//+build mage

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

const binary = "./drunkenfall"

func version() string {
	version, err := sh.OutCmd("git", "describe", "--always", "--tags")()
	if err != nil {
		log.Fatal(err)
	}
	return version
}

func buildTime() string {
	return time.Now().Format(time.RFC3339)
}

type Build mg.Namespace

// Builds the backend service
func (Build) Drunkenfall() error {
	env := make(map[string]string)

	// These are needed for static linking to work inside Docker
	env["CGO_ENABLED"] = "0"
	env["GOOS"] = "linux"

	ldflags := fmt.Sprintf(
		"-w -X main.version=%s -X main.buildtime=%s",
		version(),
		buildTime(),
	)

	return sh.RunWith(
		env,
		"go", "build", "-v",
		"-ldflags", ldflags,
		"-o", binary)
}

type Run mg.Namespace

// Runs the backend service
func (Run) Drunkenfall() error {
	mg.Deps(Build.Drunkenfall)

	return sh.Run(binary)
}

// Runs the npm frontend
func (Run) Npm() error {
	if err := os.Chdir("js"); err != nil {
		return err
	}

	defer os.Chdir("../")
	return sh.Run("npm", "run", "dev")
}

type Docker mg.Namespace

// Builds the backend service docker image
func (Docker) Drunkenfall() error {
	return sh.Run("docker-compose", "build", "drunkenfall")
}

// Builds the frontend docker image
func (Docker) Frontend() error {
	return sh.Run("docker-compose", "build", "frontend")
}

type Code mg.Namespace
