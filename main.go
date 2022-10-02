package main

import (
	"github.com/Netflix/go-env"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/faryne/api-server/api/avgle"
	"github.com/faryne/api-server/api/dmm"
	"github.com/faryne/api-server/api/hanime_api"
	"github.com/faryne/api-server/api/telegraph"
	"github.com/faryne/api-server/config"
	_ "github.com/faryne/api-server/docs"
	"github.com/faryne/api-server/service/output"
	_ "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/joho/godotenv"
	"os"
	"time"
)

//go:generate swag init
func main() {
	// initialize env
	if _, err := os.Stat("./.env"); err == nil {
		godotenv.Load("./.env")
	}
	env.UnmarshalFromEnviron(&config.Config)

	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		UnescapePath:  true,
		ErrorHandler:  output.ErrorHandler,
	})
	// setting up middleware
	app.Use(etag.New())
	app.Use(cors.New())
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals("start_time", time.Now().UnixNano())
		return ctx.Next()
	})

	// 取得 telegraph news
	app.Get("/telegraph/news", telegraph.News)
	// dmm crawler
	app.Get("/dmm/crawler", dmm.Crawler)
	// hanime newest
	app.Get("/hanime/new.rss", hanime_api.NewUpload)
	// AVGle 縮圖
	app.Use(avgle.New())

	app.Get("/*", swagger.HandlerDefault)

	app.Listen(":8080")
}
