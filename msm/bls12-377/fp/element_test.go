

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fp

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/bits"


	"testing"

	"github.com/leanovate/gopter"
	ggen "github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"github.com/stretchr/testify/require"
)

// -------------------------------------------------------------------------------------------------
// benchmarks
// most benchmarks are rudimentary and should sample a large number of random inputs
// or be run multiple times to ensure it didn't measure the fastest path of the function

var benchResElement Element

func BenchmarkElementSelect(b *testing.B) {
	var x, y Element
	x.SetRandom()
	y.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Select(i%3, &x, &y)
	}
}

func BenchmarkElementSetRandom(b *testing.B) {
	var x Element
	x.SetRandom()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = x.SetRandom()
	}
}

func BenchmarkElementSetBytes(b *testing.B) {
	var x Element
	x.SetRandom()
	bb := x.ZBytes()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchResElement.SetBytes(bb[:])
	}

}







func BenchmarkElementDouble(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Double(&benchResElement)
	}
}

func BenchmarkElementAdd(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Add(&x, &benchResElement)
	}
}

func BenchmarkElementSub(b *testing.B) {
	var x Element
	x.SetRandom()
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Sub(&x, &benchResElement)
	}
}

func BenchmarkElementNeg(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Neg(&benchResElement)
	}
}


func BenchmarkElementFromMont(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.FromMont()
	}
}

func BenchmarkElementToMont(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.ToMont()
	}
}
func BenchmarkElementSquare(b *testing.B) {
	benchResElement.SetRandom()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Square(&benchResElement)
	}
}

func BenchmarkElementSqrt(b *testing.B) {
	var a Element
	a.SetUint64(4)
	a.Neg(&a)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Sqrt(&a)
	}
}

func BenchmarkElementMul(b *testing.B) {
	x := Element{
		13224372171368877346,
		227991066186625457,
		2496666625421784173,
		13825906835078366124,
		9475172226622360569,
		30958721782860680,
	}
	benchResElement.SetOne()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Mul(&benchResElement, &x)
	}
}

func BenchmarkElementCmp(b *testing.B) {
	x := Element{
		13224372171368877346,
		227991066186625457,
		2496666625421784173,
		13825906835078366124,
		9475172226622360569,
		30958721782860680,
	}
	benchResElement = x
	benchResElement[0] = 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchResElement.Cmp(&x)
	}
}

func TestElementCmp(t *testing.T) {
	var x, y Element

	if x.Cmp(&y) != 0 {
		t.Fatal("x == y")
	}

	one := One()
	y.Sub(&y, &one)

	if x.Cmp(&y) != -1 {
		t.Fatal("x < y")
	}
	if y.Cmp(&x) != 1 {
		t.Fatal("x < y")
	}

	x = y
	if x.Cmp(&y) != 0 {
		t.Fatal("x == y")
	}

	x.Sub(&x, &one)
	if x.Cmp(&y) != -1 {
		t.Fatal("x < y")
	}
	if y.Cmp(&x) != 1 {
		t.Fatal("x < y")
	}
}
func TestElementIsRandom(t *testing.T) {
	for i := 0; i < 50; i++ {
		var x, y Element
		x.SetRandom()
		y.SetRandom()
		if x.Equal(&y) {
			t.Fatal("2 random numbers are unlikely to be equal")
		}
	}
}

