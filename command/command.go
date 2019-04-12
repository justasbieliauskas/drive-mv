package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/justasbieliauskas/drivemv/fs"
)

// Command holds environment, stdin, stdout, stderr for drivemv command.
type Command struct {
	Env []string
}

// New creates command with default setting.
func New() Command {
	return Command{
		Env: os.Environ(),
	}
}

// Run executes command with given arguments.
func (command *Command) Run(args []string) error {
	contains := envContainsAll(
		command.Env,
		"DRIVE_CLIENT_ID",
		"DRIVE_PROJECT_ID",
		"DRIVE_CLIENT_SECRET",
		"DRIVE_ACCESS_TOKEN",
		"DRIVE_REFRESH_TOKEN",
		"DRIVE_TOKEN_EXPIRY",
	)
	if !contains {
		return errors.New("Environment variables missing")
	}
	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("Error while opening file \"%s\": %v", args[0], err)
	}
	defer file.Close()
	root, err := fs.NewRoot(command.Env)
	if err != nil {
		return fmt.Errorf("Error while obtaining drive root: %v", err)
	}
	_, err = root.UploadFile(file, args[0])
	return err
}

func envContainsAll(env []string, vars ...string) bool {
	contains := false
	for _, varName := range vars {
		for _, envVariable := range env {
			envNameAndValue := strings.Split(envVariable, "=")
			if envNameAndValue[0] == varName {
				contains = true
				break
			}
		}
		if !contains {
			break
		}
	}
	return contains
}
