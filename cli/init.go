package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

type InitCmd struct {
	Url string `kong:"help='Specify a URL from which to retrieve a manifest repository.',short='u'"`
	Manifest string `kong:"help='Select a manifest file within the repository.',short='m',default='default.yaml'"`
	Branch string `kong:"help='Specify a revision, that is, a particular manifest-branch.',short='b'"`
}

func (r *InitCmd) Run(globals *Globals) error {
	log.Debug().Str("Url", r.Url).Str("Manifest", r.Manifest).Str("Branch", r.Branch).Msg("init")

	if _, err := os.Stat(filepath.Join(".repo", "manifests")); !os.IsNotExist(err) {
		// log.Error().Err(err).Msg("'.repo/manifests' already exists")

		return errors.Wrap(err, "'.repo/manifests' already exists")
	}

	err := os.MkdirAll(filepath.Join(".repo"), os.ModePerm); if err != nil {
		// log.Error().Err(err).Msg("Could not create '.repo directory")

		return errors.Wrap(err, "Could not create '.repo directory")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Minute)
	defer cancel()

	manifestsDir := filepath.Join(".repo", "manifests")

	cmdClone := exec.CommandContext(ctx, "git", "clone", r.Url, manifestsDir)
	cmdClone.Env = os.Environ()

	if err := cmdClone.Run(); err != nil {
		// log.Error().Err(err).Msg("Could not clone manifests repository")

		return errors.Wrap(err, "Could not clone manifests repository")
	}

	cmdCheckout := exec.CommandContext(ctx, "git", "checkout", r.Branch)
	cmdCheckout.Dir = manifestsDir
	cmdCheckout.Env = os.Environ()

	if err := cmdCheckout.Run(); err != nil {
		// log.Error().Err(err).Msgf("Could not checkout branch '%s'", r.Branch)

		return errors.Wrapf(err, "Could not checkout branch '%s'", r.Branch)
	}

	oldSymLink := filepath.Join("manifests", r.Manifest)

	if err := os.Chdir(".repo"); err != nil {
		return errors.Wrap(err, "Could not current directory to 'repo' directory")
	}

	if err := os.Symlink(oldSymLink, "manifest.yaml"); err != nil {
		return errors.Wrapf(err, "Could not create symlink to manifest '%s'", r.Manifest)
	}

	return nil
}