package bls12

// #include "relic_core.h"
// #include "relic_epx.h"
// ep2_t _ep2_new() { ep2_t t; ep2_new(t); return t; }
// void _ep2_free(ep2_t t) { ep2_free(t); }
// void _ep2_add(ep2_t r, const ep2_t p, const ep2_t q) { ep2_add(r, p, q); }
// void _ep2_neg(ep2_t r, const ep2_t p) { ep2_neg(r, p); }
// int ep2_y_is_higher(const ep2_t ep2);
import "C"
import "errors"

// EP2 is a point in G2 backed by a relic ep2_t.
//
// EP2 requires manual memory management.
type EP2 struct {
	t C.ep2_t
}

func NewEP2() *EP2 {
	ep2 := &EP2{C._ep2_new()}
	checkError()
	return ep2
}

func (ep2 *EP2) Close() {
	C._ep2_free(ep2.t)
}

func (ep2 *EP2) SetZero() *EP2 {
	C.ep2_set_infty(ep2.t)
	return ep2
}

func (ep2 *EP2) SetOne() *EP2 {
	C.ep2_curve_get_gen(ep2.t)
	if C.ep2_is_infty(ep2.t) == 1 {
		panic("G == 0")
	}
	return ep2
}

func (ep2 *EP2) Add(a *EP2) *EP2 {
	C._ep2_add(ep2.t, ep2.t, a.t)
	return ep2
}

func (ep2 *EP2) Equal(a *EP2) bool {
	return C.ep2_cmp(ep2.t, a.t) == C.CMP_EQ
}

const (
	Fq2ElementSize     = 96
	G2CompressedSize   = Fq2ElementSize
	G2UncompressedSize = 2 * Fq2ElementSize
)

// EncodeUncompressed encodes a point according to ebfull/pairing bls12_381
// serialization into a byte slice of length G2UncompressedSize.
func (ep2 *EP2) EncodeUncompressed() []byte {
	bin := make([]byte, 2*Fq2ElementSize+1)
	res := bin[1:]

	if C.ep2_is_infty(ep2.t) == 1 {
		res[0] |= serializationInfinity
		return res
	}

	C.ep2_write_bin((*C.uint8_t)(&bin[0]), C.int(len(bin)), ep2.t, 0)
	checkError()

	return swapLimbs(make([]byte, 0, 2*Fq2ElementSize), res)
}

// EncodeCompressed encodes a point according to ebfull/pairing bls12_381
// serialization into a byte slice of length G2CompressedSize.
func (ep2 *EP2) EncodeCompressed() []byte {
	bin := make([]byte, Fq2ElementSize+1)
	res := bin[1:]

	if C.ep2_is_infty(ep2.t) == 1 {
		res[0] |= serializationInfinity | serializationCompressed
		return res
	}

	C.ep2_norm(ep2.t, ep2.t)
	C.ep2_write_bin((*C.uint8_t)(&bin[0]), C.int(len(bin)), ep2.t, 1)
	checkError()

	res = swapLimbs(make([]byte, 0, Fq2ElementSize), res)

	if C.ep2_y_is_higher(ep2.t) == 1 {
		res[0] |= serializationBigY
	}

	res[0] |= serializationCompressed
	return res
}

// swapLimbs appends in to out after swapping Fq2Element limbs, and
// returns out. relic uses c0||c1 for c0 + c1 * u, instead of c1||c0.
func swapLimbs(out, in []byte) []byte {
	if len(in)%Fq2ElementSize != 0 {
		panic("wrong Fq2Element size")
	}
	for len(in) > 0 {
		out = append(out, in[Fq2ElementSize/2:Fq2ElementSize]...)
		out = append(out, in[:Fq2ElementSize/2]...)
		in = in[Fq2ElementSize:]
	}
	return out
}

// DecodeUncompressed decodes a point according to ebfull/pairing bls12_381
// serialization from a byte slice of length G2UncompressedSize.
func (ep2 *EP2) DecodeUncompressed(in []byte) (*EP2, error) {
	if len(in) != G2UncompressedSize {
		return nil, errors.New("wrong encoded point size")
	}
	if in[0]&serializationCompressed != 0 {
		return nil, errors.New("point is compressed")
	}
	if in[0]&serializationBigY != 0 {
		return nil, errors.New("high Y bit improperly set")
	}

	bin := make([]byte, 1, 2*Fq2ElementSize+1)
	bin[0] = 4
	bin = swapLimbs(bin, in)
	bin[Fq2ElementSize/2+1] &= serializationMask

	if in[0]&serializationInfinity != 0 {
		for i := range bin[1:] {
			if bin[1+i] != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}
		C.ep2_set_infty(ep2.t)
		return ep2, nil
	}

	C.ep2_read_bin(ep2.t, (*C.uint8_t)(&bin[0]), C.int(len(bin)))
	checkError()
	return ep2, nil
}

// DecodeCompressed decodes a point according to ebfull/pairing bls12_381
// serialization from a byte slice of length G2CompressedSize.
func (ep2 *EP2) DecodeCompressed(in []byte) (*EP2, error) {
	if len(in) != G2CompressedSize {
		return nil, errors.New("wrong encoded point size")
	}
	if in[0]&serializationCompressed == 0 {
		return nil, errors.New("point isn't compressed")
	}

	bin := make([]byte, 1, Fq2ElementSize+1)
	bin[0] = 2
	bin = swapLimbs(bin, in)
	bin[Fq2ElementSize/2+1] &= serializationMask

	if in[0]&serializationInfinity != 0 {
		if in[0]&serializationBigY != 0 {
			return nil, errors.New("high Y bit improperly set")
		}
		for i := range bin[1:] {
			if bin[1+i] != 0 {
				return nil, errors.New("invalid infinity encoding")
			}
		}
		C.ep2_set_infty(ep2.t)
		return ep2, nil
	}

	C.ep2_norm(ep2.t, ep2.t)
	C.ep2_read_bin(ep2.t, (*C.uint8_t)(&bin[0]), C.int(len(bin)))
	checkError()

	if C.ep2_y_is_higher(ep2.t) == 0 {
		if in[0]&serializationBigY != 0 {
			C._ep2_neg(ep2.t, ep2.t)
		}
	} else {
		if in[0]&serializationBigY == 0 {
			C._ep2_neg(ep2.t, ep2.t)
		}
	}
	return ep2, nil
}
