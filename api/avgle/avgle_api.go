package avgle

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"strings"
)

func New() fiber.Handler {
	return proxy.Balancer(proxy.Config{
		Next: func(ctx *fiber.Ctx) bool {
			url := string(ctx.Request().URI().RequestURI())
			// 路徑不含 /avgle 則略過
			if !strings.Contains(url, "/avgle") {
				return true
			}
			return false
		},
		Servers: []string{
			"static-clst.avgle.com",
		},
		ModifyRequest: func(ctx *fiber.Ctx) error {
			url := string(ctx.Request().URI().RequestURI())
			if strings.Contains(url, "/avgle") {
				ctx.Request().URI().SetPath(strings.Replace(url, "avgle", "", -1))
			}
			ctx.Request().Header.Add("Host", "static-clst.avgle.com")
			ctx.Request().Header.Add("Referer", "https://avgle.com/video/CwgSmSuuTSk/%E7%B5%B6%E5%AF%BE%E3%81%AB%E3%83%8A%E3%83%9E%E3%81%A7%E9%80%A3%E5%B0%84%E3%81%95%E3%81%9B%E3%81%A6%E3%81%8F%E3%82%8C%E3%82%8B%E9%80%A3%E7%B6%9A%E4%B8%AD%E5%87%BA%E3%81%97j-%E5%AD%A6%E5%9C%92%E7%A5%AD%E3%82%BD%E3%83%BC%E3%83%97-sim-062-1")
			ctx.Request().Header.Add("Referrer-Policy", "Origin")
			return nil
		},
	})

}
