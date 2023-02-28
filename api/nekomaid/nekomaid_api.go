package nekomaid

import (
	"github.com/faryne/api-server/service/nekomaid"
	"github.com/faryne/api-server/service/nekomaid/nico"
	"github.com/faryne/api-server/service/nekomaid/pixiv"
	"github.com/faryne/api-server/service/nekomaid/tinami"
	"github.com/faryne/api-server/service/output"
	"github.com/gofiber/fiber/v2"
)

type retrieveRequest struct {
	Site      nekomaid.Site `form:"site" validate:"required,oneof=pixiv nico tinami"`
	ArtworkId string        `form:"artwork_id" validate:"required"`
}

// Retrieve 執行作品資訊解析與圖片抓取
// @Summary 執行作品資訊解析與圖片抓取
// @Produce json
// @Accept x-www-form-urlencoded
// @Tags Nekomaid
// @Param site formData string true "string enums" Enums(pixiv, nico, tinami) "網站類型，必須為 pixiv / nico / tinami 之一"
// @Param artwork_id formData string true "作品ID"
// @Success
// @Router /nekomaid/retrieve.json [POST]
func Retrieve(c *fiber.Ctx) error {
	var req retrieveRequest
	if err := c.BodyParser(&req); err != nil {
		return output.New(400, "", map[string]string{
			"error": err.Error(),
		})
	}
	var processor nekomaid.RetrieverInterface
	switch req.Site {
	case nekomaid.Pixiv:
		processor = pixiv.New()
	case nekomaid.Nico:
		processor = nico.New()
	case nekomaid.Tinami:
		processor = tinami.New()
	}
	artwork, err := processor.Get(req.ArtworkId)
	if err != nil {
		return output.New(500, "", map[string]string{
			"error": err.Error(),
		})
	}
	return output.New(200, "", artwork)
}
