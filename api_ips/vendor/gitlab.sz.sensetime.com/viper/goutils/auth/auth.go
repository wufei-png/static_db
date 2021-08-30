// +build !nolic

package auth

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	// XXX

	"github.com/juju/ratelimit"
	log "github.com/sirupsen/logrus"

	lh "gitlab.bj.sensetime.com/sys-dev/sdk-license-go-plugin/apiv1"
	lhl "gitlab.bj.sensetime.com/sys-dev/sdk-license-go-plugin/loader"
)

const (
	featureEncryptionKey = "iva_feature_encryption_key"
	QuotaKeyNodeTPS      = "node-tps"
	licenseAct           = "ca_private"

	activationCode = "activation_code"

	okCode    = "000000"
	fatalCode = "000001"
)

var (
	defaultConstNames = []string{
		featureEncryptionKey,
	}
	// should keep sync with licensesdk
	defaultHeartBeatInterval int32 = 28 // sec
	defaultHeartBeatTimeout  int32 = 13 // http timeout, sec, defaultHeartBeatTimeout * 2 < defaultHeartBeatInterval
)

// SetHeartbeatOption ...
func SetHeartbeatOption(interval, timeout int32) {
	if timeout*2 > interval {
		log.Error("license: invalid heartbeat option")
		return
	}
	if interval < 15 {
		defaultHeartBeatInterval = 15
	}
	if timeout < 5 {
		defaultHeartBeatTimeout = 5
	}
	defaultHeartBeatInterval = interval
	defaultHeartBeatTimeout = timeout
	log.Infof("license: heartbeat interval : [%d]s,heartBeat timeout : [%d]s", defaultHeartBeatInterval, defaultHeartBeatTimeout)
}

type limiterRegistry struct {
	mu sync.Mutex

	m map[string][]*ReloadableLimiter
}

func newLimiterRegistry() *limiterRegistry {
	return &limiterRegistry{
		m: make(map[string][]*ReloadableLimiter, 8),
	}
}

func (l *limiterRegistry) Add(name string, r *ReloadableLimiter) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if rl, ok := l.m[name]; ok {
		l.m[name] = append(rl, r)
	} else {
		l.m[name] = []*ReloadableLimiter{r}
	}
}

func (l *limiterRegistry) RangeUpdate(f func(name string) int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for k, v := range l.m {
		newRate := f(k)
		if newRate <= 0 {
			log.Warn("license: invalid ratelimit on update: ", k)
			continue
		}
		for i := range v {
			if v[i].Reset(int64(newRate)) {
				log.Info("license: ratelimit updated: ", k)
			}
		}
	}
}

type Authorizor struct {
	componentName string
	productName   string
	uuid          string

	cachedDongleInfo atomic.Value

	// monitor file change
	licChecksum string
	licPath     string

	quotaManager lh.QuotaManagerAPI // nolint

	curVersion int32
	status     chan QuotaStatus
	doneCh     chan struct{}

	caReq CARequest

	constValuePair atomic.Value
	limiters       *limiterRegistry
}

type constsPairMap map[string]interface{}

func nodeQuotaKey(sn string) string {
	return sn + "-nodes"
}

// nolint
func convertStatus(s lh.QuotaStatus) QuotaStatus {
	limits := make(map[string]QuotaLimit, len(s.Limits))
	for k, v := range s.Limits {
		limits[k] = QuotaLimit{
			Current: v.Current,
			Free:    v.Free,
		}
	}
	constValuePair := make(map[string]interface{}, len(s.Consts))
	for k, v := range s.Consts {
		constValuePair[k] = v
	}
	return QuotaStatus{
		Err: Error{
			Code:    s.Err.Code,
			Message: s.Err.Message,
		},
		Version: s.Version,
		Limits:  limits,
		Consts:  constValuePair,
	}
}

