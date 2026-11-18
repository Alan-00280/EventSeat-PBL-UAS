package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName        string
	AppEnv         string
	AppHost        string
	AppPort        string
	AllowedOrigins []string
	SwaggerEnabled bool

	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxConns        int32
	DBMinConns        int32
	DBMaxConnLifetime time.Duration
	DBMaxConnIdleTime time.Duration

	MongoURI    string
	MongoDBName string

	RedisHost          string
	RedisPort          string
	RedisPassword      string
	RedisDB            int
	WorkerConcurrency  int
	OutboxPollInterval time.Duration
	OutboxBatchSize    int

	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	JWTIssuer     string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

func LoadEnv() (*Env, error) {
	_ = godotenv.Load()

	appEnv := strings.ToLower(getEnv("APP_ENV", "development"))
	if appEnv != "development" && appEnv != "production" && appEnv != "test" {
		return nil, fmt.Errorf("konfigurasi tidak valid: APP_ENV harus bernilai development, production, atau test (dapat: %q)", appEnv)
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if len(jwtSecret) < 32 {
		return nil, fmt.Errorf("konfigurasi tidak valid: JWT_SECRET wajib diisi minimal 32 karakter")
	}

	dbUser := getEnv("DB_USER", "")
	dbName := getEnv("DB_NAME", "")
	if dbUser == "" || dbName == "" {
		return nil, fmt.Errorf("konfigurasi tidak valid: DB_USER dan DB_NAME wajib diisi")
	}

	dbMaxConns := getEnvAsInt("DB_MAX_CONNS", 20)
	dbMinConns := getEnvAsInt("DB_MIN_CONNS", 5)
	if dbMaxConns < 1 || dbMinConns < 0 || dbMinConns > dbMaxConns {
		return nil, fmt.Errorf("konfigurasi tidak valid: DB_MIN_CONNS (%d) dan DB_MAX_CONNS (%d) tidak konsisten", dbMinConns, dbMaxConns)
	}

	mongoURI := getEnv("MONGO_URI", "")
	mongoDBName := getEnv("MONGO_DB_NAME", "")
	if mongoURI == "" || mongoDBName == "" {
		return nil, fmt.Errorf("konfigurasi tidak valid: MONGO_URI dan MONGO_DB_NAME wajib diisi")
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	if redisHost == "" || redisPort == "" {
		return nil, fmt.Errorf("konfigurasi tidak valid: REDIS_HOST dan REDIS_PORT wajib diisi")
	}

	originsRaw := getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	var origins []string
	for _, item := range strings.Split(originsRaw, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return nil, fmt.Errorf("konfigurasi tidak valid: ALLOWED_ORIGINS minimal berisi satu origin")
	}

	return &Env{
		AppName:        getEnv("APP_NAME", "EventSeat API"),
		AppEnv:         appEnv,
		AppHost:        getEnv("APP_HOST", "0.0.0.0"),
		AppPort:        getEnv("APP_PORT", "8080"),
		AllowedOrigins: origins,
		SwaggerEnabled: getEnvAsBool("SWAGGER_ENABLED", appEnv != "production"),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            dbUser,
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            dbName,
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxConns:        int32(dbMaxConns),
		DBMinConns:        int32(dbMinConns),
		DBMaxConnLifetime: getEnvAsDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute),
		DBMaxConnIdleTime: getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),

		MongoURI:    mongoURI,
		MongoDBName: mongoDBName,

		RedisHost:          redisHost,
		RedisPort:          redisPort,
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
		WorkerConcurrency:  getEnvAsInt("WORKER_CONCURRENCY", 10),
		OutboxPollInterval: getEnvAsDuration("OUTBOX_POLL_INTERVAL", 2*time.Second),
		OutboxBatchSize:    getEnvAsInt("OUTBOX_BATCH_SIZE", 50),

		JWTSecret:     jwtSecret,
		JWTAccessTTL:  getEnvAsDuration("JWT_ACCESS_TTL_MINUTES", 15*time.Minute),
		JWTRefreshTTL: getEnvAsDuration("JWT_REFRESH_TTL_DAYS", 7) * 24 * time.Hour,
		JWTIssuer:     getEnv("JWT_ISSUER", "uas_backend_prak"),

		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 1025),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@eventseat.campus.ac.id"),
	}, nil
}

func (e *Env) ServerAddr() string {
	return fmt.Sprintf("%s:%s", e.AppHost, e.AppPort)
}

func (e *Env) PostgresDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(e.DBUser, e.DBPassword),
		Host:   fmt.Sprintf("%s:%s", e.DBHost, e.DBPort),
		Path:   e.DBName,
	}
	q := u.Query()
	q.Set("sslmode", e.DBSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func (e *Env) RedisAddr() string {
	return fmt.Sprintf("%s:%s", e.RedisHost, e.RedisPort)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		if trimmed := strings.TrimSpace(val); trimmed != "" {
			return trimmed
		}
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return fallback
}
