package main

import (
	"errors"
	"github.com/Netflix/go-env"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/faryne/crawler-server/config"
	_ "github.com/faryne/crawler-server/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// initialize env
	if _, err := os.Stat("./.env"); err == nil {
		godotenv.Load("./.env")
	}
	env.UnmarshalFromEnviron(&config.Config)

	app := fiber.New()
	app.Get("/parse", func(ctx *fiber.Ctx) error {
		return errors.New("It works")
	})
	app.Get("/*", swagger.HandlerDefault)

	app.Listen(":" + config.Config.Port)
}
