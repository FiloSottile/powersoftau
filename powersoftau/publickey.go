package powersoftau

import (
	"crypto/rand"

	"golang.org/x/crypto/blake2b"

	"github.com/FiloSottile/powersoftau/bls12"
)

type PublicKey struct {
	Tau, Alpha, Beta struct {
		S     *bls12.EP
		Sx    *bls12.EP
		SxG2x *bls12.EP2
	}
}

type PrivateKey struct {
	Tau, Alpha, Beta []byte
}

func NewKeypair(digest []byte) (*PublicKey, *PrivateKey) {
	pub, priv := &PublicKey{}, &PrivateKey{}

	priv.Tau = randomScalar()
	priv.Alpha = randomScalar()
	priv.Beta = randomScalar()

	gen := func(x []byte, personalization byte) struct {
		S     *bls12.EP
		Sx    *bls12.EP
		SxG2x *bls12.EP2
	} {
		s := randomScalar()
		S := (&bls12.EP{}).ScalarBaseMult(s)
		Sx := S.Copy().ScalarMult(x)
		h, _ := blake2b.New512(nil)
		h.Write([]byte{personalization})
		h.Write(digest)
		h.Write(S.EncodeUncompressed())
		h.Write(Sx.EncodeUncompressed())
		SxG2x := HashToG2(h.Sum(nil)).ScalarMult(x)
		return struct {
			S     *bls12.EP
			Sx    *bls12.EP
			SxG2x *bls12.EP2
		}{
			S: S, Sx: Sx, SxG2x: SxG2x,
		}
	}

	pub.Tau = gen(priv.Tau, 0)
	pub.Alpha = gen(priv.Alpha, 1)
	pub.Beta = gen(priv.Beta, 2)

	return pub, priv
}

func randomScalar() []byte {
	for {
		s := make([]byte, 32)
		if _, err := rand.Read(s); err != nil {
			panic(err)
		}
		if bls12.IsScalar(s) {
			return s
		}
	}
}
