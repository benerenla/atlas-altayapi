package repository

import (
	"context"
	"log/slog"

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
		err = player.UnmarshalJSON([]byte(cachedData))
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
		return err
	}
	slog.Info("✅ Kayıt tamamlandı", "user", req.Username)
	return nil
}

func (r *PlayerRepository) IsRegistered(ctx context.Context, uuid models.IsRegisteredRequest) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Player{}).Where("uuid = ?", uuid.UUID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