// NewAuthorizorFromFile used to build an auth client to fetch auth info from license-ca server
// licensePath: local license file path
// componentName: service node name,like engine-timespace-feature-db,engine-static-feature-db...
// resourceCost: quota cost per-service-node,could be [0,1<<31 - 1)
func NewAuthorizorFromFile(licensePath, componentName, clientID, product, uuid string, req CARequest) (*Authorizor, error) {
	absPath, err := filepath.Abs(licensePath)
	if err != nil {
		return nil, err
	}
	licenseContent, err := ioutil.ReadFile(absPath) // #nosec
	if err != nil {
		return nil, err
	}
	return newAuthorizorWithPath(licenseContent, absPath, componentName, clientID, product, uuid, req)
}

func NewAuthorizor(licenseContent []byte, componentName, clientID, product, uuid string, req CARequest) (*Authorizor, error) {
	return newAuthorizorWithPath(licenseContent, "", componentName, clientID, product, uuid, req)
}

func newAuthorizorWithPath(licenseContent []byte, path string, componentName, clientID, product, uuid string, req CARequest) (*Authorizor, error) {
	// nolint
	if ver, err := lhl.GetSDKAPIVersion(); err != nil || ver < 4 {
		panic("incompatiable SDK version or bad license lib")
	}
	if componentName == "" {
		log.Fatal("license: empty service name")
	}
	if req.ResourceCost < 0 {
		log.Fatal("license: error params resourceCost,should be greater than 0 but pramas is ", req.ResourceCost)
	}

	// deep copy
	newReq := CARequest{
		ResourceCost: req.ResourceCost,
		ConstNames:   make([]string, 0, len(req.ConstNames)+len(defaultConstNames)),
	}
	newReq.ConstNames = append(newReq.ConstNames, req.ConstNames...)
	newReq.ConstNames = append(newReq.ConstNames, defaultConstNames...)
	newReq.Capabilities = append(newReq.Capabilities, req.Capabilities...)

	log.Infof("license: init authorizor with: service=%s, clientID=%s, product=%s", componentName, clientID, product)
	handler, err := lhl.NewLicenseHandlerAPI(licenseContent, product, uuid) // nolint
	if err != nil {
		return nil, err
	}

	if !handler.CheckExpiration() { // nolint
		return nil, ErrLicenseExpired
	}

	if !handler.CheckActivationType(licenseAct) { // nolint
		log.Info("license: standalone mode")
	}

	quotaManager, err := lhl.NewQuotaManagerAPI(clientID, handler, defaultHeartBeatInterval, defaultHeartBeatTimeout) // nolint
	if err != nil {
		return nil, err
	}

	auth := &Authorizor{
		componentName: componentName,
		productName:   product,
		uuid:          uuid,

		licPath:     path,
		licChecksum: checksumLic(licenseContent),

		caReq:        newReq,
		quotaManager: quotaManager,
		status:       make(chan QuotaStatus, 5),
		doneCh:       make(chan struct{}),

		limiters: newLimiterRegistry(),
	}
	err = auth.init()
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func (a *Authorizor) init() error {
	// quotas:        make(map[string]int32),
	quotaReq := make(map[string]int32, 1)
	if a.caReq.ResourceCost > 0 {
		quotaReq[nodeQuotaKey(a.componentName)] = a.caReq.ResourceCost
	}

	a.curVersion = a.quotaManager.UpdateRequestQuota(quotaReq, a.caReq.ConstNames, a.caReq.Capabilities) // nolint

	// initial quota check already done in Run()
	qs, err := a.quotaManager.Run() // nolint
	if err != nil {
		log.Error("license: failed to init CA license: ", err)
		return err
	}

	// load initial consts
	a.loadConstsFromStatus(convertStatus(qs))

	go a.runForward()
	return nil
}

func (a *Authorizor) detectLicenseUpdate() {
	if a.licPath == "" {
		return
	}

	licenseContent, err := ioutil.ReadFile(a.licPath)
	if err != nil {
		log.Debug("license: failed to reload license: ", err)
		return
	}
	newChecksum := checksumLic(licenseContent)
	if a.licChecksum == newChecksum {
		return
	}

	handler, err := lhl.NewLicenseHandlerAPI(licenseContent, a.productName, a.uuid) // nolint
	if err != nil {
		log.Warn("license: failed to load new license: ", err)
		return
	}

	if !handler.CheckExpiration() { // nolint
		log.Warn("license: new license expired")
		return
	}

	if !handler.CheckActivationType(licenseAct) { // nolint
		log.Info("license: reloaded with standalone mode")
	}

	// load it!
	a.quotaManager.UpdateLicenseHandler(handler) // nolint
	a.licChecksum = newChecksum
	log.Info("license: new license loaded, checksum: ", newChecksum)
}

// nolint
func (a *Authorizor) getLicenseHandler() lh.LicenseHandlerAPI {
	return a.quotaManager.GetLicenseHandler()
}

func (a *Authorizor) loadConstsFromStatus(val QuotaStatus) {
	if val.Err.Code != okCode {
		return
	}
	newM := make(constsPairMap)
	for key, value := range val.Consts {
		switch value.(type) {
		case string:
			if key == featureEncryptionKey {
				encryptionKey, err := checkEncryptionKey(value)
				if err != nil {
					log.Fatal("fetch encryptionKey failed: ", err)
				} else if len(encryptionKey) == 0 {
					newM[key] = value
				} else if len(encryptionKey) == 16 {
					newM[key] = encryptionKey
				}
			} else {
				newM[key] = value
			}

		case int32:
			newM[key] = value
		}
	}

	a.constValuePair.Store(newM)
	a.limiters.RangeUpdate(func(name string) int {
		l, _ := a.GetQuotaLimit(name)
		return l
	})
}

func (a *Authorizor) runForward() {
	defer close(a.doneCh)
	defer log.Info("license: auth loop done")
	defer close(a.status)

	ticker := time.NewTicker(time.Duration(defaultHeartBeatInterval/2+1) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case v, ok := <-a.quotaManager.Status(): // nolint
			if !ok {
				return
			}
			val := convertStatus(v)
			var d DongleInfo
			d.HardwareID = v.DongleID
			d.HardwareTime = v.DongleTime
			a.cachedDongleInfo.Store(d)
			a.loadConstsFromStatus(val)
			select {
			case a.status <- val:
			default:
				<-a.status
				a.status <- val
			}
		case <-ticker.C:
			a.detectLicenseUpdate()
		}
	}
}

