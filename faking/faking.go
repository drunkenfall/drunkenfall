package faking

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	roman "github.com/StefanSchroeder/Golang-Roman"
)

// Cracklib is a random-pop interface to the cracklib wordlist
type Cracklib []string

var avatardir = os.ExpandEnv("$GOPATH/src/github.com/drunkenfall/drunkenfall/js/static/avatars/fake/")
var avatars []os.FileInfo

// ai is the avatar index of the current avatar
var ai int

func init() {
	var err error
	log.Printf("faking: Loaded %d cracklib strings", len(words))

	avatars, err = ioutil.ReadDir(avatardir)
	if err != nil {
		log.Printf("No fake avatars available. Faking will not work.")
	}

	// Shuffle the avatars so that they are not repeated
	for i := range avatars {
		j := rand.Intn(i + 1)
		avatars[i], avatars[j] = avatars[j], avatars[i]
	}

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
		p := []string{"von ", "Mc", "of ", "the ", "De", "La"}
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
	f := []string{"Faking", "Testing", "Durnkern", "Meta", "Jalo", "Poser", "Digging", words.Title()}
	prefix := f[rand.Intn(len(f))]

	// Add a roman numeral between 20 and 50
	x := roman.Roman(rand.Intn(30) + 20)

	// Fake a subtitle
	var subtitle string
	switch rand.Intn(9) {
	case 0:
		subtitle = fmt.Sprintf("The %s of %s", words.Title(), words.Title())
	case 1:
		subtitle = fmt.Sprintf("%s and %s", words.Title(), words.Title())
	case 2:
		subtitle = fmt.Sprintf("%s and %s for %s", words.Title(), words.Title(), words.Title())
	case 3:
		subtitle = fmt.Sprintf("A Song of %s and %s", words.Title(), words.Title())
	case 4:
		subtitle = fmt.Sprintf("To The %s", words.Title())
	case 5:
		subtitle = words.Title()
	case 6:
		subtitle = fmt.Sprintf("%s %s", words.Title(), words.Title())
	case 7:
		subtitle = fmt.Sprintf("In a World of %s", words.Title())
	case 8:
		subtitle = fmt.Sprintf("%s Against %s", words.Title(), words.Title())
	}

	title := fmt.Sprintf("%sFall %s: %s", prefix, x, subtitle)

	// Add question marks every now and again?
	if rand.Intn(5)%5 == 0 {
		title = fmt.Sprintf("%s?", title)
	}

	return title, strings.ToLower(x)
}

// FakeAvatar returns a fake avatar URL
func FakeAvatar() string {
	avatar := avatars[ai%len(avatars)]
	ai++

	// Add a slash so that it makes sense as a URL
	return "/" + avatardir + avatar.Name()
}

func percentTrue(n int) bool {
	return rand.Intn(100) <= n
}
