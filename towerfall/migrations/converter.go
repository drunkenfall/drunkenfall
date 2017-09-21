package migrations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/drunkenfall/drunkenfall/towerfall"
	"github.com/pkg/errors"
)

type Tournaments map[string][]byte

// Convert runs a migration.
//
// It expects an executable script to be found at <dir>/convert. First
// Convert produces pretty printed JSON blobs, one per tournament, at
// <dir>/in/<id>. Then it executes the convert script, which is
// expected to take the infile as the first argument and the outfile
// as a second. The convert script should then do its conversion and
// produce <dir>/out/<id>.
//
// The script could technically work with stdin/stdout instead, but
// producing files lets us do easy diffing to verify the patches.
func Convert(db *bolt.DB, dir string) (Tournaments, error) {
	p, _ := filepath.Abs(db.Path())
	dir = filepath.Join(filepath.Dir(p), "towerfall", "migrations", dir)

	os.MkdirAll(filepath.Join(dir, "in"), 0755)
	os.MkdirAll(filepath.Join(dir, "out"), 0755)
	ret := make(Tournaments)

	// Grab all the current tournaments and put them into <dir>/<id>, as
	// pretty printed JSON strings
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(towerfall.TournamentKey)
		err := b.ForEach(func(id []byte, data []byte) error {

			var pretty bytes.Buffer
			err := json.Indent(&pretty, data, "", "  ")
			if err != nil {
				return errors.Wrap(err, "pretty-printing json failed")
			}

			f, err := os.Create(filepath.Join(dir, "in", string(id)))
			if err != nil {
				return errors.Wrap(err, "opening the file failed")
			}
			defer f.Close()

			_, err = f.Write(pretty.Bytes())
			if err != nil {
				return errors.Wrap(err, "writing the file failed")
			}

			return nil
		})

		return err
	})

	if err != nil {
		return ret, errors.Wrap(err, "producing input files failed")
	}

	// Run the converter
	cmd := exec.Command(filepath.Join(dir, "convert"))
	output, err := cmd.Output()
	log.Printf("script output:\n%s", output)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			log.Print(fmt.Sprintf("exit %d:\n\n %s\n", waitStatus.ExitStatus(), exitError.Stderr))
		}

		return ret, errors.Wrap(err, "script failed")
	}

	// Load the results from the converter into the return value
	files, err := ioutil.ReadDir(filepath.Join(dir, "out"))
	if err != nil {
		return ret, errors.Wrap(err, "loading result files failed")
	}

	for _, f := range files {
		fn := filepath.Join(dir, "out", f.Name())
		key := filepath.Base(fn)
		file, err := os.Open(fn)
		if err != nil {
			return ret, errors.Wrap(err, fmt.Sprintf("could not load file %s", fn))
		}

		b, err := ioutil.ReadAll(file)
		if err != nil {
			return ret, errors.Wrap(err, fmt.Sprintf("could not read file %s", fn))
		}

		ret[key] = b
	}

	return ret, nil
}
