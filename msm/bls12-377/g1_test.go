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
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"

	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestG1EdwardsExtended(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[BLS12-377] check conversion Short Weierstrass affine to/From twisted Edwards extended (Z=1)", prop.ForAll(
		func(a fp.Element) bool {
			var g G1Jac
			var p, _p G1Affine
			var s big.Int
			a.ToBigIntRegular(&s)
			g.ScalarMultiplication(&g1Gen, &s)
			p.FromJacobian(&g)
			var q G1EdExtended
			q.FromAffineSW(&p)
			_p.FromExtendedEd(&q)

			return p.Equal(&_p)

		},
		GenFp(),
	))

	properties.Property("[BLS12-377] Check arithmetic in twisted Edwards extended", prop.ForAll(
		func(a fp.Element) bool {
			var g G1Jac
			var p, p1, p2 G1Affine
			var s big.Int
			a.ToBigIntRegular(&s)
			g.ScalarMultiplication(&g1Gen, &s)
			p.FromJacobian(&g)
			var q, q1, q2 G1EdExtended
			q.FromAffineSW(&p)
			q1.DedicatedDouble(&q)
			q1.UnifiedAdd(&q)
			q2.DedicatedDouble(&q)
			var _q G1EdMSM
			_q.FromExtendedEd(&q)
			q2.UnifiedMixedAdd(&_q)
			p1.FromExtendedEd(&q1)
			p2.FromExtendedEd(&q2)

			three := big.NewInt(3)
			g.ScalarMultiplication(&g, three)
			p.FromJacobian(&g)

			return p.Equal(&p1) && p.Equal(&p2)

		},
		GenFp(),
	))
	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG1AffineIsOnCurve(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[BLS12-377] g1Gen (affine) should be on the curve", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2 G1Affine
			op1.FromJacobian(&g1Gen)
			op2.Set(&op1)
			op2.Y.Mul(&op2.Y, &a)
			return op1.IsOnCurve() && !op2.IsOnCurve()
		},
		GenFp(),
	))

	properties.Property("[BLS12-377] g1Gen (Jacobian) should be on the curve", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2, op3 G1Jac
			op1.Set(&g1Gen)
			op3.Set(&g1Gen)

			op2 = fuzzG1Jac(&g1Gen, a)
			op3.Y.Mul(&op3.Y, &a)
			return op1.IsOnCurve() && op2.IsOnCurve() && !op3.IsOnCurve()
		},
		GenFp(),
	))

	properties.Property("[BLS12-377] IsInSubGroup and MulBy subgroup order should be the same", prop.ForAll(
		func(a fp.Element) bool {
			var op1, op2 G1Jac
			op1 = fuzzG1Jac(&g1Gen, a)
			_r := fr.Modulus()
			op2.ScalarMultiplication(&op1, _r)
			return op1.IsInSubGroup() && op2.Z.IsZero()
		},
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG1AffineConversions(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[BLS12-377] Affine representation should be independent of the Jacobian representative", prop.ForAll(
		func(a fp.Element) bool {
			g := fuzzG1Jac(&g1Gen, a)
			var op1 G1Affine
			op1.FromJacobian(&g)
			return op1.X.Equal(&g1Gen.X) && op1.Y.Equal(&g1Gen.Y)
		},
		GenFp(),
	))



	properties.Property("[BLS12-377] Jacobian representation should be the same as the affine representative", prop.ForAll(
		func(a fp.Element) bool {
			var g G1Jac
			var op1 G1Affine
			op1.X.Set(&g1Gen.X)
			op1.Y.Set(&g1Gen.Y)

			var one fp.Element
			one.SetOne()

			g.FromAffine(&op1)

			return g.X.Equal(&g1Gen.X) && g.Y.Equal(&g1Gen.Y) && g.Z.Equal(&one)
		},
		GenFp(),
	))

	properties.Property("[BLS12-377] Converting affine symbol for infinity to Jacobian should output correct infinity in Jacobian", prop.ForAll(
		func() bool {
			var g G1Affine
			g.X.SetZero()
			g.Y.SetZero()
			var op1 G1Jac
			op1.FromAffine(&g)
			var one, zero fp.Element
			one.SetOne()
			return op1.X.Equal(&one) && op1.Y.Equal(&one) && op1.Z.Equal(&zero)
		},
	))




	properties.Property("[BLS12-377] [Jacobian] Two representatives of the same class should be equal", prop.ForAll(
		func(a, b fp.Element) bool {
			op1 := fuzzG1Jac(&g1Gen, a)
			op2 := fuzzG1Jac(&g1Gen, b)
			return op1.Equal(&op2)
		},
		GenFp(),
		GenFp(),
	))
	properties.Property("[BLS12-377] BatchJacobianToAffineG1 and FromJacobian should output the same result", prop.ForAll(
		func(a, b fp.Element) bool {
			g1 := fuzzG1Jac(&g1Gen, a)
			g2 := fuzzG1Jac(&g1Gen, b)
			var op1, op2 G1Affine
			op1.FromJacobian(&g1)
			op2.FromJacobian(&g2)
			baseTableAff := BatchJacobianToAffineG1([]G1Jac{g1, g2})
			return op1.Equal(&baseTableAff[0]) && op2.Equal(&baseTableAff[1])
		},
		GenFp(),
		GenFp(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG1AffineOps(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 10

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	properties.Property("[BLS12-377] [Jacobian] Add should call double when having adding the same point", prop.ForAll(
		func(a, b fp.Element) bool {
			fop1 := fuzzG1Jac(&g1Gen, a)
			fop2 := fuzzG1Jac(&g1Gen, b)
			var op1, op2 G1Jac
			op1.Set(&fop1).AddAssign(&fop2)
			op2.Double(&fop2)
			return op1.Equal(&op2)
		},
		GenFp(),
		GenFp(),
	))

	properties.Property("[BLS12-377] [Jacobian] Adding the opposite of a point to itself should output inf", prop.ForAll(
		func(a, b fp.Element) bool {
			fop1 := fuzzG1Jac(&g1Gen, a)
			fop2 := fuzzG1Jac(&g1Gen, b)
			fop2.Neg(&fop2)
			fop1.AddAssign(&fop2)
			return fop1.Equal(&g1Infinity)
		},
		GenFp(),
		GenFp(),
	))

	properties.Property("[BLS12-377] [Jacobian] Adding the inf to a point should not modify the point", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzG1Jac(&g1Gen, a)
			fop1.AddAssign(&g1Infinity)
			var op2 G1Jac
			op2.Set(&g1Infinity)
			op2.AddAssign(&g1Gen)
			return fop1.Equal(&g1Gen) && op2.Equal(&g1Gen)
		},
		GenFp(),
	))

	

	properties.Property("[BLS12-377] [Jacobian] Addmix the negation to itself should output 0", prop.ForAll(
		func(a fp.Element) bool {
			fop1 := fuzzG1Jac(&g1Gen, a)
			fop1.Neg(&fop1)
			var op2 G1Affine
			op2.FromJacobian(&g1Gen)
			fop1.AddMixed(&op2)
			return fop1.Equal(&g1Infinity)
		},
		GenFp(),
	))

	properties.Property("[BLS12-377] scalar multiplication (double and add) should depend only on the scalar mod r", prop.ForAll(
		func(s fr.Element) bool {

			r := fr.Modulus()
			var g G1Jac
			g.mulGLV(&g1Gen, r)

			var scalar, blindedScalar, rminusone big.Int
			var op1, op2, op3, gneg G1Jac
			rminusone.SetUint64(1).Sub(r, &rminusone)
			op3.mulWindowed(&g1Gen, &rminusone)
			gneg.Neg(&g1Gen)
			s.ToBigIntRegular(&scalar)
			blindedScalar.Mul(&scalar, r).Add(&blindedScalar, &scalar)
			op1.mulWindowed(&g1Gen, &scalar)
			op2.mulWindowed(&g1Gen, &blindedScalar)

			return op1.Equal(&op2) && g.Equal(&g1Infinity) && !op1.Equal(&g1Infinity) && gneg.Equal(&op3)

		},
		genScalar,
	))

	properties.Property("[BLS12-377] scalar multiplication (GLV) should depend only on the scalar mod r", prop.ForAll(
		func(s fr.Element) bool {

			r := fr.Modulus()
			var g G1Jac
			g.mulGLV(&g1Gen, r)

			var scalar, blindedScalar, rminusone big.Int
			var op1, op2, op3, gneg G1Jac
			rminusone.SetUint64(1).Sub(r, &rminusone)
			op3.ScalarMultiplication(&g1Gen, &rminusone)
			gneg.Neg(&g1Gen)
			s.ToBigIntRegular(&scalar)
			blindedScalar.Mul(&scalar, r).Add(&blindedScalar, &scalar)
			op1.ScalarMultiplication(&g1Gen, &scalar)
			op2.ScalarMultiplication(&g1Gen, &blindedScalar)

			return op1.Equal(&op2) && g.Equal(&g1Infinity) && !op1.Equal(&g1Infinity) && gneg.Equal(&op3)

		},
		genScalar,
	))

	properties.Property("[BLS12-377] GLV and Double and Add should output the same result", prop.ForAll(
		func(s fr.Element) bool {

			var r big.Int
			var op1, op2 G1Jac
			s.ToBigIntRegular(&r)
			op1.mulWindowed(&g1Gen, &r)
			op2.mulGLV(&g1Gen, &r)
			return op1.Equal(&op2) && !op1.Equal(&g1Infinity)

		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestG1AffineCofactorCleaning(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("[BLS12-377] Clearing the cofactor of a random point should set it in the r-torsion", prop.ForAll(
		func() bool {
			var a, x, b fp.Element
			a.SetRandom()

			x.Square(&a).Mul(&x, &a).Add(&x, &bCurveCoeff)

			for x.Legendre() != 1 {
				a.SetRandom()

				x.Square(&a).Mul(&x, &a).Add(&x, &bCurveCoeff)

			}

			b.Sqrt(&x)
			var point, pointCleared, infinity G1Jac
			point.X.Set(&a)
			point.Y.Set(&b)
			point.Z.SetOne()
			pointCleared.ClearCofactor(&point)
			infinity.Set(&g1Infinity)
			return point.IsOnCurve() && pointCleared.IsInSubGroup() && !pointCleared.Equal(&infinity)
		},
	))
	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestG1AffineBatchScalarMultiplication(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzzShort
	}

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	// size of the multiExps
	const nbSamples = 10

	properties.Property("[BLS12-377] BatchScalarMultiplication should be consistent with individual scalar multiplications", prop.ForAll(
		func(mixer fr.Element) bool {
			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer).
					FromMont()
			}

			result := BatchScalarMultiplicationG1(&g1GenAff, sampleScalars[:])

			if len(result) != len(sampleScalars) {
				return false
			}

			for i := 0; i < len(result); i++ {
				var expectedJac G1Jac
				var expected G1Affine
				var b big.Int
				expectedJac.mulGLV(&g1Gen, sampleScalars[i].ToBigInt(&b))
				expected.FromJacobian(&expectedJac)
				if !result[i].Equal(&expected) {
					return false
				}
			}
			return true
		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// ------------------------------------------------------------
// benches

func BenchmarkG1JacIsInSubGroup(b *testing.B) {
	var a G1Jac
	a.Set(&g1Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.IsInSubGroup()
	}

}

func BenchmarkG1AffineBatchScalarMultiplication(b *testing.B) {
	// ensure every words of the scalars are filled
	var mixer fr.Element
	mixer.SetString("7716837800905789770901243404444209691916730933998574719964609384059111546487")

	const pow = 15
	const nbSamples = 1 << pow

	var sampleScalars [nbSamples]fr.Element

	for i := 1; i <= nbSamples; i++ {
		sampleScalars[i-1].SetUint64(uint64(i)).
			Mul(&sampleScalars[i-1], &mixer).
			FromMont()
	}

	for i := 5; i <= pow; i++ {
		using := 1 << i

		b.Run(fmt.Sprintf("%d points", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				_ = BatchScalarMultiplicationG1(&g1GenAff, sampleScalars[:using])
			}
		})
	}
}

func BenchmarkG1JacScalarMultiplication(b *testing.B) {

	var scalar big.Int
	r := fr.Modulus()
	scalar.SetString("5243587517512619047944770508185965837690552500527637822603658699938581184513", 10)
	scalar.Add(&scalar, r)

	var doubleAndAdd G1Jac

	b.Run("double and add", func(b *testing.B) {
		b.ResetTimer()
		for j := 0; j < b.N; j++ {
			doubleAndAdd.mulWindowed(&g1Gen, &scalar)
		}
	})

	var glv G1Jac
	b.Run("GLV", func(b *testing.B) {
		b.ResetTimer()
		for j := 0; j < b.N; j++ {
			glv.mulGLV(&g1Gen, &scalar)
		}
	})

}

func BenchmarkG1AffineCofactorClearing(b *testing.B) {
	var a G1Jac
	a.Set(&g1Gen)
	for i := 0; i < b.N; i++ {
		a.ClearCofactor(&a)
	}
}

func BenchmarkG1JacAdd(b *testing.B) {
	var a G1Jac
	a.Double(&g1Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.AddAssign(&g1Gen)
	}
}

func BenchmarkG1JacAddMixed(b *testing.B) {
	var a G1Jac
	a.Double(&g1Gen)

	var c G1Affine
	c.FromJacobian(&g1Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.AddMixed(&c)
	}

}

func BenchmarkG1JacDouble(b *testing.B) {
	var a G1Jac
	a.Set(&g1Gen)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.DoubleAssign()
	}

}

func BenchmarkG1EdExtDedicatedDouble(b *testing.B) {
	var g G1Jac
	g.Double(&g1Gen)
	var p G1Affine
	p.FromJacobian(&g)
	var q G1EdExtended
	q.FromAffineSW(&p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.DedicatedDouble(&q)
	}
}

func BenchmarkG1EdExtUnifiedAdd(b *testing.B) {
	var g G1Jac
	g.Double(&g1Gen)
	var p1, p2 G1Affine
	p1.FromJacobian(&g)
	p2.FromJacobian(&g1Gen)
	var q1, q2 G1EdExtended
	q1.FromAffineSW(&p1)
	q2.FromAffineSW(&p2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q1.UnifiedAdd(&q2)
	}
}

func BenchmarkG1EdExtUnifiedMixedAdd(b *testing.B) {
	var g G1Jac
	g.Double(&g1Gen)
	var p1, p2 G1Affine
	p1.FromJacobian(&g)
	p2.FromJacobian(&g1Gen)
	var q1, q2 G1EdExtended
	q1.FromAffineSW(&p1)
	q2.FromAffineSW(&p2)
	var _q2 G1EdMSM
	_q2.FromExtendedEd(&q2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q1.UnifiedMixedAdd(&_q2)
	}
}

func fuzzG1Jac(p *G1Jac, f fp.Element) G1Jac {
	var res G1Jac
	res.X.Mul(&p.X, &f).Mul(&res.X, &f)
	res.Y.Mul(&p.Y, &f).Mul(&res.Y, &f).Mul(&res.Y, &f)
	res.Z.Mul(&p.Z, &f)
	return res
}


func TestBatchAffineConv(t *testing.T) {
	var points [5]G1Affine

	points[0] = g1GenAff
	for i := 1; i < len(points); i++ {
		points[i].Add(&points[i-1], &points[i-1])
	}

	var p1, p2 []G1EdExtended
	p1 = make([]G1EdExtended, 5)

	for i := 0; i < len(points); i++ {
		p1[i].FromAffineSW(&points[i])
	}

	p2 = BatchFromAffineSW(points[:])

	if !reflect.DeepEqual(p1, p2) {
		t.Fatal("invalid batch conversion")
	}
}
