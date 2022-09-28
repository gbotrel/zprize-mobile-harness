package bls12377

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
)

// BLS12-377: A Barreto--Lynn--Scott curve of embedding degree k=12 with seed x₀=9586122913090633729
// 𝔽r: r=8444461749428370424248824938781546531375899335154063827935233455917409239041 (x₀⁴-x₀²+1)
// 𝔽p: p=258664426012969094010652733694893533536393512754914660539884262666720468348340822774968888139573360124440321458177 ((x₀-1)² ⋅ r(x₀)/3+x₀)
// (E/𝔽p): Y²=X³+1
// (Eₜ/𝔽p²): Y² = X³+1/u (D-type twist)
// r ∣ #E(Fp) and r ∣ #Eₜ(𝔽p²)
// Extension fields tower:
//     𝔽p²[u] = 𝔽p/u²+5
//     𝔽p⁶[v] = 𝔽p²/v³-u
//     𝔽p¹²[w] = 𝔽p⁶/w²-v
// optimal Ate loop size: x₀

// bCurveCoeff b coeff of the curve Y²=X³+b
var bCurveCoeff fp.Element

// generators of the r-torsion group, resp. in ker(pi-id), ker(Tr)
var g1Gen G1Jac

var g1GenAff G1Affine

// point at infinity
var g1Infinity G1Jac

// optimal Ate loop counter
var loopCounter [64]int8

var thirdRootOneG1 fp.Element
var four fp.Element

// seed x₀ of the curve
var xGen big.Int

// glvBasis stores R-linearly independent vectors (a,b), (c,d)
// in ker((u,v) → u+vλ[r]), and their determinant
var glvBasis ecc.Lattice

// conversion to twisted Edwards form
// -x²+y² = 1+(-d/a)x²y²
// a = -2√3+3, d = -2√3-3
var sqrtThree fp.Element
var invSqrtMinusA fp.Element
var dCurveCoeffDouble fp.Element

func init() {

	bCurveCoeff.SetUint64(1)

	g1Gen.X.SetString("81937999373150964239938255573465948239988671502647976594219695644855304257327692006745978603320413799295628339695")
	g1Gen.Y.SetString("241266749859715473739788878240585681733927191168601896383759122102112907357779751001206799952863815012735208165030")
	g1Gen.Z.SetString("1")

	// (X,Y,Z) = (1,1,0)
	g1Infinity.X.SetOne()
	g1Infinity.Y.SetOne()

	g1GenAff.FromJacobian(&g1Gen)

	// binary decomposition of x₀ little endian
	loopCounter = [64]int8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1}

	// x₀
	xGen.SetString("9586122913090633729", 10)

	thirdRootOneG1.SetString("80949648264912719408558363140637477264845294720710499478137287262712535938301461879813459410945")

	var lambdaGLV big.Int
	lambdaGLV.SetString("91893752504881257701523279626832445440", 10) //(x₀²-1)
	_r := fr.Modulus()
	ecc.PrecomputeLattice(_r, &lambdaGLV, &glvBasis)

	four.SetUint64(4)

	// conversion to twisted Edwards form
	// √3
	sqrtThree = fp.Element{
		4588006732144632292,
		14697816095396418986,
		15095485345557306380,
		15246065856125797005,
		14023251964588091418,
		94960888053355880,
	}
	// 1/√-a = 1/√{-2√3+3}
	invSqrtMinusA = fp.Element{
		4253980672159819453,
		10543543389341820547,
		9740544029435732375,
		9468753515685554864,
		16322658805964220949,
		99169878199756585,
	}
	// -d/a = 2*(7+4√3)
	dCurveCoeffDouble = fp.Element{
		17459187486984439494,
		16924566387240380666,
		1102674864729460645,
		6167408650503764774,
		10085461570702649662,
		105911569874903690,
	}

}
