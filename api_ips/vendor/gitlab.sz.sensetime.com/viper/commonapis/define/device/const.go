package device

// Viper支持的加速卡的软件Runtime的列表
const (
	RuntimeTRT2    = "trt2"
	RuntimeTRT3    = "trt3"
	RuntimeTRT5    = "trt5"
	RuntimeTRT6    = "trt6"
	RuntimeTRT7    = "trt7"
	RuntimeAMD64   = "amd64"
	RuntimeARM64   = "arm64"
	RuntimeAtlas   = "atlas"
	RuntimeST      = "st100"
	RuntimeNNIE11  = "nnie11"
	RuntimeNeuware = "neuware"
)

// Viper支持的加速卡的硬件类型的列表
const (
	HardwareNVGTX1080  = "nv_gtx1080"
	HardwareNVP4       = "nv_p4"
	HardwareNVGTX2080  = "nv_gtx2080"
	HardwareNVT4       = "nv_t4"
	HardwareNVRTX4000  = "nv_rtx4000"
	HardwareIntel      = "intel"
	HardwareAMD        = "amd"
	HardwareARM        = "arm"
	HardwareHWAtlas300 = "hw_atlas300"
	HardwareHWAtlas500 = "hw_atlas500"
	HardwareST         = "st"
	HardwareHI3559A    = "hi_3559a"
	HardwareMLU270     = "mlu270"
	HardwareMLU220Edge = "mlu220edge"
)
