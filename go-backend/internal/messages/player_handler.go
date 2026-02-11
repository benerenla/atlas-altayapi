package messages

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/benerenla/best-plugin/internal/models"
	"github.com/benerenla/best-plugin/internal/repository"
	"github.com/benerenla/best-plugin/utils"
	nats "github.com/nats-io/nats.go"
)

type AuthHandler struct {
	Conn *nats.Conn
}

func NewAuthHandler(conn *nats.Conn) *AuthHandler {
	return &AuthHandler{Conn: conn}
}

func (m *AuthHandler) RegisterHandlers(repo *repository.PlayerRepository) {
	m.Conn.Subscribe("mc.player.register", func(msg *nats.Msg) {
		m.handlePlayerRegister(msg, repo)
	})
	m.Conn.Subscribe("mc.player.is_registered", func(msg *nats.Msg) {
		m.handlePlayerIsRegistered(msg, repo)
	})
	m.Conn.Subscribe("mc.player.login", func(msg *nats.Msg) {
		m.handlePlayerLogin(msg, repo)
	})
	m.Conn.Subscribe("mc.player.verify_email", func(msg *nats.Msg) {
		m.handleEmailVerification(msg, repo)
	})
	m.Conn.Subscribe("mail.send_verification", func(msg *nats.Msg) {
		m.handleSendMail(msg, repo)
	})
	m.Conn.Subscribe("mc.player.verify", func(msg *nats.Msg) {
		m.handleVerifyCode(msg, repo)
	})
	m.Conn.Subscribe("mc.player.is_logged_in", func(msg *nats.Msg) {
		m.handlePlayerIsLoggedIn(msg, repo)
	})
	m.Conn.Subscribe("mc.player.set_logged_in", func(msg *nats.Msg) {
		m.setPlayerIsLoggedIn(msg, repo)
	})
	m.Conn.Subscribe("mc.player.is_verifed", func(msg *nats.Msg) {
		m.IsVerified(msg, repo)
	})
}
func (m *AuthHandler) handlePlayerRegister(msg *nats.Msg, repo *repository.PlayerRepository) {
	var p models.RegisterPlayerRequest
	json.Unmarshal(msg.Data, &p)

	slog.Info("🔍 Kayıt isteği alındı", "uuid", p.UUID, "username", p.Username)

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

func (m *AuthHandler) handlePlayerIsRegistered(msg *nats.Msg, repo *repository.PlayerRepository) {
	var p models.IsRegisteredRequest
	err := json.Unmarshal(msg.Data, &p)
	if err != nil {
		slog.Error("❌ JSON unmarshal başarısız", "error", err)
		m.Conn.Publish(msg.Reply, []byte("false"))
		return
	}

	// Repository'yi kullanarak veriyi işle
	player, err := repo.IsRegistered(context.Background(), p)
	if err != nil {
		slog.Error("❌ Oyuncu kontrol edilemedi", "error", err)
		m.Conn.Publish(msg.Reply, []byte("false"))
		return
	}

	if !player {
		slog.Info("❌ Oyuncu kayıtlı değil", "user", p.UUID)
		m.Conn.Publish(msg.Reply, []byte("false"))
		return
	}
	msg.Respond([]byte("true"))
	slog.Info("✅ Oyuncu kayıtlı", "user", p.UUID)
}
func (m *AuthHandler) handlePlayerLogin(msg *nats.Msg, repo *repository.PlayerRepository) {
	var p models.LoginPlayerRequest
	json.Unmarshal(msg.Data, &p)

	player, err := repo.LoginPlayer(context.Background(), &p)
	if err != nil {
		slog.Error("❌ Oyuncu giriş yapılamadı", "error", err)

		msg.Respond([]byte("ERROR_LOGIN"))
		return
	}

	// Repository'yi kullanarak veriyi işle
	// ... (Kayıt veya Giriş mantığı burada çalışır)
	slog.Info("✅ Oyuncu işlendi", "user", player.Username)
	msg.Respond([]byte("SUCCESS"))
}
func (m *AuthHandler) handleEmailVerification(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.MailPayload
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ MailPayload JSON unmarshal başarısız", "error", err)
		return
	}
	if req.Email != "" {
		// E-posta adresi sağlanmış, doğrulama kodu oluştur ve kaydet
		code := utils.GenerateSecureCode(6) // 6 haneli doğrulama kodu oluştur
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := repo.SetEmail(ctx, req.UUID, &req.Email)
		if err != nil {
			slog.Error("❌ E-posta adresi kaydedilemedi", "uuid", req.UUID, "email", req.Email, "error", err)
			return
		}

		err = repo.SaveVerificationCode(ctx, req.UUID, code)

		if err != nil {
			slog.Error("❌ Doğrulama kodu kaydedilemedi", "uuid", req.UUID, "error", err)
			return
		}
		payload := models.MailPayload{
			Email:    req.Email,
			Username: req.Username,
			Code:     code,
		}
		data, _ := json.Marshal(payload)
		m.Conn.Publish("mail.send_verification", data)
		slog.Info("✅ Doğrulama kodu oluşturuldu ve e-posta gönderildi", "uuid", req.UUID, "email", req.Email)
		msg.Respond([]byte("SUCCESS"))
	}
}

