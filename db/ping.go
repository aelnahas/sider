package db

import "context"

type pingCmd struct {
	echo string
}

func (p *pingCmd) Read(args []string, _ map[string]any) error {
	p.echo = "PONG"
	if len(args) > 0 {
		p.echo = args[0]
	}

	return nil
}

func (p *pingCmd) Execute(ctx context.Context) (any, error) {
	return p.echo, nil
}
