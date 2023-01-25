package nekomaid

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/faryne/api-server/models/artwork"
	"github.com/gofiber/fiber/v2/utils"
	"image"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Site string

const (
	Pixiv  Site = "pixiv"
	Nico   Site = "nico"
	Tinami Site = "tinami"
)

const Home = "https://nekomaid.web.app"

const PreviewUrlPattern = Home + "#/%s/%s/%s"

var domains = []string{
	"http://pcdn1.ha2.tw",
	"http://pcdn2.ha2.tw",
	"http://pcdn2.ha2.tw",
	"http://cdn-pixiv.maid.tw",
	"http://cdn-pixiv.maid.im",
}

// RetrieverInterface 作為介面控制需實作的項目
type RetrieverInterface interface {
	// Login 執行登入
	Login() error
	Get(id string) (*artwork.Artwork, error)
}

func UploadImage(artworkId string, reader *http.Response, idx int) (artwork.Image, string, error) {
	var o = artwork.Image{}
	var thumb = "" // 縮圖網址
	img, format, err := image.Decode(reader.Body)
	if err != nil {
		return o, thumb, err
	}
	// 將 image 物件轉換為 bytes，並且算出 hashId
	b := new(bytes.Buffer)
	if err := png.Encode(b, img); err != nil {
		return o, thumb, err
	}
	m := md5.New()
	if _, err := io.Copy(m, b); err != nil {
		return o, thumb, err
	}
	hashId := hex.EncodeToString(m.Sum(b.Bytes()))[0:5]
	// 解出副檔名
	o.Ext = format
	if strings.ToLower(format) == "jpeg" { // 碰到是 jpeg 時，副檔名改為 jpg
		o.Ext = "jpg"
	}
	o.Size = reader.ContentLength
	if o.Size <= 0 {
		// @TODO：tinami 的檔案長度需要額外處理
		o.Size = int64(b.Cap())
	}

	o.Height = int64(img.Bounds().Dy())
	o.Width = int64(img.Bounds().Dx())
	o.Mime = utils.GetMIME(format)
	o.Index = int64(idx)
	o.FileId = artworkId
	o.KeyId = hashId
	// 處理 Raw 的網址內容，避免重要資訊暴露
	imageUrl, _ := url.Parse(reader.Request.URL.String())
	values := imageUrl.Query()
	values.Del("api_key")
	imageUrl.RawQuery = values.Encode()
	o.Raw = imageUrl.String()
	o.Original = o.Raw

	// 處理縮圖
	if idx == 0 {
		thumbName := fmt.Sprintf("%s_%s_thumb.%s", artworkId, hashId, o.Ext)
		thumb = getDomain() + "/" + thumbName // 設定縮圖完整網址
		var width, height = 120, 0
		if img.Bounds().Dx() < img.Bounds().Dy() {
			width = 0
			height = 120
		}
		newImage := imaging.Resize(img, width, height, imaging.Lanczos)
		fp, err := os.Create("./" + thumbName)
		if err != nil {
			return o, thumb, err
		}
		// 呼叫 S3 upload
		if err := png.Encode(fp, newImage); err != nil {
			return o, thumb, err
		}
		defer fp.Close()
	}

	// 計算圖片真實路徑
	filenamePattern := "%s_%s.%s"
	var filename = fmt.Sprintf(filenamePattern, artworkId, o.KeyId, o.Ext)
	if idx > 0 {
		filename = fmt.Sprintf(filenamePattern, artworkId, o.KeyId+"_p"+strconv.Itoa(idx), o.Ext)
	}
	o.Filename = filename
	o.Url = getDomain() + "/" + filename
	// 呼叫 S3 upload
	fp, _ := os.Create("./" + filename)
	png.Encode(fp, img)
	defer fp.Close()

	return o, thumb, nil
}

func getDomain() string {
	return domains[rand.Intn(len(domains))]
}