func TestElementIsUint64(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("reduce should output a result smaller than modulus", prop.ForAll(
		func(v uint64) bool {
			var e Element
			e.SetUint64(v)

			if !e.IsUint64() {
				return false
			}

			return e.Uint64() == v
		},
		ggen.UInt64(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestElementNegZero(t *testing.T) {
	var a, b Element
	b.SetZero()
	for a.IsZero() {
		a.SetRandom()
	}
	a.Neg(&b)
	if !a.IsZero() {
		t.Fatal("neg(0) != 0")
	}
}

// -------------------------------------------------------------------------------------------------
// Gopter tests
// most of them are generated with a template

const (
	nbFuzzShort = 200
	nbFuzz      = 1000
)

// special values to be used in tests
var staticTestValues []Element

func init() {
	staticTestValues = append(staticTestValues, Element{}) // zero
	staticTestValues = append(staticTestValues, One())     // one
	staticTestValues = append(staticTestValues, rSquare)   // r²
	var e, one Element
	one.SetOne()
	e.Sub(&qElement, &one)
	staticTestValues = append(staticTestValues, e) // q - 1
	e.Double(&one)
	staticTestValues = append(staticTestValues, e) // 2

	{
		a := qElement
		a[0]--
		staticTestValues = append(staticTestValues, a)
	}
	staticTestValues = append(staticTestValues, Element{0})
	staticTestValues = append(staticTestValues, Element{0, 0})
	staticTestValues = append(staticTestValues, Element{1})
	staticTestValues = append(staticTestValues, Element{0, 1})
	staticTestValues = append(staticTestValues, Element{2})
	staticTestValues = append(staticTestValues, Element{0, 2})

	{
		a := qElement
		a[5]--
		staticTestValues = append(staticTestValues, a)
	}
	{
		a := qElement
		a[5]--
		a[0]++
		staticTestValues = append(staticTestValues, a)
	}

	{
		a := qElement
		a[5] = 0
		staticTestValues = append(staticTestValues, a)
	}

}

func TestElementReduce(t *testing.T) {
	testValues := make([]Element, len(staticTestValues))
	copy(testValues, staticTestValues)


	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := genFull()

	properties.Property("reduce should output a result smaller than modulus", prop.ForAll(
		func(a Element) bool {
			b := a
			_reduceGeneric(&b)
			return b.smallerThanModulus()
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestElementEqual(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("x.Equal(&y) iff x == y; likely false for random pairs", prop.ForAll(
		func(a testPairElement, b testPairElement) bool {
			return a.element.Equal(&b.element) == (a.element == b.element)
		},
		genA,
		genB,
	))

	properties.Property("x.Equal(&y) if x == y", prop.ForAll(
		func(a testPairElement) bool {
			b := a.element
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestElementBytes(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("SetBytes(Bytes()) should stay constant", prop.ForAll(
		func(a testPairElement) bool {
			var b Element
			bytes := a.element.ZBytes()
			b.ZSetBytes(bytes[:])
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}


func mulByConstant(z *Element, c uint8) {
	var y Element
	y.SetUint64(uint64(c))
	z.Mul(z, &y)
}


func TestElementLegendre(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("legendre should output same result than big.Int.Jacobi", prop.ForAll(
		func(a testPairElement) bool {
			return a.element.Legendre() == big.Jacobi(&a.bigint, Modulus())
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestElementBitLen(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("BitLen should output same result than big.Int.BitLen", prop.ForAll(
		func(a testPairElement) bool {
			return a.element.FromMont().BitLen() == a.bigint.BitLen()
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}



func TestElementLexicographicallyLargest(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("element.Cmp should match LexicographicallyLargest output", prop.ForAll(
		func(a testPairElement) bool {
			var negA Element
			negA.Neg(&a.element)

			cmpResult := a.element.Cmp(&negA)
			lResult := a.element.LexicographicallyLargest()

			if lResult && cmpResult == 1 {
				return true
			}
			if !lResult && cmpResult != 1 {
				return true
			}
			return false
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

}

func TestElementAdd(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("Add: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			d.Set(&a.element)

			c.Add(&a.element, &b.element)
			a.element.Add(&a.element, &b.element)
			b.element.Add(&d, &b.element)

			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Add: operation result must match big.Int result", prop.ForAll(
		func(a, b testPairElement) bool {
			{
				var c Element

				c.Add(&a.element, &b.element)

				var d, e big.Int
				d.Add(&a.bigint, &b.bigint).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]Element, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int
				r.ToBigIntRegular(&rb)

				var c Element
				c.Add(&a.element, &r)
				d.Add(&a.bigint, &rb).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}
			return true
		},
		genA,
		genB,
	))

	properties.Property("Add: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPairElement) bool {
			var c Element

			c.Add(&a.element, &b.element)

			return c.smallerThanModulus()
		},
		genA,
		genB,
	))

	specialValueTest := func() {
		// test special values against special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for _, b := range testValues {

				var bBig, d, e big.Int
				b.ToBigIntRegular(&bBig)

				var c Element
				c.Add(&a, &b)
				d.Add(&aBig, &bBig).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("Add failed special test values")
				}
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

func TestElementSub(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("Sub: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			d.Set(&a.element)

			c.Sub(&a.element, &b.element)
			a.element.Sub(&a.element, &b.element)
			b.element.Sub(&d, &b.element)

			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Sub: operation result must match big.Int result", prop.ForAll(
		func(a, b testPairElement) bool {
			{
				var c Element

				c.Sub(&a.element, &b.element)

				var d, e big.Int
				d.Sub(&a.bigint, &b.bigint).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]Element, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int
				r.ToBigIntRegular(&rb)

				var c Element
				c.Sub(&a.element, &r)
				d.Sub(&a.bigint, &rb).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}
			return true
		},
		genA,
		genB,
	))

	properties.Property("Sub: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPairElement) bool {
			var c Element

			c.Sub(&a.element, &b.element)

			return c.smallerThanModulus()
		},
		genA,
		genB,
	))

	specialValueTest := func() {
		// test special values against special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for _, b := range testValues {

				var bBig, d, e big.Int
				b.ToBigIntRegular(&bBig)

				var c Element
				c.Sub(&a, &b)
				d.Sub(&aBig, &bBig).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("Sub failed special test values")
				}
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

func TestElementMul(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genB := gen()

	properties.Property("Mul: having the receiver as operand should output the same result", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			d.Set(&a.element)

			c.Mul(&a.element, &b.element)
			a.element.Mul(&a.element, &b.element)
			b.element.Mul(&d, &b.element)

			return a.element.Equal(&b.element) && a.element.Equal(&c) && b.element.Equal(&c)
		},
		genA,
		genB,
	))

	properties.Property("Mul: operation result must match big.Int result", prop.ForAll(
		func(a, b testPairElement) bool {
			{
				var c Element

				c.Mul(&a.element, &b.element)

				var d, e big.Int
				d.Mul(&a.bigint, &b.bigint).Mod(&d, Modulus())

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}

			// fixed elements
			// a is random
			// r takes special values
			testValues := make([]Element, len(staticTestValues))
			copy(testValues, staticTestValues)

			for _, r := range testValues {
				var d, e, rb big.Int
				r.ToBigIntRegular(&rb)

				var c Element
				c.Mul(&a.element, &r)
				d.Mul(&a.bigint, &rb).Mod(&d, Modulus())

				// checking generic impl against asm path
				var cGeneric Element
				_mulGeneric(&cGeneric, &a.element, &r)
				if !cGeneric.Equal(&c) {
					// need to give context to failing error.
					return false
				}

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					return false
				}
			}
			return true
		},
		genA,
		genB,
	))

	properties.Property("Mul: operation result must be smaller than modulus", prop.ForAll(
		func(a, b testPairElement) bool {
			var c Element

			c.Mul(&a.element, &b.element)

			return c.smallerThanModulus()
		},
		genA,
		genB,
	))

	properties.Property("Mul: assembly implementation must be consistent with generic one", prop.ForAll(
		func(a, b testPairElement) bool {
			var c, d Element
			c.Mul(&a.element, &b.element)
			_mulGeneric(&d, &a.element, &b.element)
			return c.Equal(&d)
		},
		genA,
		genB,
	))

	specialValueTest := func() {
		// test special values against special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			for _, b := range testValues {

				var bBig, d, e big.Int
				b.ToBigIntRegular(&bBig)

				var c Element
				c.Mul(&a, &b)
				d.Mul(&aBig, &bBig).Mod(&d, Modulus())

				// checking asm against generic impl
				var cGeneric Element
				_mulGeneric(&cGeneric, &a, &b)
				if !cGeneric.Equal(&c) {
					t.Fatal("Mul failed special test values: asm and generic impl don't match")
				}

				if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
					t.Fatal("Mul failed special test values")
				}
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}



func TestElementSquare(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Square: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPairElement) bool {

			var b Element

			b.Square(&a.element)
			a.element.Square(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("Square: operation result must match big.Int result", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Square(&a.element)

			var d, e big.Int
			d.Mul(&a.bigint, &a.bigint).Mod(&d, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("Square: operation result must be smaller than modulus", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Square(&a.element)
			return c.smallerThanModulus()
		},
		genA,
	))

	specialValueTest := func() {
		// test special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c Element
			c.Square(&a)

			var d, e big.Int
			d.Mul(&aBig, &aBig).Mod(&d, Modulus())

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("Square failed special test values")
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}


func TestElementSqrt(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Sqrt: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPairElement) bool {

			b := a.element

			b.Sqrt(&a.element)
			a.element.Sqrt(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("Sqrt: operation result must match big.Int result", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Sqrt(&a.element)

			var d, e big.Int
			d.ModSqrt(&a.bigint, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("Sqrt: operation result must be smaller than modulus", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Sqrt(&a.element)
			return c.smallerThanModulus()
		},
		genA,
	))

	specialValueTest := func() {
		// test special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c Element
			c.Sqrt(&a)

			var d, e big.Int
			d.ModSqrt(&aBig, Modulus())

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("Sqrt failed special test values")
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

func TestElementDouble(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Double: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPairElement) bool {

			var b Element

			b.Double(&a.element)
			a.element.Double(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("Double: operation result must match big.Int result", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Double(&a.element)

			var d, e big.Int
			d.Lsh(&a.bigint, 1).Mod(&d, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("Double: operation result must be smaller than modulus", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Double(&a.element)
			return c.smallerThanModulus()
		},
		genA,
	))

	specialValueTest := func() {
		// test special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c Element
			c.Double(&a)

			var d, e big.Int
			d.Lsh(&aBig, 1).Mod(&d, Modulus())

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("Double failed special test values")
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

func TestElementNeg(t *testing.T) {
	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Neg: having the receiver as operand should output the same result", prop.ForAll(
		func(a testPairElement) bool {

			var b Element

			b.Neg(&a.element)
			a.element.Neg(&a.element)
			return a.element.Equal(&b)
		},
		genA,
	))

	properties.Property("Neg: operation result must match big.Int result", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Neg(&a.element)

			var d, e big.Int
			d.Neg(&a.bigint).Mod(&d, Modulus())

			return c.FromMont().ToBigInt(&e).Cmp(&d) == 0
		},
		genA,
	))

	properties.Property("Neg: operation result must be smaller than modulus", prop.ForAll(
		func(a testPairElement) bool {
			var c Element
			c.Neg(&a.element)
			return c.smallerThanModulus()
		},
		genA,
	))

	specialValueTest := func() {
		// test special values
		testValues := make([]Element, len(staticTestValues))
		copy(testValues, staticTestValues)

		for _, a := range testValues {
			var aBig big.Int
			a.ToBigIntRegular(&aBig)
			var c Element
			c.Neg(&a)

			var d, e big.Int
			d.Neg(&aBig).Mod(&d, Modulus())

			if c.FromMont().ToBigInt(&e).Cmp(&d) != 0 {
				t.Fatal("Neg failed special test values")
			}
		}
	}

	properties.TestingRun(t, gopter.ConsoleReporter(false))
	specialValueTest()

}

func TestElementSetInt64(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("z.SetInt64 must match z.SetString", prop.ForAll(
		func(a testPairElement, v int64) bool {
			c := a.element
			d := a.element

			c.SetInt64(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, ggen.Int64(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestElementSetInterface(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()
	genInt := ggen.Int
	genInt8 := ggen.Int8
	genInt16 := ggen.Int16
	genInt32 := ggen.Int32
	genInt64 := ggen.Int64

	genUint := ggen.UInt
	genUint8 := ggen.UInt8
	genUint16 := ggen.UInt16
	genUint32 := ggen.UInt32
	genUint64 := ggen.UInt64

	properties.Property("z.SetInterface must match z.SetString with int8", prop.ForAll(
		func(a testPairElement, v int8) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genInt8(),
	))

	properties.Property("z.SetInterface must match z.SetString with int16", prop.ForAll(
		func(a testPairElement, v int16) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genInt16(),
	))

	properties.Property("z.SetInterface must match z.SetString with int32", prop.ForAll(
		func(a testPairElement, v int32) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genInt32(),
	))

	properties.Property("z.SetInterface must match z.SetString with int64", prop.ForAll(
		func(a testPairElement, v int64) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genInt64(),
	))

	properties.Property("z.SetInterface must match z.SetString with int", prop.ForAll(
		func(a testPairElement, v int) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genInt(),
	))

	properties.Property("z.SetInterface must match z.SetString with uint8", prop.ForAll(
		func(a testPairElement, v uint8) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genUint8(),
	))

	properties.Property("z.SetInterface must match z.SetString with uint16", prop.ForAll(
		func(a testPairElement, v uint16) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genUint16(),
	))

	properties.Property("z.SetInterface must match z.SetString with uint32", prop.ForAll(
		func(a testPairElement, v uint32) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genUint32(),
	))

	properties.Property("z.SetInterface must match z.SetString with uint64", prop.ForAll(
		func(a testPairElement, v uint64) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genUint64(),
	))

	properties.Property("z.SetInterface must match z.SetString with uint", prop.ForAll(
		func(a testPairElement, v uint) bool {
			c := a.element
			d := a.element

			c.SetInterface(v)
			d.SetString(fmt.Sprintf("%v", v))

			return c.Equal(&d)
		},
		genA, genUint(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))

	{
		assert := require.New(t)
		var e Element
		r, err := e.SetInterface(nil)
		assert.Nil(r)
		assert.Error(err)

		var ptE *Element
		var ptB *big.Int

		r, err = e.SetInterface(ptE)
		assert.Nil(r)
		assert.Error(err)
		ptE = new(Element).SetOne()
		r, err = e.SetInterface(ptE)
		assert.NoError(err)
		assert.True(r.IsOne())

		r, err = e.SetInterface(ptB)
		assert.Nil(r)
		assert.Error(err)

	}
}

func TestElementNewElement(t *testing.T) {
	assert := require.New(t)

	t.Parallel()

	e := NewElement(1)
	assert.True(e.IsOne())

	e = NewElement(0)
	assert.True(e.IsZero())
}


func TestElementFromMont(t *testing.T) {

	t.Parallel()
	parameters := gopter.DefaultTestParameters()
	if testing.Short() {
		parameters.MinSuccessfulTests = nbFuzzShort
	} else {
		parameters.MinSuccessfulTests = nbFuzz
	}

	properties := gopter.NewProperties(parameters)

	genA := gen()

	properties.Property("Assembly implementation must be consistent with generic one", prop.ForAll(
		func(a testPairElement) bool {
			c := a.element
			d := a.element
			c.FromMont()
			_fromMontGeneric(&d)
			return c.Equal(&d)
		},
		genA,
	))

	properties.Property("x.FromMont().ToMont() == x", prop.ForAll(
		func(a testPairElement) bool {
			c := a.element
			c.FromMont().ToMont()
			return c.Equal(&a.element)
		},
		genA,
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

func TestElementJSON(t *testing.T) {
	assert := require.New(t)

	type S struct {
		A Element
		B [3]Element
		C *Element
		D *Element
	}

	// encode to JSON
	var s S
	s.A.SetString("-1")
	s.B[2].SetUint64(42)
	s.D = new(Element).SetUint64(8000)

	encoded, err := json.Marshal(&s)
	assert.NoError(err)
	const expected = "{\"A\":-1,\"B\":[0,0,42],\"C\":null,\"D\":8000}"
	assert.Equal(expected, string(encoded))

	// decode valid
	var decoded S
	err = json.Unmarshal([]byte(expected), &decoded)
	assert.NoError(err)

	assert.Equal(s, decoded, "element -> json -> element round trip failed")

	// decode hex and string values
	withHexValues := "{\"A\":\"-1\",\"B\":[0,\"0x00000\",\"0x2A\"],\"C\":null,\"D\":\"8000\"}"

	var decodedS S
	err = json.Unmarshal([]byte(withHexValues), &decodedS)
	assert.NoError(err)

	assert.Equal(s, decodedS, " json with strings  -> element  failed")

}

type testPairElement struct {
	element Element
	bigint  big.Int
}

func gen() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var g testPairElement

		g.element = Element{
			genParams.NextUint64(),
			genParams.NextUint64(),
			genParams.NextUint64(),
			genParams.NextUint64(),
			genParams.NextUint64(),
			genParams.NextUint64(),
		}
		if qElement[5] != ^uint64(0) {
			g.element[5] %= (qElement[5] + 1)
		}

		for !g.element.smallerThanModulus() {
			g.element = Element{
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
			}
			if qElement[5] != ^uint64(0) {
				g.element[5] %= (qElement[5] + 1)
			}
		}

		g.element.ToBigIntRegular(&g.bigint)
		genResult := gopter.NewGenResult(g, gopter.NoShrinker)
		return genResult
	}
}

func genFull() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {

		genRandomFq := func() Element {
			var g Element

			g = Element{
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
				genParams.NextUint64(),
			}

			if qElement[5] != ^uint64(0) {
				g[5] %= (qElement[5] + 1)
			}

			for !g.smallerThanModulus() {
				g = Element{
					genParams.NextUint64(),
					genParams.NextUint64(),
					genParams.NextUint64(),
					genParams.NextUint64(),
					genParams.NextUint64(),
					genParams.NextUint64(),
				}
				if qElement[5] != ^uint64(0) {
					g[5] %= (qElement[5] + 1)
				}
			}

			return g
		}
		a := genRandomFq()

		var carry uint64
		a[0], carry = bits.Add64(a[0], qElement[0], carry)
		a[1], carry = bits.Add64(a[1], qElement[1], carry)
		a[2], carry = bits.Add64(a[2], qElement[2], carry)
		a[3], carry = bits.Add64(a[3], qElement[3], carry)
		a[4], carry = bits.Add64(a[4], qElement[4], carry)
		a[5], _ = bits.Add64(a[5], qElement[5], carry)

		genResult := gopter.NewGenResult(a, gopter.NoShrinker)
		return genResult
	}
}
