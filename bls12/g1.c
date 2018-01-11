#include "relic.h"
#include "relic_fp.h"
#include "relic_ep.h"

// Helpers operate on ep_t because we can't pass fp_t though cgo.
// https://github.com/relic-toolkit/relic/issues/60

int ep_y_is_higher(const ep_t ep) {
    uint8_t a[FP_BYTES], b[FP_BYTES];
    fp_t other;

    fp_write_bin(a, FP_BYTES, ep->y);

    fp_new(other);
    fp_copy(other, ep->y);
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

void monty_reduce(uint8_t *bin, int len) {
    fp_st a;
    fp_read_bin(a, bin, len);

    dv_t t;
    dv_new(t);
    dv_zero(t, 2 * FP_DIGS);

    dv_copy(t, a, FP_DIGS);
    fp_rdc_monty_basic(a, t);

    dv_free(t);

    fp_write_bin(bin, len, a);
}
