package cli

import (
	"fmt"
)

type PruneCmd struct {
	Projects []string `kong:"arg,optional,name='1st project',help='List of projects'"`
}

func (r *PruneCmd) Run(globals *Globals) error {
	fmt.Println("prune", r.Projects, len(r.Projects))
	return nil
}

