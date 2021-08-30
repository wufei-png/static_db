package auth

import "errors"

var (
	ErrServiceDisabled = errors.New("license: service disabled by")
	ErrBadFeatureKey   = errors.New("license: bad feature key")
	ErrLicenseExpired  = errors.New("license: expired")
	ErrBadQuotaLimit   = errors.New("license: bad quota limit")
	ErrNoQuotaLimit    = errors.New("license: no quota limit")
	ErrBadNodeQuota    = errors.New("license: bad nodes quota")
	ErrQuotaRenew      = errors.New("license: quota renew error")
	ErrConstFetch      = errors.New("license: const value fetch error")
)

const (
	AuthErrExitCode = 1
)
