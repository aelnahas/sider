package resp

import (
	"errors"
	"fmt"
)

var (
	ErrNotABulkString = errors.New("invalid syntax, input is not a valid resp bulk string")
	ErrOutOfBound     = errors.New("index out of bound")
)

type ErrUnexpectedSymbol struct {
	Wanted rune
	Got    rune
}

func (e ErrUnexpectedSymbol) Error() string {
	return fmt.Sprintf("unexpected symbol found. Wanted (%c) but got (%c)", e.Wanted, e.Got)
}

type ErrUnknownCommand struct {
	Name string
}

func (e ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command '%s'", e.Name)
}
