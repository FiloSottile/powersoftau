package powersoftau

import (
	"bytes"
	"encoding/binary"

	"github.com/FiloSottile/powersoftau/bls12"
	"github.com/FiloSottile/powersoftau/internal/chacha20"
)

/*

The Rust hash_to_g2 implementation, which we have to match to pass
verification, uses the Rand trait as implemented by ChaChaRng.

Here is a reversed spec:

	1. Use a 32-byte digest as a ChaCha20 key [hash_to_g2]

	2. Pick a random field element x = c0 + c1 * u [Fq2::Rand]
		2.1. Pick a random c0 [Fq::Rand]
			2.1.1. Extract 12 random little-endian uint32 from the
				   ChaCha20 RNG, arrange them into little-endian
				   pairs as uint64 and interpret those in little-endian
				   order as a 384-bit number [FqRepr::Rand]
				   [Rng::next_u64] [ChaChaRng::next_u32]

				   The resulting big-endian byte order is like this:

	... 19 18 17 16 23 22 21 20 11 10 9 8 15 14 13 12 3 2 1 0 7 6 5 4

			2.1.2. Mask away the 3 top bits [FqRepr::Rand]
			2.1.3. If the result is not lower than the field
				   modulus [Fq::is_valid], go back to 2.1.1
			2.1.4. Perform a Montgomery reduction [Fq::into_repr]
			       [G2Uncompressed::from_affine] [Fq::mont_reduce]
		2.2. Pick a random c1, like in 2.1

	3. Pick a random flag by extracting a little-endian uint32
	   from the RNG and checking if the LSB is 1 [bool::Rand]

	4. Compute y [G2Affine::get_point_from_x]
		4.1. Compute ±y = sqrt(x^3 + b)
		4.2. If no square root exists, go back to 2
		4.3. Select the higher (modulo the field modulus) of
			 ±y if the flag at 3 is set, the lower otherwise

	5. Scale p = (x, y) by the curve cofactor [G2Affine::scale_by_cofactor]
		5.1. Perform the scalar multiplication cofactor×p

	6. If p is zero (the point at infinity) go back to 2

	7. Return p [G2::Rand]

*/

func HashToG2(digest *[32]byte) *bls12.EP2 {
	rng := chacha20.NewRng(digest)

	p := bls12.NewEP2()
	for {
		c0 := extractFieldElement(rng)
		c1 := extractFieldElement(rng)
		greater := extractBool(rng)

		// Use point deserialization instead of reimplementing lexicographic y ordering.
		buf := make([]byte, bls12.G2CompressedSize)
		copy(buf, c1[:])
		copy(buf[48:], c0[:])
		buf[0] |= 1 << 7 // serializationCompressed
		if greater {
			buf[0] |= 1 << 5 // serializationBigY
		}

		p, err := p.DecodeCompressed(buf)
		if err != nil {
			continue
		}

		p.ScaleByCofactor()

		if p.IsZero() {
			continue
		}

		return p
	}
}

var fqModulus = [48]byte{0x1a, 0x01, 0x11, 0xea, 0x39, 0x7f, 0xe6, 0x9a, 0x4b, 0x1b, 0xa7, 0xb6, 0x43, 0x4b, 0xac, 0xd7, 0x64, 0x77, 0x4b, 0x84, 0xf3, 0x85, 0x12, 0xbf, 0x67, 0x30, 0xd2, 0xa0, 0xf6, 0xb0, 0xf6, 0x24, 0x1e, 0xab, 0xff, 0xfe, 0xb1, 0x53, 0xff, 0xff, 0xb9, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xaa, 0xab}

func extractFieldElement(rng *chacha20.Rng) [48]byte {
	for {
		var res [48]byte
		for i := 48 - 8; i >= 0; i -= 8 {
			binary.BigEndian.PutUint32(res[i:], rng.ReadUint32())
			binary.BigEndian.PutUint32(res[i+4:], rng.ReadUint32())
		}
		res[0] &= 0xff >> 3
		if bytes.Compare(res[:], fqModulus[:]) >= 0 {
			continue
		}
		bls12.FqMontgomeryReduce(res[:])
		return res
	}
}

func extractBool(rng *chacha20.Rng) bool {
	x := rng.ReadUint32()
	return x&1 == 1
}
