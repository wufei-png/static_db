// +build se

package se

import (
	"errors"
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gitlab.sz.sensetime.com/viper/gosdkwrapper/dsc"
)

var AvailableIndexTypeList = map[dsc.IndexType]bool{
	dsc.IndexPQ:   true,
	dsc.IndexDC:   true,
	dsc.IndexHNSW: true,
}

type SEIndexHandlerBuilder struct{}

type SEIndexHandler struct {
	seContext      SearchEngineContext
	seType         dsc.IndexType
	ctxMemReserved int32
	seThreadNum    int32
}

func (b *SEIndexHandlerBuilder) Build(config interface{}) (dsc.IndexHandler, error) {
	c, ok := config.(*dsc.GeneralIndexHandlerConfig)
	if !ok {
		return nil, errors.New("unknown se hanlder builder config")
	}

	if c.IsCPUModel {
		log.Warnf("se is running in cpu model, valid index model must be configured")
	}
	return NewSEIndexHandler(c)
}

func (b *SEIndexHandlerBuilder) InitEnv(config interface{}) error {
	c, ok := config.(*dsc.SEInitEnvConfig)
	if !ok {
		return errors.New("unknown se env init config")
	}
	return InitSearchEngineEnv(c.ProductName, c.LicensePath)
}

func (b *SEIndexHandlerBuilder) DestroyEnv() error {
	DestroySearchEngineEnv()
	return nil
}

func NewSEIndexHandler(config *dsc.GeneralIndexHandlerConfig) (dsc.IndexHandler, error) {
	if _, ok := AvailableIndexTypeList[config.IndexType]; !ok {
		return nil, errors.New("unsupport index type")
	}

	return &SEIndexHandler{
		seType:         config.IndexType,
		ctxMemReserved: int32(config.DeviceTemMemory),
		seThreadNum:    config.SEThreadNum,
	}, nil
}

func (handle *SEIndexHandler) BindDevice(deviceID int32) error {
	runtime.LockOSThread()
	sc, err := InitSearchEngineContext(deviceID)
	if err != nil {
		return fmt.Errorf("fail to init se context, err: %v", err)
	}

	if err = SetSearchEngineContextReservedMemory(sc, handle.ctxMemReserved); err != nil {
		if err2 := DestroySearchEngineContext(sc); err2 != nil {
			return fmt.Errorf("fail to set reserved memory of se context, err: %s; %s", err, err2)
		}
		return fmt.Errorf("fail to set reserved memory of se context, err: %s", err)
	}

	if handle.seThreadNum > 0 {
		if err = SetSearchEngineContextThreadNum(sc, handle.seThreadNum); err != nil {
			if err2 := DestroySearchEngineContext(sc); err2 != nil {
				return fmt.Errorf("fail to set thread number of se context, err: %s; %s", err, err2)
			}
			return fmt.Errorf("fail to set thread number of se context, err: %s", err)
		}
	}
	handle.seContext = sc
	return BindDevice(deviceID)
}

func (handle *SEIndexHandler) UnbindDevice() error {
	defer runtime.UnlockOSThread()
	defer UnbindDevice()
	err := DestroySearchEngineContext(handle.seContext)
	if err != nil {
		return fmt.Errorf("fail to destroy se context, err: %v", err)
	}
	return nil
}

func (handle *SEIndexHandler) InitIndex(indexConfig interface{}) (dsc.SearchIndex, error) {
	c, ok := indexConfig.(*dsc.GeneralIndexConfig)
	if !ok {
		return nil, errors.New("unknown index config for se")
	}
	c.IndexType = handle.seType
	return InitSearchEngineIndex(handle.seContext, c)
}

func (handle *SEIndexHandler) LoadIndex(filepath string) (dsc.SearchIndex, error) {
	return LoadSearchEngineIndex(handle.seContext, filepath)
}

func init() {
	//register
	dsc.IndexHandlerFactory[dsc.SE] = &SEIndexHandlerBuilder{}
}
