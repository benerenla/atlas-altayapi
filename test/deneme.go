package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type RegisterRequest struct {
	UUID     string  `json:"uuid"`
	Username string  `json:"username"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type VerifyRequest struct {
	UUID     string  `json:"uuid"`
	Code     string  `json:"code"`
	Email    *string `json:"email,omitempty"`
	Username string  `json:"username,omitempty"`
}

func main() {
	fmt.Println("🔌 NATS Sunucusuna Bağlanılıyor...")
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		fmt.Printf("❌ Kritik Hata: NATS'e bağlanılamadı! Docker çalışıyor mu? (%v)\n", err)
		return
	}
	defer nc.Close()
	fmt.Println("✅ Bağlantı Başarılı! Simülasyon Başlıyor...\n")

	regReq := RegisterRequest{
		UUID: "9e78dd1c-6d63-3ff4-a3bc-ee8258fcb42b",
	}
	regData, _ := json.Marshal(regReq)

	respMsg, err := nc.Request("mc.player.is_registered", regData, 3*time.Second)
	if err != nil {
		fmt.Printf("❌ Kayıt isteği başarısız oldu! (%v)\n", err)
		return
	}

	fmt.Printf("✅ Sunucudan Gelen Yanıt: %s\n", string(respMsg.Data))
}
