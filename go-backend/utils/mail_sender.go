package utils

import (
	"fmt"
	"log/slog"

	"gopkg.in/gomail.v2"
)

const (
	SMTPHost     = "smtp.hostinger.com" // SMTP sunucunuzun adresi (örn: "smtp.gmail.com")
	SMTPPort     = 465                 // SMTP sunucunuzun portu (örn: 587 veya 465)
	SMTPUsername = "" // Mail adresiniz (örn: "info@oshnetwork.shop gibi)
	SMTPPassword = "," // Mail adresinizin şifresi
)

func SendVerificationEmail(toEmail string, username string, code string) {
	// 1. Yeni bir "Boş Mektup" kağıdı alıyoruz
	m := gomail.NewMessage()

	// 2. Zarfın üzerine adresleri yazıyoruz
	m.SetHeader("From", SMTPUsername)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Osh Network - Doğrulama Kodu")

	// 3. Mektubun içeriğini (HTML) hazırlıyoruz
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Doğrulama</title>
	<style>
        @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;800&display=swap');
        body { margin: 0; padding: 0; background-color: #f4f4f9; font-family: 'Poppins', sans-serif !important; }
    </style>
</head>
<body style="margin: 0; padding: 0; background-color: #f4f4f9;">
    <div style="max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 15px rgba(0,0,0,0.1);">
        
        <div style="background: linear-gradient(135deg, #5f07c4ff 0%%, #250561ff 100%%); padding: 40px 20px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 700; letter-spacing: 1px;">OSH NETWORK</h1>
            <p style="color: rgba(255,255,255,0.9); margin-top: 10px; font-size: 16px;">Hesap Doğrulama İşlemi</p>
        </div>

        <div style="padding: 40px 30px; text-align: center;">
            <h2 style="color: #2d3436; font-size: 22px; margin-bottom: 20px;">Hoş Geldin, <span style="color: #5f07c4ff;">%s</span>!</h2>
            
            <p style="color: #636e72; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">
                Sunucumuza katılmana çok sevindik. Hesabını güvene almak ve maceraya başlamak için aşağıdaki doğrulama kodunu kullanman gerekiyor.
            </p>

            <div style="background-color: #f8f6ff; border: 2px dashed #5f07c4ff; border-radius: 8px; padding: 20px; margin: 0 auto 30px auto; display: inline-block;">
                <span style="font-size: 36px; font-weight: 800; color: #5f07c4ff; letter-spacing: 8px; display: block;">%s</span>
            </div>

            <p style="color: #636e72; font-size: 14px;">
                Bu kodu oyun içinde <b>/onayla [kod]</b> şeklinde yazabilirsin.
            </p>
        </div>

        <div style="background-color: #f9f9fc; padding: 20px; text-align: center; border-top: 1px solid #eee;">
            <p style="color: #646e72ff; font-size: 12px; margin: 0;">
                Bu kod güvenlik nedeniyle <b>5 dakika</b> içinde geçerliliğini yitirecektir.<br>
                Eğer bu isteği sen yapmadıysan, bu maili görmezden gelebilirsin.
                Yardım için <a style="color: #5f07c4ff" href="oshnetwork.shop/discord"><b>Discord</b></a>
                sunucumuza katılabilirsin.
            </p>
        </div>
    </div>
</body>
</html>`, username, code)
	m.SetBody("text/html", htmlBody)

	// 4. Postaneye (SMTP Sunucusuna) giden yolu tarif ediyoruz
	d := gomail.NewDialer(SMTPHost, SMTPPort, SMTPUsername, SMTPPassword)

	// 5. Kapıyı çal, içeri gir, mektubu ver ve çık
	if err := d.DialAndSend(m); err != nil {
		slog.Error("❌ Mail gönderilemedi", "hata", err)
	} else {
		slog.Info("📧 Mail başarıyla uçtu!", "kime", toEmail)
	}
}

func SendWelcomeMessage(toEmail string, username string) {
	// 1. Yeni bir "Boş Mektup" kağıdı alıyoruz
	m := gomail.NewMessage()

	// 2. Zarfın üzerine adresleri yazıyoruz
	m.SetHeader("From", SMTPUsername)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Osh Network - Sunucumuza Hoşgeldin")

	// 3. Mektubun içeriğini (HTML) hazırlıyoruz
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;600;800&display=swap');
        body { margin: 0; padding: 0; background-color: #f4f4f9; font-family: 'Poppins', sans-serif !important; }
    </style>
</head>
<body>
    <div style="max-width: 600px; margin: 40px auto; background-color: #f4f4f9; border-radius: 16px; overflow: hidden; box-shadow: 0 10px 30px rgba(108, 92, 231, 0.2);">
        
        <div style="background: linear-gradient(135deg, #5f07c4ff 0%%, #250561ff 100%%); padding: 40px 20px; text-align: center;">
            <h1 style="color: #ffffff; margin: 0; font-size: 28px; font-weight: 800; letter-spacing: 1px;">HESAP ONAYLANDI!</h1>
        </div>

        <div style="padding: 40px 30px; text-align: center;">
            <h2 style="color: #2d3436; font-size: 22px; margin-bottom: 20px;">Selam, <span style="color: #5f07c4ff;">%s</span>!</h2>
            
            <p style="color: #636e72; font-size: 16px; line-height: 1.6;">
                E-posta adresin başarıyla doğrulandı. Artık <b style="color: #5f07c4ff;">Osh Network</b> sunucusuna giriş yapabilirsin.
            </p>

            <div style="margin: 30px 0;">
                <p style="font-size: 18px; color: #2d3436; font-weight: 600;">Seni Neler Bekliyor?</p>
                <ul style="text-align: left; display: inline-block; color: #636e72;">
                    <li> Güvenli Hesap</li>
                    <li> VIP Çekilişlerine Katılım</li>
                    <li> Özel Etkinlik Bildirimleri</li>
                </ul>
            </div>

            <a href="https://oshnetwork.com" style="background-color: #5f07c4ff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 50px; font-weight: bold; display: inline-block;">Sitemize Göz At</a>

        </div>
    </div>
</body>
</html>`, username)
	m.SetBody("text/html", htmlBody)

	// 4. Postaneye (SMTP Sunucusuna) giden yolu tarif ediyoruz
	d := gomail.NewDialer(SMTPHost, SMTPPort, SMTPUsername, SMTPPassword)

	// 5. Kapıyı çal, içeri gir, mektubu ver ve çık
	if err := d.DialAndSend(m); err != nil {
		slog.Error("❌ Mail gönderilemedi", "hata", err)
	} else {
		slog.Info("📧 Mail başarıyla uçtu!", "kime", toEmail)
	}
}
