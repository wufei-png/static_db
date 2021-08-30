package monitor

import (
	"fmt"

	"gitlab.sz.sensetime.com/viper/commonapis/define/device"
)

const (
	CPUDEVICE = "cpu"
	HIDEVICE  = "hidevice"
	CUDA      = "cuda"
	STPU      = "stpu_device"
	ATLAS     = "acldevice"
	CAMBRICON = "cambricon"
)

const (
	OneMB = 1024 * 1024
)

var DeviceMonitors = map[string]DeviceMonitor{CPUDEVICE: &CPUMonitor{}}

var DefaultDevice = CPUDEVICE

var DefaultMonitor = DeviceMonitor(&CPUMonitor{})

type SDKProperties struct {
	SDKMajor int
	SDKMinor int
}

type DeviceProperties struct {
	HardwareMajor int
	HardwareMinor int
	DeviceName    string
}

type DeviceMonitor interface {
	Init() error
	BindDevice(id int) error
	GetDeviceCount() (int, error)
	// GetMemInfo return mem info in bytes
	GetMemInfo(id int) (free int, total int, err error)
	GetSDKProperties() (*SDKProperties, error)
	GetDeviceProperties(id int) (*DeviceProperties, error)
	UnbindDevice(id int)
	UnInit() error
}

const (
	P4 = "Tesla P4"
	// 部分驱动下(如367.48), P4的name为"Graphics Device"
	GraphicsDevice     = "Graphics Device"
	T4                 = "Tesla T4"
	RTX4000            = "Quadro RTX 4000"
	GTX1060            = "GeForce GTX 1060"
	GTX10606GB         = "GeForce GTX 1060 6GB"
	GTX1660            = "GeForce GTX 1660"
	GTX1070            = "GeForce GTX 1070"
	GTX2080            = "GeForce RTX 2080"
	NvidiaGTX1080      = "GeForce GTX 1080"
	NvidiaGTX1080Ti    = "GeForce GTX 1080 Ti"
	NvidiaRTX2070Super = "GeForce RTX 2070 SUPER"
	HI3559A            = "hi_3559a" // device.HardwareHI3559A

	StpuDevice       = "stpu_device"
	StpuDeviceSV114A = "SV-114A"
	StpuDeviceSV116E = "SV-116E"

	AscendDevice310 = "Ascend310"

	MLU270     = "MLU270-S4"
	MLU220Edge = "MLU220-EDGE"
)

const (
	CUDA10Version = 10000
	CUDA11Version = 11000
)

type RuntimeInfo struct {
	Runtime  string
	Hardware string
}

func GetDefaultDeviceName() string {
	probeList := []string{CUDA, HIDEVICE, ATLAS, STPU}
	for _, dev := range probeList {
		v, ok := DeviceMonitors[dev]
		if !ok {
			continue
		}
		if c, err := v.GetDeviceCount(); err == nil && c > 0 {
			return dev
		}
	}
	return ""
}

var (
	DeviceName2CommonDeviceName = map[string]string{
		GraphicsDevice: device.HardwareNVP4,

		GTX1060:         device.HardwareNVP4,
		GTX10606GB:      device.HardwareNVP4,
		GTX1660:         device.HardwareNVT4,
		GTX1070:         device.HardwareNVP4,
		NvidiaGTX1080:   device.HardwareNVP4,
		NvidiaGTX1080Ti: device.HardwareNVP4,

		NvidiaRTX2070Super: device.HardwareNVT4,
		GTX2080:            device.HardwareNVT4,

		P4:               device.HardwareNVP4,
		T4:               device.HardwareNVT4,
		RTX4000:          device.HardwareNVT4,
		ATLAS:            device.HardwareHWAtlas300,
		AscendDevice310:  device.HardwareHWAtlas300,
		STPU:             device.HardwareST,
		StpuDeviceSV114A: device.HardwareST,
		StpuDeviceSV116E: device.HardwareST,
		MLU220Edge:       device.HardwareMLU220Edge,
		MLU270:           device.HardwareMLU270,
	}
)

func GetRuntimeAndHardware(deviceName string) (RuntimeInfo, error) {
	const defaultInfo = "default"
	info := RuntimeInfo{
		Runtime:  defaultInfo,
		Hardware: defaultInfo,
	}

	// 认为在deviceName为空的情况下为CPU模式, 其模型的Runtime与Hardware都为default
	if len(deviceName) == 0 {
		return info, nil
	}

	v, ok := DeviceMonitors[deviceName]
	if !ok {
		return info, fmt.Errorf("unknown device: %s", deviceName)
	}

	switch deviceName {
	case CUDA:
		sprop, err := v.GetSDKProperties()
		if err != nil {
			return info, fmt.Errorf("get %s version err: %v", deviceName, err)
		}
		if sprop.SDKMajor >= CUDA11Version {
			info.Runtime = device.RuntimeTRT7
		} else if sprop.SDKMajor >= CUDA10Version {
			info.Runtime = device.RuntimeTRT5

		} else {
			info.Runtime = device.RuntimeTRT2
		}

		dprop, err := v.GetDeviceProperties(0)
		if err != nil {
			return info, fmt.Errorf("get %s device property err: %v", deviceName, err)
		}
		// t4
		if dprop.HardwareMajor >= 7 {
			info.Hardware = device.HardwareNVT4
		} else {
			info.Hardware = device.HardwareNVP4
		}

	case HIDEVICE:
		sprop, err := v.GetSDKProperties()
		if err != nil {
			return info, fmt.Errorf("get %s version err: %v", deviceName, err)
		}
		if sprop.SDKMajor == 11 {
			info.Runtime = device.RuntimeNNIE11
		}
		dprop, err := v.GetDeviceProperties(0)
		if err != nil {
			return info, fmt.Errorf("get %s device property err: %v", deviceName, err)
		}
		info.Hardware = dprop.DeviceName
	case ATLAS:
		info.Runtime = device.RuntimeAtlas
		info.Hardware = device.HardwareHWAtlas300
	case STPU:
		info.Runtime = device.RuntimeST
		info.Hardware = device.HardwareST
	case CAMBRICON:
		dprop, err := v.GetDeviceProperties(0)
		if err != nil {
			return info, fmt.Errorf("get %s device property err: %v", deviceName, err)
		}
		info = RuntimeInfo{Runtime: device.RuntimeNeuware}
		if hardware, ok := DeviceName2CommonDeviceName[dprop.DeviceName]; ok {
			info.Hardware = hardware
		} else {
			return info, fmt.Errorf("get %s device CommonDeviceName err: %v", deviceName, err)
		}
	}

	return info, nil
}
