package main

import (
	"log/slog"
	"os"

	"github.com/benerenla/best-plugin/internal/database"
	"github.com/benerenla/best-plugin/internal/messages"
	"github.com/lmittmann/tint"
)

func main() {

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05", // Saat:Dakika:Saniye formatı
		NoColor:    false,
	}))

	slog.SetDefault(logger)
	database.NewRedis(database.RedisConfig{Host: "auth-redis", Port: 6379})
	messages.NewMessage("nats://auth-nats:4222")
	database.NewMysql("root:root_password@tcp(auth-mysql:3306)/server_auth?charset=utf8mb4&parseTime=True&loc=Local")
}
