package cli

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func FindRepoDir() (string, error) {
	for cwd, err := os.Getwd(); cwd != string(filepath.Separator); cwd = filepath.Dir(cwd) {

		if err != nil {
			return "", errors.Wrap(err, "Could not get current working dir")
		}

		if _, err := os.Stat(filepath.Join(cwd, ".repo")); os.IsNotExist(err) {
			log.Debug().Str("cwd", cwd).Msg("'.repo' not found in current dir, searching upwards")
		} else {
			log.Debug().Str("cwd", cwd).Msg("'.repo' found")

			return cwd, nil
		}

	}

	return "", errors.New("Path is not part of a 'Repo' project")
}
