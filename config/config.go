package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type DBConfig struct {
	DBDriver             string
	DBHost               string
	DBPort               string
	DBUser               string
	DBName               string
	DBPassword           string
	TSLConfig            string
	AllowNativePasswords bool
	MultiStatements      bool
	MaxOpenConns         uint
	MaxIdleConns         uint
	ConnMaxLifetime      uint
}

type RedisConfig struct {
	RedisHost       string
	RedisPort       string
	RedisPassword   string
	RedisDBNum      uint8
	RedisExpiration uint
}

type WebsocketConfig struct {
	WSURL               string
	WSApiToken          string
	WSReconnectInterval uint
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

// Shifter Services
type MailServiceConfig struct {
	MailServiceAdr string
}

type AuthServiceConfig struct {
	AuthServiceAdr string
}

type GeoApifyConfig struct {
	APIUrl string
	APIKey string
}

type GooglePlacesConfig struct {
	GoogleAPIKey       string
	GoogleSearchRadius uint
}

type SearchConfig struct {
	SearchRadius int
}

type AppConfig struct {
	DBConfig
	RedisConfig
	WebsocketConfig
	LoggingConfig
	SecurityConfig
	SSLConfig
	SMTPConfig
	MailServiceConfig
	AuthServiceConfig
	GeoApifyConfig
	GooglePlacesConfig
	SearchConfig
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
			DBDriver:             getEnv("DB_DRIVER", "mysql"),
			DBHost:               getEnv("DB_HOST", ""),
			DBPort:               getEnv("DB_PORT", "3306"),
			DBUser:               getEnv("DB_USER", ""),
			DBPassword:           getEnv("DB_PASSWORD", ""),
			DBName:               getEnv("DB_NAME", ""),
			TSLConfig:            getEnv("DB_TLS_CONFIG", "true"),
			AllowNativePasswords: getEnvAsBool("DB_ALLOW_NATIVE_PASSWORDS", true),
			MultiStatements:      getEnvAsBool("DB_MULTI_STATEMENTS", false),
			MaxOpenConns:         uint(getEnvAsInt("DB_MAX_OPEN_CONNS", 20)),
			MaxIdleConns:         uint(getEnvAsInt("DB_MAX_IDLE_CONNS", 5)),
			ConnMaxLifetime:      uint(getEnvAsInt("DB_CONN_MAX_LIFETIME", 5)),
		},
		RedisConfig: RedisConfig{
			RedisHost:       getEnv("REDIS_HOST", ""),
			RedisPort:       getEnv("REDIS_PORT", "6378"),
			RedisPassword:   getEnv("REDIS_PASSWORD", ""),
			RedisDBNum:      uint8(getEnvAsInt("REDIS_DB_NUM", 0)),
			RedisExpiration: uint(getEnvAsInt("REDIS_EXPIRATION", 0)),
		},
		SMTPConfig: SMTPConfig{
			SMTPHost:     getEnv("SMTP_HOST", ""),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		},
		WebsocketConfig: WebsocketConfig{
			WSURL:               getEnv("WS_URL", ""),
			WSApiToken:          getEnv("WS_API_TOKEN", ""),
			WSReconnectInterval: uint(getEnvAsInt("WS_RECONNECT_INTERVAL", 5)),
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
		MailServiceConfig: MailServiceConfig{
			MailServiceAdr: getEnv("MAIL_SERVICE_ADR", ""),
		},
		AuthServiceConfig: AuthServiceConfig{
			AuthServiceAdr: getEnv("AUTH_SERVICE_ADR", ""),
		},
		GeoApifyConfig: GeoApifyConfig{
			APIUrl: getEnv("GEOAPIFY_LINK", ""),
			APIKey: getEnv("GEAPIFY_API_TOKEN", ""),
		},
		GooglePlacesConfig: GooglePlacesConfig{
			GoogleAPIKey:       getEnv("GOOGLE_PLACES_API_TOKEN", ""),
			GoogleSearchRadius: uint(getEnvAsInt("GOOGLE_SEARCH_RADIUS", 50000)),
		},
		SearchConfig: SearchConfig{
			SearchRadius: getEnvAsInt("SEARCH_BIAS_RADIUS", 21000),
		},
		AppMode:     getEnv("APP_MODE", "devs"),
		AppLanguage: getEnv("APP_LANG", "en"),
		AppTimezone: getEnv("APP_TIMEZONE", "Asia/Jakarta"),
		AppPort:     getEnv("APP_PORT", ""),
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

// Helper to read an environment variable into a slice of a specific type or return default value.
func getEnvAsSlice[T any](name string, defaultVal []T, sep string) []T {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	vals := strings.Split(valStr, sep)
	result := make([]T, len(vals))

	for i, v := range vals {
		switch any(result).(type) {
		case []string:
			result[i] = any(v).(T)
		case []int:
			intVal, _ := strconv.Atoi(v)
			result[i] = any(intVal).(T)
		case []bool:
			boolVal, _ := strconv.ParseBool(v)
			result[i] = any(boolVal).(T)
		default:
			return defaultVal
		}
	}

	return result
}
