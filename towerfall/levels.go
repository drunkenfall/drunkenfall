package towerfall

import "math/rand"

type Levels map[string][]string

// Because of reasons, there is a difference between the level names
// in the app and in the game. This map translates the app names to
// the game names. Sorry.
var realNames map[string]string

func init() {
	realNames = make(map[string]string)
	realNames["sacred"] = "SacredGround" // Won't happen, but thorough
	realNames["twilight"] = "TwilightSpire"
	realNames["backfire"] = "Backfire"
	realNames["flight"] = "Flight"
	realNames["mirage"] = "Mirage"
	realNames["thornwood"] = "Thornwood"
	realNames["frostfang"] = "FrostfangKeep"
	realNames["moonstone"] = "Moonstone"
	realNames["kingscourt"] = "KingsCourt"
	realNames["sunken"] = "SunkenCity"
	realNames["towerforge"] = "TowerForge"
	realNames["ascension"] = "Ascension"
	realNames["amaranth"] = "TheAmaranth"
	realNames["dreadwood"] = "Dreadwood"
	realNames["darkfang"] = "Darkfang"
	realNames["cataclysm"] = "Cataclysm"
}

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
