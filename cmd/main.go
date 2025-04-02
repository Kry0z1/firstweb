package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Kry0z1/firstweb/internal/database"
	"github.com/Kry0z1/firstweb/internal/digitalclock"
	"github.com/Kry0z1/firstweb/internal/middleware/auth"
	"github.com/Kry0z1/firstweb/internal/urlshortener"
)

type Config struct {
	Port               int `json:"port"`
	AccessTokenExpires int `json:"access_token_expires"`
}

var (
	Cfg                  Config
	secretKey, algorithm string
)

func parseConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Failed to read config: %w", err)
	}

	err = json.NewDecoder(file).Decode(&Cfg)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal config: %w", err)
	}

	return nil
}

func init() {
	godotenv.Load()

	secretKey = os.Getenv("SECRET_KEY")
	algorithm = os.Getenv("JWT_ALGO")

	err := parseConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	database.Init()
}

func main() {
	tokenizer, err := auth.NewJWTTokenizer(secretKey, algorithm, time.Duration(Cfg.AccessTokenExpires*int(time.Minute)))
	if err != nil {
		log.Fatal("Invalid secret key: %w", err)
	}

	r := gin.Default()
	r.Use(auth.CheckAuth(tokenizer))

	shortenHandler, redirectHandler := urlshortener.GetHandlers(urlshortener.GetRAMStorage())

	r.POST("/shorten", shortenHandler)
	r.GET("/redirect/:key", redirectHandler)

	r.GET("/digitalclock", digitalclock.Handler)

	r.POST("/login", auth.LoginForToken(tokenizer))

	r.POST("/user/create", database.CreateUserHandler)

	var _ database.User

	r.Run(fmt.Sprintf(":%d", Cfg.Port))
}
