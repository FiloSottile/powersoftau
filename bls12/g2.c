#include "relic.h"
#include "relic_fp.h"
#include "relic_epx.h"

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
