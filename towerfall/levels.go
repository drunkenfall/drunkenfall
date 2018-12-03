package towerfall

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
