package main

// Backend'e gidecek veri yapıları


/*
func main() {
	// 1. NATS Bağlantısı
	fmt.Println("🔌 NATS Sunucusuna Bağlanılıyor...")
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		fmt.Printf("❌ Kritik Hata: NATS'e bağlanılamadı! Docker çalışıyor mu? (%v)\n", err)
		return
	}
	defer nc.Close()
	fmt.Println("✅ Bağlantı Başarılı! Simülasyon Başlıyor...\n")

	// Klavye okuyucusu
	reader := bufio.NewReader(os.Stdin)

	// ---------------------------------------------------------
	// ADIM 1: BİLGİLERİ TOPLA
	// ---------------------------------------------------------
	fmt.Println("📝 Lütfen Kayıt Bilgilerini Giriniz:")

	fmt.Print("👉 Kullanıcı Adı (User ID): ")
	uuid, _ := reader.ReadString('\n')
	uuid = strings.TrimSpace(uuid)


	// Kullanıcı Adı
	fmt.Print("👉 Kullanıcı Adı (User ID): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// E-Posta
	fmt.Print("👉 E-Posta Adresi: ")
	emailInput, _ := reader.ReadString('\n')
	emailInput = strings.TrimSpace(emailInput)

	// Şifre
	fmt.Print("👉 Şifre: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)


	// ---------------------------------------------------------
	// ADIM 2: KAYIT İSTEĞİ GÖNDER (REGISTER)
	// ---------------------------------------------------------
	fmt.Println("\n⏳ Sunucuya kayıt isteği gönderiliyor...")

	regReq := RegisterRequest{
		UUID:     uuid,
		Username: username,
		Password: password,
	}

	regData, _ := json.Marshal(regReq)

	// Backend'e NATS üzerinden soruyoruz
	respMsg, err := nc.Request("mc.player.register", regData, 3*time.Second)
	if err != nil {
		fmt.Printf("❌ Sunucu Cevap Vermedi (Timeout): %v\n", err)
		return
	}

	responseStr := string(respMsg.Data)
	fmt.Printf("📨 Sunucu Cevabı: %s\n", responseStr)

	if responseStr != "SUCCESS" {
		fmt.Println("❌ Kayıt başarısız oldu, işlem durduruluyor.")
		return
	}
	newRegData := RegisterRequest{
		UUID:     username,
		Username: username,
		Email:    &emailInput,
	}

	emailData, _ := json.Marshal(newRegData)

	emailMessage, err := nc.Request("mc.player.verify_email", emailData, 3*time.Second)
	if err != nil {
		fmt.Printf("❌ Doğrulama isteği sırasında hata: %v\n", err)
		return
	}
	fmt.Println("✅ Kayıt Başarılı! Mail kutuna doğrulama kodu gönderildi.", emailMessage.Data)
	fmt.Println("---------------------------------------------------------")

	// ---------------------------------------------------------
	// ADIM 3: DOĞRULAMA KODU GİR (VERIFY)
	// ---------------------------------------------------------
	fmt.Print("🔑 Lütfen Mailinize Gelen Kodu Giriniz: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)

	fmt.Println("⏳ Kod doğrulanıyor...")

	verifyReq := VerifyRequest{
		UUID:     uuid,
		Code:     code,
		Email:    &emailInput,
		Username: username,
	}
	verifyData, _ := json.Marshal(verifyReq)

	// Backend'e doğrulama isteği at
	verifyResp, err := nc.Request("mc.player.verify", verifyData, 3*time.Second)
	if err != nil {
		fmt.Printf("❌ Doğrulama sırasında hata: %v\n", err)
		return
	}

	finalResponse := string(verifyResp.Data)

	if finalResponse == "SUCCESS" {
		fmt.Println("\n🎉 TEBRİKLER! Hesap başarıyla doğrulandı ve aktifleştirildi.")
	} else if finalResponse == "INVALID_OR_EXPIRED" {
		fmt.Println("\n⚠️ HATA: Girdiğin kod yanlış veya süresi dolmuş.")
	} else {
		fmt.Printf("\n⚠️ Bilinmeyen Durum: %s\n", finalResponse)
	}
}

// 	data := `{"uuid": "550e8400-e29b-41d4-a716-446655440021", "username": "testuser1", "email": "testuser1@example.com"}`
*/
