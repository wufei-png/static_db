package common

import "fmt"

type SDKError struct {
	module string
	code   int
}

func (e *SDKError) Error() string {
	return fmt.Sprintf("%s error: %d", e.module, e.code)
}

func (e *SDKError) Code() int {
	return e.code
}

func NewSDKError(module string, code int) error {
	if code == 0 {
		return nil
	}
	return &SDKError{
		module: module,
		code:   code,
	}
}
