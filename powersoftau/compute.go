package powersoftau

import (
	"math/big"
	"sync"

	"github.com/FiloSottile/powersoftau/bls12"
)

func (c *Challenge) Compute(processes int) {
	pub, priv := NewKeypair(c.ChallengeHash[:])
	c.PublicKey = pub

	r := (&big.Int{}).SetBytes(bls12.ScalarOrder())

	tau, alpha, beta := &big.Int{}, &big.Int{}, &big.Int{}
	tau.SetBytes(priv.Tau)
	alpha.SetBytes(priv.Alpha)
	beta.SetBytes(priv.Beta)

	computeRange := func(a, b int) {
		k, ka, kb := &big.Int{}, &big.Int{}, &big.Int{}
		k.Exp(tau, big.NewInt(int64(a)), r)

		for i := a; i < b; i++ {
			c.Accumulator.TauG1[i].ScalarMult(k.Bytes())
			if i < TauPowers {
				c.Accumulator.TauG2[i].ScalarMult(k.Bytes())
				ka.Mul(k, alpha).Mod(ka, r)
				c.Accumulator.AlphaTau[i].ScalarMult(ka.Bytes())
				kb.Mul(k, beta).Mod(kb, r)
				c.Accumulator.BetaTau[i].ScalarMult(kb.Bytes())
			}

			k.Mul(k, tau).Mod(k, r)
		}
	}

	chunk := 1 << 10
	work := make(chan struct{ a, b int })

	var wg sync.WaitGroup
	for i := 0; i < processes; i++ {
		wg.Add(1)
		go func() {
			for job := range work {
				computeRange(job.a, job.b)
			}
			wg.Done()
		}()
	}

	for i := 0; i < TauPowersG1; i += chunk {
		a, b := i, i+chunk
		if b > TauPowersG1 {
			b = TauPowersG1
		}
		work <- struct{ a, b int }{a, b}
	}
	close(work)
	wg.Wait()

	c.Accumulator.BetaG2.ScalarMult(priv.Beta)
}
