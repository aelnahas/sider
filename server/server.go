package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/aelnahas/sider/db"
	"github.com/aelnahas/sider/pubsub"
	"github.com/aelnahas/sider/resp"
)

const (
	ConfigDefaultPort     = 6379
	ConfigDefaultHostName = "0.0.0.0"
)

type Config struct {
	Port     uint
	HostName string
}

type Connection struct {
	config Config
	store  *db.DB
	broker pubsub.Broker
}

type Option = func(*Config)

func WithPort(port uint) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithHostName(hostname string) Option {
	return func(c *Config) {
		c.HostName = hostname
	}
}

func NewConnection(opts ...Option) *Connection {

	config := &Config{
		Port:     ConfigDefaultPort,
		HostName: ConfigDefaultHostName,
	}

	for _, opt := range opts {
		opt(config)
	}

	store := db.NewDB()

	broker := pubsub.NewBroker()

	return &Connection{
		config: *config,
		store:  store,
		broker: broker,
	}
}

func (c *Connection) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.config.HostName, c.config.Port))
	if err != nil {
		return fmt.Errorf("could not start redis server: %w", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			slog.Error("error openning connection", "error", err)
			return fmt.Errorf("error opening connection: %w", err)
		}

		go c.handleConn(conn)
	}
}

func (c *Connection) handleConn(conn net.Conn) error {
	var buf [512]byte
	request := bytes.NewBuffer(nil)
	var response []byte

	slog.Info("new incomming connection")

	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err != io.EOF {
				slog.Warn("connection closed")
				return nil
			}
			continue
		}

		request.Write(buf[0:n])
		slog.Info("received", "raw", request.String())

		cmd, err := resp.Parse(request)
		if err != nil {
			response = resp.Encode(err)
			_, err = conn.Write(response)
			if err != io.EOF {
				return nil
			}

			return err
		}

		var result any
		if cmd.IsPubSubCMD {
			c.executePubSubCmd(cmd, conn)
		} else {
			result, err = c.store.Execute(context.Background(), cmd.Name, cmd.Args, cmd.Options)
			var err error
			defer func() {
				err2 := conn.Close()
				if err2 != nil {
					slog.Error("could not close current connection", "error", err2)
					err = errors.Join(err, err2)
				}
			}()
		}
		if err != nil {
			response = resp.Encode(err)
		} else {
			response = resp.Encode(result)
		}
		_, err = conn.Write(response)
		if err != nil {
			if err != io.EOF {
				return nil
			}
			return err
		}
		slog.Info("sent", "raw", string(response))
	}
}

func (c *Connection) executePubSubCmd(cmd *resp.RawCommand, conn net.Conn) {
	id := conn.RemoteAddr().String()

	switch cmd.Name {
	case resp.CmdSub:
		c.broker.Subscribe(id, cmd.Args, conn)
	case resp.CmdUnSub:
		c.broker.Unsubscribe(id, cmd.Args)
	default:
		defer conn.Close()
		count := c.broker.Publish(cmd.Args[0], cmd.Args[1])
		res := resp.Encode(count)
		conn.Write(res)
	}
}
