package client

import (
	"context"
	"reflect"
	"sync"

	"github.com/sirupsen/logrus"

	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
)

const (
	selectorTagKey   = "filter"
	selectorTagValue = "replace"
)

type DownloaderRecord struct {
	Err      error
	Model    *api.Model
	ModelRef string
	RealPath string
}

type Downloader struct {
	ctx   context.Context
	opts  ModelFetchOpts
	mutex sync.Mutex

	cache  *Cache
	record []DownloaderRecord
}

func (d *Downloader) GetRecord() []DownloaderRecord {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.record
}

type DownloaderOption func(*Downloader)

func ModelFetchOptionSet(m ModelFetchOpts) DownloaderOption {
	return func(d *Downloader) {
		d.opts = m
	}
}

func ModelLocalCacheSet(c *Cache) DownloaderOption {
	return func(d *Downloader) {
		d.cache = c
	}
}

func (d *Downloader) Handle(ref string) string {
	model, path, err := d.cache.GetModelLocalPath(d.ctx, ref, d.opts)
	record := DownloaderRecord{
		Err:      err,
		Model:    model,
		ModelRef: ref,
		RealPath: path,
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.record = append(d.record, record)

	return record.RealPath
}

func NewDownloader(ctx context.Context, opt ...DownloaderOption) *Downloader {
	// default options
	options := []DownloaderOption{
		ModelLocalCacheSet(GetDefaultModelCache()),
	}
	options = append(options, opt...)

	d := &Downloader{ctx: ctx}
	for _, o := range options {
		o(d)
	}

	if d.cache == nil {
		logrus.Fatalln("model cache is nil")
	}

	return d
}

// SelectorHandler 为Selector选定struct field之后的处理逻辑
type SelectorHandler interface {
	Handle(string) string
}

// Selector 按照约定的StructTag选定struct field, 处理逻辑会交给SelectorHandler
type Selector struct {
	group   *sync.WaitGroup
	ticket  int // 用来控制并发
	bucket  chan struct{}
	handler SelectorHandler
	tagKey  string
	tagVal  string
}

func (s *Selector) recursive(val reflect.Value) {
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return
		}
		s.recursive(val.Elem())

	case reflect.Slice:
		if val.IsNil() {
			return
		}

		for i := 0; i < val.Len(); i++ {
			s.recursive(val.Index(i))
		}

	case reflect.Interface:
		if val.IsNil() {
			return
		}

		s.recursive(val.Elem())

	case reflect.Map:
		if val.IsNil() {
			return
		}

		for _, key := range val.MapKeys() {
			if key.Kind() != reflect.String {
				continue
			}

			item := val.MapIndex(key)
			s.recursive(item)
			val.SetMapIndex(key, item) // 覆盖map的value
		}

	case reflect.Struct:
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			vField := val.Field(i)
			tField := typ.Field(i)

			if tField.Tag.Get(s.tagKey) == s.tagVal {
				if vField.Kind() != reflect.String {
					logrus.Warnf("structure field tag hit but value type not String, actually: %v", vField.Kind())
					continue
				}
				if !vField.CanSet() {
					logrus.Warn("structure field tag hit but value can not be set")
					continue
				}

				s.group.Add(1)
				go func() {
					defer func() {
						s.bucket <- struct{}{}
						s.group.Done()
					}()

					<-s.bucket
					rst := s.handler.Handle(vField.String())
					vField.Set(reflect.ValueOf(rst))
				}()
			} else {
				s.recursive(vField)
			}
		}
	}
}

// Do函数内部会改变入参的值, 确保入参为指针类型
func (s *Selector) Do(src interface{}) {
	if src == nil {
		return
	}

	val := reflect.ValueOf(src)
	if val.Kind() != reflect.Ptr {
		logrus.Error("params must be pointer", val.Kind(), src)
		return
	}

	s.recursive(val)
	s.group.Wait()
	logrus.Infoln("selector done")
}

type SelectorOption func(*Selector)

func ConcurrencySet(c int) SelectorOption {
	return func(s *Selector) {
		s.ticket = c
	}
}

func HandlerSet(h SelectorHandler) SelectorOption {
	return func(s *Selector) {
		s.handler = h
	}
}

func StructureTagSet(key, value string) SelectorOption {
	return func(s *Selector) {
		s.tagKey = key
		s.tagVal = value
	}
}

func NewSelector(opt ...SelectorOption) *Selector {
	// default options
	options := []SelectorOption{
		ConcurrencySet(1),
		StructureTagSet(selectorTagKey, selectorTagValue),
	}
	options = append(options, opt...)

	s := new(Selector)
	for _, o := range options {
		o(s)
	}

	if s.ticket < 1 || s.handler == nil {
		logrus.Fatalln("option wrong")
	}

	s.group = new(sync.WaitGroup)
	s.bucket = make(chan struct{}, s.ticket)
	for i := 0; i < s.ticket; i++ {
		s.bucket <- struct{}{}
	}

	return s
}

// DownloadBySelectorDefault downloads models given by fields with specific tags of a struct
// and replace field's value with the real path of the model referenced one by one,
// returns information of models downloaded.
// Attention: local handler and downloader will be created internally.
func DownloadBySelectorDefault(src interface{}) []DownloaderRecord {
	d := NewDownloader(context.Background())
	s := NewSelector(HandlerSet(d))

	s.Do(src)
	return d.GetRecord()
}
