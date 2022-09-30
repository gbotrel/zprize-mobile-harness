// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package bls12377

import (
	"math/big"
	"runtime"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
)

// G1Affine point in affine coordinates
type G1Affine struct {
	X, Y fp.Element
}

// G1Jac is a point with fp.Element coordinates
type G1Jac struct {
	X, Y, Z fp.Element
}

// G1EdExtended point in extended coordinates on a twisted Edwards curve (x=X/Z, y=Y/Z, x*y=T/Z)
type G1EdExtended struct {
	X, Y, Z, T fp.Element
}

// G1EdMSM point in custom affine coordinates on a twisted Edwards curve (y-x=X, y+x=Y, 2d*x*y=T)
type G1EdMSM struct {
	X, Y, T fp.Element
}

// -------------------------------------------------------------------------------------------------
// Affine

// Set sets p to the provided point
func (p *G1Affine) Set(a *G1Affine) *G1Affine {
	p.X, p.Y = a.X, a.Y
	return p
}

// IsZero returns true if p=0 false otherwise
func (p *G1Affine) IsZero() bool {
	var one fp.Element
	one.SetOne()
	return p.X.IsZero() && p.Y.Equal(&one)
}

// ScalarMultiplication computes and returns p = a ⋅ s
func (p *G1Affine) ScalarMultiplication(a *G1Affine, s *big.Int) *G1Affine {
	var _p G1Jac
	_p.FromAffine(a)
	_p.mulGLV(&_p, s)
	p.FromJacobian(&_p)
	return p
}

// ScalarMultiplicationAffine computes and returns p = a ⋅ s
// Takes an affine point and returns a Jacobian point (useful for KZG)
func (p *G1Jac) ScalarMultiplicationAffine(a *G1Affine, s *big.Int) *G1Jac {
	p.FromAffine(a)
	p.mulGLV(p, s)
	return p
}

// Add adds two point in affine coordinates.
// This should rarely be used as it is very inefficient compared to Jacobian
func (p *G1Affine) Add(a, b *G1Affine) *G1Affine {
	var p1, p2 G1Jac
	p1.FromAffine(a)
	p2.FromAffine(b)
	p1.AddAssign(&p2)
	p.FromJacobian(&p1)
	return p
}

// Sub subs two point in affine coordinates.
// This should rarely be used as it is very inefficient compared to Jacobian
func (p *G1Affine) Sub(a, b *G1Affine) *G1Affine {
	var p1, p2 G1Jac
	p1.FromAffine(a)
	p2.FromAffine(b)
	p1.SubAssign(&p2)
	p.FromJacobian(&p1)
	return p
}

// Equal tests if two points (in Affine coordinates) are equal
func (p *G1Affine) Equal(a *G1Affine) bool {
	return p.X.Equal(&a.X) && p.Y.Equal(&a.Y)
}

// Neg sets p to -a
func (p *G1Affine) Neg(a *G1Affine) *G1Affine {
	p.X = a.X
	p.Y.Neg(&a.Y)
	return p
}

// FromJacobian rescales a point in Jacobian coord in z=1 plane
func (p *G1Affine) FromJacobian(p1 *G1Jac) *G1Affine {

	var a, b fp.Element

	if p1.Z.IsZero() {
		p.X.SetZero()
		p.Y.SetZero()
		return p
	}

	a.Inverse(&p1.Z)
	b.Square(&a)
	p.X.Mul(&p1.X, &b)
	p.Y.Mul(&p1.Y, &b).Mul(&p.Y, &a)

	return p
}

// String returns the string representation of the point or "O" if it is infinity
func (p *G1Affine) String() string {
	if p.IsInfinity() {
		return "O"
	}
	return "E([" + p.X.String() + "," + p.Y.String() + "])"
}

// IsInfinity checks if the point is infinity
// in affine, it's encoded as (0,0)
// (0,0) is never on the curve for j=0 curves
func (p *G1Affine) IsInfinity() bool {
	return p.X.IsZero() && p.Y.IsZero()
}

