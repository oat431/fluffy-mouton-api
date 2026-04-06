package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	CodeExpiration time.Duration
}

func GetEmailConfig() *Config {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587
	}
	return &Config{
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       port,
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPassword:   os.Getenv("SMTP_PASS"),
		CodeExpiration: time.Minute,
	}
}
