package cli

import (
	"fmt"
)

type StatusCmd struct {
	Projects []string `arg optional name:"projects" help:"List of projects"`
}

func (r *StatusCmd) Run(globals *Globals) error {
	fmt.Println("status", r.Projects)
	return nil
}
