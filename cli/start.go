package cli

import (
	"fmt"
)

type StartCmd struct {
	Branch   string   `kong:"arg,help='Provides a short description of the change you are trying to make to the projects.',default='default'"`
	Projects []string `kong:"arg,optional,name='1st project',help='Specifies which projects participate in this topic branch.'"`
}

func (r *StartCmd) Run(globals *Globals) error {
	fmt.Println("start", r.Branch, r.Projects, len(r.Projects))
	return nil
}