func (a *Authorizor) IsCAPrivate() bool {
	return a.quotaManager != nil
}

func (a *Authorizor) ComponentName() string {
	return a.componentName
}

func (a *Authorizor) GetLimit(n string) interface{} {
	return a.getLicenseHandler().GetLimit(n) // nolint
}

func (a *Authorizor) getCurrentConstPair() constsPairMap {
	v := a.constValuePair.Load()
	if v == nil {
		return nil
	}
	return v.(constsPairMap)
}

// nolint
func (a *Authorizor) GetEncryptionKey() []byte {
	// in dongle mode, key saved in consts
	// in soft mode, key saved in client lic

	//remote const value first
	cvp := a.getCurrentConstPair()
	if value, ok := cvp[featureEncryptionKey]; ok && value != nil {
		return value.([]byte)
	}

	handler := a.getLicenseHandler()
	//then local file
	val := handler.GetLimit(featureEncryptionKey) // nolint
	if val == nil {
		val = handler.GetLimit("const_" + featureEncryptionKey)
	}
	if val != nil {
		log.Println("license: overriding with local license file")
		encryptionKey, err := checkEncryptionKey(val)
		if err != nil {
			log.Fatal("license: bad feature key")
		}
		return encryptionKey
	}

	log.Fatal("license: invalid license key")
	return nil
}

func (a *Authorizor) IsStatusFatal(status QuotaStatus) bool { // nolint
	return status.Err.Code == fatalCode
}

func (a *Authorizor) Status() <-chan QuotaStatus { // nolint
	return a.status
}

func (a *Authorizor) Close() {
	if !a.IsCAPrivate() {
		return
	}
	log.Info("license: waiting quota manager close")
	a.quotaManager.Close() // nolint
	<-a.doneCh
}

