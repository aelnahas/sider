package executor

type Command interface {
	Execute() (any, error)
}
