package powersoftau

import (
	"errors"
	"io"
	"os"

	"github.com/FiloSottile/powersoftau/bls12"
	"golang.org/x/crypto/blake2b"
)

var (
	TauPowers     = 1 << 10 // 1 << 21
	TauPowersG1   = TauPowers<<1 - 1
	ChallengeSize = TauPowersG1*bls12.G1UncompressedSize + // G1 powers
		TauPowers*bls12.G2UncompressedSize + // G2 powers
		TauPowers*bls12.G1UncompressedSize + // alpha powers
		TauPowers*bls12.G1UncompressedSize + // beta powers
		bls12.G2UncompressedSize + // beta
		blake2b.Size
	PublicKeySize = 3*bls12.G2UncompressedSize + 6*bls12.G1UncompressedSize
	ResponseSize  = TauPowersG1*bls12.G1CompressedSize + // G1 powers
		TauPowers*bls12.G2CompressedSize + // G2 powers
		TauPowers*bls12.G1CompressedSize + // alpha powers
		TauPowers*bls12.G1CompressedSize + // beta powers
		bls12.G2CompressedSize + // beta
		blake2b.Size + PublicKeySize
)

type Challenge struct {
	PreviousHash [blake2b.Size]byte
	Accumulator  *Accumulator
	Hash         [blake2b.Size]byte
}

func ReadChallenge(filename string) (*Challenge, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if fi.Size() != int64(ChallengeSize) {
		return nil, errors.New("the challenge file has the wrong size")
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	h, _ := blake2b.New512(nil)
	r := io.TeeReader(f, h)

	c := &Challenge{}
	if _, err := io.ReadFull(r, c.PreviousHash[:]); err != nil {
		return nil, err
	}
	if _, err := ReadAccumulator(r, false); err != nil {
		return nil, err
	}
	h.Sum(c.Hash[:0])

	return c, nil
}

type Accumulator struct {
	TauG1    []*bls12.EP
	TauG2    []*bls12.EP2
	AlphaTau []*bls12.EP
	BetaTau  []*bls12.EP
	BetaG2   *bls12.EP2
}

func ReadAccumulator(r io.Reader, compressed bool) (*Accumulator, error) {
	a := &Accumulator{}
	var err error
	a.TauG1, err = readG1Slice(r, TauPowersG1, compressed)
	if err != nil {
		return nil, err
	}
	a.TauG2, err = readG2Slice(r, TauPowers, compressed)
	if err != nil {
		return nil, err
	}
	a.AlphaTau, err = readG1Slice(r, TauPowers, compressed)
	if err != nil {
		return nil, err
	}
	a.BetaTau, err = readG1Slice(r, TauPowers, compressed)
	if err != nil {
		return nil, err
	}
	pp, err := readG2Slice(r, 1, compressed)
	if err != nil {
		return nil, err
	}
	a.BetaG2 = pp[0]
	return a, nil
}

func readG1Slice(r io.Reader, n int, compressed bool) ([]*bls12.EP, error) {
	var buf []byte
	if compressed {
		buf = make([]byte, bls12.G1CompressedSize)
	} else {
		buf = make([]byte, bls12.G1UncompressedSize)
	}
	var res []*bls12.EP
	for i := 0; i < n; i++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		p := &bls12.EP{}
		var err error
		if compressed {
			p, err = p.DecodeCompressed(buf)
		} else {
			p, err = p.DecodeUncompressed(buf)
		}
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func readG2Slice(r io.Reader, n int, compressed bool) ([]*bls12.EP2, error) {
	var buf []byte
	if compressed {
		buf = make([]byte, bls12.G2CompressedSize)
	} else {
		buf = make([]byte, bls12.G2UncompressedSize)
	}
	var res []*bls12.EP2
	for i := 0; i < n; i++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		p := bls12.NewEP2()
		var err error
		if compressed {
			p, err = p.DecodeCompressed(buf)
		} else {
			p, err = p.DecodeUncompressed(buf)
		}
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (a *Accumulator) WriteTo(w io.Writer, compressed bool) error {
	if err := writeG1Slice(w, a.TauG1, compressed); err != nil {
		return err
	}
	if err := writeG2Slice(w, a.TauG2, compressed); err != nil {
		return err
	}
	if err := writeG1Slice(w, a.AlphaTau, compressed); err != nil {
		return err
	}
	if err := writeG1Slice(w, a.BetaTau, compressed); err != nil {
		return err
	}
	if err := writeG2Slice(w, []*bls12.EP2{a.BetaG2}, compressed); err != nil {
		return err
	}
	return nil
}

func writeG1Slice(w io.Writer, s []*bls12.EP, compressed bool) error {
	for _, p := range s {
		var buf []byte
		if compressed {
			buf = p.EncodeCompressed()
		} else {
			buf = p.EncodeUncompressed()
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

func writeG2Slice(w io.Writer, s []*bls12.EP2, compressed bool) error {
	for _, p := range s {
		var buf []byte
		if compressed {
			buf = p.EncodeCompressed()
		} else {
			buf = p.EncodeUncompressed()
		}
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
