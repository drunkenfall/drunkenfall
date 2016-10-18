package main

import (
	"fmt"
	roman "github.com/StefanSchroeder/Golang-Roman"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Cracklib is a random-pop interface to the cracklib wordlist
type Cracklib []string

var words Cracklib

func init() {
	cracklibLocation := os.Getenv("DF_CRACKLIB_LOCATION")
	defaultCracklibLocation := "/usr/share/cracklib/cracklib-small"
	if cracklibLocation == "" {
		log.Print("Cracklib location set to default: " + defaultCracklibLocation)

		cracklibLocation = defaultCracklibLocation
	}
	out, err := ioutil.ReadFile(cracklibLocation)
	if err != nil {
		log.Print("Cracklib failed to load. Set location by env variable DF_CRACKLIB_LOCATION")
		return
	}
	words = removeApostrophes(strings.Split(string(out), "\n"))
	log.Printf("Loaded %d cracklib strings", len(words))

	// Also reset the random...
	rand.Seed(time.Now().UnixNano())
}

// Random returns a random string from cracklib
func (c *Cracklib) Random() string {
	return words[rand.Intn(len(words))]
}

// Title returns a random string from cracklib in Title Case
func (c *Cracklib) Title() string {
	return strings.Title(c.Random())
}

// Upper returns a random string from cracklib in UPPERCASE
func (c *Cracklib) Upper() string {
	return strings.ToUpper(c.Random())
}

// FakeName returns a fake player name
func FakeName() string {
	var prefix string
	// Add a last name prefix 80% of the time
	if percentTrue(80) {
		p := []string{"von ", "Mc", "of ", "the ", "De"}
		prefix = p[rand.Intn(len(p))]
	}

	return fmt.Sprintf(
		"%s %s%s",
		words.Title(),
		prefix,
		words.Title(),
	)
}

// FakeNick returns a fake player nickname
func FakeNick() string {
	var tag, prefix, suffix, number, char string

	// Add a clan tag
	if percentTrue(25) {
		content := words.Upper()

		if percentTrue(50) {
			// Half of the time, pick three random uppercase characters.
			b := make([]rune, 3)
			letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			content = string(b)
		}

		tag = fmt.Sprintf("[%s]", content)
	}

	// Add a prefix
	if percentTrue(20) {
		p := []string{"sexy", "super", "muscle", "ultra", "mega"}
		prefix = p[rand.Intn(len(p))]
	}

	// Add a suffix
	if percentTrue(33) {
		s := []string{"gamer", "dude", "girl", "man", "lady", ""}
		suffix = s[rand.Intn(len(s))]
	}

	// Add a number between 250 and 1000
	if percentTrue(10) {
		number = strconv.Itoa(rand.Intn(750) + 250)
	}

	// Add a random character at the end
	if percentTrue(10) {
		c := []string{"-", "_", "^", "$", "#", "`", ">"}
		char = c[rand.Intn(len(c))]
	}

	return fmt.Sprintf(
		"%s%s%s%s%s%s",
		tag,
		prefix,
		words.Title(),
		suffix,
		number,
		char,
	)
}

// FakeTournamentTitle returns a fake tournament title and a numeral
func FakeTournamentTitle() (string, string) {
	f := []string{"Faking", "Testing", "Durnkern", "Meta", "Jalo", "Poser"}
	prefix := f[rand.Intn(len(f))]

	// Add a roman numeral between 20 and 400
	x := roman.Roman(rand.Intn(380) + 20)

	// Fake a subtitle
	var subtitle string
	switch rand.Intn(5) {
	case 0:
		subtitle = fmt.Sprintf("The %s of %s", words.Title(), words.Title())
	case 1:
		subtitle = fmt.Sprintf("%s and %s", words.Title(), words.Title())
	case 2:
		subtitle = fmt.Sprintf("%s and %s for %s", words.Title(), words.Title(), words.Title())
	case 3:
		subtitle = fmt.Sprintf("Barry Trotter and the %s of %s", words.Title(), words.Title())
	case 4:
		subtitle = words.Title()
	}

	title := fmt.Sprintf("%sFall %s: %s", prefix, x, subtitle)
	return title, strings.ToLower(x)
}

// FakeAvatar returns a fake avatar URL
func FakeAvatar() string {
	dir := "static/avatars/fake/"
	avatars, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("No fake avatars available")
	}

	avatar := avatars[rand.Intn(len(avatars))]

	// Add a slash so that it makes sense as a URL
	return "/" + dir + avatar.Name()
}

func percentTrue(n int) bool {
	return rand.Intn(100) <= n
}

func removeApostrophes(c Cracklib) Cracklib {
	var vsf Cracklib
	for _, v := range c {
		if !strings.Contains(v, "'") {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
