package main

import (
	"os"

	"log/slog"

	"github.com/aelnahas/sider/cmd"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	if err := cmd.Execute(); err != nil {
		slog.Error("error orccured during execution", "err", err)
		panic(err)
	}
}
