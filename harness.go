package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	bls12377 "github.com/gbotrel/zprize-mobile-harness/msm/bls12-377"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
	"github.com/icza/gox/timex"
)

// benchmarkMSM follows closely the test harness
func benchmarkMSM(outDir string, instances []Instance, nbIterations int) ([]time.Duration, error) {
	var results []time.Duration

	fTimes, err := os.OpenFile(filepath.Join(outDir, "resulttimes.txt"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer fTimes.Close()

	fResults, err := os.Create(filepath.Join(outDir, "gnark_result.txt"))
	if err != nil {
		return nil, err
	}
	defer fResults.Close()

	var buf [48]byte

	// for each instance we run the number of iterations and save results and measured time
	for vID, instance := range instances {

		var totalDuration time.Duration
		var instanceResult bls12377.G1Jac
		var tmp bls12377.G1Affine

		for i := 0; i < nbIterations; i++ {
			start := time.Now()
			instanceResult.MultiExp(instance.Points, instance.Scalars, bls12377.BestC(len(instance.Scalars)))
			end := time.Since(start)
			totalDuration += end

			// write result to file
			tmp.FromJacobian(&instanceResult)
			buf = tmp.Bytes()
			if _, err := fResults.Write(buf[:]); err != nil {
				return nil, err
			}

			// write time elapsed to file
			if _, err := fTimes.WriteString(fmt.Sprintf("[gnark] iteration %d: %s\n", i+1, timex.Round(end, 2))); err != nil {
				return nil, err
			}
		}

		mean := totalDuration / time.Duration(nbIterations)
		if _, err := fTimes.WriteString(fmt.Sprintf("[gnark] Mean across all iterations (vector %d): %s\n", vID, timex.Round(mean, 2))); err != nil {
			return nil, err
		}

		results = append(results, mean)
	}

	totalMean := mean(results)

	if _, err := fTimes.WriteString(fmt.Sprintf("[gnark] Mean across all vectors: %s\n", timex.Round(totalMean, 2))); err != nil {
		return nil, err
	}

	return results, nil

}

func mean(times []time.Duration) time.Duration {
	total := time.Duration(0)
	for _, t := range times {
		total += t
	}
	total /= time.Duration(len(times))
	return total
}

type Instance struct {
	Points  []bls12377.G1Affine
	Scalars []fr.Element
}

func NewRandomInstance(nbElements int) Instance {
	var r Instance
	r.Points = make([]bls12377.G1Affine, nbElements)
	r.Scalars = make([]fr.Element, nbElements)

	randomPoints(r.Points)
	randomScalars(r.Scalars)

	return r
}

func randomScalars(scalars []fr.Element) {
	for i := 0; i < len(scalars); i++ {
		scalars[i].SetRandom()
	}
}

func randomPoints(points []bls12377.G1Affine) {
	rr := rand.New(rand.NewSource(time.Now().Unix()))
	var r, bound big.Int
	bound.SetString("8444461749428370424248824938781546531375899335154063827935233455917409239039", 10)

	var g1Gen bls12377.G1Affine
	g1Gen.X.SetString("81937999373150964239938255573465948239988671502647976594219695644855304257327692006745978603320413799295628339695")
	g1Gen.Y.SetString("241266749859715473739788878240585681733927191168601896383759122102112907357779751001206799952863815012735208165030")

	for i := 0; i < len(points); i++ {
		for r.IsUint64() && r.Uint64() == 0 {
			r.Rand(rr, &bound)
		}
		points[i].ScalarMultiplication(&g1Gen, &r)
	}
}
