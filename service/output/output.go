package output

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"math"
	"time"
)

// ResponseInterface 標準輸出介面
type ResponseInterface interface {
	error // 需實作 error.Error() 方法
	Code() int
	Output() interface{}
}

// StdOutput 標準輸出物件
type StdOutput struct {
	Code     int         `json:"code"`       // HTTP Code，例如 `200` / `400` 等
	CostTime float64     `json:"cost_time" ` // API 執行花費時間，單位為秒。e.g. `1.234`
	Data     interface{} `json:"data"`       // 回傳內容
	Ip       string      `json:"ip"`         // ip
	Uri      string      `json:"uri"`        // 呼叫 api 路徑，e.g. `/v1/api-token`
	Method   string      `json:"method"`     // HTTP 方法
	Message  string      `json:"message"`    // 執行完成訊息
}

// appOutput
type appOutput struct {
	code    int
	message string
	output  interface{}
}

// ErrorHandler 給 fiber 使用預設的 error handler
func ErrorHandler(c *fiber.Ctx, err error) error {
	var message = err.Error() // 預設從 error 取得輸出文字訊息
	var httpCode = 500        // 預設回傳的 httpcode

	// 計算執行時間
	var costTime = float64(-1)
	var endTime = float64(time.Now().UnixMilli())
	if startTime := c.Locals("start_time"); startTime != nil {
		costTime = (endTime - float64(startTime.(int64))) / float64(time.Microsecond)
	}

	// 取出 ip
	var ip = c.IP()
	if len(c.IPs()) > 0 {
		ip = c.IPs()[0]
	}

	if resp, ok := err.(ResponseInterface); ok {
		if content, ok := resp.Output().(io.Reader); ok {
			return c.Status(resp.Code()).SendStream(content, -1)
		}
		// 傳入 String 當成字串處理
		if content, ok := resp.Output().(string); ok {
			return c.Status(resp.Code()).SendString(content)
		}
		// 否則按照一般 JSON 輸出處理
		return c.Status(resp.Code()).JSON(StdOutput{
			Code:     resp.Code(),
			Ip:       ip,
			CostTime: math.Round(costTime*10000) / 10000,
			Data:     resp.Output(),
			Uri:      c.Request().URI().String(),
			Method:   c.Method(),
			Message:  err.Error(),
		})
	}

	return c.Status(httpCode).SendString(message)
}

// New 預設標準輸出方法
func New(code int, message string, output interface{}) ResponseInterface {
	var o = appOutput{
		code:    code,
		message: message,
		output:  output,
	}
	return &o
}

func (a *appOutput) Error() string {
	return a.message
}

func (a *appOutput) Code() int {
	return a.code
}

func (a *appOutput) Output() interface{} {
	return a.output
}
