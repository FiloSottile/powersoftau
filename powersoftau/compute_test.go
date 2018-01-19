package powersoftau_test

import (
	"testing"

	"github.com/FiloSottile/powersoftau/powersoftau"
)

func TestShortChallenge(t *testing.T) {
	ch, err := powersoftau.ReadChallenge("testdata/challenge")
	if err != nil {
		t.Fatal(err)
	}

	ch.Compute()

	hash, err := powersoftau.WriteResponse("testdata/response", ch)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%x", hash)
}