// IsOnCurve returns true if p in on the curve
func (p *G1Affine) IsOnCurve() bool {
	var point G1Jac
	point.FromAffine(p)
	return point.IsOnCurve() // call this function to handle infinity point
}

// IsInSubGroup returns true if p is in the correct subgroup, false otherwise
func (p *G1Affine) IsInSubGroup() bool {
	var _p G1Jac
	_p.FromAffine(p)
	return _p.IsInSubGroup()
}

// -------------------------------------------------------------------------------------------------
// Jacobian

// Set sets p to the provided point
func (p *G1Jac) Set(a *G1Jac) *G1Jac {
	p.X, p.Y, p.Z = a.X, a.Y, a.Z
	return p
}

// Equal tests if two points (in Jacobian coordinates) are equal
func (p *G1Jac) Equal(a *G1Jac) bool {

	if p.Z.IsZero() && a.Z.IsZero() {
		return true
	}
	_p := G1Affine{}
	_p.FromJacobian(p)

	_a := G1Affine{}
	_a.FromJacobian(a)

	return _p.X.Equal(&_a.X) && _p.Y.Equal(&_a.Y)
}

// Neg computes -G
func (p *G1Jac) Neg(a *G1Jac) *G1Jac {
	*p = *a
	p.Y.Neg(&a.Y)
	return p
}

// SubAssign subtracts two points on the curve
func (p *G1Jac) SubAssign(a *G1Jac) *G1Jac {
	var tmp G1Jac
	tmp.Set(a)
	tmp.Y.Neg(&tmp.Y)
	p.AddAssign(&tmp)
	return p
}

// AddAssign point addition in montgomery form
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#addition-add-2007-bl
func (p *G1Jac) AddAssign(a *G1Jac) *G1Jac {

	// p is infinity, return a
	if p.Z.IsZero() {
		p.Set(a)
		return p
	}

	// a is infinity, return p
	if a.Z.IsZero() {
		return p
	}

	var Z1Z1, Z2Z2, U1, U2, S1, S2, H, I, J, r, V fp.Element
	Z1Z1.Square(&a.Z)
	Z2Z2.Square(&p.Z)
	U1.Mul(&a.X, &Z2Z2)
	U2.Mul(&p.X, &Z1Z1)
	S1.Mul(&a.Y, &p.Z).
		Mul(&S1, &Z2Z2)
	S2.Mul(&p.Y, &a.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U1.Equal(&U2) && S1.Equal(&S2) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &U1)
	I.Double(&H).
		Square(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &S1).Double(&r)
	V.Mul(&U1, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	S1.Mul(&S1, &J).Double(&S1)
	p.Y.Sub(&p.Y, &S1)
	p.Z.Add(&p.Z, &a.Z)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &Z2Z2).
		Mul(&p.Z, &H)

	return p
}

// AddMixed point addition
// http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-madd-2007-bl
func (p *G1Jac) AddMixed(a *G1Affine) *G1Jac {

	//if a is infinity return p
	if a.IsInfinity() {
		return p
	}
	// p is infinity, return a
	if p.Z.IsZero() {
		p.X = a.X
		p.Y = a.Y
		p.Z.SetOne()
		return p
	}

	var Z1Z1, U2, S2, H, HH, I, J, r, V fp.Element
	Z1Z1.Square(&p.Z)
	U2.Mul(&a.X, &Z1Z1)
	S2.Mul(&a.Y, &p.Z).
		Mul(&S2, &Z1Z1)

	// if p == a, we double instead
	if U2.Equal(&p.X) && S2.Equal(&p.Y) {
		return p.DoubleAssign()
	}

	H.Sub(&U2, &p.X)
	HH.Square(&H)
	I.Double(&HH).Double(&I)
	J.Mul(&H, &I)
	r.Sub(&S2, &p.Y).Double(&r)
	V.Mul(&p.X, &I)
	p.X.Square(&r).
		Sub(&p.X, &J).
		Sub(&p.X, &V).
		Sub(&p.X, &V)
	J.Mul(&J, &p.Y).Double(&J)
	p.Y.Sub(&V, &p.X).
		Mul(&p.Y, &r)
	p.Y.Sub(&p.Y, &J)
	p.Z.Add(&p.Z, &H)
	p.Z.Square(&p.Z).
		Sub(&p.Z, &Z1Z1).
		Sub(&p.Z, &HH)

	return p
}

