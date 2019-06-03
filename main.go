package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"

	"gorepo/cli"
)

type CLI struct {
	cli.Globals

	Init cli.InitCmd `kong:"cmd,help='Installs Repo in the current directory.'"`
	Sync cli.SyncCmd `kong:"cmd,help='Downloads new changes and updates the working files in your local environment.'"`
	// Upload UploadCmd `cmd help:"For the specified projects, Repo compares the local branches to the remote branches updated during the last Repo sync."`
	// Diff DiffCmd `cmd help:"Shows outstanding changes between the commit and the working tree."`
	// Download DownloadCmd `cmd help:"Downloads the specified change from the review system and makes it available in your project's local working directory."`
	// Forall ForallCmd `cmd help:"Executes the given shell command in each project."`
	Prune  cli.PruneCmd  `kong:"cmd,help='Prunes (deletes) topics that are already merged.'"`
	Start  cli.StartCmd  `cmd help:"Begins a new branch for development, starting from the revision specified in the manifest."`
	Status cli.StatusCmd `cmd help:"Compares the working tree to the staging area (index) and the most recent commit on this branch (HEAD) in each project specified."`
}

var (
	log zerolog.Logger
)

func init() {
	log = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Str("component", "main").
		Logger()
}

func main() {
	cli := CLI{
		Globals: cli.Globals{
			Version: cli.VersionFlag("0.1.0"),
		},
	}

	ctx := kong.Parse(&cli,
		kong.Name("gorepo"),
		kong.Description("Go 'repo' tool"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": "0.1.0",
		})

	if cli.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	err := ctx.Run(&cli.Globals)

	ctx.FatalIfErrorf(err)
}