// WaitQuotaRenew Deprecated: Will Not support, do not use.
// nolint
func (a *Authorizor) WaitQuotaRenew() error {
	if !a.IsCAPrivate() {
		return nil
	}
	for r := range a.status { // nolint
		if r.Err.Code != okCode {
			return ErrQuotaRenew
		}
		if r.Version >= a.curVersion {
			return nil
		}
	}
	return ErrQuotaRenew
}

func convertNumber(val interface{}) (int, bool) {
	//licensesdk use json format
	switch tv := val.(type) {
	case int:
		return tv, true
	case int32:
		return int(tv), true
	case int64:
		return int(tv), true
	case float32:
		return int(tv), true
	case float64:
		return int(tv), true
	case string:
		if v, err := strconv.Atoi(tv); err == nil {
			return v, true
		}
	}
	return 0, false
}

// GetQuotaLimit local config file is first choice,then is remote license-ca server's value
func (a *Authorizor) GetQuotaLimit(name string) (int, error) {
	key := a.componentName + "-" + name
	val, ok := a.GetConstValueByKey(key)
	if !ok {
		return 0, ErrBadQuotaLimit // nolint
	}
	return val, nil
}

// NewRateLimiter Deprecated, use NewReloadableLimiter instead
func (a *Authorizor) NewRateLimiter(name string) (*ratelimit.Bucket, int, error) {
	r, err := a.GetQuotaLimit(name)
	if err != nil {
		return nil, 0, err
	}
	if r <= 0 {
		return nil, 0, ErrBadQuotaLimit
	}
	interval := time.Second / time.Duration(r)
	if interval <= 0 {
		interval = time.Nanosecond
	}
	return ratelimit.NewBucket(interval, int64(r)), r, nil
}

func (a *Authorizor) NewReloadableLimiter(name string) (*ReloadableLimiter, int, error) {
	r, err := a.GetQuotaLimit(name)
	if err != nil {
		return nil, 0, err
	}
	if r <= 0 {
		return nil, 0, ErrBadQuotaLimit
	}
	limiter := NewReloadableLimiter(int64(r))
	a.limiters.Add(name, limiter)
	return limiter, r, nil
}

//GetConstValueByKey used to get const value not regular
func (a *Authorizor) GetConstValueByKey(constKey string) (int, bool) {
	handler := a.getLicenseHandler()
	//fetch const value from local lic file first
	//nolint
	if value := handler.GetLimit(constKey); value != nil {
		log.Infof("license: get const value [%s] without const prefix from local lic file success", constKey)
		return convertNumber(value)
	}
	//nolint
	if value := handler.GetLimit("const_" + constKey); value != nil {
		log.Infof("license: get const value [const_%s] with const prefix from local lic file success", constKey)
		return convertNumber(value)
	}

	//if constKey is fetched before,then return
	cvp := a.getCurrentConstPair()
	if value, ok := cvp[constKey]; ok && value != nil {
		return convertNumber(value)
	}
	return 0, false
}

// GetExtraInfo ...
func (a *Authorizor) GetExtraInfo() ([]byte, error) {
	extraInfo, err := a.quotaManager.GetExtraInfo(activationCode) // nolint
	if err != nil {
		return nil, err
	}
	return extraInfo, nil
}

//GetConstStringByKey used to get const string not regular
func (a *Authorizor) GetConstStringByKey(constKey string) (string, bool) {
	handler := a.getLicenseHandler()
	//fetch const value from local lic file first
	//nolint
	if value := handler.GetLimit(constKey); value != nil {
		return convertString(value)
	}
	//nolint
	if value := handler.GetLimit("const_" + constKey); value != nil {
		return convertString(value)
	}
	//if constKey is fetched before,then return
	cvp := a.getCurrentConstPair()
	if value, ok := cvp[constKey]; ok && value != nil {
		return convertString(value)
	}
	return "", false
}

func convertString(val interface{}) (string, bool) {
	switch tv := val.(type) {
	case string:
		return tv, true
	}
	return "", false
}

// GetDongleInfo ...
func (a *Authorizor) GetDongleInfo() (DongleInfo, bool) {
	if val, ok := a.cachedDongleInfo.Load().(DongleInfo); ok {
		return val, true
	}
	return DongleInfo{}, false
}
