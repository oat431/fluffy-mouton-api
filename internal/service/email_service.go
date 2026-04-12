package service

import (
	"fmt"
	"net/smtp"
	"oat431/fluffy-mouton/internal/config"
)

type SMTPService struct {
	config *config.Config
}

func NewSMTPService(cfg *config.Config) *SMTPService {
	return &SMTPService{config: cfg}
}

func (s *SMTPService) SendMail(to string) error {
	subject := "Subject: Fluffy Mouton Account Verification \n"
	body := `
		If you are receiving this email, it means that your account has been successfully created. \n
		Please verify your email address to complete the registration process.\n
		Thank you for joining Fluffy Mouton!
	`

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUser,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPUser,
		[]string{to},
		message,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMTPService) SendVerificationEmail(to, token string) error {
	subject := "Subject: Fluffy Mouton Account Verification\n"
	body := fmt.Sprintf(`
		If you are receiving this email, it means that your account has been successfully created.
		Please verify your email address by clicking the link here:http://localhost:3000/verify-email?token=%s
		Thank you for joining Fluffy Mouton!
	`, token)

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUser,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPUser,
		[]string{to},
		message,
	)
	if err != nil {
		return err
	}

	return nil
}
