package common

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math"
)

const Float32Bytes = 4

func Base64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	} else {
		return dst[:n], nil
	}
}

func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func EncodeFloat32(src []float32) []byte {
	dst := make([]byte, len(src)*Float32Bytes)
	offset := 0
	for _, f := range src {
		binary.LittleEndian.PutUint32(dst[offset:], math.Float32bits(f))
		offset += Float32Bytes
	}
	return dst
}

func DecodeFloat32(src []byte) []float32 {
	dst := make([]float32, len(src)/Float32Bytes)
	offset := 0
	for i := range dst {
		b := binary.LittleEndian.Uint32(src[offset:])
		offset += Float32Bytes
		dst[i] = math.Float32frombits(b)
	}
	return dst
}
func ReadLines(r io.Reader, delim byte) chan []byte {
	br := bufio.NewReader(r)
	ch := make(chan []byte, 8)

	go func() {
		defer close(ch)
		for {
			b, err := br.ReadBytes(delim)
			if err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Print("failed to read: %s", err)
				}
			}
			//log.Debugf("read %d bytes from file", len(b))
			ch <- b
		}
	}()
	return ch
}

// fp, err := os.Open(file)
func ChangeTobase(path string) string {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	res := base64.StdEncoding.EncodeToString(srcByte)
	return res
}
