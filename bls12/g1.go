package bls12

// #include "relic_core.h"
// #include "relic_ep.h"
// #include "relic_bn.h"
// void _ep_add(ep_t r, const ep_t p, const ep_t q) { ep_add(r, p, q); }
// void _ep_neg(ep_t r, const ep_t p) { ep_neg(r, p); }
// void _ep_mul(ep_t r, const ep_t p, const bn_t k) { ep_mul(r, p, k); }
// void _fp_rdc_monty(fp_t c, dv_t a) { fp_rdc_monty(c, a); };
// int ep_y_is_higher(const ep_t);
// void monty_reduce(uint8_t *bin, int len);
// bn_t _bn_new();
// void _bn_free(bn_t t);
import "C"
import (
	"errors"
)

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

func (ep *EP) Copy() *EP {
	a := &EP{}
	C.ep_copy(&a.st, &ep.st)
	return a
}

func (ep *EP) ScalarMult(s []byte) *EP {
	bn := C._bn_new()
	defer C._bn_free(bn)
	C.bn_read_bin(bn, (*C.uint8_t)(&s[0]), C.int(len(s)))
	checkError()
	C._ep_mul(&ep.st, &ep.st, bn)
	checkError()
	return ep
}

func (ep *EP) ScalarBaseMult(s []byte) *EP {
	bn := C._bn_new()
	defer C._bn_free(bn)
	C.bn_read_bin(bn, (*C.uint8_t)(&s[0]), C.int(len(s)))
	checkError()
	C.ep_mul_gen(&ep.st, bn)
	checkError()
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

	if C.ep_y_is_higher(&ep.st) == 1 {
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

	if C.ep_y_is_higher(&ep.st) == 0 {
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

func FqMontgomeryReduce(b []byte) {
	C.monty_reduce((*C.uint8_t)(&b[0]), C.int(len(b)))
	checkError()
}
