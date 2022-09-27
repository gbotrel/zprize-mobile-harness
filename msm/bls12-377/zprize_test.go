package bls12377

import (
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
	"github.com/leanovate/gopter"
)

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
