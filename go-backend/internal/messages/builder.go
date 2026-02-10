package messages

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/benerenla/best-plugin/internal/models"
	"github.com/benerenla/best-plugin/internal/repository"
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
func (m *MessageBuilder) RegisterHandlers(repo *repository.PlayerRepository) {
	m.Coon.Subscribe("mc.player.register", func(msg *nats.Msg) {
		m.handlePlayerRegister(msg, repo)
	})
	m.Coon.Subscribe("mc.player.is_registered", func(msg *nats.Msg) {
		m.handlePlayerIsRegistered(msg, repo)
	})
}
func (m *MessageBuilder) handlePlayerRegister(msg *nats.Msg, repo *repository.PlayerRepository) {
	var p models.RegisterPlayerRequest
	json.Unmarshal(msg.Data, &p)

	newPlayer := models.Player{
		UUID:     p.UUID,
		Username: p.Username,
		Password: p.Password,
	}

	// Repository'yi kullanarak veriyi işle
	err := repo.RegisterPlayer(context.Background(), &newPlayer)

	if err != nil {
		slog.Error("❌ Oyuncu işlenemedi", "error", err)
		msg.Respond([]byte("ERROR_DATABASE"))
		return
	}

	// ... (Kayıt veya Giriş mantığı burada çalışır)
	slog.Info("✅ Oyuncu işlendi", "user", p.Username)
	msg.Respond([]byte("SUCCESS"))
}

func (m *MessageBuilder) handlePlayerIsRegistered(msg *nats.Msg, repo *repository.PlayerRepository) {
	var p models.IsRegisteredRequest
	err := json.Unmarshal(msg.Data, &p)
	if err != nil {
		slog.Error("❌ JSON unmarshal başarısız", "error", err)
		m.Coon.Publish(msg.Reply, []byte("false"))
		return
	}

	// Repository'yi kullanarak veriyi işle
	player, err := repo.IsRegistered(context.Background(), p)
	if err != nil {
		slog.Error("❌ Oyuncu kontrol edilemedi", "error", err)
		m.Coon.Publish(msg.Reply, []byte("false"))
		return
	}

	if !player {
		slog.Info("❌ Oyuncu kayıtlı değil", "user", p.UUID)
		m.Coon.Publish(msg.Reply, []byte("false"))
		return
	}
	msg.Respond([]byte("true"))
	slog.Info("✅ Oyuncu kayıtlı", "user", p.UUID)
}
