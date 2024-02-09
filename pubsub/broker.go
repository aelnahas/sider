package pubsub

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/aelnahas/sider/resp"
)

type MessageKind int

const (
	MessageKindUnknown MessageKind = iota
	MessageKindSubscribe
	MessageKindMessage
)

func (mk MessageKind) String() string {
	switch mk {
	case MessageKindMessage:
		return "message"
	case MessageKindSubscribe:
		return "subscribed"
	default:
		return "unknown"
	}
}

type Message struct {
	Kind      MessageKind
	Topic     string
	Data      any
	Pattern   string
	Timestamp time.Time
}

type Topic struct {
	Name     string
	Messages []string
}

type Broker interface {
	Connect(id string, conn net.Conn)
	Disconnect(id string)
	Subscribe(id string, topics []string, conn net.Conn)
	Unsubscribe(id string, topic []string)
	Publish(topic string, data any) int
}

type client struct {
	id   string
	conn net.Conn

	messages chan *Message

	writer bufio.Writer
	reader bufio.Reader
}

type broker struct {
	clients       map[string]*client
	subscriptions map[string]map[string]struct{}
}

var _ Broker = new(broker)

func NewBroker() *broker {
	return &broker{
		clients:       make(map[string]*client),
		subscriptions: make(map[string]map[string]struct{}),
	}
}

func (b *broker) Connect(id string, conn net.Conn) {
	c := newClient(id, conn)
	c.Listen()

	b.clients[id] = c
}

func (b *broker) Disconnect(id string) {
	client, ok := b.clients[id]
	if !ok {
		return
	}

	client.Close()
	delete(b.clients, id)
}

func (b *broker) Subscribe(id string, topics []string, conn net.Conn) {
	_, ok := b.clients[id]
	if !ok {
		b.Connect(id, conn)
	}

	client := b.clients[id]
	for _, topic := range topics {
		b.subTopic(id, topic)

		client.messages <- &Message{Topic: topic, Timestamp: time.Now(), Data: "subscribed", Kind: MessageKindSubscribe}
	}

}

func (b *broker) subTopic(id string, topic string) {

	if _, ok := b.subscriptions[topic]; !ok {
		b.subscriptions[topic] = make(map[string]struct{})
	}

	b.subscriptions[topic][id] = struct{}{}
}

func (b *broker) Unsubscribe(id string, topics []string) {
	if _, ok := b.clients[id]; !ok {
		return
	}

	for _, topic := range topics {
		b.unSubTopic(id, topic)
	}
}

func (b *broker) unSubTopic(id, topic string) {
	if _, ok := b.subscriptions[topic]; !ok {
		return
	}

	delete(b.subscriptions[topic], id)
}

func (b *broker) Publish(topic string, data any) int {

	sub := b.subscriptions[topic]

	count := 0
	for id := range sub {
		client, ok := b.clients[id]
		if !ok {
			slog.Error("client not found", "id", id, "topic", topic)
			break
		}

		msg := &Message{
			Kind:      MessageKindMessage,
			Topic:     topic,
			Timestamp: time.Now(),
			Data:      data,
		}
		client.messages <- msg
		count++
	}
	return count
}

func newClient(id string, conn net.Conn) *client {
	client := &client{
		id:   id,
		conn: conn,

		messages: make(chan *Message),
		writer:   *bufio.NewWriter(conn),
		reader:   *bufio.NewReader(conn),
	}
	client.Listen()
	return client
}

func (c *client) Listen() {
	//go c.Read()
	go c.SendMessages()
}

func (c *client) Read() {
	for {
		var buf [512]byte
		str, err := c.reader.Read(buf[0:])
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		slog.Info("received", "data", str)

	}

}

func (c *client) SendMessages() {
	for msg := range c.messages {
		slog.Info("publishing message", "topic", msg.Topic, "data", msg.Data, "kind", fmt.Sprint(msg.Kind))

		res := resp.EncodeArray(
			msg.Kind.String(),
			msg.Topic,
			msg.Data,
		)

		_, err := c.writer.Write(res)
		fmt.Println("error ", err)
		if err != nil {
			slog.Error("could not publish message", "error", err, "topic", msg.Topic, "id", c.id)
			c.writer.Write(resp.Encode(err))
			break
		}

		if err := c.writer.Flush(); err != nil {
			slog.Error("could not publish message", "error", err, "topic", msg.Topic, "id", c.id)
			c.writer.Write(resp.Encode(err))
			break
		}
	}
	slog.Info("closed client writer thread", "id", c.id)
}

func (c *client) Close() {
	c.conn.Close()
	slog.Info("closed connection", "id", c.id)
}
