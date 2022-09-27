package bls12377

// SizeOfG1AffineCompressed represents the size in bytes that a G1Affine need in binary form, compressed
const SizeOfG1AffineCompressed = 48

// SizeOfG1AffineUncompressed represents the size in bytes that a G1Affine need in binary form, uncompressed
const SizeOfG1AffineUncompressed = SizeOfG1AffineCompressed * 2

// ϕ assigns p to ϕ(a) where ϕ: (x,y) → (w x,y), and returns p
// where w is a third root of unity in 𝔽p
func (p *G1Jac) phi(a *G1Jac) *G1Jac {
	p.Set(a)
	p.X.Mul(&p.X, &thirdRootOneG1)
	return p
}

// IsInSubGroup returns true if p is on the r-torsion, false otherwise.
// Z[r,0]+Z[-lambdaG1Affine, 1] is the kernel
// of (u,v)->u+lambdaG1Affinev mod r. Expressing r, lambdaG1Affine as
// polynomials in x, a short vector of this Zmodule is
// 1, x². So we check that p+x²ϕ(p)
// is the infinity.
func (p *G1Jac) IsInSubGroup() bool {

	var res G1Jac
	res.phi(p).
		ScalarMultiplication(&res, &xGen).
		ScalarMultiplication(&res, &xGen).
		AddAssign(p)

	return res.IsOnCurve() && res.Z.IsZero()

}

// IsInSubGroup returns true if p is in the correct subgroup, false otherwise
func (p *G1Affine) IsInSubGroup() bool {
	var _p G1Jac
	_p.FromAffine(p)
	return _p.IsInSubGroup()
}
