package towerfall

import "math/rand"

type Levels map[string][]string

// NewLevels returns a newly shuffled set of levels
//
// This returns a new discrete set of levels so that tournaments won't
// re-use them.
func NewLevels() Levels {
	levels := make(Levels)

	levels[playoff] = []string{
		"twilight",
		"backfire",
		"flight",
		"mirage",
		"thornwood",
		"frostfang",
		"moonstone",
		"kingscourt",
	}

	levels[semi] = []string{
		"sunken",
		"towerforge",
		"ascension",
		"amaranth",
		"dreadwood",
		"darkfang",
	}

	levels[final] = []string{
		"cataclysm",
	}

	shuffle(levels[playoff])
	shuffle(levels[semi])
	shuffle(levels[final])

	return levels
}

// shuffle ...
func shuffle(s []string) {
	n := len(s)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		s[r], s[i] = s[i], s[r]
	}
}