// Double doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G1Jac) Double(q *G1Jac) *G1Jac {
	p.Set(q)
	p.DoubleAssign()
	return p
}

// DoubleAssign doubles a point in Jacobian coordinates
// https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-3.html#doubling-dbl-2007-bl
func (p *G1Jac) DoubleAssign() *G1Jac {

	var XX, YY, YYYY, ZZ, S, M, T fp.Element

	XX.Square(&p.X)
	YY.Square(&p.Y)
	YYYY.Square(&YY)
	ZZ.Square(&p.Z)
	S.Add(&p.X, &YY)
	S.Square(&S).
		Sub(&S, &XX).
		Sub(&S, &YYYY).
		Double(&S)
	M.Double(&XX).Add(&M, &XX)
	p.Z.Add(&p.Z, &p.Y).
		Square(&p.Z).
		Sub(&p.Z, &YY).
		Sub(&p.Z, &ZZ)
	T.Square(&M)
	p.X = T
	T.Double(&S)
	p.X.Sub(&p.X, &T)
	p.Y.Sub(&S, &p.X).
		Mul(&p.Y, &M)
	YYYY.Double(&YYYY).Double(&YYYY).Double(&YYYY)
	p.Y.Sub(&p.Y, &YYYY)

	return p
}

// ScalarMultiplication computes and returns p = a ⋅ s
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G1Jac) ScalarMultiplication(a *G1Jac, s *big.Int) *G1Jac {
	return p.mulGLV(a, s)
}

// String returns canonical representation of the point in affine coordinates
func (p *G1Jac) String() string {
	_p := G1Affine{}
	_p.FromJacobian(p)
	return _p.String()
}

// FromAffine sets p = Q, p in Jacobian, Q in affine
func (p *G1Jac) FromAffine(Q *G1Affine) *G1Jac {
	if Q.IsInfinity() {
		p.Z.SetZero()
		p.X.SetOne()
		p.Y.SetOne()
		return p
	}
	p.Z.SetOne()
	p.X.Set(&Q.X)
	p.Y.Set(&Q.Y)
	return p
}

// IsOnCurve returns true if p in on the curve
func (p *G1Jac) IsOnCurve() bool {
	var left, right, tmp fp.Element
	left.Square(&p.Y)
	right.Square(&p.X).Mul(&right, &p.X)
	tmp.Square(&p.Z).
		Square(&tmp).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &p.Z).
		Mul(&tmp, &bCurveCoeff)
	right.Add(&right, &tmp)
	return left.Equal(&right)
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

// mulWindowed computes a 2-bits windowed scalar multiplication
func (p *G1Jac) mulWindowed(a *G1Jac, s *big.Int) *G1Jac {

	var res G1Jac
	var ops [3]G1Jac

	res.Set(&g1Infinity)
	ops[0].Set(a)
	ops[1].Double(&ops[0])
	ops[2].Set(&ops[0]).AddAssign(&ops[1])

	b := s.Bytes()
	for i := range b {
		w := b[i]
		mask := byte(0xc0)
		for j := 0; j < 4; j++ {
			res.DoubleAssign().DoubleAssign()
			c := (w & mask) >> (6 - 2*j)
			if c != 0 {
				res.AddAssign(&ops[c-1])
			}
			mask = mask >> 2
		}
	}
	p.Set(&res)

	return p

}

// ϕ assigns p to ϕ(a) where ϕ: (x,y) → (w x,y), and returns p
// where w is a third root of unity in 𝔽p
func (p *G1Jac) phi(a *G1Jac) *G1Jac {
	p.Set(a)
	p.X.Mul(&p.X, &thirdRootOneG1)
	return p
}

