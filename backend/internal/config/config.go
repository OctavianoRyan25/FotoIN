package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	Port            string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
	Database        DatabaseConfig
	Redis           RedisConfig
	S3              S3Config
	Google          GoogleConfig
	SMTP            SMTPConfig
	Duitku          DuitkuConfig
	Midtrans        MidtransConfig
}

type DatabaseConfig struct {
	Host        string
	User        string
	Password    string
	Name        string
	Port        string
	SSLMode     string
	Environment string
}

type RedisConfig struct {
	URL string
}

type S3Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type DuitkuConfig struct {
	Environment   string
	APIKey        string
	MerchantCode  string
	PaymentMethod string
	CallbackURL   string
	ReturnURL     string
}

type MidtransConfig struct {
	Environment string
	ServerKey   string
}

func LoadConfig() Config {
	backendURL := getEnv("BACKEND_URL", "http://localhost:8080")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	callbackURL := getEnv("DUIDKU_CALLBACK_URL", backendURL+"/api/v1/webhooks/payment")
	returnURL := getEnv("DUIDKU_RETURN_URL", frontendURL+"/payment/{order_id}")
	redirectURL := getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")
	port := getEnv("APP_PORT", "8080")

	return Config{
		Port:            port,
		AllowedOrigins:  parseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS")),
		ShutdownTimeout: 10 * time.Second,
		Database: DatabaseConfig{
			Host:        os.Getenv("DB_HOST"),
			User:        os.Getenv("DB_USER"),
			Password:    os.Getenv("DB_PASSWORD"),
			Name:        os.Getenv("DB_NAME"),
			Port:        os.Getenv("DB_PORT"),
			SSLMode:     os.Getenv("DB_SSLMODE"),
			Environment: os.Getenv("APP_ENV"),
		},
		Redis: RedisConfig{
			URL: getEnv("URL_REDIS", "redis://localhost:6379/0"),
		},
		S3: S3Config{
			Region:    getEnv("AWS_REGION", os.Getenv("S3_REGION")),
			AccessKey: getEnv("AWS_ACCESS_KEY_ID", os.Getenv("S3_ACCESS_KEY")),
			SecretKey: getEnv("AWS_SECRET_ACCESS_KEY", os.Getenv("S3_SECRET_KEY")),
			Bucket:    getEnv("AWS_S3_BUCKET", os.Getenv("S3_BUCKET")),
			Endpoint:  getEnv("AWS_S3_ENDPOINT", os.Getenv("S3_ENDPOINT")),
		},
		Google: GoogleConfig{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  redirectURL,
		},
		SMTP: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: getEnv("SMTP_USER", os.Getenv("SMTP_USERNAME")),
			Password: getEnv("SMTP_PASS", os.Getenv("SMTP_PASSWORD")),
			From:     getEnv("SMTP_FROM", os.Getenv("SMTP_USER")),
		},
		Duitku: DuitkuConfig{
			Environment:   os.Getenv("DUIDKU_ENV"),
			APIKey:        os.Getenv("DUIDKU_APIKEY"),
			MerchantCode:  os.Getenv("DUIDKU_MERCHANTCODE"),
			PaymentMethod: getEnv("DUIDKU_PAYMENT_METHOD", "SP"),
			CallbackURL:   callbackURL,
			ReturnURL:     returnURL,
		},
		Midtrans: MidtransConfig{
			Environment: os.Getenv("MIDTRANS_ENV"),
			ServerKey:   os.Getenv("MIDTRANS_SERVER_KEY"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func parseAllowedOrigins(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{"*"}
	}

	origins := strings.Split(value, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
