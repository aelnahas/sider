package db

import (
	"context"
	"fmt"
)

type Command interface {
	Read(args []string, options map[string]any) error
	Execute(ctx context.Context) (any, error)
}

type DB struct {
	store *memory
}

type Expiration struct {
	Type    string
	Present bool
	TTL     any

	SetOnKeyExists    bool
	SetOnKeyNotExists bool
	KeepTTL           bool
}

func NewDB() *DB {
	store := newMemory()
	return &DB{store: store}
}

func (d *DB) Execute(ctx context.Context, name string, args []string, opts map[string]any) (any, error) {
	cmd, err := d.getCommand(name)
	if err != nil {
		return nil, err
	}

	if err := cmd.Read(args, opts); err != nil {
		return nil, err
	}
	return cmd.Execute(ctx)
}

func (d *DB) getCommand(name string) (Command, error) {
	switch name {
	case "SET":
		return &setCmd{store: d}, nil
	case "GET":
		return &getCmd{store: d}, nil
	case "PING":
		return &pingCmd{}, nil
	case "ECHO":
		return &pingCmd{}, nil
	case "EXISTS":
		return &existsCmd{store: d}, nil
	case "DEL":
		return &delCmd{store: d}, nil
	default:
		return nil, fmt.Errorf("unknown command %s", name)
	}

}
