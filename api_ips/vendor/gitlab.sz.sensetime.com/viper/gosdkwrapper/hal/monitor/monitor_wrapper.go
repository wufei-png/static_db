package monitor

import (
	"fmt"
	"runtime"
	"time"
)

const requestEnqueueTimeout = 10 * time.Second

type DeviceHandler struct {
	deviceID            int
	m                   DeviceMonitor
	isSingleThreadModel bool
	requestChan         chan *generalRequest
}

type generalRequest struct {
	request  interface{}
	response chan *generalResponse
}

type generalResponse struct {
	response interface{}
	err      error
}

type getDeviceCountRequest struct{}
type getDeviceCountResponse struct {
	deviceNum int
}

type getMemInfoRequest struct{}
type getMemInfoReponse struct {
	free  int
	total int
}

type getSDKPropertiesRequest struct{}
type getSDKPropertiesReponse struct {
	sdkProperties *SDKProperties
}

type getDevicePropertiesRequest struct{}
type getDevicePropertiesReponse struct {
	deviceProperties *DeviceProperties
}

// NewDeviceHandler init a monitor wrapper, if isSingleThreadModel set to true, handler will init a locked thread to read device info, and handler should invoke Close() to exit safety
func NewDeviceHandler(deviceID int, deviceName string, isForceCPUModel bool, isSingleThreadModel bool) (*DeviceHandler, error) {
	d := DeviceHandler{deviceID: deviceID, isSingleThreadModel: isSingleThreadModel}
	if isForceCPUModel {
		d.m = &CPUMonitor{}
	} else {
		monitor, ok := DeviceMonitors[deviceName]
		if !ok {
			return nil, fmt.Errorf("no device monitor %v", deviceName)
		}
		d.m = monitor
	}

	if isSingleThreadModel {
		err := d.threadInit()
		if err != nil {
			return nil, err
		}
	}
	return &d, nil
}

func (d *DeviceHandler) threadInit() error {
	d.requestChan = make(chan *generalRequest, 1)
	isReady := make(chan error, 1)
	go d.readDeviceLoop(isReady)
	err := <-isReady
	if err != nil {
		return err
	}
	return nil
}

func (d *DeviceHandler) readDeviceLoop(isReady chan<- error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := d.BindDevice(d.deviceID); err != nil {
		isReady <- err
		return
	}
	defer d.UnBindDevice(d.deviceID)
	isReady <- nil
	for req := range d.requestChan {
		//TODO: add timeout handle
		res, err := d.handleRequest(req)
		req.response <- &generalResponse{response: res, err: err}
	}
}

func (d *DeviceHandler) sendRequest(request interface{}) (interface{}, error) {
	responseChannel := make(chan *generalResponse, 1)
	req := &generalRequest{
		request:  request,
		response: responseChannel,
	}
	tick := time.NewTimer(requestEnqueueTimeout)
	select {
	case d.requestChan <- req:
	case <-tick.C:
		return nil, fmt.Errorf("backend busy")
	}
	tick.Stop()
	generalResponse := <-responseChannel
	if generalResponse.err != nil {
		return nil, generalResponse.err
	}
	return generalResponse.response, nil
}

func (d *DeviceHandler) handleRequest(generalRequest *generalRequest) (interface{}, error) {
	if generalRequest == nil {
		return nil, fmt.Errorf("invalid request")
	}
	switch generalRequest.request.(type) {
	case *getDeviceCountRequest:
		n, err := d.m.GetDeviceCount()
		return &getDeviceCountResponse{deviceNum: n}, err
	case *getMemInfoRequest:
		f, t, err := d.m.GetMemInfo(d.deviceID)
		return &getMemInfoReponse{free: f, total: t}, err
	case *getSDKPropertiesRequest:
		p, err := d.m.GetSDKProperties()
		return &getSDKPropertiesReponse{sdkProperties: p}, err
	case *getDevicePropertiesRequest:
		p, err := d.m.GetDeviceProperties(d.deviceID)
		return &getDevicePropertiesReponse{deviceProperties: p}, err
	default:
		return nil, fmt.Errorf("unknown request type")
	}
}

func (d *DeviceHandler) BindDevice(id int) error {
	return d.m.BindDevice(id)
}

func (d *DeviceHandler) UnBindDevice(id int) {
	d.m.UnbindDevice(id)
}

func (d *DeviceHandler) GetDeviceCount() (int, error) {
	if d.isSingleThreadModel {
		res, err := d.sendRequest(&getDeviceCountRequest{})
		if err != nil {
			return 0, err
		}
		r, ok := res.(*getDeviceCountResponse)
		if !ok {
			return 0, fmt.Errorf("incorrect response")
		}
		return r.deviceNum, nil
	}
	return d.m.GetDeviceCount()
}
func (d *DeviceHandler) GetMemInfo(id int) (free int, total int, err error) {
	if d.isSingleThreadModel {
		res, err := d.sendRequest(&getMemInfoRequest{})
		if err != nil {
			return 0, 0, err
		}
		r, ok := res.(*getMemInfoReponse)
		if !ok {
			return 0, 0, fmt.Errorf("incorrect response")
		}
		return r.free, r.total, nil
	}
	return d.m.GetMemInfo(id)
}

func (d *DeviceHandler) GetSDKProperties() (*SDKProperties, error) {
	if d.isSingleThreadModel {
		res, err := d.sendRequest(&getSDKPropertiesRequest{})
		if err != nil {
			return nil, err
		}
		r, ok := res.(*getSDKPropertiesReponse)
		if !ok {
			return nil, fmt.Errorf("incorrect response")
		}
		return r.sdkProperties, nil
	}
	return d.m.GetSDKProperties()
}

func (d *DeviceHandler) GetDeviceProperties(id int) (*DeviceProperties, error) {
	if d.isSingleThreadModel {
		res, err := d.sendRequest(&getDevicePropertiesRequest{})
		if err != nil {
			return nil, err
		}
		r, ok := res.(*getDevicePropertiesReponse)
		if !ok {
			return nil, fmt.Errorf("incorrect response")
		}
		return r.deviceProperties, nil
	}
	return d.m.GetDeviceProperties(id)
}

func (d *DeviceHandler) Close() error {
	if d.isSingleThreadModel {
		//XXX should we check double close?
		close(d.requestChan)
	}

	return d.m.UnInit()
}
