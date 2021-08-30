package auth

// keep sync with sdk-license-go

// QuotaLimit carries current and Free quota
type QuotaLimit struct {
	Current int32
	Free    int32
}

// QuotaStatus carries current quota limits, const limits, version and error status
type QuotaStatus struct {
	Err     Error
	Version int32
	Limits  map[string]QuotaLimit
	Consts  map[string]interface{} // value type will be int32 or string
}

// DongleInfo ...
type DongleInfo struct {
	HardwareTime int64
	HardwareID   string
}

// Error defines errors from sdk to users
type Error struct {
	Code    string
	Message string
}

type CARequest struct {
	ResourceCost int32
	ConstNames   []string
	Capabilities []string
}

const (
	ViperDefaultProductName = "IVA-VIPER"
	ViperCADefaultUUID      = ""
)
