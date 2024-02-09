package db

import "context"

type getCmd struct {
	key string

	store *DB
}

func (g *getCmd) Read(args []string, _ map[string]any) error {
	g.key = args[0]

	return nil
}

func (g *getCmd) Execute(ctx context.Context) (any, error) {
	return g.store.store.get(ctx, g.key)
}
