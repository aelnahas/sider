package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"log/slog"

	"github.com/aelnahas/sider/server"
	"github.com/urfave/cli/v2"
)

const (
	pidFile = "sider.lock"
)

func Execute() error {

	daemon := false
	var port uint

	app := &cli.App{
		Name:                 "sider",
		Usage:                "a key/value db that mimics redis",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "starts sider service",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "daemon",
						Value:       false,
						Aliases:     []string{"d"},
						Usage:       "run service in daemon mode",
						Destination: &daemon,
					},
					&cli.UintFlag{
						Name:        "Port",
						Value:       server.ConfigDefaultPort,
						Aliases:     []string{"p"},
						Usage:       "set port",
						Destination: &port,
					},
				},
				Action: func(c *cli.Context) error {
					if daemon {
						bgCmd := exec.Command("sider", "start")
						if err := bgCmd.Start(); err != nil {
							return err
						}
						slog.Info("starting server...", "PID", bgCmd.Process.Pid)
						err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", bgCmd.Process.Pid)), 0666)
						if err != nil {
							return err
						}
						daemon = false
						return nil
					}
					conn := server.NewConnection(server.WithPort(port))
					return conn.Start()

				},
			},
			{
				Name:    "stop",
				Aliases: []string{"x"},
				Usage:   "stops sider service",
				Action: func(c *cli.Context) error {
					pid, err := os.ReadFile(pidFile)
					if err != nil {
						return err
					}

					slog.Info("stopping server ...", "PID", pid)

					killCmd := exec.Command("kill", string(pid))

					if err := killCmd.Start(); err != nil {
						return err
					}

					return os.Remove(pidFile)

				},
			},
		},
	}

	return app.Run(os.Args)
}
