package database

import (
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func NewMysql(dsn string) (*Storage, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		slog.Error("❌ Database bağlantısı başarısız.", "error", err)
		return nil, err
	}
	slog.Info("✅ Database bağlantısı başarılı",
		"host", "localhost",
		"db", "minecraft",
		"latency", time.Since(time.Now()),
	)
	return &Storage{DB: db}, nil
}
