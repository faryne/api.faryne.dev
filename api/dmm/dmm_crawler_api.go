package dmm

import (
	"errors"
	"github.com/faryne/api-server/models/apireq"
	"github.com/faryne/api-server/service/dmm"
	"github.com/faryne/api-server/service/output"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Summary 使用 DMM 爬蟲爬出頁面指定資料
// @Produce json
// @Accept json
// @Tags DMM
// @Param url query string true "DMM目標頁面網址"
// @Param type query string true "string enums" Enums(video, actress) "爬取頁面類型"
// @Router /dmm/crawler [GET]
func Crawler(ctx *fiber.Ctx) error {
	var r apireq.DMMCrawlerRequest
	if err := ctx.QueryParser(&r); err != nil {
		return errors.New("")
	}
	v := validator.New()
	if err := v.Struct(r); err != nil {
		return err
	}
	resp, err := dmm.GetData(&r)
	if err != nil {
		return err
	}
	return output.New(200, "OK", resp)
}
