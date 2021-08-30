package client

import (
	"errors"
	"strings"

	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/localcache"
)

var (
	ErrInvalidModelRef  = errors.New("invalid model ref")
	ErrModelNotFound    = localcache.ErrFileNotFound
	ErrModelLocked      = localcache.ErrFileLocked
	ErrModelExists      = errors.New("model already exists")
	ErrInvalidModelBlob = errors.New("invalid model blob")
	ErrOfflineMode      = errors.New("work in offline mode")
	ErrContextDone      = errors.New("context done")
	ErrRetryLimitDone   = errors.New("retry limit exhausted")
)

// Error defines client error
type Error struct {
	Op  string
	Err error
}

// Error returns an error message
func (e *Error) Error() string { return e.Op + ": " + e.Err.Error() }

func ParseModelRef(ref string) (*api.ModelPath, error) {
	arr := strings.Split(ref, "/")
	if len(arr) != 5 {
		return nil, ErrInvalidModelRef
	}
	mt, ok := api.ModelType_value[arr[0]]
	if !ok {
		return nil, ErrInvalidModelRef
	}
	return &api.ModelPath{
		Type:     api.ModelType(mt),
		SubType:  arr[1],
		Runtime:  arr[2],
		Hardware: arr[3],
		Name:     arr[4],
	}, nil
}
