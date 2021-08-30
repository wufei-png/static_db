// +build hal_cuda11 hal_cuda10 hal_cuda8

package cuda

import (
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/gpu/cudaruntime"
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/monitor"
)

type cudaDeviceMonitorImpl struct{}

func (c cudaDeviceMonitorImpl) Init() error {
	return nil
}

func (c cudaDeviceMonitorImpl) BindDevice(id int) error {
	return cudaruntime.SetDevice(id)
}

func (c cudaDeviceMonitorImpl) GetDeviceCount() (int, error) {
	return cudaruntime.DeviceGetCount()

}

func (c cudaDeviceMonitorImpl) GetMemInfo(id int) (all int, used int, err error) {
	return cudaruntime.GetMemInfo()
}

func (c cudaDeviceMonitorImpl) GetSDKProperties() (*monitor.SDKProperties, error) {
	version, err := cudaruntime.GetCUDAVersion()
	if err != nil {
		return nil, err
	}

	return &monitor.SDKProperties{
		SDKMajor: version,
	}, nil
}

func (c cudaDeviceMonitorImpl) GetDeviceProperties(id int) (*monitor.DeviceProperties, error) {
	devProp, err := cudaruntime.GetCUDADeviceProperties(id)
	if err != nil {
		return nil, err
	}

	return &monitor.DeviceProperties{
		HardwareMajor: devProp.Major,
		HardwareMinor: devProp.Minor,
		DeviceName:    devProp.Name,
	}, nil
}

func (c cudaDeviceMonitorImpl) UnbindDevice(id int) {
}

func (c cudaDeviceMonitorImpl) UnInit() error {
	return nil
}

func init() {
	cm := &cudaDeviceMonitorImpl{}
	if err := cm.Init(); err != nil {
		panic(err)
	}
	monitor.DeviceMonitors[monitor.CUDA] = cm
	monitor.DefaultMonitor = cm
	monitor.DefaultDevice = monitor.CUDA
}
