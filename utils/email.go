package utils

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"github.com/k3a/html2text"
	"github.com/tonybobo/auth-template/config"
	"github.com/tonybobo/auth-template/models"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

func SendEmail(user *models.DBResponse, data *EmailData, temp *template.Template, templateName string) error {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Error Loading Config", err)
	}

	from := config.EmailFrom
	smtpUser := config.SMTPUser
	smtpPass := config.SMTPPass
	to := user.Email
	smptHost := config.SMTPHost
	smtpPort := config.SMTPPort

	var body bytes.Buffer

	if err := temp.ExecuteTemplate(&body, templateName, &data); err != nil {
		log.Fatal("error parsing templates", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)

	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smptHost, smtpPort, smtpUser, smtpPass)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil

}
