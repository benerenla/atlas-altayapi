package messages

import (
	"log/slog"

	nats "github.com/nats-io/nats.go"
)

type MessageBuilder struct {
	Coon *nats.Conn
}

func NewMessage(url string) *MessageBuilder {
	nc, err := nats.Connect(url)

	if err != nil {
		panic(err)
	}
	slog.Info("✅ NATS bağlantısı başarılı", "url", url)

	return &MessageBuilder{Coon: nc}
}

func (m *MessageBuilder) GetConnection() *nats.Conn {
	return m.Coon
}