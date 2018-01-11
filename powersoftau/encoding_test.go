package powersoftau_test

import (
	"testing"

	"github.com/FiloSottile/powersoftau/powersoftau"
)

func TestReadChallenge(t *testing.T) {
	_, err := powersoftau.ReadChallenge("testdata/challenge")
	if err != nil {
		t.Fatal(err)
	}
}
