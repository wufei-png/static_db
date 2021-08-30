package dsc

// IndexHandlerBuilder Type
const (
	GSC = "gsc"
	SE  = "se"
)

var IndexHandlerFactory = map[string]IndexHandlerBuilder{}

type IndexHandlerBuilder interface {
	InitEnv(config interface{}) error
	DestroyEnv() error
	Build(config interface{}) (IndexHandler, error)
}

type IndexHandler interface {
	BindDevice(deviceID int32) error
	UnbindDevice() error
	InitIndex(indexConfig interface{}) (SearchIndex, error)
	LoadIndex(cfilepath string) (SearchIndex, error)
}

type GeneralIndexHandlerConfig struct {
	IsCPUModel      bool
	DeviceTemMemory uint64
	IndexType       IndexType
	SEThreadNum     int32
}

type SEInitEnvConfig struct {
	ProductName string
	LicensePath string
}
