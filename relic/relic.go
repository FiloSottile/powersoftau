package relic

// #cgo CFLAGS: -I${SRCDIR}/relic/include -I${SRCDIR}/build/include
// #cgo LDFLAGS: ${SRCDIR}/build/lib/librelic_s.a
// #include "relic_core.h"
// #include "relic_err.h"
// #include "relic_ep.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// int _y_is_higher(const ep_t);
import "C"
import (
	"errors"
	"os"
)

func init() {
	C.core_init()
	C.ep_param_set(C.B12_381)
}

// With CHECK on, the program exits on the second uncaught(?) error,
// and there are functions like ep_read_bin that will cause two errors
// in a row without returning.
//
// With CHECK off there is no err_get_msg.
//
// Basically there's nothing we can do beyond keeping CHECK on, so that
// we see log+exit, and treat all errors as irrecoverable. YOLO.
//
// But anyway, if by mistake we cause one error and not two, we need
// to detonate ourselves. Sigh.
//
// Ah, and https://github.com/relic-toolkit/relic/issues/59.

func checkError() {
	if C.err_get_code() != C.STS_OK {
		var e *C.err_t
		var msg **C.char
		C.err_get_msg(e, msg)
		// errors.New(C.GoString(*msg))
		os.Exit(int(*e))
	}
}

// EP is a point in G1 backed by a relic ep_st.
type EP struct {
	st C.ep_st
}

func (ep *EP) SetZero() *EP {
	C.ep_set_infty(&ep.st)
	return ep
}

func (ep *EP) SetOne() *EP {
	C.ep_curve_get_gen(&ep.st)
	if C.ep_is_infty(&ep.st) == 1 {
		panic("G == 0")
	}
	return ep
}

func (ep *EP) Add(a *EP) *EP {
	C._ep_add(&ep.st, &ep.st, &a.st)
	return ep
}

func (ep *EP) Equal(a *EP) bool {
	return C.ep_cmp(&ep.st, &a.st) == C.CMP_EQ
}

const (
	FqElementSize      = 48
	G1CompressedSize   = FqElementSize
	G1UncompressedSize = 2 * FqElementSize

	Fq2ElementSize     = 96
	G2CompressedSize   = Fq2ElementSize
	G2UncompressedSize = 2 * Fq2ElementSize
)

// https://github.com/ebfull/pairing/tree/master/src/bls12_381#serialization
const (
	serializationMask       = (1 << 5) - 1
	serializationCompressed = 1 << 7
	serializationInfinity   = 1 << 6
	serializationBigY       = 1 << 5
)

// EncodeUncompressed encodes a point according to ebfull/pairing bls12_381
// serialization into a byte slice of length G1UncompressedSize.
func (ep *EP) EncodeUncompressed() []byte {
	bin := make([]byte, 2*FqElementSize+1)
	res := bin[1:]

	if C.ep_is_infty(&ep.st) == 1 {
		res[0] |= serializationInfinity
		return res
	}

	C.ep_write_bin((*C.uint8_t)(&bin[0]), C.int(len(bin)), &ep.st, 0)
	checkError()

	return res
}

// EncodeCompressed encodes a point according to ebfull/pairing bls12_381
// serialization into a byte slice of length G1CompressedSize.
func (ep *EP) EncodeCompressed() []byte {
	bin := make([]byte, FqElementSize+1)
	res := bin[1:]

	if C.ep_is_infty(&ep.st) == 1 {
		res[0] |= serializationInfinity | serializationCompressed
		return res
	}

	C.ep_norm(&ep.st, &ep.st)
	C.ep_write_bin((*C.uint8_t)(&bin[0]), C.int(len(bin)), &ep.st, 1)
	checkError()

	if C._y_is_higher(&ep.st) == 1 {
		res[0] |= serializationBigY
	}

	res[0] |= serializationCompressed
	return res
}

// DecodeUncompressed decodes a point according to ebfull/pairing bls12_381
// serialization from a byte slice of length G1UncompressedSize.
func (ep *EP) DecodeUncompressed(in []byte) (*EP, error) {
	if len(in) != G1UncompressedSize {
		return nil, errors.New("wrong encoded point size")
	}
	if in[0]&serializationCompressed != 0 {
		return nil, errors.New("point is compressed")
	}
	if in[0]&serializationBigY != 0 {
		return nil, errors.New("high Y bit improperly set")
	}

	bin := make([]byte, 2*FqElementSize+1)
	copy(bin[1:], in)
	bin[0] = 4
	bin[1] &= serializationMask

	if in[0]&serializationInfinity != 0 {
		for i := range bin[1:] {
			if bin[1+i] != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}
		C.ep_set_infty(&ep.st)
		return ep, nil
	}

	C.ep_read_bin(&ep.st, (*C.uint8_t)(&bin[0]), C.int(len(bin)))
	checkError()
	return ep, nil
}

// DecodeCompressed decodes a point according to ebfull/pairing bls12_381
// serialization from a byte slice of length G1CompressedSize.
func (ep *EP) DecodeCompressed(in []byte) (*EP, error) {
	if len(in) != G1CompressedSize {
		return nil, errors.New("wrong encoded point size")
	}
	if in[0]&serializationCompressed == 0 {
		return nil, errors.New("point isn't compressed")
	}

	bin := make([]byte, FqElementSize+1)
	copy(bin[1:], in)
	bin[0] = 2
	bin[1] &= serializationMask

	if in[0]&serializationInfinity != 0 {
		if in[0]&serializationBigY != 0 {
			return nil, errors.New("high Y bit improperly set")
		}
		for i := range bin[1:] {
			if bin[1+i] != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}
		C.ep_set_infty(&ep.st)
		return ep, nil
	}

	C.ep_norm(&ep.st, &ep.st)
	C.ep_read_bin(&ep.st, (*C.uint8_t)(&bin[0]), C.int(len(bin)))
	checkError()

	if C._y_is_higher(&ep.st) == 0 {
		if in[0]&serializationBigY != 0 {
			C._ep_neg(&ep.st, &ep.st)
		}
	} else {
		if in[0]&serializationBigY == 0 {
			C._ep_neg(&ep.st, &ep.st)
		}
	}
	return ep, nil
}
