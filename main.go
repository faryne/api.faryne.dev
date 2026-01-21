package main

import (
	"github.com/Netflix/go-env"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/faryne/api-server/api/avgle"
	"github.com/faryne/api-server/api/dmm"
	"github.com/faryne/api-server/api/hanime_api"
	"github.com/faryne/api-server/api/nekomaid"
	"github.com/faryne/api-server/api/telegraph"
	"github.com/faryne/api-server/config"
	_ "github.com/faryne/api-server/docs"
	"github.com/faryne/api-server/service/output"
	_ "github.com/goccy/go-json"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/joho/godotenv"
	goapp "github.com/maxence-charriere/go-app/v9/pkg/app"
	"net/http"
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
	// nekomaid
	nekomaidGroup := app.Group("/nekomaid")
	nekomaidGroup.Post("/retrieve.json", nekomaid.Retrieve)
	// AVGle 縮圖
	app.Use(avgle.New())

	app.Use(adaptor.HTTPHandler(&goapp.Handler{}))

	app.Get("/*", swagger.HandlerDefault)

	//app.Listen(":8080")
	// 將 fiber app 轉換為 http.Handler 以便可以使用 GAE 上的資源
	http.ListenAndServe(":8080", adaptor.FiberApp(app))
}
