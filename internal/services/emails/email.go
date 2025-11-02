package emails

import (
	"fmt"
	"log"
	"rio-go-model/configs"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Sender      string
	To          string
	Subject     string
	Body        string
	Host        string
	Port        int
	AppPassword string
}

func NewEmailService(name string, email string, message string) *EmailService {
	settings := configs.GetSettings()
	// Expect the following envs set in settings: SENDER_EMAIL, EMAIL_HOST, EMAIL_PORT, EMAIL_APP_PASSWORD
	sender := settings.EmailSender
	host := settings.EmailHost
	port := settings.EmailPort
	appPassword := settings.EmailAppPassword

	subject := fmt.Sprintf("New message from %s", name)
	body := fmt.Sprintf("Name: %s\nEmail: %s\n\nMessage:\n%s", name, email, message)

	return &EmailService{Sender: sender, To: email, Subject: subject, Body: body, Host: host, Port: port, AppPassword: appPassword}
}

func (e *EmailService) SendEmail() error {
	if e.Sender == "" || e.AppPassword == "" {
		log.Println("Email credentials not configured; skipping send")
		return fmt.Errorf("email credentials not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.Sender)
	m.SetHeader("To", e.To, configs.GlobalSettings.EmailTo)
	m.SetHeader("Subject", e.Subject)
	m.SetBody("text/plain", e.Body)

	d := gomail.NewDialer(e.Host, e.Port, e.Sender, e.AppPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending email:", err)
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