// mulGLV computes the scalar multiplication using a windowed-GLV method
// see https://www.iacr.org/archive/crypto2001/21390189.pdf
func (p *G1Jac) mulGLV(a *G1Jac, s *big.Int) *G1Jac {

	var table [15]G1Jac
	var res G1Jac
	var k1, k2 fr.Element

	res.Set(&g1Infinity)

	// table[b3b2b1b0-1] = b3b2 ⋅ ϕ(a) + b1b0*a
	table[0].Set(a)
	table[3].phi(a)

	// split the scalar, modifies ±a, ϕ(a) accordingly
	k := ecc.SplitScalar(s, &glvBasis)

	if k[0].Sign() == -1 {
		k[0].Neg(&k[0])
		table[0].Neg(&table[0])
	}
	if k[1].Sign() == -1 {
		k[1].Neg(&k[1])
		table[3].Neg(&table[3])
	}

	// precompute table (2 bits sliding window)
	// table[b3b2b1b0-1] = b3b2 ⋅ ϕ(a) + b1b0 ⋅ a if b3b2b1b0 != 0
	table[1].Double(&table[0])
	table[2].Set(&table[1]).AddAssign(&table[0])
	table[4].Set(&table[3]).AddAssign(&table[0])
	table[5].Set(&table[3]).AddAssign(&table[1])
	table[6].Set(&table[3]).AddAssign(&table[2])
	table[7].Double(&table[3])
	table[8].Set(&table[7]).AddAssign(&table[0])
	table[9].Set(&table[7]).AddAssign(&table[1])
	table[10].Set(&table[7]).AddAssign(&table[2])
	table[11].Set(&table[7]).AddAssign(&table[3])
	table[12].Set(&table[11]).AddAssign(&table[0])
	table[13].Set(&table[11]).AddAssign(&table[1])
	table[14].Set(&table[11]).AddAssign(&table[2])

	// bounds on the lattice base vectors guarantee that k1, k2 are len(r)/2 or len(r)/2+1 bits long max
	// this is because we use a probabilistic scalar decomposition that replaces a division by a right-shift
	k1.SetBigInt(&k[0]).FromMont()
	k2.SetBigInt(&k[1]).FromMont()

	// we don't target constant-timeness so we check first if we increase the bounds or not
	maxBit := k1.BitLen()
	if k2.BitLen() > maxBit {
		maxBit = k2.BitLen()
	}
	hiWordIndex := (maxBit - 1) / 64

	// loop starts from len(k1)/2 or len(k1)/2+1 due to the bounds
	for i := hiWordIndex; i >= 0; i-- {
		mask := uint64(3) << 62
		for j := 0; j < 32; j++ {
			res.Double(&res).Double(&res)
			b1 := (k1[i] & mask) >> (62 - 2*j)
			b2 := (k2[i] & mask) >> (62 - 2*j)
			if b1|b2 != 0 {
				s := (b2<<2 | b1)
				res.AddAssign(&table[s-1])
			}
			mask = mask >> 2
		}
	}

	p.Set(&res)
	return p
}

// ClearCofactor maps a point in curve to r-torsion
func (p *G1Affine) ClearCofactor(a *G1Affine) *G1Affine {
	var _p G1Jac
	_p.FromAffine(a)
	_p.ClearCofactor(&_p)
	p.FromJacobian(&_p)
	return p
}

// ClearCofactor maps a point in E(Fp) to E(Fp)[r]
func (p *G1Jac) ClearCofactor(a *G1Jac) *G1Jac {
	// cf https://eprint.iacr.org/2019/403.pdf, 5
	var res G1Jac
	res.ScalarMultiplication(a, &xGen).Neg(&res).AddAssign(a)
	p.Set(&res)
	return p

}

