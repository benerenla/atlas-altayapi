package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	metric "github.com/benerenla/best-plugin/internal/metrics"
	"github.com/benerenla/best-plugin/internal/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PlayerRepository struct {
	// Bu struct, oyuncu verilerini yönetmek için gerekli yöntemleri içerebilir.
	// Örneğin: GetPlayerByID, CreatePlayer, UpdatePlayer, DeletePlayer gibi yöntemler.
	db    *gorm.DB
	redis *redis.Client
}

func NewPlayerRepository(db *gorm.DB, redisClient *redis.Client) *PlayerRepository {
	return &PlayerRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, uuid string) (*models.Player, error) {
	var player models.Player
	cacheKey := "player:" + uuid

	// Öncelikle Redis'te oyuncu verisini kontrol et
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// Redis'te veri bulundu, JSON'dan Player struct'ına dönüştür
		var player models.Player
		err = json.Unmarshal([]byte(cachedData), &player)
		if err == nil {
			return &player, nil
		}
	}

	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&player).Error; err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *PlayerRepository) RegisterPlayer(ctx context.Context, req *models.Player) error {
	// Kayıt işlemi burada yapılır
	// Örneğin, yeni bir oyuncu oluşturup veritabanına kaydedebilirsiniz.
	// Kayıt işlemi başarılı olduktan sonra, oyuncu verisini Redis'e kaydedebilirsiniz.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err := r.db.WithContext(ctx).Create(&models.Player{
		UUID:     req.UUID,
		Username: req.Username,
		Password: string(hashedPassword),
	}).Error; err != nil {
		slog.Info("❌ Kayıt tamamlanamadı", "user", req.Username, "error", err)
		return err
	}
	slog.Info("✅ Kayıt tamamlandı", "user", req.Username)
	metric.TotalRegisters.Inc() // Grafana'da çizgi yukarı çıkacak!
	return nil
}

func (r *PlayerRepository) IsRegistered(ctx context.Context, uuid models.IsRegisteredRequest) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Player{}).Where("uuid = ?", uuid.UUID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PlayerRepository) LoginPlayer(ctx context.Context, req *models.LoginPlayerRequest) (*models.Player, error) {
	var player models.Player

	if err := r.db.WithContext(ctx).Where("uuid = ?", req.UUID).First(&player).Error; err != nil {
		return nil, err
	}

	newPlayer := models.Player{
		UUID:          player.UUID,
		Username:      player.Username,
		Level:         player.Level,
		XP:            player.XP,
		Coins:         player.Coins,
		Rank:          player.Rank,
		DiscordID:     player.DiscordID,
		Verified:      player.Verified,
		CurrentServer: player.CurrentServer,
	}

	data, err := json.Marshal(&newPlayer)

	if err != nil {
		slog.Error("❌ Redis'e oyuncu verisi kaydedilemedi", "error", err)
	} else {
		// Redis'e oyuncu verisini kaydet
		err = r.redis.Set(ctx, "player:"+newPlayer.UUID, data, 0).Err()
		if err != nil {
			slog.Error("❌ Redis'e oyuncu verisi kaydedilemedi", "error", err)
		}
	}
	// Şifre doğrulaması yap
	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(req.Password)); err != nil {
		return nil, err
	}
	return &newPlayer, nil
}

func (r *PlayerRepository) UpdatePlayerLastSeen(ctx context.Context, player *models.Player) error {
	// Oyuncu verisini güncellemek için bu yöntemi kullanabilirsiniz.
	// Örneğin, oyuncunun son görüldüğü zamanı güncelleyebilirsiniz.
	return r.db.WithContext(ctx).Model(player).Update("last_seen", player.LastSeen).Error
}

func (r *PlayerRepository) UpdatePlayer(ctx context.Context, player *models.Player) error {
	// Oyuncu verisini güncellemek için bu yöntemi kullanabilirsiniz.
	// Örneğin, oyuncunun son görüldüğü zamanı güncelleyebilirsiniz.
	return r.db.WithContext(ctx).Save(player).Error
}

// Kodu Redis'e kaydeder (5 Dakika Süreli)
func (r *PlayerRepository) SaveVerificationCode(ctx context.Context, uuid string, code string) error {
	// Key: "verify:uuid" -> Value: "123456" (TTL: 5 Dakika)
	return r.redis.Set(ctx, "verify:"+uuid, code, 5*time.Minute).Err()
}

// Kodu doğrular ve kullanıcıyı onaylar
func (r *PlayerRepository) VerifyEmail(ctx context.Context, uuid string, codeInput string) (bool, error) {
	key := "verify:" + uuid

	val, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return false, nil // Kod bulunamadı (Süresi dolmuş veya hiç yok)
	}

	if val != codeInput {
		return false, nil // Kod yanlış
	}

	err = r.db.WithContext(ctx).Model(&models.Player{}).
		Where("uuid = ?", uuid).
		Update("verified", true).Error

	if err != nil {
		slog.Error("❌ MySQL güncelleme hatası", "err", err)
		return false, err
	}

	r.redis.Del(ctx, key)

	r.redis.Del(ctx, "player:"+uuid)

	return true, nil
}

func (r *PlayerRepository) IsEmailVerified(ctx context.Context, uuid string) (bool, error) {
	var player models.Player
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&player).Error; err != nil {
		return false, err
	}

	return player.Verified, nil
}

func (r *PlayerRepository) GetVerificationCode(ctx context.Context, uuid string) (string, error) {
	key := "verify:" + uuid
	return r.redis.Get(ctx, key).Result()
}

func (r *PlayerRepository) VerifyPlayer(ctx context.Context, uuid string) error {
	err := r.db.WithContext(ctx).Model(&models.Player{}).
		Where("uuid = ?", uuid).
		Update("verified", true).Error

	if err != nil {
		slog.Error("❌ MySQL güncelleme hatası", "err", err)
		return err
	}
	r.redis.Del(ctx, "player:"+uuid)
	return nil
}

func (r *PlayerRepository) SetEmail(ctx context.Context, uuid string, email *string) error {
	return r.db.WithContext(ctx).Model(&models.Player{}).
		Where("uuid = ?", uuid).
		Update("email", email).Error
}

func (r *PlayerRepository) GetEmail(ctx context.Context, uuid string) (*string, error) {
	var player models.Player
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&player).Error; err != nil {
		return nil, err
	}
	return player.Email, nil
}

func (r *PlayerRepository) SetLoggedIn(ctx context.Context, uuid string, isLoggedIn bool) error {
	return r.redis.Set(ctx, "session:"+uuid, isLoggedIn, 24*time.Hour).Err()
}

func (r *PlayerRepository) IsLoggedIn(ctx context.Context, uuid string) (bool, error) {
	val, err := r.redis.Get(ctx, "session:"+uuid).Result()
	if err != nil {
		return false, nil // Oturum bilgisi bulunamadı
	}
	return val == "true", nil
}
func (r *PlayerRepository) IsVerıfedEmail(ctx context.Context, uuid string) (bool, error) {
	var player models.Player
	if err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&player).Error; err != nil {
		return false, err
	}
	return player.Verified, nil
}
