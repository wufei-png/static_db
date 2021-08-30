package feature

import (
	"crypto/aes"
	"crypto/cipher"
)

type Decoder struct {
	block cipher.Block
}

func NewDecoder(key []byte) (*Decoder, error) {
	dec := &Decoder{}
	if len(key) > 0 {
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		dec.block = block
	}

	return dec, nil
}

func (d *Decoder) Decode(pf PersistedFeature) (RawFeature, error) {
	if pf.header.MagicNumber != FeatureMagicNumber {
		return RawFeature{}, ErrInvalidMagic
	}
	if pf.header.Flags&ET1 != 0 {
		return RawFeature{}, ErrUnsupportedEncoding
	}
	blob := pf.blob
	if uint32(len(blob)) != pf.header.DataLen {
		return RawFeature{}, ErrBadLength
	}
	if pf.header.Flags&ET2 != 0 {
		if d.block == nil {
			return RawFeature{}, ErrUnsupportedEncoding
		}
		if len(blob) < aes.BlockSize {
			return RawFeature{}, ErrCorrupted
		}
		if len(blob)%aes.BlockSize != 0 {
			return RawFeature{}, ErrCorrupted
		}
		iv := blob[:aes.BlockSize]
		blob = blob[aes.BlockSize:]
		mode := cipher.NewCBCDecrypter(d.block, iv)
		dst := make([]byte, len(blob))
		mode.CryptBlocks(dst, blob)
		blob = dst
	}

	var raw []float32
	if pf.header.Flags&Float16 != 0 {
		if len(blob)&0x1 != 0 {
			return RawFeature{}, ErrCorrupted
		}
		dim := len(blob) / 2
		if uint32(dim) != pf.header.Dim {
			return RawFeature{}, ErrBadDimension
		}
		raw = make([]float32, dim)
		GetFloat16To32(blob, raw)
	} else {
		if len(blob)&0x3 != 0 {
			return RawFeature{}, ErrCorrupted
		}
		dim := len(blob) / 4
		if uint32(dim) != pf.header.Dim {
			return RawFeature{}, ErrBadDimension
		}
		raw = make([]float32, dim)
		GetFloat32(blob, raw)
	}
	h := pf.header
	h.DataLen = uint32(sizeOfFloat32 * len(raw))
	h.Flags = 0
	return RawFeature{
		Header: h,
		Raw:    raw,
	}, nil
}
