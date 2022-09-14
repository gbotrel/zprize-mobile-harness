package bls12377

import (
	"math/big"

	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
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

// point at infinity
var g1Infinity G1Jac

// optimal Ate loop counter
var loopCounter [64]int8

var thirdRootOneG1 fp.Element

// seed x₀ of the curve
var xGen big.Int

func init() {

	bCurveCoeff.SetUint64(1)

	g1Gen.X.SetString("81937999373150964239938255573465948239988671502647976594219695644855304257327692006745978603320413799295628339695")
	g1Gen.Y.SetString("241266749859715473739788878240585681733927191168601896383759122102112907357779751001206799952863815012735208165030")
	g1Gen.Z.SetString("1")

	// (X,Y,Z) = (1,1,0)
	g1Infinity.X.SetOne()
	g1Infinity.Y.SetOne()

	// binary decomposition of x₀ little endian
	loopCounter = [64]int8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1}

	// x₀
	xGen.SetString("9586122913090633729", 10)

	thirdRootOneG1.SetString("80949648264912719408558363140637477264845294720710499478137287262712535938301461879813459410945")
}
