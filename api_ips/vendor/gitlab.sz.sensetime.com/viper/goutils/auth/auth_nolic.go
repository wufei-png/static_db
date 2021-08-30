// +build nolic

package auth

import (
	"time"

	"github.com/juju/ratelimit"
)

const (
	intUnlimitValue = 1e6
)

const (
	QuotaKeyNodeTPS = "node-tps"
	pocSignatureKey = "iva_poc_signature_key"
	testPrivateKey  = `-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDfLLIAo71dR+dmnHNIRqay2UoblDowbetOzCXpWM/FuLMk8cAT
V711BqGodvKCQ3sPodQYIDozrFJaHL3yvEYpkohQk3fB8E5wpmaTOTMAStixjHOM
Wh7De3DHuD6xDHcgdnxycL3JXAdrrvxxnLjthoWpiqQ0cWAe1aOUdDhlswIDAQAB
AoGBAIJv/RmimesLO5QdnMOBh6zAky/LUrc7h2xmpUHdRpKpgQr2wOyNO45KcfGk
k9jO1/4q92uiamIJNZi1c8+LevV52r3+7msigFGRTApv/6ryQ7fsqdiqaxv5XSDu
T/MJ+vpJhn+XaJmwKx1dXGidkL73QSwVgVz4QpNM+HNB73IBAkEA848knXocUnAu
H9u6cs1hTuK94A3nLGvenMnC1N7rgh2IHS05cnqEQsgA9pxPzfd+M01Hp+XdB6zW
erM3KCZuiQJBAOqS/04FfibhjDE5mv5c4t/YlXZG0AfRpfZwF9iRYUbjTC6H1jM/
TlfhArlFyU56qvDXDNQcUeUa7m1AxvxHg1sCQQDJRFfAnqD66uLixsPrjJbGBo9b
sIGBbt51+DC9kj1Rt6+8VJvtYxsayIYrRH6aONQb9tepAkXFyukuWhzRW/jpAkEA
kA38GbcH1OrYhHZi9+ilL2C70OoF4XdfW2tVSZtmSPlXlI+4/LnY6+D0IkF4Mejl
R4hZiX1m0bLrT07jpJRqVwJBAKFqSLQ6QImzePGX+uJ8zI0DLudQ0aSM4T8mmwjW
eDsWosFa2rngjg5AfJeyE56L92lY8Ob/LWR6aKooCyMCOJY=
-----END RSA PRIVATE KEY-----,1766160000`
)

var (
	tokens = []byte(
		`eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJjb20uc2Vuc2V0aW1lLnZpcGVyLnNlbmFyaW8udnBzLmV4cCIsImNobiI6NjQsImV4cCI6MTU4NjI1Njg4NSwiaWRzIjpbIjk4NzY1MTIzNDAiLCIxIl0sIm5iZiI6MTU4NDUyODg4NX0.SqcFMw3yGw57eACtUqbrkpB7lqVEKElyFqgpCJ8pDgOHENHO3ydFVh2M5jrbVFwLY4LAj4phW5qIYrSQDKhafI-yLLxpVDc2AhT-guSnCkamdN92fKUjwY-bZslqPy3vgMVgMJm8CfXXh7kePBYsD5jqDV8m218RFtWhMBIs5lY,eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJjb20uc2Vuc2V0aW1lLnZpcGVyLnNlbmFyaW8uaXBzLmV4cCIsImV4cCI6MTU4NjI1Njg4NSwiaWRzIjpbIjk4NzY1MTIzNDAiLCIxIl0sIm5iZiI6MTU4NDUyODg4NSwicXBzIjoxNjB9.WAxzkg1T_8VA2Our8wcC3rHvDCPCPLQJgi1HelniJPkBbtnFnn036Ebg8ydX-OLQmyWILg7KZ1B1x-cGOw7Evt6dzbAh0Ho29mheK2n-3ZnsEKjbzmjJMIANZxpRK9Mpk98nrmGG6FUDZc0n05k1bnaUhl8IL-qhJSM2Lk7lFRk,eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJjb20uc2Vuc2V0aW1lLnZpcGVyLnNlbmFyaW8uaXBzIiwiZXhwIjoxNjE3NzkyODg1LCJpZHMiOlsiOTg3NjUxMjM0MCIsIjEiXSwibmJmIjoxNTg0NTI4ODg1LCJxcHMiOjE2MH0.n0Xma_XRG5JY2QD108B2X3SKV8QGNRUQ3t220LQzIVEnG8iWo__vzHPfPcppB1Tjlcvya15sYHvsoFlus1_S0FjbrhAUJPj3uK3Cwv7KyGlv2my3cQgeUkFSa4REwN8_K0xg4bAT0vyOtdsDZov02zIphu2ldMKemafLbvvCyVk,eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJjb20uc2Vuc2V0aW1lLnZpcGVyLnNlbmFyaW8udnBzIiwiY2huIjo2NCwiZXhwIjoxNjE1OTc4NDg1LCJpZHMiOlsiOTg3NjUxMjM0MCIsImRvbmdsZV8xMjMiXSwibmJmIjoxNTg0NTI4ODg0fQ.HTFtWCJwax5-LdbQkZuRz0QHT3L_PPt4zYeibr1vZ393YHajL-dTmZ6SWSIKhQO8GYR5nwobZ1RlSgjD7PlG_JD3kFa8bN_JUmuir7-jUvo3L8GGgfNfw1WQfvUP3MMu4jqzHLAt4g9JGwAS7G94GtNEuFOuCzpfaOwN9ThqXD0`)
)

