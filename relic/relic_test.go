package relic_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/FiloSottile/powersoftau/relic"
)

func TestVector(t *testing.T) {
	t.Run("G1Uncompressed", func(t *testing.T) {
		var (
			data = readFile(t, "testdata/g1_uncompressed_valid_test_vectors.dat")
			ep   = (&relic.EP{}).SetZero()
			a    = &relic.EP{}
			one  = (&relic.EP{}).SetOne()
			d    = data
		)
		for i := 0; i < 1000; i++ {
			t.Logf("%d <- %x", i, d[:relic.G1UncompressedSize])
			_, err := a.DecodeUncompressed(d[:relic.G1UncompressedSize])
			if err != nil {
				t.Errorf("%d: failed decoding: %v", i, err)
			}
			if !ep.Equal(a) {
				t.Errorf("%d: different point", i)
			}
			buf := ep.EncodeUncompressed()
			t.Logf("%d -> %x", i, buf)
			if !bytes.Equal(buf, d[:relic.G1UncompressedSize]) {
				t.Errorf("%d: different encoding", i)
			}
			d = d[relic.G1UncompressedSize:]
			ep.Add(one)
		}
	})
	t.Run("G1Compressed", func(t *testing.T) {
		var (
			data = readFile(t, "testdata/g1_compressed_valid_test_vectors.dat")
			ep   = (&relic.EP{}).SetZero()
			a    = &relic.EP{}
			one  = (&relic.EP{}).SetOne()
			d    = data
		)
		for i := 0; i < 1000; i++ {
			t.Logf("%d <- %x", i, d[:relic.G1CompressedSize])
			_, err := a.DecodeCompressed(d[:relic.G1CompressedSize])
			if err != nil {
				t.Errorf("%d: failed decoding: %v", i, err)
			}
			if !ep.Equal(a) {
				t.Errorf("%d: different point", i)
			}
			buf := ep.EncodeCompressed()
			t.Logf("%d -> %x", i, buf)
			if !bytes.Equal(buf, d[:relic.G1CompressedSize]) {
				t.Errorf("%d: different encoding", i)
			}
			d = d[relic.G1CompressedSize:]
			ep.Add(one)
		}
	})
}

func readFile(t *testing.T, name string) []byte {
	t.Helper()
	res, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
