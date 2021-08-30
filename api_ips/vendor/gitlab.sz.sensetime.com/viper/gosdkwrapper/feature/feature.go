package feature

import (
	"encoding/binary"
	"errors"
)

// FeatureMagicNumber is 'SOF1' (Sensetime Object Feature) in little endian
const FeatureMagicNumber = 0x31464f53

const MaxFeatureBlobSize = 1 << 20

type Flag uint32

const (
	Float16 Flag = 1 << 0
	ET1          = 1 << 1 // garble table, name is deliberately obfuscated
	ET2          = 1 << 2 // AES128-CBC, prefix with 16bytes IV
	ET3          = 1 << 3 // transform flag, which will be set when feature is processed with transformation model
)

var (
	ErrBadLength           = errors.New("bad feature data length")
	ErrInvalidMagic        = errors.New("invalid feature magic")
	ErrCorrupted           = errors.New("corrupted feature")
	ErrBadDimension        = errors.New("bad feature dimension")
	ErrUnsupportedEncoding = errors.New("unsupported feature encoding")
)

// Header is the header of feature which stored by us
type Header struct {
	MagicNumber uint32
	Version     int32
	// DataLen is data bytes size, in case for data padding
	DataLen uint32
	// Dim is feature dimension
	Dim        uint32
	objectType int32
	Flags      Flag
	reserved1  int32
	reserved2  int32
}

const HeaderSize = 32

type PersistedFeature struct {
	header Header
	blob   []byte
}

func (p PersistedFeature) Header() Header {
	return p.header
}

func (p PersistedFeature) Blob() []byte {
	return p.blob
}

func (p PersistedFeature) ByteSize() int {
	return HeaderSize + len(p.blob)
}

func (p PersistedFeature) Bytes() []byte {
	b := make([]byte, HeaderSize+len(p.blob))
	binary.LittleEndian.PutUint32(b[0:], p.header.MagicNumber)
	binary.LittleEndian.PutUint32(b[4:], uint32(p.header.Version))
	binary.LittleEndian.PutUint32(b[8:], p.header.DataLen)
	binary.LittleEndian.PutUint32(b[12:], p.header.Dim)
	binary.LittleEndian.PutUint32(b[16:], uint32(p.header.objectType))
	binary.LittleEndian.PutUint32(b[20:], uint32(p.header.Flags))
	binary.LittleEndian.PutUint32(b[24:], uint32(p.header.reserved1))
	binary.LittleEndian.PutUint32(b[28:], uint32(p.header.reserved2))
	copy(b[HeaderSize:], p.blob)
	return b
}

func (p PersistedFeature) Dimension() int {
	return int(p.header.Dim)
}

func (p PersistedFeature) IsTransformed() bool {
	return p.header.Flags&ET3 != 0
}

func persistedBlobSize(dim int, flag Flag) int {
	flen := dim * 4
	if flag&Float16 != 0 {
		flen /= 2
	}
	if flag&ET2 != 0 {
		flen += 16
	}
	return flen
}

func NewPersistedFeatureFromBytes(b []byte, clone bool) (PersistedFeature, error) {
	if len(b) < HeaderSize {
		return PersistedFeature{}, ErrBadLength
	}
	var header Header
	header.MagicNumber = binary.LittleEndian.Uint32(b[0:])
	if header.MagicNumber != FeatureMagicNumber {
		return PersistedFeature{}, ErrInvalidMagic
	}
	header.Version = int32(binary.LittleEndian.Uint32(b[4:]))
	header.DataLen = binary.LittleEndian.Uint32(b[8:])
	header.Dim = binary.LittleEndian.Uint32(b[12:])
	header.objectType = int32(binary.LittleEndian.Uint32(b[16:]))
	header.Flags = Flag(binary.LittleEndian.Uint32(b[20:]))
	header.reserved1 = int32(binary.LittleEndian.Uint32(b[24:]))
	header.reserved2 = int32(binary.LittleEndian.Uint32(b[28:]))
	if header.DataLen > MaxFeatureBlobSize {
		return PersistedFeature{}, ErrBadLength
	}
	if len(b) < int(header.DataLen+HeaderSize) {
		return PersistedFeature{}, ErrBadLength
	}
	if header.Flags&ET1 != 0 {
		return PersistedFeature{}, ErrUnsupportedEncoding
	}
	if header.Dim <= 0 {
		return PersistedFeature{}, ErrBadDimension
	}
	var blob []byte
	if clone {
		blob = make([]byte, len(b)-HeaderSize)
		copy(blob, b[HeaderSize:])
	} else {
		blob = b[HeaderSize:]
	}
	return PersistedFeature{
		header: header,
		blob:   blob,
	}, nil
}

// RawFeature is in-memory representation of object feature.
type RawFeature struct {
	Header Header
	Raw    []float32
}

func NewRawFeatureFromFloat32(version int32, raw []float32) RawFeature {
	if raw == nil {
		panic("feature: nil raw feature")
	}
	return RawFeature{
		Header: Header{
			MagicNumber: FeatureMagicNumber,
			Version:     version,
			DataLen:     uint32(len(raw) * sizeOfFloat32),
			Dim:         uint32(len(raw)),
		},
		Raw: raw,
	}
}

func NewRawFeatureFromFloat32Bytes(version int32, raw []byte) (RawFeature, error) {
	if len(raw)&0x3 != 0 {
		return RawFeature{}, ErrBadLength
	}
	f32 := make([]float32, len(raw)/4)
	GetFloat32(raw, f32)
	return RawFeature{
		Header: Header{
			MagicNumber: FeatureMagicNumber,
			Version:     version,
			DataLen:     uint32(len(f32) * sizeOfFloat32),
			Dim:         uint32(len(f32)),
		},
		Raw: f32,
	}, nil
}

func (r RawFeature) Dimension() int {
	return len(r.Raw)
}

func (r RawFeature) Norm2() float32 {
	return r.Dot(r)
}

func (r RawFeature) IsNormalized() bool {
	s := r.Norm2()
	return s > 0.8 && s < 1.2
}

func (r RawFeature) Dot(v RawFeature) float32 {
	if len(r.Raw) != len(v.Raw) {
		panic("feature: dimension mismatch")
	}
	var s float32
	for i := range r.Raw {
		s += r.Raw[i] * v.Raw[i]
	}
	return s
}
