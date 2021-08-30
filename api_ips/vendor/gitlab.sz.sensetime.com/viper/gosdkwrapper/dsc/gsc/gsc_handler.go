// +build gsc

package gsc

import (
	"errors"
	"runtime"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/dsc"
)

var AvailableIndexTypeList = map[dsc.IndexType]bool{
	dsc.IndexIVFPQ: true,
}

const (
	// 100MB default temp device mem for gpu faiss
	defaultTempGPUMemorySize = 100 << 20
)

type GSCIndexHandlerBuilder struct{}

type GSCIndexHandler struct {
	gpuDevice   *GpuDevice
	gpuResource *GpuResource
	isCPUModel  bool
	temMemory   uint64
}

func (b *GSCIndexHandlerBuilder) Build(config interface{}) (dsc.IndexHandler, error) {
	c, ok := config.(*dsc.GeneralIndexHandlerConfig)
	if !ok {
		return nil, errors.New("unknown gsc hanlder builder config")
	}
	if !c.IsCPUModel {
		device, err := GetGpuDevice()
		if err != nil {
			return nil, err
		}
		return NewGSCIndexHandler(c, device)
	}
	return NewGSCIndexHandler(c, nil)
}

func (b *GSCIndexHandlerBuilder) InitEnv(config interface{}) error {
	return nil
}

func (b *GSCIndexHandlerBuilder) DestroyEnv() error {
	return nil
}

func NewGSCIndexHandler(config *dsc.GeneralIndexHandlerConfig, device *GpuDevice) (*GSCIndexHandler, error) {
	if _, ok := AvailableIndexTypeList[config.IndexType]; !ok {
		return nil, errors.New("Unsupport index type")
	}
	if !config.IsCPUModel {
		if device == nil {
			return nil, errors.New("Device not init")
		}
		if config.DeviceTemMemory <= 0 {
			config.DeviceTemMemory = defaultTempGPUMemorySize
		}
	}
	return &GSCIndexHandler{gpuDevice: device, isCPUModel: config.IsCPUModel, temMemory: config.DeviceTemMemory}, nil
}

func (handle *GSCIndexHandler) BindDevice(deviceID int32) error {
	if !handle.isCPUModel {
		runtime.LockOSThread()
		gpuResource, err := handle.gpuDevice.CreateGpuResource(handle.temMemory)
		if err != nil {
			return errors.New("CreateGpuResource fail, err: " + err.Error())
		}
		handle.gpuResource = gpuResource
	}
	return nil
}

func (handle *GSCIndexHandler) UnbindDevice() error {
	if handle.isCPUModel || handle.gpuResource == nil {
		return nil
	}
	defer runtime.UnlockOSThread()
	return handle.gpuResource.DestroyGpuResource()
}

func (handle *GSCIndexHandler) InitIndex(indexConfig interface{}) (dsc.SearchIndex, error) {
	c, ok := indexConfig.(*dsc.GeneralIndexConfig)
	if !ok {
		return nil, errors.New("unknown index config for gsc")
	}
	var gscIndex dsc.SearchIndex
	var err error
	if handle.isCPUModel {
		gscIndex, err = InitCpuIndex(c)
	} else {
		if handle.gpuResource == nil {
			return nil, errors.New("InitIndex fail, gpuResource do not init")
		}
		gscIndex, err = handle.gpuResource.InitGpuIndex(c)
	}
	return gscIndex, err
}

func (handle *GSCIndexHandler) LoadIndex(filepath string) (dsc.SearchIndex, error) {
	var gscIndex dsc.SearchIndex
	var err error
	if handle.isCPUModel {
		gscIndex, err = LoadCpuIndex(filepath)
	} else {
		if handle.gpuResource == nil {
			return nil, errors.New("LoadIndex fail, gpuResource do not init")
		}
		gscIndex, err = handle.gpuResource.LoadGpuIndex(filepath)
	}
	return gscIndex, err
}

func init() {
	//register
	dsc.IndexHandlerFactory[dsc.GSC] = &GSCIndexHandlerBuilder{}
}
