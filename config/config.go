package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
}

type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDBNum    uint8
}

type WebsocketConfig struct {
	WSURL      string
	WSApiToken string
}

type LoggingConfig struct {
	ServerLogPath string
	ErrLogPath    string
	LogMaxSize    int
	LogMaxBackup  int
	LogMaxAge     int
	LogCompress   bool
}

type SecurityConfig struct {
	AESKey string
	JWTKey string

	// JWT life span in hour(s), default to 1 hour
	JWTLifeSpan uint32

	// Password minimal length, default to 8
	PasswordMinLength uint32
}

type SSLConfig struct {
	KeyPath  string
	CertPath string
}

type AppConfig struct {
	DBConfig
	RedisConfig
	WebsocketConfig
	LoggingConfig
	SecurityConfig
	SSLConfig
	AppMode     string
	AppLanguage string
	AppTimezone string
	AppPort     string
	AppHost     string
	AppRoot     string
	DebugMode   bool
}

var config *AppConfig

func New(envPath string) *AppConfig {

	if err := godotenv.Load(envPath); err != nil {
		log.Println("Failed to locate .env file, program will proceed with provided env if any is provided")
	}

	config = &AppConfig{
		DBConfig: DBConfig{
			DBDriver:   getEnv("DB_DRIVER", "mysql"),
			DBHost:     getEnv("DB_HOST", ""),
			DBPort:     getEnv("DB_PORT", "3306"),
			DBUser:     getEnv("DB_USER", ""),
			DBPassword: getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", ""),
		},
		RedisConfig: RedisConfig{
			RedisHost:     getEnv("REDIS_HOST", ""),
			RedisPort:     getEnv("REDIS_PORT", "6378"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
			RedisDBNum:    uint8(getEnvAsInt("REDIS_DB_NUM", 0)),
		},
		WebsocketConfig: WebsocketConfig{
			WSURL:      getEnv("WS_URL", ""),
			WSApiToken: getEnv("WS_API_TOKEN", ""),
		},
		LoggingConfig: LoggingConfig{
			ServerLogPath: getEnv("LOG_PATH", "./log/server.log"),
			ErrLogPath:    getEnv("ERR_LOG_PATH", "./log/error.log"),
			LogMaxSize:    getEnvAsInt("LOG_MAX_SIZE", 30),
			LogMaxBackup:  getEnvAsInt("LOG_MAX_BACKUP", 5),
			LogMaxAge:     getEnvAsInt("LOG_MAX_AGE", 30),
			LogCompress:   getEnvAsBool("LOG_COMPRESS", true),
		},
		SecurityConfig: SecurityConfig{
			AESKey:            getEnv("AES_KEY", ""),
			JWTKey:            getEnv("JWT_KEY", ""),
			JWTLifeSpan:       uint32(getEnvAsInt("JWT_LIFE_SPAN", 1)),
			PasswordMinLength: uint32(getEnvAsInt("PASSWORD_MIN_LENGTH", 8)),
		},
		SSLConfig: SSLConfig{
			KeyPath:  getEnv("KEY_PATH", ""),
			CertPath: getEnv("CERT_PATH", ""),
		},
		AppMode:     getEnv("APP_MODE", "devs"),
		AppLanguage: getEnv("APP_LANG", "en"),
		AppTimezone: getEnv("APP_TIMEZONE", "Asia/Jakarta"),
		AppPort:     getEnv("APP_PORT", "47000"),
		AppHost:     getEnv("APP_HOST", ""),
		AppRoot:     getEnv("APP_ROOT", "/api/v1"),
		DebugMode:   getEnvAsBool("DEBUG", false),
	}

	return config
}

func GetConfig() *AppConfig {
	return config
}

// Simple helper function to read an environment or return a default value.
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value.
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
