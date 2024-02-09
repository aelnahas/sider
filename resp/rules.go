package resp

import (
	"time"
)

type rule struct {
	minArgCount int
	maxArgCount int
	argType     int
	hasOptions  bool
	options     map[string]optionSyntax
	isPubSubCmd bool
}

type optionSyntax struct {
	name     string
	dataType any
}

const (
	argTypeUnknown = iota
	argTypeVar
	argTypeOptional
	argTypeRequired
)

var (
	ruleSet = rule{
		minArgCount: 2,
		maxArgCount: 2,
		argType:     argTypeRequired,
		hasOptions:  true,
		options: map[string]optionSyntax{
			"EX": {
				name:     "EX",
				dataType: time.Second,
			},
			"PX": {
				name:     "PX",
				dataType: time.Millisecond,
			},
			"EXAT": {
				name:     "EXAT",
				dataType: time.Time{},
			},
			"PXAT": {
				name:     "PXAT",
				dataType: time.Time{},
			},
			"NX": {
				name:     "NX",
				dataType: true,
			},
			"XX": {
				name:     "PX",
				dataType: true,
			},
			"KEEPTTL": {
				name:     "KEEPTTL",
				dataType: true,
			},
			"GET": {
				name:     "GET",
				dataType: true,
			},
		},
	}

	ruleGet = rule{
		minArgCount: 1,
		maxArgCount: 1,
		argType:     argTypeRequired,
		hasOptions:  false,
	}

	rulePing = rule{
		minArgCount: 0,
		maxArgCount: 1,
		argType:     argTypeOptional,
		hasOptions:  false,
	}

	ruleEcho = rule{
		minArgCount: 1,
		maxArgCount: 1,
		argType:     argTypeOptional,
		hasOptions:  false,
	}

	ruleDel = rule{
		minArgCount: 1,
		argType:     argTypeVar,
		hasOptions:  false,
	}

	ruleExists = rule{
		minArgCount: 1,
		argType:     argTypeVar,
		hasOptions:  false,
	}

	ruleSub = rule{
		minArgCount: 1,
		argType:     argTypeVar,
		hasOptions:  false,
		isPubSubCmd: true,
	}

	ruleUnSub = rule{
		minArgCount: 1,
		argType:     argTypeVar,
		hasOptions:  false,
		isPubSubCmd: true,
	}

	rulePub = rule{
		minArgCount: 2,
		maxArgCount: 2,
		argType:     argTypeRequired,
		hasOptions:  false,
		isPubSubCmd: true,
	}
)

var rules map[string]rule = map[string]rule{
	"SET":    ruleSet,
	"GET":    ruleGet,
	"PING":   rulePing,
	"ECHO":   ruleEcho,
	"DEL":    ruleDel,
	"EXISTS": ruleExists,
	CmdSub:   ruleSub,
	CmdPub:   rulePub,
	CmdUnSub: ruleUnSub,
}
