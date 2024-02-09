package resp

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type RawCommand struct {
	Name        string
	Args        []string
	Options     map[string]any
	IsPubSubCMD bool
}

func Parse(input io.Reader) (*RawCommand, error) {
	scanner := NewScanner(input)

	args := make([]string, 0)

	token, lit, err := scanner.Next()
	if err != nil {
		return nil, err
	}

	if token == TokenArg {
		return nil, fmt.Errorf("expected token to be an arg")
	}
	name := strings.ToUpper(lit)

	cmdRule, ok := rules[name]
	if !ok {
		return nil, ErrUnknownCommand{Name: name}
	}

	optionName := ""
	options := make(map[string]any)
	for scanner.HasNext() {
		token, lit, err := scanner.Next()
		if err != nil {
			return nil, err
		}

		if token == TokenArg {

			if len(args) < cmdRule.maxArgCount || cmdRule.argType == argTypeVar {
				args = append(args, lit)
			} else if !cmdRule.hasOptions {
				return nil, fmt.Errorf("cmd %s does not accept options", name)
			} else if optionName != "" {
				optSyntax := cmdRule.options[optionName]
				val, err := parseOptValue(optSyntax.dataType, optionName, lit)

				if err != nil {
					return nil, err
				}
				options[optionName] = val
				optionName = ""
			} else {
				optionName = strings.ToUpper(lit)
				opt, found := cmdRule.options[optionName]
				if !found {
					return nil, fmt.Errorf("syntax error command %s does not support option %s", name, lit)
				}

				_, ok := opt.dataType.(bool)
				if ok {
					options[optionName] = true
					optionName = ""
				}
			}
		}
	}

	if len(args) < cmdRule.minArgCount {
		return nil, fmt.Errorf("syntax err command %s is missing required args", name)
	}

	if optionName != "" {
		return nil, fmt.Errorf("syntax err option %s is missing a value", optionName)
	}

	return &RawCommand{Name: name, Args: args, Options: options, IsPubSubCMD: cmdRule.isPubSubCmd}, nil
}

func parseOptValue(dataType any, name, val string) (any, error) {

	switch dataType := dataType.(type) {
	case string:
		return val, nil
	case int:
		return strconv.Atoi(val)
	case time.Time:
		return time.Parse(time.RFC3339, val)
	case time.Duration:

		unit := "s"
		if dataType == time.Millisecond {
			unit = "ms"
		}

		val += unit
		return time.ParseDuration(val)
	default:
		return nil, fmt.Errorf("unsupported data type (%T) provided with option (%s)", dataType, name)
	}
}
