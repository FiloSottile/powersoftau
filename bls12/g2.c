#include "relic.h"
#include "relic_fp.h"
#include "relic_epx.h"

// Helpers operate on ep2_t because we can't pass fp2_t though cgo.
// https://github.com/relic-toolkit/relic/issues/60

int ep2_y_is_higher(const ep2_t ep2) {
    uint8_t a[FP_BYTES], b[FP_BYTES];
    fp_t other;

    fp_write_bin(a, FP_BYTES, ep2->y[1]);

    fp_new(other);
    fp_copy(other, ep2->y[1]);
    fp_neg(other, other);
    fp_write_bin(b, FP_BYTES, other);
    fp_free(other);

    for (int i = 0; i < FP_BYTES; i++) {
        if (a[i] > b[i]) {
            return 1;
        } else if (a[i] < b[i]) {
            return 0;
        }
    }
    return 0;
}

void ep2_scale_by_cofactor(ep2_t p) {
    bn_t k;
    bn_new(k);
    // https://github.com/FiloSottile/powersoftau/issues/3
    bn_read_str(k, "5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5", 127, 16);
    // https://github.com/relic-toolkit/relic/issues/64
    ep2_mul_basic(p, p, k);
    bn_free(k);
}

void ep2_read_x(ep2_t a, uint8_t* bin, int len) {
    a->norm = 1;
    fp_set_dig(a->z[0], 1);
    fp_zero(a->z[1]);
    fp2_read_bin(a->x, bin, len);
    fp2_zero(a->y);
}
