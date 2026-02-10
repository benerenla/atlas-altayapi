package database

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host string
	Port int
}

func NewRedis(config RedisConfig) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Password: "", // Şifre yoksa boş bırakın
		DB:       0,  // Varsayılan DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("❌ Redis bağlantısı başarısız.", "error", err)
		return
	}
	slog.Info("✅ Redis bağlantısı başarılı ")
}
