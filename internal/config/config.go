package config

import (
	"os"
	"time"
)

type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Solana     SolanaConfig
	DexScreens DexScreensConfig
	JWT        JWTConfig
	ApplePay   ApplePayConfig
}

type AppConfig struct {
	Env    string
	Port   string
	Secret string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      time.Duration
}

type SolanaConfig struct {
	RpcURL       string
	WebsocketURL string
	Network      string
}

type DexScreensConfig struct {
	ApiURL string
	ApiKey string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type ApplePayConfig struct {
	MerchantID      string `env:"APPLE_PAY_MERCHANT_ID,required"`
	CertificatePath string `env:"APPLE_PAY_CERT_PATH,required"`
	PrivateKeyPath  string `env:"APPLE_PAY_KEY_PATH,required"`
	DomainName      string `env:"APPLE_PAY_DOMAIN,required"`
}

func Load() (*Config, error) {
	return &Config{
		App: AppConfig{
			Env:    getEnv("APP_ENV", "development"),
			Port:   getEnv("APP_PORT", "8080"),
			Secret: getEnv("APP_SECRET", ""),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "memecoin_db"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
			TTL:      300 * time.Second,
		},
		Solana: SolanaConfig{
			RpcURL:       getEnv("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com"),
			WebsocketURL: getEnv("SOLANA_WEBSOCKET_URL", "wss://api.mainnet-beta.solana.com"),
			Network:      getEnv("SOLANA_NETWORK", "mainnet"),
		},
		DexScreens: DexScreensConfig{
			ApiURL: getEnv("DEXSCREENS_API_URL", "https://api.dexscreens.io/v1"),
			ApiKey: getEnv("DEXSCREENS_API_KEY", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: 24 * time.Hour,
		},
		ApplePay: ApplePayConfig{
			MerchantID:      getEnv("APPLE_PAY_MERCHANT_ID", ""),
			CertificatePath: getEnv("APPLE_PAY_CERT_PATH", ""),
			PrivateKeyPath:  getEnv("APPLE_PAY_KEY_PATH", ""),
			DomainName:      getEnv("APPLE_PAY_DOMAIN", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