func (m *AuthHandler) handleSendMail(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.MailPayload
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ MailPayload JSON unmarshal başarısız", "error", err)
		return
	}

	utils.SendVerificationEmail(req.Email, req.Username, req.Code)

}

func (m *AuthHandler) handleVerifyCode(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.MailPayload
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ MailPayload JSON unmarshal başarısız", "error", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	storedCode, err := repo.GetVerificationCode(ctx, req.UUID)
	if err != nil {
		slog.Error("❌ Doğrulama kodu alınamadı", "uuid", req.UUID, "error", err)
		msg.Respond([]byte("ERROR"))
		return
	}
	if storedCode == req.Code {
		// Kod doğrulandı, oyuncuyu onayla
		err := repo.VerifyPlayer(ctx, req.UUID)
		if err != nil {
			slog.Error("❌ Oyuncu doğrulanamadı", "uuid", req.UUID, "error", err)
			msg.Respond([]byte("ERROR"))
			return
		}
		utils.SendWelcomeMessage(req.Email, req.Username)
		slog.Info("✅ Oyuncu doğrulandı", "uuid", req.UUID)
		msg.Respond([]byte("SUCCESS"))
	} else {
		slog.Info("❌ Doğrulama kodu yanlış", "uuid", req.UUID)
		msg.Respond([]byte("INVALID_OR_EXPIRED"))
	}
}

func (m *AuthHandler) setPlayerIsLoggedIn(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.SessionValidateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ SessionValidateRequest JSON unmarshal başarısız", "error", err)
		msg.Respond([]byte("false"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := repo.SetLoggedIn(ctx, req.UUID, true)
	if err != nil {
		slog.Error("❌ Oturum durumu kontrol edilemedi", "uuid", req.UUID, "error", err)
		msg.Respond([]byte("false"))
		return
	}
	slog.Info("✅ Oyuncu oturum açıldı", "uuid", req.UUID)
	msg.Respond([]byte("true"))
}
func (m *AuthHandler) handlePlayerIsLoggedIn(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.SessionValidateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ SessionValidateRequest JSON unmarshal başarısız", "error", err)
		msg.Respond([]byte("false"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	isLoggedIn, err := repo.IsLoggedIn(ctx, req.UUID)
	if err != nil {
		slog.Error("❌ Oturum durumu kontrol edilemedi", "uuid", req.UUID, "error", err)
		msg.Respond([]byte("false"))
		return
	}
	slog.Info("✅ Oyuncu oturum durumu kontrol edildi", "uuid", req.UUID, "isLoggedIn", isLoggedIn)
	msg.Respond([]byte(strconv.FormatBool(isLoggedIn)))
}

func (m *AuthHandler) IsVerified(msg *nats.Msg, repo *repository.PlayerRepository) {
	var req models.SessionValidateRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		slog.Error("❌ SessionValidateRequest JSON unmarshal başarısız", "error", err)
		msg.Respond([]byte("false"))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	isVerified, err := repo.IsEmailVerified(ctx, req.UUID)
	if err != nil {
		slog.Error("❌ E-posta doğrulama durumu kontrol edilemedi", "uuid", req.UUID, "error", err)
		msg.Respond([]byte("false"))
		return
	}
	slog.Info("✅ E-posta doğrulama durumu kontrol edildi", "uuid", req.UUID, "isVerified", isVerified)
	msg.Respond([]byte(strconv.FormatBool(isVerified)))
}
