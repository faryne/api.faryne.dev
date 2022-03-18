package telegraph

import (
	"errors"
	"github.com/faryne/api-server/models/apireq"
	"github.com/faryne/api-server/service/output"
	t "github.com/faryne/api-server/service/telegraph"
	"github.com/gofiber/fiber/v2"
)

// @Summary 取得 Telegraph 文章列表
// @Produce json
// @Accept json
// @Tags Telegraph
// @Param page query int false "頁碼"
// @Param limit query int false "每頁 N 筆"
// @Success 200 {object} output.StdOutput{code=apiresp.TelegraphPagesList} "OK"
// @Router /telegraph/news [GET]
func News(ctx *fiber.Ctx) error {
	var req = apireq.CommonRequest{
		Page:  1,
		Limit: 10,
	}
	account, err := t.New()
	if err != nil {
		return err
	}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.New("")
	}
	pages, err := account.GetPageList(req.Page, req.Limit)
	return output.New(200, "ok", pages)
}
