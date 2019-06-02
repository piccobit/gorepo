package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gorepo/manifest"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type SyncCmd struct {
	SwitchBack bool `kong:"help='Switch specified projects back to the manifest revision.',short='d'"`
	SyncToGoodBuild bool `kong:"help='Sync to a known good build as specified by the manifest-server element in the current manifest.',short='s'"`
	Force bool `kong:"help='Proceed with syncing other projects even if a project fails to sync.',short='f'"`
	Projects []string `kong:"arg,optional,name='1st project',help='List of projects to sync'"`
}

var (
	manifestData manifest.Manifest
)

func (r *SyncCmd) Run(globals *Globals) error {
	log.Debug().
		Strs("Projects",r.Projects).
		Bool("SwitchBack", r.SwitchBack).
		Bool("SyncToGoodBuild", r.SyncToGoodBuild).
		Bool("Force", r.Force).Msg("sync")

	repoDir, err := FindRepoDir()

	if err != nil {
		return errors.Wrap(err, "Could not find '.repo' dir")
	}

	manifestFile, err := os.Open(filepath.Join(repoDir, ".repo", "manifest.yaml"))

	if err != nil {
		return errors.Wrap(err, "Could not open manifest file")
	}

	bytes, err := ioutil.ReadAll(manifestFile)

	if err != nil {
		return errors.Wrap(err, "Could not read manifest file")
	}

    err = yaml.Unmarshal(bytes, &manifestData)

    if err != nil {
	    return errors.Wrap(err, "Could not unmarshal manifest file")
    }

    log.Debug().Interface("manifestData", manifestData).Msg("sync")

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Minute)
	defer cancel()

    for _, project := range manifestData.Manifest.Projects {
    	log.Debug().Str("name", project.Name).Str("path", project.Path).Msg("sync")

	    if _, err := os.Stat(filepath.Join(repoDir, project.Path)); os.IsNotExist(err) {
		    log.Debug().Str("name", project.Name).Msg("sync - project does not exist, cloning it")

		    var cloneDir string

		    if len(project.Path) != 0 {
		    	cloneDir = project.Path
		    } else {
			    cloneDir = project.Name
		    }

		    var remoteUrl string

		    remote := manifestData.FindRemote(project.Remote)
		    remoteUrl = remote.Fetch

		    log.Debug().Str("cloneDir", cloneDir).Str("remoteUrl", remoteUrl).Msg("sync - Cloning the project")
		    cmdClone := exec.CommandContext(ctx, "git", "clone", remoteUrl, cloneDir)
		    cmdClone.Dir = repoDir
		    cmdClone.Env = os.Environ()

		    if err := cmdClone.Run(); err != nil {
			    // log.Error().Err(err).Msg("Could not clone project repository")

			    return errors.Wrap(err, "Could not clone project repository")
		    }
	    } else {
		    log.Debug().Str("name", project.Name).Msg("sync - project exists, syncing it")

		    cmdRemoteUpdate := exec.CommandContext(ctx, "git", "remote", "update")
		    cmdRemoteUpdate.Dir = filepath.Join(repoDir, project.Path)
		    cmdRemoteUpdate.Env = os.Environ()

		    if err := cmdRemoteUpdate.Run(); err != nil {
			    // log.Error().Err(err).Msg("Could not update remote")

			    return errors.Wrap(err, "Could not update remote")
		    }

		    var revision string

		    remote := manifestData.FindRemote(project.Remote)

		    if len(project.Revision) != 0 {
			    revision = project.Revision
		    } else if len(remote.Revision) != 0 {
			    revision = remote.Revision
		    } else if len(manifestData.Manifest.Default.Revision) != 0 {
			    revision = manifestData.Manifest.Default.Revision
		    } else {
		    	log.Debug().Msg("Could not find any revision, therefore using 'master")

		    	revision = "master"
		    }

		    origin := fmt.Sprintf("origin/%s", revision)

		    log.Debug().Str("origin", origin).Msg("sync - origin")

		    cmdRebase := exec.CommandContext(ctx, "git", "rebase", origin)
		    cmdRebase.Dir = filepath.Join(repoDir, project.Path)
		    cmdRebase.Env = os.Environ()

		    if err := cmdRebase.Run(); err != nil {
			    // log.Error().Err(err).Msg("Could not rebase origin")

			    return errors.Wrap(err, "Could not rebase origin")
		    }
	    }
    }

	return nil
}

