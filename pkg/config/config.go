package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

type Config struct {
	AppName              string
	Port                 string
	AccessTokenSecretKey string
	AccessTokenDuratin   int

	Postgres
	Redis
}

type Postgres struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

type Redis struct {
	Host string
	Port string
	Pass string
	DB   int
}

func GetConfig() *Config {
	return &Config{
		AppName:              getStr("APP_NAME"),
		Port:                 getStr("PORT"),
		AccessTokenSecretKey: getStr("ACCESS_TOKEN_SECRET_KEY"),
		AccessTokenDuratin:   getInt("ACCESS_TOKEN_DURATION"),
		Postgres: Postgres{
			Host: getStr("POSTGRES_HOST"),
			Port: getStr("POSTGRES_PORT"),
			User: getStr("POSTGRES_USER"),
			Pass: getStr("POSTGRES_PASS"),
			DB:   getStr("POSTGRES_DB"),
		},
		Redis: Redis{
			Host: getStr("REDIS_HOST"),
			Port: getStr("REDIS_PORT"),
			Pass: getStr("REDIS_PASS"),
			DB:   getInt("REDIS_DB"),
		},
	}
}

func getStr(name string) string {
	return os.Getenv(name)
}

func getInt(name string) int {
	val, _ := strconv.Atoi(os.Getenv(name))
	return val
}
