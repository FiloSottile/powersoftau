# Modified from https://github.com/relic-toolkit/relic/issues/55#issuecomment-332984301

param = -0xd201000000010000

# /* x = -(2^63 + 2^62 + 2^60 + 2^57 + 2^48 + 2^16). */
assert(param == -((2**63) + (2**62) + (2**60) + (2**57) + (2**48) + (2**16)))

def r(x):
	return (x**4) - (x**2) + 1

def q(x):
	# /* p = (x^2 - 2x + 1) * (x^4 - x^2 + 1)/3 + x. */
	return (((x - 1) ** 2) * ((x**4) - (x**2) + 1) // 3) + x

def g1_h(x):
	return ((x-1)**2) // 3

def g2_h(x):
	# (x^8 - 4x^7 + 5x^6 - 4x^4 + 6x^3 - 4x^2 - 4x + 13) / 9
	return ((x**8) - (4 * (x**7)) + (5 * (x**6)) - (4 * (x**4)) + (6 * (x**3)) - (4 * (x**2)) - (4*x) + 13) // 9

q = q(param)
r = r(param)

assert(r == 0x73EDA753299D7D483339D80809A1D80553BDA402FFFE5BFEFFFFFFFF00000001) # B12_P381_R
assert(q == 0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab)

Fq = GF(q)

ec = EllipticCurve(Fq, [0, 4])

def psqrt(v):
	assert(not v.is_zero())
	a = sqrt(v)
	b = -a
	if a < b:
		return a
	else:
		return b

assert(g1_h(param) == 0x396C8C005555E1568C00AAAB0000AAAB) # B12_P381_H

for x in range(0,100):
	rhs = Fq(x)^3 + 4
	if rhs.is_square():
		y = psqrt(rhs)
		p = ec(x, y) * g1_h(param)
		if (not p.is_zero()) and (p * r).is_zero():
			Px, Py = p.xy()
			assert(Px == 0x17F1D3A73197D7942695638C4FA9AC0FC3688C4F9774B905A14E3A3F171BAC586C55E83FF97A1AEFFB3AF00ADB22C6BB)
			assert(Py == 0x08B3F481E3AAA0F1A09E30ED741D8AE4FCF5E095D5D00AF600DB18CB2C04B3EDD03CC744A2888AE40CAA232946C5E7E1)
			break

Fqx.<j> = PolynomialRing(Fq, 'j')
Fq2.<i> = GF(q^2, modulus=j^2 + 1)

ec2 = EllipticCurve(Fq2, [0, (4 * (1 + i))])

assert(ec2.order() == (r * g2_h(param)))
assert(g2_h(param) == 0x5d543a95414e7f1091d50792876a202cd91de4547085abaa68a205b2e5a7ddfa628f1cb4d9e82ef21537e293a6691ae1616ec6e786f0c70cf1c38e31c7238e5)

for x in range(0,100):
	rhs = (Fq2(x))^3 + (4 * (1 + i))
	if rhs.is_square():
		y = psqrt(rhs)
		p = ec2(Fq2(x), y) * g2_h(param)
		if (not p.is_zero()) and (p * r).is_zero():
			Px, Py = p.xy()
			assert(Px == 0x13E02B6052719F607DACD3A088274F65596BD0D09920B61AB5DA61BBDC7F5049334CF11213945D57E5AC7D055D042B7E*i + 0x024AA2B2F08F0A91260805272DC51051C6E47AD4FA403B02B4510B647AE3D1770BAC0326A805BBEFD48056C8C121BDB8)
			assert(Py == 0x0606C4A02EA734CC32ACD2B02BC28B99CB3E287E85A763AF267492AB572E99AB3F370D275CEC1DA1AAA9075FF05F79BE*i + 0x0CE5D527727D6E118CC9CDC6DA2E351AADFD9BAA8CBDD3A76D429A695160D12C923AC9CC3BACA289E193548608B82801)
			break

print "All parameters check out!"