// BatchJacobianToAffineG1 converts points in Jacobian coordinates to Affine coordinates
// performing a single field inversion (Montgomery batch inversion trick).
func BatchJacobianToAffineG1(points []G1Jac) []G1Affine {
	result := make([]G1Affine, len(points))
	zeroes := make([]bool, len(points))
	accumulator := fp.One()

	// batch invert all points[].Z coordinates with Montgomery batch inversion trick
	// (stores points[].Z^-1 in result[i].X to avoid allocating a slice of fr.Elements)
	for i := 0; i < len(points); i++ {
		if points[i].Z.IsZero() {
			zeroes[i] = true
			continue
		}
		result[i].X = accumulator
		accumulator.Mul(&accumulator, &points[i].Z)
	}

	var accInverse fp.Element
	accInverse.Inverse(&accumulator)

	for i := len(points) - 1; i >= 0; i-- {
		if zeroes[i] {
			// do nothing, (X=0, Y=0) is infinity point in affine
			continue
		}
		result[i].X.Mul(&result[i].X, &accInverse)
		accInverse.Mul(&accInverse, &points[i].Z)
	}

	// batch convert to affine.
	Execute(len(points), func(start, end int) {
		for i := start; i < end; i++ {
			if zeroes[i] {
				// do nothing, (X=0, Y=0) is infinity point in affine
				continue
			}
			var a, b fp.Element
			a = result[i].X
			b.Square(&a)
			result[i].X.Mul(&points[i].X, &b)
			result[i].Y.Mul(&points[i].Y, &b).
				Mul(&result[i].Y, &a)
		}
	})

	return result
}

// BatchScalarMultiplicationG1 multiplies the same base by all scalars
// and return resulting points in affine coordinates
// uses a simple windowed-NAF like exponentiation algorithm
func BatchScalarMultiplicationG1(base *G1Affine, scalars []fr.Element) []G1Affine {

	// approximate cost in group ops is
	// cost = 2^{c-1} + n(scalar.nbBits+nbChunks)

	nbPoints := uint64(len(scalars))
	min := ^uint64(0)
	bestC := 0
	for c := 2; c < 18; c++ {
		cost := uint64(1 << (c - 1))
		nbChunks := uint64(fr.Limbs * 64 / c)
		if (fr.Limbs*64)%c != 0 {
			nbChunks++
		}
		cost += nbPoints * ((fr.Limbs * 64) + nbChunks)
		if cost < min {
			min = cost
			bestC = c
		}
	}
	c := uint64(bestC) // window size
	nbChunks := int(fr.Limbs * 64 / c)
	if (fr.Limbs*64)%c != 0 {
		nbChunks++
	}
	mask := uint64((1 << c) - 1) // low c bits are 1
	msbWindow := uint64(1 << (c - 1))

	// precompute all powers of base for our window
	// note here that if performance is critical, we can implement as in the msmX methods
	// this allocation to be on the stack
	baseTable := make([]G1Jac, (1 << (c - 1)))
	baseTable[0].Set(&g1Infinity)
	baseTable[0].AddMixed(base)
	for i := 1; i < len(baseTable); i++ {
		baseTable[i] = baseTable[i-1]
		baseTable[i].AddMixed(base)
	}

	pScalars, _ := partitionScalars(scalars, c, false, runtime.NumCPU())

	// compute offset and word selector / shift to select the right bits of our windows
	selectors := make([]selector, nbChunks)
	for chunk := 0; chunk < nbChunks; chunk++ {
		jc := uint64(uint64(chunk) * c)
		d := selector{}
		d.index = jc / 64
		d.shift = jc - (d.index * 64)
		d.mask = mask << d.shift
		d.multiWordSelect = (64%c) != 0 && d.shift > (64-c) && d.index < (fr.Limbs-1)
		if d.multiWordSelect {
			nbBitsHigh := d.shift - uint64(64-c)
			d.maskHigh = (1 << nbBitsHigh) - 1
			d.shiftHigh = (c - nbBitsHigh)
		}
		selectors[chunk] = d
	}
	// convert our base exp table into affine to use AddMixed
	baseTableAff := BatchJacobianToAffineG1(baseTable)
	toReturn := make([]G1Jac, len(scalars))

	// for each digit, take value in the base table, double it c time, voilà.
	Execute(len(pScalars), func(start, end int) {
		var p G1Jac
		for i := start; i < end; i++ {
			p.Set(&g1Infinity)
			for chunk := nbChunks - 1; chunk >= 0; chunk-- {
				s := selectors[chunk]
				if chunk != nbChunks-1 {
					for j := uint64(0); j < c; j++ {
						p.DoubleAssign()
					}
				}

				bits := (pScalars[i][s.index] & s.mask) >> s.shift
				if s.multiWordSelect {
					bits += (pScalars[i][s.index+1] & s.maskHigh) << s.shiftHigh
				}

				if bits == 0 {
					continue
				}

				// if msbWindow bit is set, we need to substract
				if bits&msbWindow == 0 {
					// add
					p.AddMixed(&baseTableAff[bits-1])
				} else {
					// sub
					t := baseTableAff[bits & ^msbWindow]
					t.Neg(&t)
					p.AddMixed(&t)
				}
			}

			// set our result point
			toReturn[i] = p

		}
	})
	toReturnAff := BatchJacobianToAffineG1(toReturn)
	return toReturnAff
}

