package telegraph

import (
	"github.com/faryne/api-server/service/output"
	t "github.com/faryne/api-server/service/telegraph"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// @Summary 取得 Telegraph 文章列表
// @Produce json
// @Accept json
// @Tags Telegraph
// @Router /telegraph/news [GET]
func News(ctx *fiber.Ctx) error {
	account, err := t.New()
	if err != nil {
		return err
	}
	var offset = 0
	var limit = 100
	if page := ctx.Query("page"); page != "" {
		if p, err := strconv.ParseInt(page, 10, 64); err == nil {
			offset = (int(p) - 1) * limit
		}
	}
	pages, err := account.GetPageList(offset, limit)
	return output.New(200, "ok", pages)
}
