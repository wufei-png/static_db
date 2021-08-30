package utils

/*
#cgo CFLAGS: -std=c99 -fPIC -mavx -fopenmp
#include <immintrin.h>

float dot_product(const float *a, const float *b, int n) {
    float v[8];
    float prod = 0;
    int step = 8;   // avx2 __m256 can hold 8 floats
    int epoch = n / step;
    for (int i = 0; i < epoch; i++) {
            __m256 vec1 = _mm256_loadu_ps(a + i * step); // from memory to ymm1
            __m256 vec2 = _mm256_loadu_ps(b + i * step); // from memory to ymm2
            __m256 prod_sum = _mm256_dp_ps(vec1, vec2, 0xFF); // dot product
            _mm256_storeu_ps(v, prod_sum); // form ymm3 to memory
            prod += (v[0] + v[4]); // get the result
    }
    for (int i = epoch*step; i < n; i++){
		prod += a[i] * b[i];
    }
    return prod;
}
*/
import "C"
