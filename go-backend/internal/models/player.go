package models

import (
	"encoding/json"
	"time"
)

type Player struct {
	// gorm:"primaryKey" -> Birincil anahtar
	// gorm:"uniqueIndex" -> Hızlı arama için indeksleme
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID     string `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
	Username string `gorm:"type:varchar(16);uniqueIndex;not null" json:"username"`
	LastIP   string `gorm:"type:varchar(45)" json:"last_ip"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // Şifre alanı JSON çıktısında gösterilmez

	// Verifiction
	Email      *string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	Verified   bool    `gorm:"default:false" json:"verified"`
	VerifyCode *string `gorm:"type:varchar(255)" json:"-"`
	DiscordID  *string `gorm:"type:varchar(20);uniqueIndex" json:"discord_id"`

	// Oyun Verileri
	Coins int64  `gorm:"default:0" json:"coins"`
	Level int    `gorm:"default:1" json:"level"`
	XP    int64  `gorm:"default:0" json:"xp"`
	Rank  string `gorm:"type:varchar(20);default:'player'" json:"rank"`

	// Kalıcı olmayan (Sadece Redis/Memory) veriler
	// gorm:"-" etiketi GORM'a bu alanı MySQL'e kaydetmemesini söyler
	CurrentServer   string `gorm:"-" json:"current_server"`
	IsOnline        bool   `gorm:"-" json:"is_online"`
	IsRegistered    bool   `gorm:"-" json:"is_registered"`
	IsAuthenticated bool   `gorm:"-" json:"is_authenticated"`

	// Zaman Bilgileri
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	LastSeen  *time.Time `json:"last_seen"`

	// Bans
	IsBanned   bool       `gorm:"default:false" json:"is_banned"`
	BanReason  string     `gorm:"type:varchar(255)" json:"ban_reason"`
	BanExpires *time.Time `json:"ban_expires"`
}

type RegisterPlayerRequest struct {
	UUID     string `json:"uuid" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type IsRegisteredRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

type LoginPlayerRequest struct {
	UUID     string `json:"uuid" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SessionValidateRequest struct {
	UUID string `json:"uuid" binding:"required"`
}

func (p *Player) toJSON() ([]byte, error) {
	return json.Marshal(p)
}

func playerFromJSON(data []byte) (*Player, error) {
	var player Player
	err := json.Unmarshal(data, &player)
	return &player, err
}
