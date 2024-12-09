package repository

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"tanya_dokter_app/app/middlewares"
	"tanya_dokter_app/app/models"
	"tanya_dokter_app/app/reqres"
	"tanya_dokter_app/config"
	"text/template"
	"time"

	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type email struct {
	to      []string
	subject string
	body    string
}

func SignIn(email, password string) (user models.GlobalUser, token string, err error) {
	err = config.DB.
		Where("email = '" + strings.ToLower(email) + "'").First(&user).Error
	if err != nil {
		return
	}
	err = middlewares.VerifyPassword(password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		err = errors.New("incorrect password")
		return
	}
	if user.EmailVerifiedAt.IsZero() {
		err = errors.New("please verify your email")
		return
	}
	token, err = middlewares.AuthMakeToken(user)
	if err != nil {
		return
	}
	return
}

func GetVerificationToken(request *reqres.EmailRequest) (data string, err error) {
	user, err := GetUserByEmail(request.Email)
	if err != nil {
		log.Println("Failed to get user:", err)
		return
	}
	rand.Seed(time.Now().UnixNano())
	pin := fmt.Sprintf("%06d", rand.Intn(1000000))

	go sendEmailAuthentification(pin, request.Email, user.Fullname, "/assets/html/email-verification.html", "Verifikasi Email")

	data = "Email sent successfully."
	return
}

func sendEmailAuthentification(pin, emailTo string, fullname, templateHTML string, subjectTitle string) {
	log.Println("Sending Email")

	htmlFile, err := os.ReadFile(config.RootPath() + templateHTML)
	if err != nil {
		log.Println("Failed to read file:", err)
		return
	}

	type Data struct {
		AppName  string
		PIN      string
		Email    string
		Logo     string
		Fullname string
	}

	mailHost, mailPort, mailUsername, mailPassword, _, err := GetEmailHostConfiguration()
	if err != nil {
		log.Println("Failed to get host configuration:", err)
	} else {
		log.Println("mailHost:", mailHost)
		log.Println("mailPort:", mailPort)
		log.Println("mailUsername:", mailUsername)
		log.Println("mailPassword:", mailPassword)
	}

	fill := &Data{
		AppName:  config.LoadConfig().AppName,
		PIN:      pin,
		Fullname: fullname,
		Email:    emailTo,
	}

	logoBase64, err := GetBase64Logo()
	if err != nil {
		log.Println("Failed to load logo:", err)
	} else {
		fill.Logo = logoBase64
	}

	// Parse the HTML template
	tmpl, err := template.New("emailTemplate").Parse(string(htmlFile))
	if err != nil {
		log.Println("Failed to parse template:", err)
		return
	}

	// Fill in the template with the data
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, fill)
	if err != nil {
		log.Println("Failed to fill in template:", err)
		return
	}

	// Send the email with the filled-in template
	email := &email{
		to:      []string{emailTo},
		subject: subjectTitle,
		body:    tpl.String(),
	}

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", mailHost, mailPort),
		smtp.PlainAuth("", mailUsername, mailPassword, mailHost),
		mailUsername,
		email.to,
		email.buildMessage(),
	)

	if err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("Success to send email to", emailTo)
	}

}

func EmailVerification(request *reqres.TokenRequest) (data string, err error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}

	user, err := GetUserByEmail(request.Email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	user.EmailVerifiedAt = null.TimeFrom(time.Now().In(location))
	UpdateUser(user)

	data = "Email verification success"
	return data, nil
}

func (e *email) buildMessage() []byte {
	header := ""
	header += fmt.Sprintf("From: %s\r\n", config.LoadConfig().MailUsername)
	header += fmt.Sprintf("To: %s\r\n", strings.Join(e.to, ";"))
	header += fmt.Sprintf("Subject: %s\r\n", e.subject)
	header += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	return []byte(header + e.body)
}

func GetEmailHostConfiguration() (mailHost string, mailPort int, mailUsername string, mailPassword string, encryption string, err error) {
	mailHost = config.LoadConfig().MailHost
	mailPort = config.LoadConfig().MailPort
	mailUsername = config.LoadConfig().MailUsername
	mailPassword = config.LoadConfig().MailPassword
	encryption = config.LoadConfig().MailEncryption

	return mailHost, mailPort, mailUsername, mailPassword, encryption, err
}

func GetLogo() ([]byte, error) {
	// Path ke file logo di dalam folder assets
	logoPath := "./assets/logo.png"

	// Membaca file logo
	logo, err := ioutil.ReadFile(logoPath)
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func GetBase64Logo() (string, error) {
	logo, err := GetLogo()
	if err != nil {
		return "", err
	}

	// Mengonversi logo ke base64
	encodedLogo := base64.StdEncoding.EncodeToString(logo)
	return encodedLogo, nil
}

func SendResetPassword(pin, email string) (err error) {

	user, err := GetUserByEmail(strings.ToLower(email))
	if err != nil {
		return
	}

	go sendEmailAuthentification(pin, user.Email, user.Fullname, "/assets/html/reset-password.html", "Permintaan Lupa Password")

	return
}

func ResetPassword(pin, password, email string) (user models.GlobalUser, tokenOutput string, err error) {

	decoded, err := base64.StdEncoding.DecodeString(pin)
	if err != nil {
		return
	}

	newdecoded, err := base64.StdEncoding.DecodeString(string(decoded))
	if err != nil {
		return
	}
	parts := strings.Split(string(newdecoded), "&")
	err = config.DB.
		Where("email = '" + parts[0] + "'").First(&user).Error
	if err != nil {
		return
	}
	if user.ID >= 1 {
		newPassword := middlewares.BcryptPassword(password)
		user.Password = newPassword
		err = config.DB.Save(&user).Error
		if err != nil {
			return
		}

	}

	tokenOutput, err = middlewares.AuthMakeToken(user)
	if err != nil {
		return
	}
	return
}

func SaveResetRequest(email, pin string, expiresAt time.Time) error {
	resetRequest := reqres.ResetPasswordRequest{
		Email:     strings.ToLower(email),
		Pin:       pin,
		ExpiresAt: expiresAt,
	}
	if err := config.DB.Create(&resetRequest).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePassword(email, newPassword string) error {
	var user models.GlobalUser

	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword := middlewares.BcryptPassword(newPassword)
	user.Password = hashedPassword

	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}

	err = config.DB.Where("email = ?", email).Delete(&models.ResetPasswordRequest{}).Error
	if err != nil {
		return err
	}

	return nil
}

func ValidatePin(email, pin string) error {
	var resetRequest reqres.ResetPasswordRequest
	err := config.DB.
		Where("email = ? AND pin = ?", email, pin).
		Where("expires_at > ?", time.Now()).
		First(&resetRequest).Error
	if err != nil {
		return errors.New("invalid email or pin")
	}
	return nil
}
