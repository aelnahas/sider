package db

import "context"

type delCmd struct {
	keys  []string
	store *DB
}

func (d *delCmd) Read(args []string, _ map[string]any) error {
	d.keys = args
	return nil
}

func (d *delCmd) Execute(ctx context.Context) (any, error) {
	count := 0

	for _, key := range d.keys {
		if d.store.store.exists(ctx, key) {
			count++
		}
		if err := d.store.store.del(ctx, key); err != nil {
			return count, err
		}

	}

	return count, nil
}
