#include "relic.h"
#include "relic_fp.h"
#include "relic_ep.h"

// Helpers operate on ep_t because we can't pass fp_t though cgo.
// https://github.com/relic-toolkit/relic/issues/60

int _y_is_higher(const ep_t ep) {
    fp_t other;
    fp_new(other);
    fp_copy(other, ep->y);
    fp_neg(other, other);
    int res = fp_cmp(ep->y, other);
    fp_free(other);
    if (res == CMP_GT) {
        return 1;
    } else {
        return 0;
    }
}
