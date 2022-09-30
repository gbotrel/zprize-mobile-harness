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
	"math/bits"
	"runtime"
	"sync"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func TestMultiExpG1(t *testing.T) {

	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = 2
	} else {
		parameters.MinSuccessfulTests = nbFuzzShort
	}

	properties := gopter.NewProperties(parameters)

	genScalar := GenFr()

	// size of the multiExps
	const nbSamples = 73

	// multi exp points
	var samplePoints [nbSamples]G1Affine
	var g G1Jac
	g.Set(&g1Gen)
	for i := 1; i <= nbSamples; i++ {
		samplePoints[i-1].FromJacobian(&g)
		g.AddAssign(&g1Gen)
	}

	// sprinkle some points at infinity
	samplePoints[0].setInfinity()
	samplePoints[17].setInfinity()
	samplePoints[34].setInfinity()
	samplePoints[3].setInfinity()
	samplePoints[72].setInfinity()

	// final scalar to use in double and add method (without mixer factor)
	// n(n+1)(2n+1)/6  (sum of the squares from 1 to n)
	var scalar big.Int
	scalar.SetInt64(nbSamples)
	scalar.Mul(&scalar, new(big.Int).SetInt64(nbSamples+1))
	scalar.Mul(&scalar, new(big.Int).SetInt64(2*nbSamples+1))
	scalar.Div(&scalar, new(big.Int).SetInt64(6))

	// ensure a multiexp that's splitted has the same result as a non-splitted one..
	properties.Property("[G1] Multi exponentation (c=16) should be consistent with splitted multiexp", prop.ForAll(
		func(mixer fr.Element) bool {
			var samplePointsLarge [nbSamples * 13]G1Affine
			for i := 0; i < 13; i++ {
				copy(samplePointsLarge[i*nbSamples:], samplePoints[:])
			}

			var r16, splitted1, splitted2 G1EdExtended

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples * 13]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer).
					FromMont()
			}

			scalars16, _ := partitionScalars(sampleScalars[:], 16, false, runtime.NumCPU())
			r16.msmC16(BatchFromAffineSWC(samplePoints[:]), scalars16, true)

			splitted1.MultiExp(BatchFromAffineSWC(samplePointsLarge[:]), sampleScalars[:], ecc.MultiExpConfig{NbTasks: 128})
			splitted2.MultiExp(BatchFromAffineSWC(samplePointsLarge[:]), sampleScalars[:], ecc.MultiExpConfig{NbTasks: 51})
			return r16.Equal(&splitted1) && r16.Equal(&splitted2)
		},
		genScalar,
	))

	cRange := []uint64{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	if testing.Short() {
		// test only "odd" and "even" (ie windows size divide word size vs not)
		cRange = []uint64{5, 16}
	}

	properties.Property(fmt.Sprintf("[G1] Multi exponentation (c in %v) should be consistent with sum of square", cRange), prop.ForAll(
		func(mixer fr.Element) bool {

			var expected G1Affine

			// compute expected result with double and add
			var finalScalar, mixerBigInt big.Int
			finalScalar.Mul(&scalar, mixer.ToBigIntRegular(&mixerBigInt))
			expected.ScalarMultiplication(&g1GenAff, &finalScalar)

			// mixer ensures that all the words of a fpElement are set
			var sampleScalars [nbSamples]fr.Element

			for i := 1; i <= nbSamples; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer).
					FromMont()
			}

			results := make([]G1EdExtended, len(cRange)+1)
			for i, c := range cRange {
				scalars, _ := partitionScalars(sampleScalars[:], c, false, runtime.NumCPU())
				msmInnerg1EdExtended(&results[i], int(c), BatchFromAffineSWC(samplePoints[:]), scalars, false)
				if c == 16 {
					// split the first chunk
					msmInnerg1EdExtended(&results[len(results)-1], 16, BatchFromAffineSWC(samplePoints[:]), scalars, true)
				}
			}
			for i := 1; i < len(results); i++ {
				if !results[i].Equal(&results[i-1]) {
					return false
				}
			}
			return true
		},
		genScalar,
	))

	// note : this test is here as we expect to have a different multiExp than the above bucket method
	// for small number of points
	properties.Property("[G1] Multi exponentation (<50points) should be consistant with sum of square", prop.ForAll(
		func(mixer fr.Element) bool {

			var g G1Jac
			g.Set(&g1Gen)

			// mixer ensures that all the words of a fpElement are set
			samplePoints := make([]G1Affine, 30)
			sampleScalars := make([]fr.Element, 30)

			for i := 1; i <= 30; i++ {
				sampleScalars[i-1].SetUint64(uint64(i)).
					Mul(&sampleScalars[i-1], &mixer).
					FromMont()
				samplePoints[i-1].FromJacobian(&g)
				g.AddAssign(&g1Gen)
			}

			var op1MultiExp G1Affine
			op1MultiExp.MultiExp(samplePoints, sampleScalars, ecc.MultiExpConfig{})

			var finalBigScalar fr.Element
			var finalBigScalarBi big.Int
			var op1ScalarMul G1Affine
			finalBigScalar.SetUint64(9455).Mul(&finalBigScalar, &mixer)
			finalBigScalar.ToBigIntRegular(&finalBigScalarBi)
			op1ScalarMul.ScalarMultiplication(&g1GenAff, &finalBigScalarBi)

			return op1ScalarMul.Equal(&op1MultiExp)
		},
		genScalar,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func BenchmarkMultiExpG1(b *testing.B) {

	const (
		pow       = (bits.UintSize / 2) - (bits.UintSize / 8) // 24 on 64 bits arch, 12 on 32 bits
		nbSamples = 1 << pow
	)

	var (
		samplePoints  [nbSamples]G1Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG1(samplePoints[:])

	var testPoint G1Affine

	for i := 5; i <= pow; i++ {
		using := 1 << i

		b.Run(fmt.Sprintf("%d points", using), func(b *testing.B) {
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				testPoint.MultiExp(samplePoints[:using], sampleScalars[:using], ecc.MultiExpConfig{})
			}
		})
	}
}

func BenchmarkMultiExpG1Reference1(b *testing.B) {
	// G1EdExtended <- MSM(G1EdCustom[:])
	const nbSamples = 1 << 16

	var (
		samplePoints  [nbSamples]G1EdMSM
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG12(samplePoints[:])

	var testPoint G1EdExtended

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		testPoint.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
	}
}

func BenchmarkMultiExpG1Reference2(b *testing.B) {
	// G1Affine <- MSM(G1EdCustom[:])
	const nbSamples = 1 << 16

	var (
		samplePoints  [nbSamples]G1EdMSM
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG12(samplePoints[:])

	var testPoint G1EdExtended
	var testPointAff G1Affine

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		testPoint.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
		testPointAff.FromExtendedEd(&testPoint)
	}
}

func BenchmarkMultiExpG1Reference3(b *testing.B) {
	// G1Affine <- MSM(G1Affine[:])
	const nbSamples = 1 << 16

	var (
		samplePoints  [nbSamples]G1Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG1(samplePoints[:])

	var testPoint G1Affine

	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		testPoint.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
	}
}

func BenchmarkManyMultiExpG1Reference(b *testing.B) {
	const nbSamples = 1 << 20

	var (
		samplePoints  [nbSamples]G1Affine
		sampleScalars [nbSamples]fr.Element
	)

	fillBenchScalars(sampleScalars[:])
	fillBenchBasesG1(samplePoints[:])

	var t1, t2, t3 G1Affine
	b.ResetTimer()
	for j := 0; j < b.N; j++ {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			t1.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t2.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		go func() {
			t3.MultiExp(samplePoints[:], sampleScalars[:], ecc.MultiExpConfig{})
			wg.Done()
		}()
		wg.Wait()
	}
}

// WARNING: this return points that are NOT on the curve and is meant to be use for benchmarking
// purposes only. We don't check that the result is valid but just measure "computational complexity".
//
// Rationale for generating points that are not on the curve is that for large benchmarks, generating
// a vector of different points can take minutes. Using the same point or subset will bias the benchmark result
// since bucket additions in extended jacobian coordinates will hit doubling algorithm instead of add.
func fillBenchBasesG1(samplePoints []G1Affine) {
	var r big.Int
	r.SetString("340444420969191673093399857471996460938405", 10)
	samplePoints[0].ScalarMultiplication(&samplePoints[0], &r)

	one := samplePoints[0].X
	one.SetOne()

	for i := 1; i < len(samplePoints); i++ {
		samplePoints[i].X.Add(&samplePoints[i-1].X, &one)
		samplePoints[i].Y.Sub(&samplePoints[i-1].Y, &one)
	}
}

func fillBenchBasesG12(samplePoints []G1EdMSM) {
	samplePoints[0].X.SetRandom()
	samplePoints[0].Y.SetRandom()
	samplePoints[0].T.SetRandom()

	one := samplePoints[0].X
	one.SetOne()

	for i := 1; i < len(samplePoints); i++ {
		samplePoints[i].X.Add(&samplePoints[i-1].X, &one)
		samplePoints[i].Y.Sub(&samplePoints[i-1].Y, &one)
		samplePoints[i].T.Neg(&samplePoints[i-1].T).Add(&samplePoints[i].T, &samplePoints[i].X)
	}
}

func fillBenchScalars(sampleScalars []fr.Element) {
	// ensure every words of the scalars are filled
	var mixer fr.Element
	mixer.SetString("7716837800905789770901243404444209691916730933998574719964609384059111546487")
	for i := 1; i <= len(sampleScalars); i++ {
		sampleScalars[i-1].SetUint64(uint64(i)).
			Mul(&sampleScalars[i-1], &mixer).
			FromMont()
	}
}
