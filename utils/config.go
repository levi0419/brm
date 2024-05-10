package utils

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/brm/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {

	
	
	DBDriver             string
	DBSource             string

	ServerAddress 		  string			
	
	TokenSymmetricKey    string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	
	
}

func NewConfig() *Config {
	var envFile string
	flag.StringVar(&envFile, "env", "", "Env Variable File")
	flag.Parse()
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			logger.Error("error: ", zap.Error(err))
		}
	}
	return &Config{}
}

// LoadEnvConfig reads configuration from environment variables.
func (c *Config) LoadEnvConfig() error {

	atd := os.Getenv("ACCESS_TOKEN_DURATION")
	rtd := os.Getenv("REFRESH_TOKEN_DURATION")

	c.AccessTokenDuration, _ = time.ParseDuration(atd)
	c.RefreshTokenDuration, _ = time.ParseDuration(rtd)

	
	c.DBDriver = os.Getenv("DB_DRIVER")
	c.DBSource = os.Getenv("DB_SOURCE")
	c.ServerAddress = os.Getenv("SERVER_ADDRESS")
	
	c.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")

	
	return nil
}

func (c *Config) SanityCheck() {
	envProps := []string{
		"ACCESS_TOKEN_DURATION",
		"REFRESH_TOKEN_DURATION",

		"DB_DRIVER",
		"DB_SOURCE",

		"SERVER_ADDRESS",
		

		"TOKEN_SYMMETRIC_KEY",
		
		
	}

	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
