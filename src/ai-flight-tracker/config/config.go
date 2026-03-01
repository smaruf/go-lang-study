package config

import (
	"os"
	"strconv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port               int
	Debug              bool
	SecretKey          string
	CORSOrigins        string
	OpenSkyURL         string
	OpenSkyUsername    string
	OpenSkyPassword    string
	OpenSkyTimeout     int
	OllamaURL          string
	OllamaModel        string
	OllamaTimeout      int
	CacheTTL           int
	RateLimitAsk       int // requests per minute
	RateLimitFlights   int // requests per minute
	LivePricingEnabled bool
	PricingAPIKey      string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:               getEnvInt("PORT", 8080),
		Debug:              getEnvBool("DEBUG", false),
		SecretKey:          getEnvStr("SECRET_KEY", "change-me-in-production"),
		CORSOrigins:        getEnvStr("CORS_ORIGINS", "*"),
		OpenSkyURL:         getEnvStr("OPENSKY_URL", "https://opensky-network.org/api"),
		OpenSkyUsername:    getEnvStr("OPENSKY_USERNAME", ""),
		OpenSkyPassword:    getEnvStr("OPENSKY_PASSWORD", ""),
		OpenSkyTimeout:     getEnvInt("OPENSKY_TIMEOUT", 10),
		OllamaURL:          getEnvStr("OLLAMA_URL", "http://localhost:11434/api/generate"),
		OllamaModel:        getEnvStr("OLLAMA_MODEL", "llama2"),
		OllamaTimeout:      getEnvInt("OLLAMA_TIMEOUT", 30),
		CacheTTL:           getEnvInt("CACHE_TTL", 300),
		RateLimitAsk:       getEnvInt("RATE_LIMIT_ASK", 10),
		RateLimitFlights:   getEnvInt("RATE_LIMIT_FLIGHTS", 30),
		LivePricingEnabled: getEnvBool("LIVE_PRICING_ENABLED", false),
		PricingAPIKey:      getEnvStr("PRICING_API_KEY", ""),
	}
}

func getEnvStr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	return v == "true" || v == "1" || v == "yes"
}
