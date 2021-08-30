package feature

import (
	"encoding/binary"
	"math"

	"github.com/h2so5/half"
)

const sizeOfFloat32 = 4

func PutFloat32(buf []byte, fs []float32) int {
	off := 0
	for _, v := range fs {
		binary.LittleEndian.PutUint32(buf[off:], math.Float32bits(v))
		off += 4
	}
	return off
}

func PutFloat16(buf []byte, fs []half.Float16) int {
	off := 0
	for _, v := range fs {
		binary.LittleEndian.PutUint16(buf[off:], uint16(v))
		off += 2
	}
	return off
}

func PutFloat16From32(buf []byte, fs []float32) int {
	off := 0
	for _, v := range fs {
		x := half.NewFloat16(v)
		binary.LittleEndian.PutUint16(buf[off:], uint16(x))
		off += 2
	}
	return off
}

func GetFloat32(buf []byte, fs []float32) int {
	off := 0
	for i := 0; i < len(fs); i++ {
		v := binary.LittleEndian.Uint32(buf[off:])
		fs[i] = math.Float32frombits(v)
		off += 4
	}
	return off
}

func GetFloat16(buf []byte, fs []half.Float16) int {
	off := 0
	for i := 0; i < len(fs); i++ {
		v := binary.LittleEndian.Uint16(buf[off:])
		fs[i] = half.Float16(v)
		off += 2
	}
	return off
}

func GetFloat16To32(buf []byte, fs []float32) int {
	off := 0
	for i := 0; i < len(fs); i++ {
		v := binary.LittleEndian.Uint16(buf[off:])
		fs[i] = half.Float16(v).Float32()
		off += 2
	}
	return off
}