// -------------------------------------------------------------------------------------------------
// Extended coordinates on twisted Edwards

// FromAffine sets p = a, p in twisted Edwards (extended), a in Short Weierstrass (affine)
func (p *G1EdExtended) FromAffineSW(a *G1Affine) *G1EdExtended {

	var d1, d2, one fp.Element
	one.SetOne()

	d1.Mul(&a.Y, &invSqrtMinusA)
	d2.Add(&a.X, &one).
		Add(&d2, &sqrtThree)

	inv := fp.BatchInvert([]fp.Element{d1, d2})

	p.X.Add(&a.X, &one).
		Mul(&p.X, &inv[0])
	p.Y.Add(&a.X, &one).
		Sub(&p.Y, &sqrtThree).
		Mul(&p.Y, &inv[1])

	p.Z.SetOne()

	p.T.Mul(&p.X, &p.Y)

	return p
}

// BatchFromAffineSW converts a_i from affine short Weierstrass to extended twisted Edwards
// performing a single field inversion (Montgomery batch inversion trick).
func BatchFromAffineSW(a []G1Affine) []G1EdExtended {

	p := make([]G1EdExtended, len(a))
	d := make([]fp.Element, 2*len(a))

	var one fp.Element
	one.SetOne()

	Execute(len(a), func(start, end int) {
		for i := start; i < end; i++ {
			d[i].Mul(&a[i].Y, &invSqrtMinusA)
			d[i+len(a)].Add(&a[i].X, &one).
				Add(&d[i+len(a)], &sqrtThree)
		}
	})

	inv := fp.BatchInvert(d)

	Execute(len(a), func(start, end int) {
		for i := start; i < end; i++ {
			p[i].X.Add(&a[i].X, &one).
				Mul(&p[i].X, &inv[i])
			p[i].Y.Add(&a[i].X, &one).
				Sub(&p[i].Y, &sqrtThree).
				Mul(&p[i].Y, &inv[i+len(a)])

			p[i].Z.SetOne()

			p[i].T.Mul(&p[i].X, &p[i].Y)
		}
	})

	return p
}

// BatchFromAffineSWC converts a_i from affine short Weierstrass to custom twisted Edwards
// performing a single field inversion (Montgomery batch inversion trick).
func BatchFromAffineSWC(a []G1Affine) []G1EdMSM {

	p := make([]G1EdMSM, len(a))
	d := make([]fp.Element, 2*len(a))

	var one fp.Element
	one.SetOne()

	Execute(len(a), func(start, end int) {
		for i := start; i < end; i++ {
			d[i].Mul(&a[i].Y, &invSqrtMinusA)
			d[i+len(a)].Add(&a[i].X, &one).
				Add(&d[i+len(a)], &sqrtThree)
		}
	})

	inv := fp.BatchInvert(d)

	Execute(len(a), func(start, end int) {
		var x, y, t fp.Element
		for i := start; i < end; i++ {
			x.Add(&a[i].X, &one).
				Mul(&x, &inv[i])
			y.Add(&a[i].X, &one).
				Sub(&y, &sqrtThree).
				Mul(&y, &inv[i+len(a)])

			t.Mul(&x, &y).Mul(&t, &dCurveCoeffDouble)
			p[i].X.Sub(&y, &x)
			p[i].Y.Add(&y, &x)
			p[i].T = t
		}
	})

	return p
}

