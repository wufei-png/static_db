package client

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/localcache"
)

const DefaultTimeout = 10 * time.Minute

type Manager struct {
	endpoint   string
	httpClient *http.Client
}

func NewManager(endpoint string, timeout time.Duration) *Manager {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   DefaultTimeout,
			KeepAlive: DefaultTimeout,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	}
	return &Manager{
		endpoint:   endpoint,
		httpClient: &http.Client{Transport: transport, Timeout: timeout},
	}
}

func (m *Manager) checkOffline() bool {
	return m.endpoint == ""
}

func (m *Manager) DeleteModel(mp *api.ModelPath) error {
	if m.checkOffline() {
		return ErrOfflineMode
	}
	url, err := localcache.ModelPathToFilename(mp, "")
	if err != nil {
		return err
	}
	metaURL := fmt.Sprintf("%s/v1/models/%s", m.endpoint, url)
	req, _ := http.NewRequest("DELETE", metaURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("delete model failed: %v", resp.StatusCode)
}

func (m *Manager) GetModelMeta(mp *api.ModelPath) (*localcache.MetaFile, error) {
	if m.checkOffline() {
		return nil, ErrOfflineMode
	}
	url, err := localcache.ModelPathToFilename(mp, "")
	if err != nil {
		return nil, err
	}
	metaURL := fmt.Sprintf("%s/v1/models/%s", m.endpoint, url)
	resp, err := http.Get(metaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint
	if resp.StatusCode == 404 {
		return nil, ErrModelNotFound
	}
	if resp.StatusCode == 200 {
		var r api.ModelGetResponse
		if err := jsonpb.Unmarshal(resp.Body, &r); err != nil {
			return nil, err
		}
		return localcache.NewMetaFileFromBytes(r.GetRawMeta())
	}
	return nil, fmt.Errorf("get model meta failed: %v", resp.StatusCode)
}

func (m *Manager) GetModelBlob(model *api.Model) (io.ReadCloser, error) {
	if m.checkOffline() {
		return nil, ErrOfflineMode
	}
	blobURL := fmt.Sprintf("%s/v1/blobs/%s", m.endpoint, model.GetChecksum())
	resp, err := m.httpClient.Get(blobURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close() // nolint
		return nil, fmt.Errorf("get model blob failed: %v", resp.StatusCode)
	}

	return resp.Body, nil
}

func (m *Manager) ListModel() ([]*api.Model, error) {
	if m.checkOffline() {
		return nil, ErrOfflineMode
	}
	metaURL := fmt.Sprintf("%s/v1/models", m.endpoint)
	resp, err := http.Get(metaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint
	if resp.StatusCode == 200 {
		var out api.ModelListResponse
		if err := jsonpb.Unmarshal(resp.Body, &out); err != nil {
			return nil, err
		}
		return out.GetModels(), nil
	}
	return nil, fmt.Errorf("get model failed: %v", resp.StatusCode)
}

func (m *Manager) UploadModel(model *api.Model, path string, overwrite bool) error {
	if m.checkOffline() {
		return ErrOfflineMode
	}
	checksum, size, err := localcache.ChecksumFile(path)
	if err != nil {
		return &Error{"checksum", err}
	}
	log.Info(path, " checksum: ", checksum, ", size: ", size)

	model.Checksum = checksum
	model.Size = size
	if model.Oid == "" {
		model.Oid = checksum
	}
	req := &api.ModelNewRequest{
		Overwrite: overwrite,
		Model:     model,
	}
	mar := jsonpb.Marshaler{OrigName: true}
	buf := bytes.NewBuffer(nil)
	err = mar.Marshal(buf, req)
	if err != nil {
		return &Error{"marshal", err}
	}

	f, err := os.Open(path)
	if err != nil {
		return &Error{"open model", err}
	}
	defer f.Close() // nolint

	blobURL := fmt.Sprintf("%s/v1/blobs/%s", m.endpoint, checksum)
	resp, err := m.httpClient.Post(blobURL, "application/octet-stream", f)
	// TODO retry
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint
	switch resp.StatusCode {
	case http.StatusOK:
		log.Info("blob ", checksum, " uploaded")
	case http.StatusConflict:
		log.Info("blob ", checksum, " already exists")
	default:
		return fmt.Errorf("failed to upload blob: %v", resp.StatusCode)
	}

	metaURL := fmt.Sprintf("%s/v1/models", m.endpoint)
	resp, err = m.httpClient.Post(metaURL, "application/json", buf)
	// TODO retry
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint
	switch resp.StatusCode {
	case http.StatusOK:
		log.Info("model uploaded: ", model.GetModelPath())
		return nil
	case http.StatusConflict:
		log.Info("model already exists: ", model.GetModelPath())
		return ErrModelExists
	default:
		return fmt.Errorf("failed to upload meta: %v", resp.StatusCode)
	}
}

//SyncModel triggers models synchronization from minio to managers
func (m *Manager) SyncModel() error {
	if m.checkOffline() {
		return ErrOfflineMode
	}
	req := api.ModelSynchronizeRequest{}
	buf, err := marshalRequest(&req)
	if err != nil {
		return err
	}

	URL := fmt.Sprintf("%s/v1/models/synchronize", m.endpoint)
	resp, err := http.Post(URL, "application/json", buf)
	if err != nil {
		log.Error("failed to trigger models synchronization, ", err)
		return err
	}
	defer resp.Body.Close() // nolint
	if resp.StatusCode >= 200 && resp.StatusCode < 299 {
		return nil
	}

	return fmt.Errorf("failed to trigger models synchronization: %v", resp.StatusCode)
}

// GetSystemInfo returns storage capacity and the last synchronization time
func (m *Manager) GetSystemInfo() (*api.GetSystemInfoResponse, error) {
	if m.checkOffline() {
		return nil, ErrOfflineMode
	}

	URL := fmt.Sprintf("%s/v1/get_system_info", m.endpoint)
	resp, err := http.Get(URL)
	if err != nil {
		log.Error("failed to get system info, ", err)
		return nil, err
	}
	defer resp.Body.Close() // nolint
	if resp.StatusCode == 200 {
		var r api.GetSystemInfoResponse
		if err := jsonpb.Unmarshal(resp.Body, &r); err != nil {
			return nil, err
		}
		return &r, nil
	}
	return nil, fmt.Errorf("failed to get system info: %v", resp.StatusCode)
}

func marshalRequest(req proto.Message) (*bytes.Buffer, error) {
	marshaler := jsonpb.Marshaler{OrigName: true}
	buf := bytes.NewBuffer(nil)
	err := marshaler.Marshal(buf, req)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
