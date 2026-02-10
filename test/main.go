package main

import (
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222") // Bilgisayardan çalıştırıyorsan localhost
	if err != nil {
		panic(err)
	}
	defer nc.Close() // Bağlantıyı düzgün kapat

	data := `{"uuid": "550e8400-e29b-41d4-a716-446655440000", "username": "testuser", "password": "testpass"}`

	msg, err := nc.Request("mc.player.register", []byte(data), 2*time.Second)
	if err != nil {
		panic(err)
	}

	// Mesajın gerçekten gönderildiğinden emin olmak için tamponu boşalt
	nc.Flush()

	println("✅ Test mesajı başarıyla gönderildi!", string(msg.Data))
}