// FromEdExtended converts a point in twisted Edwards from extended (Z=1) to custom coordinates
func (p *G1EdMSM) FromExtendedEd(q *G1EdExtended) *G1EdMSM {
	p.X.Sub(&q.Y, &q.X)               // x = y - x
	p.Y.Add(&q.Y, &q.X)               // x = y + x
	p.T.Mul(&q.T, &dCurveCoeffDouble) // t = t * (2d)

	return p
}

// FromEdExtended converts a point in twisted Edwards (extended) to short Weierstrass (affine)
func (a *G1Affine) FromExtendedEd(p *G1EdExtended) *G1Affine {

	if p.Z.IsZero() {
		a.X.SetZero()
		a.Y.SetZero()
		return a
	}

	var x, y, one, n, d1, d2, d3 fp.Element
	one.SetOne()

	d1.Set(&p.Z)
	d2.Sub(&p.Z, &p.Y)
	d3.Mul(&p.X, &invSqrtMinusA)
	inv := fp.BatchInvert([]fp.Element{d1, d2, d3})
	inv[1].Mul(&inv[1], &p.Z)
	inv[2].Mul(&inv[2], &p.Z)

	x.Mul(&d2, &inv[0])
	y.Mul(&p.Y, &inv[0])

	if x.IsZero() && y.IsOne() {
		a.X.SetZero()
		a.Y.SetZero()
		return a
	}

	if x.IsZero() && y.Neg(&y).IsOne() {
		a.X.SetString("86221475337656364670217577898297844512131170918304886846628087555573489449446940924989629379857786708146773819392") // -1/3
		a.Y.SetZero()
		return a
	}

	n.Add(&one, &y).
		Mul(&n, &inv[1]).
		Mul(&n, &sqrtThree)

	a.X.Sub(&n, &one)

	a.Y.Mul(&n, &inv[2])

	return a
}

// Set sets p to q and return it
func (p *G1EdExtended) Set(q *G1EdExtended) *G1EdExtended {
	p.X.Set(&q.X)
	p.Y.Set(&q.Y)
	p.T.Set(&q.T)
	p.Z.Set(&q.Z)
	return p
}

// setInfinity sets p to O (0:1:1:0)
func (p *G1EdExtended) setInfinity() *G1EdExtended {
	p.X.SetZero()
	p.Y.SetOne()
	p.Z.SetOne()
	p.T.SetZero()
	return p
}

// IsInfinity returns true if p=0 false otherwise
func (p *G1EdExtended) IsInfinity() bool {
	return p.X.IsZero() && p.Y.Equal(&p.Z)
}

// Equal returns true if p=q false otherwise
// If one point is on the affine chart Z=0 it returns false
func (p *G1EdExtended) Equal(q *G1EdExtended) bool {
	if p.Z.IsZero() || q.Z.IsZero() {
		return false
	}
	var pAffine, qAffine G1Affine
	pAffine.FromExtendedEd(p)
	qAffine.FromExtendedEd(q)
	return pAffine.Equal(&qAffine)
}

// Neg set p to -q
func (p *G1EdExtended) Neg(q *G1EdExtended) *G1EdExtended {
	p.Set(q)
	p.X.Neg(&p.X)
	p.T.Neg(&p.T)
	return p
}

// UnifiedMixedAdd adds any two points (p+q) in twisted Edwards extended coordinates when q.Z=1
// adapted from (re-madd):
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-madd-2008-hwcd-3
func (p *G1EdExtended) UnifiedMixedAdd(q *G1EdMSM) {
	if p.IsInfinity() {
		A := q.X
		B := q.Y

		fp.Butterfly(&B, &A)

		p.X.Double(&A)
		p.Y.Double(&B)
		p.T.Mul(&A, &B)
		p.Z = four
		return
	}

	var C, D fp.Element
	A := p.X
	B := p.Y

	C.Mul(&p.T, &q.T)
	D.Double(&p.Z)

	fp.Butterfly(&D, &C)

	fp.Butterfly(&B, &A)

	A.Mul(&A, &q.X)
	B.Mul(&B, &q.Y)

	fp.Butterfly(&B, &A)

	p.X.Mul(&A, &C)
	p.Y.Mul(&D, &B)
	p.T.Mul(&A, &B)
	p.Z.Mul(&C, &D)

}

