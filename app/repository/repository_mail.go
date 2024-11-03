package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"tanya_dokter_app/config"
	"text/template"

	"github.com/go-gomail/gomail"
)

func SendEmail(htmlTemplate, targetEmail, subjectMessage, attachmentPath string, fill any) {
	SmtpHost := config.LoadConfig().MailHost
	SmtpPort := config.LoadConfig().MailPort
	SmtpUsername := config.LoadConfig().MailUsername
	SmtpPassword := config.LoadConfig().MailPassword

	fmt.Println("Success to get host configuration. Configuration: ", fmt.Sprintf("%v:%v/%v:%v", SmtpHost, SmtpPort, SmtpUsername, SmtpPassword))

	htmlFile, err := os.ReadFile(config.RootPath() + "/assets/html/" + htmlTemplate + ".html")
	if err != nil {
		fmt.Println("Failed to read file:", err)
	} else {
		fmt.Println("Success to read file")
	}

	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(string(htmlFile))
	if err != nil {
		fmt.Println("Failed to parse template:", err)
	} else {
		fmt.Println("Success to parse template")
	}

	// Fill in the template with the data
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, fill)
	if err != nil {
		fmt.Println("Failed to fill in template:", err)
	} else {
		fmt.Println("Success to fill in template")
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", SmtpUsername)
	mailer.SetHeader("To", strings.ToLower(targetEmail))
	mailer.SetHeader("Subject", subjectMessage)
	mailer.SetBody("text/html", tpl.String())
	if attachmentPath != "" {
		mailer.Attach(attachmentPath)
	}

	d := gomail.NewDialer(SmtpHost, SmtpPort, SmtpUsername, SmtpPassword)
	if err := d.DialAndSend(mailer); err != nil {
		fmt.Println("Failed to send email. Error: ", err)
	} else {
		fmt.Println("Success to send email. Target Email: ", strings.ToLower(targetEmail))
	}

}

func BuildMessage(to []string, from, subject string, body bytes.Buffer) []byte {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", from)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(to, ";"))
	header += fmt.Sprintf("Subject: %s\r\n", subject)
	header += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	return []byte(header + body.String())
}

func AddAttachmentToMessage(message []byte, attachment []byte, attachmentPath string) []byte {
	// Buat penampung buffer
	var buffer bytes.Buffer

	// Tulis bagian header pesan
	buffer.Write(message)

	// Buat boundary string
	boundary := "boundary_string"

	// Tambahkan header untuk bagian lampiran
	attachmentHeader := fmt.Sprintf("\r\n--%s\r\nContent-Disposition: attachment; filename=\"%s\"\r\nContent-Type: application/octet-stream\r\n\r\n", boundary, filepath.Base(attachmentPath))

	// Tulis header lampiran
	buffer.WriteString(attachmentHeader)

	// Tulis isi lampiran
	buffer.Write(attachment)

	// Menambahkan footer boundary
	buffer.WriteString("\r\n--" + boundary + "--\r\n")

	return buffer.Bytes()
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func GetSMTPConfiguration() SMTPConfig {
	return SMTPConfig{
		Host:     config.LoadConfig().MailHost,
		Port:     config.LoadConfig().MailPort,
		Username: config.LoadConfig().MailUsername,
		Password: config.LoadConfig().MailPassword,
	}
}
