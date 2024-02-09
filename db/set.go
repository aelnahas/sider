package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log/slog"
)

var ErrSyntax = errors.New("syntax error")

type setCmd struct {
	key        string
	val        string
	expiration Expiration

	getOldVal bool

	store *DB
}

func (s *setCmd) Read(args []string, opts map[string]any) error {
	s.key = args[0]
	s.val = args[1]

	for key, opt := range opts {
		if err := s.readOpt(key, opt); err != nil {
			return err
		}
	}

	return nil
}

func (s *setCmd) Execute(ctx context.Context) (any, error) {
	var ret any = "OK"
	if s.getOldVal {
		oldVal, err := s.store.store.get(ctx, s.key)
		if err != nil {
			return nil, err
		}

		ret = oldVal
	}

	err := s.store.store.set(ctx, &record{key: s.key, val: s.val})

	if s.expiration.Present {
		s.startTTLBackground(ctx, ret != nil && ret != "OK")
	}

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *setCmd) startTTLBackground(ctx context.Context, keyExists bool) {
	slog.Info("starting ttl background deletion", "key-exists", keyExists, "expiration-settings", s.expiration)
	if !s.expiration.Present || keyExists && s.expiration.SetOnKeyNotExists || !keyExists && s.expiration.SetOnKeyExists {
		return
	}

	if s.expiration.Type == "EX" || s.expiration.Type == "PX" {
		slog.Info("starting a timer", "expiration-type", s.expiration.Type, "expiration-duration", s.expiration.TTL)
		s.setTimer(ctx)
		return
	}

	slog.Info("starting a deadline with timestamp settings", "expiration-type", s.expiration.Type, "expiratio-duation", s.expiration.TTL)
	s.setDeadline(ctx)
}

func (s *setCmd) setDeadline(ctx context.Context) {
	go func() {
		deadline := s.expiration.TTL.(time.Time)
		pollFreq := time.Millisecond * 500
		for time.Now().Before(deadline) {
			time.Sleep(pollFreq)
		}

		if err := s.store.store.del(ctx, s.key); err != nil {
			slog.Error("could not delete data after ttl expired", "error", err, "key", s.key)
			return
		}

		slog.Info("deleted record successfully", "key", s.key)
	}()

}

func (s *setCmd) setTimer(ctx context.Context) {
	duration := s.expiration.TTL.(time.Duration)
	timer := time.NewTicker(duration)

	go func() {
		<-timer.C
		err := s.store.store.del(ctx, s.key)

		if err != nil {
			slog.Error("could not delete data after ttl expired", "error", err, "key", s.key)
			return
		}

		slog.Info("deleted record successfully", "key", s.key)
	}()
}

func (s *setCmd) readOpt(name string, val any) error {

	switch name {
	case "EX", "PX", "EXAT", "PXAT":
		return s.setExpiration(name, val)
	case "NX", "XX":
		return s.setExpirationPolicy(name)
	case "KEEPTTL":
		s.expiration.KeepTTL = true
		return nil
	case "GET":
		s.getOldVal = true
		return nil
	default:
		return fmt.Errorf("syntax error")
	}
}

func (s *setCmd) setExpirationPolicy(name string) error {

	if name == "NX" {
		if s.expiration.SetOnKeyExists {
			return ErrSyntax
		}

		s.expiration.SetOnKeyNotExists = true
		return nil
	}

	if name == "XX" {
		if s.expiration.SetOnKeyNotExists {
			return ErrSyntax
		}

		s.expiration.SetOnKeyExists = true
		return nil
	}

	return nil

}

func (s *setCmd) setExpiration(name string, val any) error {
	if s.expiration.Present {
		return fmt.Errorf("syntax error")
	}

	s.expiration = Expiration{
		Type:    name,
		TTL:     val,
		Present: true,
	}

	return nil
}