// UnifiedMixedSub subtracts any two points (p-q) in twisted Edwards extended coordinates when q.Z=1
// adapted from (re-madd):
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-madd-2008-hwcd-3
func (p *G1EdExtended) UnifiedMixedSub(q *G1EdMSM) {
	if p.IsInfinity() {
		A := q.Y
		B := q.X

		fp.Butterfly(&B, &A)

		p.X.Double(&A)
		p.Y.Double(&B)
		p.T.Mul(&A, &B)
		p.Z = four
		return
	}

	var C, D fp.Element
	A := p.X
	B := p.Y

	C.Mul(&p.T, &q.T).
		Neg(&C)
	D.Double(&p.Z)

	fp.Butterfly(&D, &C)

	fp.Butterfly(&B, &A)

	A.Mul(&A, &q.Y)
	B.Mul(&B, &q.X)

	fp.Butterfly(&B, &A)

	p.X.Mul(&A, &C)
	p.Y.Mul(&D, &B)
	p.T.Mul(&A, &B)
	p.Z.Mul(&C, &D)

}

// UnifiedAdd adds any two points (p+q) in twisted Edwards extended coordinates
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-add-2008-hwcd-3
func (p *G1EdExtended) UnifiedAdd(q *G1EdExtended) {

	var A, B, C, D, E, F, G, H, tmp fp.Element

	tmp.Sub(&q.Y, &q.X)
	A.Sub(&p.Y, &p.X).
		Mul(&A, &tmp)
	tmp.Add(&p.Y, &p.X)
	B.Add(&q.Y, &q.X).
		Mul(&B, &tmp)
	C.Mul(&p.T, &q.T).
		Mul(&C, &dCurveCoeffDouble)
	D.Mul(&p.Z, &q.Z).
		Double(&D)

	H = B
	E = A
	fp.Butterfly(&H, &E)
	G = D
	F = C
	fp.Butterfly(&G, &F)

	p.X.Mul(&E, &F)
	p.Y.Mul(&G, &H)
	p.T.Mul(&E, &H)
	p.Z.Mul(&F, &G)

}

// UnifiedReAdd adds any two points (p+q) in twisted Edwards extended coordinates
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#addition-add-2008-hwcd-3
func (p *G1EdExtended) UnifiedReAdd(q1, q2 *G1EdExtended, aux *fp.Element) {

	var A, B, C, D, E, F, G, H, tmp fp.Element

	tmp.Sub(&q2.Y, &q2.X)
	A.Sub(&q1.Y, &q1.X).
		Mul(&A, &tmp)
	tmp.Add(&q1.Y, &q1.X)
	B.Add(&q2.Y, &q2.X).
		Mul(&B, &tmp)
	C.Mul(&q1.T, aux)
	D.Mul(&q1.Z, &q2.Z).
		Double(&D)

	H = B
	E = A
	fp.Butterfly(&H, &E)
	G = D
	F = C
	fp.Butterfly(&G, &F)

	p.X.Mul(&E, &F)
	p.Y.Mul(&G, &H)
	p.T.Mul(&E, &H)
	p.Z.Mul(&F, &G)

}

// DedicatedDouble doubles a point in twisted Edwards extended coordinates
// https://hyperelliptic.org/EFD/g1p/auto-twisted-extended-1.html#doubling-dbl-2008-hwcd
func (p *G1EdExtended) DedicatedDouble(q *G1EdExtended) {

	var A, B, C, D, E, F, G, H fp.Element

	A.Square(&q.X)
	B.Square(&q.Y)
	C.Square(&q.Z).
		Double(&C)
	D.Neg(&A)
	E.Add(&q.X, &q.Y).
		Square(&E).
		Sub(&E, &A).
		Sub(&E, &B)

	G = D
	H = B
	fp.Butterfly(&G, &H)

	F.Sub(&G, &C)

	p.X.Mul(&E, &F)
	p.Y.Mul(&G, &H)
	p.T.Mul(&H, &E)
	p.Z.Mul(&F, &G)

}
