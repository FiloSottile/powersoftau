package powersoftau

import (
	"math/big"

	"github.com/FiloSottile/powersoftau/bls12"
)

func (c *Challenge) Compute() {
	pub, priv := NewKeypair(c.Hash[:])
	c.PublicKey = pub

	r := (&big.Int{}).SetBytes(bls12.ScalarOrder())

	tau, alpha, beta := &big.Int{}, &big.Int{}, &big.Int{}
	tau.SetBytes(priv.Tau)
	alpha.SetBytes(priv.Alpha)
	beta.SetBytes(priv.Beta)

	k, ka, kb := &big.Int{}, &big.Int{}, &big.Int{}
	k.SetInt64(1)

	for i := 0; i < TauPowersG1; i++ {
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

	c.Accumulator.BetaG2.ScalarMult(priv.Beta)
}
