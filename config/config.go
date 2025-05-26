package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/Prototype-1/freelanceX_message.notification_service/email"
)

type SMTPConfig struct {
    EmailSender string `mapstructure:"EMAIL_SENDER"`
    EmailPass   string `mapstructure:"EMAIL_PASSWORD"`
    SMTPHost    string `mapstructure:"EMAIL_SMTP_HOST"`
    SMTPPort    string `mapstructure:"SMTP_PORT"`
}

type Config struct {
	MongoURI      string
	DatabaseName  string
	ServerPort    string
	UserServiceAddress string 
	InvoiceServiceAddress string
	RedisAddr     string
	SMTP         email.SMTPConfig
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on system environment variables...")
	}

	    v := viper.New()
    v.AutomaticEnv()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is required but not set")
	}

	databaseName := os.Getenv("MONGO_DB")
	if databaseName == "" {
		databaseName = "freelanceX_proposals" 
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":50055" 
	}

	redisAddr := os.Getenv("REDIS_ADDR")
if redisAddr == "" {
	log.Fatal("REDIS_ADDR is required but not set")
}

	  smtpCfg := email.SMTPConfig{
        EmailSender: v.GetString("EMAIL_SENDER"),
        EmailPass:   v.GetString("EMAIL_PASSWORD"),
        SMTPHost:    v.GetString("EMAIL_SMTP_HOST"),
        SMTPPort:    v.GetString("SMTP_PORT"),
    }

    if smtpCfg.EmailSender == "" || smtpCfg.EmailPass == "" || smtpCfg.SMTPHost == "" || smtpCfg.SMTPPort == "" {
        log.Println("Warning: SMTP config is incomplete, email sending may fail")
    }

	userServiceAddr := v.GetString("USER_SERVICE_GRPC_ADDR")
if userServiceAddr == "" {
    log.Fatal("USER_SERVICE_GRPC_ADDR is required but not set")
}

	InvoiceServiceAddress := v.GetString("INVOICE_SERVICE_GRPC_ADDR")
if InvoiceServiceAddress == "" {
    log.Fatal("INVOICE_SERVICE_GRPC_ADDR is required but not set")
}

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: databaseName,
		ServerPort:   serverPort,
		UserServiceAddress: userServiceAddr, 
		RedisAddr:    redisAddr,
		SMTP:         smtpCfg,
	}
}
