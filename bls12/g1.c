#include "relic.h"
#include "relic_fp.h"
#include "relic_ep.h"

// Helpers operate on ep_t because we can't pass fp_t though cgo.
// https://github.com/relic-toolkit/relic/issues/60

int _y_is_higher(const ep_t ep) {
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
