package main

import (
	"log/slog"
	"os"

	"github.com/benerenla/best-plugin/internal/database"
	"github.com/benerenla/best-plugin/internal/messages"
	"github.com/benerenla/best-plugin/internal/repository"
	"github.com/lmittmann/tint"
)

func main() {

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "15:04:05", // Saat:Dakika:Saniye formatı
		NoColor:    false,
	}))

	slog.SetDefault(logger)
	slog.Info("Sistem başlatılıyor..")
	redisClient := database.NewRedis(database.RedisConfig{Host: "auth-redis", Port: 6379})
	messages.NewMessage("nats://auth-nats:4222")
	mysqlClient, _ := database.NewMysql("root:root_password@tcp(auth-mysql:3306)/server_auth?charset=utf8mb4&parseTime=True&loc=Local")
	handler := messages.NewAuthHandler(messages.NewMessage("nats://auth-nats:4222").GetConnection())
	repository := repository.NewPlayerRepository(mysqlClient.DB, redisClient)

	handler.RegisterHandlers(repository)

	select {}
}
