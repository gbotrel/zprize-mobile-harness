package bls12377

import (
	"bytes"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
	"github.com/leanovate/gopter"
	"github.com/stretchr/testify/require"
)

func TestMultiExp(t *testing.T) {
	assert := require.New(t)
	scalars, err := ReadScalars("scalars")
	assert.NoError(err, "reading scalar file")
	points, err := ReadPoints("points")
	assert.NoError(err, "reading points file")
	results, err := ReadResults("result")
	assert.NoError(err, "reading results file")

	assert.True(len(scalars) == len(points))
	assert.True(len(scalars) == len(results))

	// perform the msms
	for i := 0; i < len(scalars); i++ {
		for J := 0; J < len(scalars[i]); J++ {
			scalars[i][J].FromMont()
		}
		var p G1EdExtended
		_, err = p.MultiExp(points[i], scalars[i], ecc.MultiExpConfig{})
		assert.NoError(err)
		var pr G1Affine
		pr.FromExtendedEd(&p)

		assert.True(pr.Equal(&results[i]), "msm mismatch")
	}
}

func TestSerializationScalars(t *testing.T) {
	assert := require.New(t)

	scalars, err := ReadScalars("scalars")
	assert.NoError(err, "reading scalar file")

	bScalars := SerializeScalars(scalars)
	fScalars, err := os.ReadFile("scalars")
	assert.NoError(err, "reading scalar file 2nd time")

	assert.True(bytes.Equal(bScalars, fScalars), "bad encoding of scalars")
}

func TestSerializationPoints(t *testing.T) {
	t.Skip("skipping with ed extended for now")
	// assert := require.New(t)

	// points, err := ReadPoints("points")
	// assert.NoError(err, "reading points file")

	// bPoints := SerializePoints(points)
	// fPoints, err := os.ReadFile("points")
	// assert.NoError(err, "reading points file 2nd time")

	// assert.True(bytes.Equal(bPoints, fPoints), "bad encoding of points")
}

func TestSerializationResults(t *testing.T) {
	assert := require.New(t)

	points, err := ReadResults("result")
	assert.NoError(err, "reading results file")

	bPoints := SerializeResults(points)
	fPoints, err := os.ReadFile("result")
	assert.NoError(err, "reading results file 2nd time")

	assert.True(bytes.Equal(bPoints, fPoints), "bad encoding of results")
}

const nbFuzzShort = 10
const nbFuzz = 100

// GenFp generates an Fp element
func GenFp() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var elmt fp.Element

		if _, err := elmt.SetRandom(); err != nil {
			panic(err)
		}

		return gopter.NewGenResult(elmt, gopter.NoShrinker)
	}
}

// GenFr generates an Fr element
func GenFr() gopter.Gen {
	return func(genParams *gopter.GenParameters) *gopter.GenResult {
		var elmt fr.Element

		if _, err := elmt.SetRandom(); err != nil {
			panic(err)
		}

		return gopter.NewGenResult(elmt, gopter.NoShrinker)
	}
}
