package h265parser

import (
	"github.com/nareix/bits"
)

func readBitToUint8(r *bits.GolombBitReader) (v uint8, err error) {
	var t uint
	if t, err = r.ReadBit(); err != nil {
		return
	}
	v = uint8(t)
	return
}

func readNBitToUint8(r *bits.GolombBitReader, n int) (v uint8, err error) {
	var t uint
	if t, err = r.ReadBits(n); err != nil {
		return
	}
	v = uint8(t)
	return
}

// Alogorithm is from
// http://graphics.stanford.edu/~seander/bithacks.html#IntegerLogDeBruijn
func log2(n uint32) (r int) {

	var multiplyDeBruijnBitPosition = []int{
		0, 9, 1, 10, 13, 21, 2, 29, 11, 14, 16, 18, 22, 25, 3, 30,
		8, 12, 20, 28, 15, 17, 24, 7, 19, 27, 23, 6, 26, 5, 4, 31,
	}

	// De Bruijn sequence
	const b = uint32(0x07C4ACDD)

	// first round down to one less than a power of 2
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	i := (uint32(n) * b) >> 27
	r = multiplyDeBruijnBitPosition[i]
	return
}

/*
func readExponentialGolombLong(r *bits.GolombBitReader) (v uint, err error) {
	var t uint
	if t, err = r.ReadBits(32); err != nil {
		return
	}
	log := 31 - log2(uint32(t))
	fmt.Println("XX ", t, log)
	if _, err = r.ReadBits(log); err != nil {
		return
	}
	if t, err = r.ReadBits(log + 1); err != nil {
		return
	}
	v = t - 1
	return
}
*/

func readLongUint(r *bits.GolombBitReader) (v uint, err error) {
	var b uint
	if b, err = r.ReadBits(8); err != nil {
		return
	}
	for b == 255 {
		v += b
		if b, err = r.ReadBits(8); err != nil {
			return
		}
	}
	v += b
	return
}
