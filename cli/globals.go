package cli

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
)

type VersionFlag string

var (
	log zerolog.Logger
)

func init() {
	log = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger()

}

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Println(vars["version"])
	app.Exit(0)
	return nil
}

type Globals struct {
	Config    string      `help:"Location of client config files" default:"~/.gorepo" type:"path"`
	Debug     bool        `short:"D" help:"Enable debug mode"`
	Version   VersionFlag `name:"version" help:"Print version information and quit"`
}

