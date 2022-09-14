package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	bls12377 "github.com/gbotrel/zprize-mobile-harness/msm/bls12-377"
	"github.com/icza/gox/timex"
)

var (
	fWorkingDirectory = flag.String("wd", ".", "working directory (saves output file and reads test vector file)")
	fTest             = flag.Bool("t", false, "if set, runs workingdirectory/points /scalars examples")
	fNbIterations     = flag.Uint("i", 1, "number of iterations to bench on each vector")
	fNbInstances      = flag.Uint("v", 1, "number of instances (vector) to randomly sample")
	fNbElemementsPow  = flag.Uint("n", 1, "number of elements to sample per instance (2**n)")
)

var _ atomic.Int64 // ensures we have latest Go version :-).

// This closely follows https://github.com/celo-org/zprize-mobile-harness for a clear comparaison
func main() {
	flag.Parse()

	if *fNbElemementsPow <= 0 || *fNbElemementsPow >= 20 {
		log.Fatal("nbElemesPow must be >= 1 && < 20")
	}

	// random multiple vectors
	// this is equivalent to benchmarkMSMRandomMultipleVecs in the harness
	if !*fTest {
		// sample random instances
		nbElems := 1 << *fNbElemementsPow
		instances := make([]Instance, *fNbInstances)
		for i := 0; i < len(instances); i++ {
			instances[i] = NewRandomInstance(nbElems)
		}

		// benchmark them
		res, err := benchmarkMSM(*fWorkingDirectory, instances, int(*fNbIterations))
		if err != nil {
			log.Fatal(err)
		}

		// output the total mean
		totalMean := mean(res)
		fmt.Println(timex.Round(totalMean, 2))
		return
	}

	// this is equivalent to benchmarkMSMFile

	// read wd/points and wd/scalars
	fPoints := filepath.Join(*fWorkingDirectory, "points")
	fScalars := filepath.Join(*fWorkingDirectory, "scalars")

	points, err := bls12377.ReadPoints(fPoints)
	if err != nil {
		log.Fatal(err)
	}

	scalars, err := bls12377.ReadScalars(fScalars)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(scalars); i++ {
		for j := 0; j < len(scalars[i]); j++ {
			scalars[i][j].FromMont()
		}
	}

	// sanity check
	if len(scalars) != len(points) {
		log.Fatal("sample file doesn't have same number of vectors of points / scalars")
	}

	instances := make([]Instance, len(scalars))
	for i := 0; i < len(instances); i++ {
		instances[i] = Instance{
			Points:  points[i],
			Scalars: scalars[i],
		}
		// sanity check
		if len(points[i]) != len(scalars[i]) {
			log.Fatal("sample file doesn't have same number of vectors of points / scalars")
		}
	}

	// benchmark it
	res, err := benchmarkMSM(*fWorkingDirectory, instances, int(*fNbIterations))
	if err != nil {
		log.Fatal(err)
	}

	// format the output similarly to the harness
	var sbb strings.Builder
	for _, r := range res {
		sbb.WriteString(timex.Round(r, 2).String())
		sbb.WriteString(", ")
	}
	formatted := sbb.String()
	formatted = formatted[:len(formatted)-2]
	fmt.Println(formatted)

	// compare results with the reference impl.
	bRef, err := os.ReadFile(filepath.Join(*fWorkingDirectory, "result.txt"))
	if err != nil {
		log.Fatal(err)
	}
	bGnark, err := os.ReadFile(filepath.Join(*fWorkingDirectory, "gnark_result.txt"))
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(bRef, bGnark) {
		log.Fatal("WARNING ref result != gnark result")
	}

}
