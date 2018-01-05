package relic_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/FiloSottile/powersoftau/relic"
)

func TestG1UncompressedVector(t *testing.T) {
	data := readFile(t, "testdata/g1_uncompressed_valid_test_vectors.dat")
	ep := (&relic.EP{}).SetZero()
	a := &relic.EP{}
	one := (&relic.EP{}).SetOne()
	d := data
	var v []byte
	for i := 0; i < 1000; i++ {
		_, err := a.DecodeUncompressed(d[:relic.G1UncompressedSize])
		if err != nil {
			t.Errorf("%d: failed decoding: %v", i, err)
		}
		d = d[relic.G1UncompressedSize:]
		if !ep.Equal(a) {
			t.Errorf("%d: different point", i)
		}
		v = append(v, ep.EncodeUncompressed()...)
		ep.Add(one)
	}
	if !bytes.Equal(data, v) {
		t.Error("different result")
	}
}

func readFile(t *testing.T, name string) []byte {
	t.Helper()
	res, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return res
}
