package db

import "context"

type existsCmd struct {
	store *DB
	keys  []string
}

func (e *existsCmd) Read(args []string, _ map[string]any) error {
	e.keys = args
	return nil
}

func (e *existsCmd) Execute(ctx context.Context) (any, error) {
	count := 0

	for _, key := range e.keys {
		if e.store.store.exists(ctx, key) {
			count++
		}

	}

	return count, nil

}
