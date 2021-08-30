package feature

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

type Encoder struct {
	flag Flag

	iv    []byte
	block cipher.Block
}

func NewEncoder(flag Flag, key []byte) (*Encoder, error) {
	enc := &Encoder{
		flag: flag,
	}
	if flag&ET1 != 0 {
		return nil, ErrUnsupportedEncoding
	}

	if flag&ET2 != 0 {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		enc.block = block
		enc.iv = make([]byte, aes.BlockSize)
		// 存在误报 feature/encoder.go:31: G404: Use of weak random number generator (math/rand instead of crypto/rand) (gosec)
		// nolint:gosec
		if _, err := rand.Read(enc.iv); err != nil {
			return nil, err
		}
	}

	return enc, nil
}

func (e *Encoder) encodeBlob(f32 []float32) []byte {
	size := persistedBlobSize(len(f32), e.flag)
	buf := make([]byte, size)
	headSize := len(e.iv)
	off := 0
	// reserve for IV
	if e.flag&ET2 != 0 {
		copy(buf, e.iv)
		off += headSize
	}
	if e.flag&Float16 != 0 {
		off += PutFloat16From32(buf[off:], f32)
	} else {
		off += PutFloat32(buf[off:], f32)
	}
	if off != len(buf) {
		panic("must equal")
	}
	if e.flag&ET2 != 0 {
		if len(buf)%aes.BlockSize != 0 {
			panic("ET2 padding not supported")
		}
		mode := cipher.NewCBCEncrypter(e.block, e.iv)
		mode.CryptBlocks(buf[headSize:], buf[headSize:])
	}
	return buf
}

func (e *Encoder) Encode(raw RawFeature) (PersistedFeature, error) {
	if raw.Header.MagicNumber != FeatureMagicNumber {
		return PersistedFeature{}, ErrInvalidMagic
	}
	b := e.encodeBlob(raw.Raw)
	return PersistedFeature{
		header: Header{
			MagicNumber: raw.Header.MagicNumber,
			Version:     raw.Header.Version,
			DataLen:     uint32(len(b)),
			Dim:         raw.Header.Dim,
			Flags:       e.flag,
		},
		blob: b,
	}, nil
}
