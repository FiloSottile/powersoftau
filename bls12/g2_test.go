package bls12_test

import (
	"bytes"
	"testing"

	"github.com/FiloSottile/powersoftau/bls12"
)

func TestVectorG2(t *testing.T) {
	t.Run("Uncompressed", func(t *testing.T) {
		var (
			data = readFile(t, "testdata/g2_uncompressed_valid_test_vectors.dat")
			ep2  = bls12.NewEP2().SetZero()
			a    = bls12.NewEP2()
			one  = bls12.NewEP2().SetOne()
			d    = data
		)
		defer ep2.Close()
		defer a.Close()
		defer one.Close()
		for i := 0; i < 1000; i++ {
			t.Logf("%d <- %x", i, d[:bls12.G2UncompressedSize])
			_, err := a.DecodeUncompressed(d[:bls12.G2UncompressedSize])
			if err != nil {
				t.Errorf("%d: failed decoding: %v", i, err)
			}
			if !ep2.Equal(a) {
				t.Errorf("%d: different point", i)
			}
			buf := ep2.EncodeUncompressed()
			t.Logf("%d -> %x", i, buf)
			if !bytes.Equal(buf, d[:bls12.G2UncompressedSize]) {
				t.Errorf("%d: different encoding", i)
			}
			d = d[bls12.G2UncompressedSize:]
			ep2.Add(one)
		}
	})
	t.Run("Compressed", func(t *testing.T) {
		var (
			data = readFile(t, "testdata/g2_compressed_valid_test_vectors.dat")
			ep2  = bls12.NewEP2().SetZero()
			a    = bls12.NewEP2()
			one  = bls12.NewEP2().SetOne()
			d    = data
		)
		defer ep2.Close()
		defer a.Close()
		defer one.Close()
		for i := 0; i < 1000; i++ {
			t.Logf("%d <- %x", i, d[:bls12.G2CompressedSize])
			_, err := a.DecodeCompressed(d[:bls12.G2CompressedSize])
			if err != nil {
				t.Errorf("%d: failed decoding: %v", i, err)
			}
			if !ep2.Equal(a) {
				t.Errorf("%d: different point", i)
			}
			buf := ep2.EncodeCompressed()
			t.Logf("%d -> %x", i, buf)
			if !bytes.Equal(buf, d[:bls12.G2CompressedSize]) {
				t.Errorf("%d: different encoding", i)
			}
			d = d[bls12.G2CompressedSize:]
			ep2.Add(one)
		}
	})
}
