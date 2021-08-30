package utils

/*
#cgo CFLAGS: -std=c99
#include <arm_neon.h>

// https://github.com/jgaeddert/liquid-dsp/blob/master/src/dotprod/src/dotprod_rrrf.neon.c#L40
float dot_product(const float *_h, const float *_x, int _n) {
    float32x4_t v;   // input vector
    float32x4_t h;   // coefficients vector
    float32x4_t s;   // dot product

    // load zeros into sum register
    float zeros[4] = {0,0,0,0};
    float32x4_t sum = vld1q_f32(zeros);

    // t = 4*(floor(_n/4))
    int t = (_n >> 2) << 2;

    int i;
    for (i=0; i<t; i+=4) {
        // load inputs into register (unaligned)
        v = vld1q_f32(&_x[i]);

        // load coefficients into register (aligned)
        h = vld1q_f32(&_h[i]);

        // compute multiplication
        s = vmulq_f32(h,v);

        // parallel addition
        sum = vaddq_f32(sum, s);
    }

    // unload packed array
    float w[4];
    vst1q_f32(w, sum);
    float total = w[0] + w[1] + w[2] + w[3];

    // cleanup
    for (; i<_n; i++)
        total += _x[i] * _h[i];

    return total;
}
*/
import "C"
