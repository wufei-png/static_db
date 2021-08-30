package localcache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/jsonpb"
	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
)

type MetaFile struct {
	Meta *api.Model
	ETag string
	// Path string
	Data    []byte
	HashKey string
}

func NewMetaFileFromModel(model *api.Model) (*MetaFile, error) {
	m := jsonpb.Marshaler{OrigName: true}
	buf := bytes.NewBuffer(nil)
	if err := m.Marshal(buf, model); err != nil {
		return nil, err
	}
	metaRaw := buf.Bytes()
	h := md5.New()
	h.Write(metaRaw) // nolint
	etag := hex.EncodeToString(h.Sum(nil))

	hashkey, err := ModelPathToFilePath(model.GetModelPath(), false)
	if err != nil {
		return nil, err
	}

	return &MetaFile{
		Meta:    model,
		ETag:    etag,
		Data:    metaRaw,
		HashKey: hashkey,
	}, nil
}

func NewMetaFileFromBytes(metaRaw []byte) (*MetaFile, error) {
	h := md5.New()
	h.Write(metaRaw) // nolint
	etag := hex.EncodeToString(h.Sum(nil))
	var model api.Model
	if err := jsonpb.Unmarshal(bytes.NewReader(metaRaw), &model); err != nil {
		return nil, err
	}

	hashkey, err := ModelPathToFilePath(model.GetModelPath(), false)
	if err != nil {
		return nil, err
	}

	return &MetaFile{
		Meta:    &model,
		ETag:    etag,
		HashKey: hashkey,
		// Path: fp,
		Data: metaRaw,
	}, nil

}

func NewMetaFileFromPath(fp string) (*MetaFile, error) {
	metaRaw, err := ioutil.ReadFile(fp)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	return NewMetaFileFromBytes(metaRaw)
}
