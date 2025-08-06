package configs

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type Configs struct {
	// APP
	AppName string
	AppVersion string
	AppEnv string
	AppDebug bool

	// Server
	ServerHost string
	ServerPort string

	// Database
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSLMode string

	// JWT
	JwtSecretKey string
	JwTExpiresIn time.Duration
	JWTRefreshExpiresIn time.Duration
}

var (
	cfg *Configs
	once sync.Once
)

func LoadConfig(isTest bool, envPath string) *Configs {
	once.Do(func ()  {
		// Load .env file
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		var (
			dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode string
		)

		jwtExpresIn, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
		if err != nil {
			log.Fatalf("Error parsing JWT_EXPIRES_IN: %v", err)
		}

		jwtRefreshExpresIn, err := time.ParseDuration(os.Getenv("JWT_REFRESH_EXPIRES_IN"))
		if err != nil {
			log.Fatalf("Error parsing JWT_REFRESH_EXPIRES_IN: %v", err)
		}


		if isTest {
			dbHost = os.Getenv("PG_HOST_TEST")
			dbPort = os.Getenv("PG_PORT_TEST")
			dbUser = os.Getenv("PG_USER_TEST")
			dbPass = os.Getenv("PG_PASS_TEST")
			dbName = os.Getenv("PG_DB_NAME_TEST")
			dbSSLMode = os.Getenv("PG_SSL_MODE_TEST")
		} else {
			dbHost = os.Getenv("PG_HOST")
			dbPort = os.Getenv("PG_PORT")
			dbUser = os.Getenv("PG_USER")
			dbPass = os.Getenv("PG_PASS")
			dbName = os.Getenv("PG_DB_NAME")
			dbSSLMode = os.Getenv("PG_SSL_MODE")
		
		}

		cfg = &Configs{
			AppName: os.Getenv("APP_NAME"),
			AppVersion: os.Getenv("APP_VERSION"),
			AppEnv: os.Getenv("APP_ENV"),
			AppDebug: os.Getenv("APP_DEBUG") == "true",

			ServerHost: os.Getenv("HTTP_HOST"),
			ServerPort: os.Getenv("HTTP_PORT"),

			DBHost: dbHost,
			DBPort: dbPort,
			DBUser: dbUser,
			DBPass: dbPass,
			DBName: dbName,
			DBSSLMode: dbSSLMode,

			JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
			JwTExpiresIn: jwtExpresIn,
			JWTRefreshExpiresIn: jwtRefreshExpresIn,
		}
	})

	return cfg
}

// GetDatabaseDSN returns database connection string
func (c *Configs) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}