func SetHeartbeatOption(interval, timeout int32) {
}

// Authorizor is mock impl of license CA
type Authorizor struct {
	componentName string
	quotaStatus   chan QuotaStatus
}

func NewAuthorizorFromFile(licensePath, componentName, clientID, product, uuid string, req CARequest) (*Authorizor, error) {
	return &Authorizor{
		componentName: componentName,
	}, nil
}

func NewAuthorizor(licenseContent []byte, componentName, clientID, product, uuid string, req CARequest) (*Authorizor, error) {
	return &Authorizor{
		componentName: componentName,
		quotaStatus:   make(chan QuotaStatus, 1),
	}, nil
}

func (a *Authorizor) IsCAPrivate() bool {
	return false
}

func (a *Authorizor) ComponentName() string {
	return a.componentName
}

func (a *Authorizor) CheckCapability(n string) (bool, error) {
	return true, nil
}

func (a *Authorizor) GetLimit(n string) interface{} {
	return nil
}

func (a *Authorizor) CounterCheckIn(n string) error {
	return nil
}

func (a *Authorizor) CounterCheckOut(n string) error {
	return nil
}

func (a *Authorizor) GetEncryptionKey() []byte {
	return nil
}

func (a *Authorizor) IsStatusFatal(status QuotaStatus) bool {
	return false
}

func (a *Authorizor) Status() <-chan QuotaStatus {
	return a.quotaStatus
}

func (a *Authorizor) Close() {
}

func (a *Authorizor) WaitQuotaRenew() error {
	return nil
}

func (a *Authorizor) GetQuotaLimit(name string) (int, error) {
	return intUnlimitValue, nil
}

func (a *Authorizor) NewRateLimiter(name string) (*ratelimit.Bucket, int, error) {
	return ratelimit.NewBucket(time.Nanosecond, int64(intUnlimitValue)), intUnlimitValue, nil
}

func (a *Authorizor) NewReloadableLimiter(name string) (*ReloadableLimiter, int, error) {
	return NewReloadableLimiter(int64(intUnlimitValue)), intUnlimitValue, nil
}

func (a Authorizor) GetConstValueByKey(constKey string) (int, bool) {
	return intUnlimitValue, true
}

func (a Authorizor) GetConstStringByKey(constKey string) (string, bool) {
	if constKey == pocSignatureKey {
		return testPrivateKey, true
	}
	return "", true
}

// GetDongleInfo ...
func (a *Authorizor) GetDongleInfo() (DongleInfo, bool) {
	return DongleInfo{HardwareID: "9876512340", HardwareTime: time.Now().Unix()}, true
}

// GetExtraInfo ...
func (a *Authorizor) GetExtraInfo() ([]byte, error) {
	return tokens, nil
}
