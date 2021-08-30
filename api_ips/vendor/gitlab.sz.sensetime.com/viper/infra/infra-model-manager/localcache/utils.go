package localcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"gitlab.sz.sensetime.com/viper/infra/infra-model-manager/api"
)

var (
	ErrInvalidModelPath = errors.New("invalid model path name")
	ErrInvalidOID       = errors.New("invalid oid")
	ErrInvalidChecksum  = errors.New("invalid checksum")
)

const DefaultName = "default"

var (
	regRuntime = regexp.MustCompile("^[a-zA-Z0-9_]{1,32}$")
	regName    = regexp.MustCompile("^[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,127}$")
)

func ModelPathToFilename(mp *api.ModelPath, suffix string) (string, error) {
	if !regName.MatchString(mp.GetName()) {
		return "", ErrInvalidModelPath
	}
	tn := mp.GetType().String()

	stn := mp.GetSubType()
	if !regRuntime.MatchString(stn) {
		return "", ErrInvalidModelPath
	}
	rn := mp.GetRuntime()
	if rn == "" {
		rn = DefaultName
	}
	if !regRuntime.MatchString(rn) {
		return "", ErrInvalidModelPath
	}
	hn := mp.GetHardware()
	if hn == "" {
		hn = DefaultName
	}
	if !regRuntime.MatchString(hn) {
		return "", ErrInvalidModelPath
	}
	return filepath.ToSlash(filepath.Join(tn, stn, rn, hn, mp.GetName()+suffix)), nil
}

func ModelPathToFilePath(mp *api.ModelPath, lock bool) (string, error) {
	suffix := ".json"
	if lock {
		suffix = ".lock"
	}
	return ModelPathToFilename(mp, suffix)
}

func GetBlobPath(oid, checksum string) (string, error) {
	if len(checksum) < 3 {
		return "", ErrInvalidChecksum
	}
	pre := checksum[:2]
	/*
		suffix := ""
		if checksum != "" {
			suffix = "-" + checksum
		}
	*/
	return filepath.Join(pre, checksum), nil
}

func TempSuffix() string {
	return fmt.Sprintf(".tmp%v", time.Now().UnixNano())
}

type SHA256Reader struct {
	r          io.Reader
	sha256Hash hash.Hash
}

func NewSHA256Reader(r io.Reader) *SHA256Reader {
	return &SHA256Reader{
		r:          r,
		sha256Hash: sha256.New(),
	}
}

func (r *SHA256Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if n > 0 {
		r.sha256Hash.Write(p[:n]) // nolint
	}
	return
}

func (r *SHA256Reader) SHA256HexString() string {
	return hex.EncodeToString(r.sha256Hash.Sum(nil))
}

func ChecksumFile(path string) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close() // nolint
	hasher := sha256.New()
	written, err := io.Copy(hasher, f)
	if err != nil {
		return "", 0, err
	}
	return hex.EncodeToString(hasher.Sum(nil)), written, nil
}

func CopyWithChecksum(w io.Writer, r io.Reader, n int64, checksum string) error {
	h := NewSHA256Reader(r)
	rn, err := io.CopyN(w, h, n)
	if err != nil {
		return err
	}
	if rn != n {
		return fmt.Errorf("copy size mismatch: want: %d got: %d", n, rn)
	}
	c := h.SHA256HexString()
	if c != checksum {
		return fmt.Errorf("copy checksum mismatch: want: %s got: %s", checksum, c)
	}
	return nil
}
