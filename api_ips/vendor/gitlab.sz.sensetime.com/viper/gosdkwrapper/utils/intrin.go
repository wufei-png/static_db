package utils

/*
float dot_product(const float *a, const float *b, int n);
*/
import "C"
import (
	"runtime"

	"golang.org/x/sys/cpu"
)

var avx2Supported = false

func Avx2Supported() bool {
	return avx2Supported
}

func AVX2DotProduct(a, b []float32) float32 {
	return float32(C.dot_product((*C.float)(&a[0]), (*C.float)(&b[0]), C.int(len(a))))
}

func init() {
	switch {
	case runtime.GOARCH == "arm64" && runtime.GOOS == "linux":
		avx2Supported = cpu.ARM64.HasASIMD
	case runtime.GOARCH == "amd64":
		avx2Supported = cpu.X86.HasAVX
	default:
		panic("unsupported arch or os")
	}
}
