package database

import (
	"log/slog"
	"time"

	"github.com/benerenla/best-plugin/internal/models"
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
	if !db.Migrator().HasTable(&models.Player{}) {
		slog.Info("📦 Tablolar oluşturuluyor...")
		err = db.AutoMigrate(&models.Player{})
		if err != nil {
			slog.Error("❌ Tablolar oluşturulamadı.", "error", err)
			return nil, err
		}
		slog.Info("✅ Tablolar başarıyla oluşturuldu.")
	} else {
		slog.Info("✅ Tablo zaten mevcut, migrasyon atlandı.")
	}

	if err != nil {
		slog.Error("❌ Tablolar oluşturulamadı.", "error", err)
		return nil, err
	}

	return &Storage{DB: db}, nil
}
