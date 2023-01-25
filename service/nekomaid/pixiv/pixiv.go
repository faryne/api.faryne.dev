package pixiv

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/faryne/api-server/config"
	"github.com/faryne/api-server/models/artwork"
	"github.com/faryne/api-server/pkg/storage/memcached"
	"github.com/faryne/api-server/service/nekomaid"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type instance struct {
	Token *OAuthAPIResponse
}

const (
	TokenName    = "pixiv_cookie"
	LoginUrl     = "https://oauth.secure.pixiv.net/auth/token"
	ApiUrl       = "https://app-api.pixiv.net/v1"
	RefererUrl   = "https://www.pixiv.net/artworks/%s"
	ClientId     = "MOBrBDS8blbauoSck0ZfDbtuzpyT"
	ClientSecret = "lsACyCD94FhDUtGTXi3QzcFE2uU1hqtDaKeqrdwj"
	HashSecret   = "28c1fdd170a5204386cb1313c7077b34f83e4aaf4aa829ce78c231e05b0bae2c"
)

type OAuthAPIResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type Response struct {
	Illust struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Type     string `json:"type"`
		Caption  string `json:"caption"`
		Restrict int64  `json:"restrict"`
		User     struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Tags []struct {
			Name string `json:"name"`
		} `json:"tags"`
		MetaSinglePage struct {
			OriginalImageUrl string `json:"original_image_url"`
		} `json:"meta_single_page"` // 單圖時使用
		MetaPages []struct {
			ImageUrls struct {
				Original     string `json:"original"`
				Large        string `json:"large"`
				Medium       string `json:"medium"`
				SquareMedium string `json:"square_medium"`
			} `json:"image_urls"` // 多圖時使用
		} `json:"meta_pages"`
		PageCount       int64 `json:"page_count"`
		IllustAiType    int64 `json:"illust_ai_type"`
		IllustBookStyle int64 `json:"illust_book_style"`
	} `json:"illust"`
}

func New() nekomaid.RetrieverInterface {
	var i = instance{}
	return &i
}

func (i *instance) Login() error {
	// 初始化 memcache：如果取不到 token 的話準備要求新的 token
	var keyToken = "token_pixiv"
	var keyExpire = time.Second * 3500 // 放 memcache 3500 秒
	m := memcached.New(memcached.Environment(config.Config.Environmet))
	tokenContent, err := m.Get(keyToken)
	if err != nil {
		return err
	}
	// 把內容從 json string 轉換為 struct
	if tokenContent != nil || len(tokenContent) > 0 {
		var token OAuthAPIResponse
		json.Unmarshal(tokenContent, &token)
		i.Token = &token
		return nil
	}
	c := http.Client{}
	var u = url.Values{}
	u.Add("client_id", ClientId)
	u.Add("client_secret", ClientSecret)
	u.Add("get_secure_url", "1")
	u.Add("username", config.Config.Pixiv.Username)
	u.Add("password", config.Config.Pixiv.Password)
	u.Add("grant_type", "refresh_token")

	req, err := http.NewRequest(http.MethodPost, LoginUrl, strings.NewReader(u.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("app-os", "ios")
	req.Header.Set("user-agent", "PixivIOSApp/7.13.3 (iOS 14.6; iPhone13,2)")
	req.Header.Set("app-os-version", "14.6")

	dt := time.Now().UTC().Format("2006-01-02T15:04:05+00:00")
	req.Header.Set("x-client-time", dt)
	str := dt + HashSecret
	has := md5.Sum([]byte(str))
	md5str1 := fmt.Sprintf("%x", has)

	req.Header.Set("x-client-hash", md5str1)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http code should be 200, but got : %d \n", resp.StatusCode)
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var token OAuthAPIResponse
	json.Unmarshal(output, &token)
	defer resp.Body.Close()

	// 將 token 內容存在 memcached 中
	i.Token = &token
	m.Set(keyToken, output, keyExpire)
	return nil
}

func (i *instance) Get(id string) (*artwork.Artwork, error) {
	if err := i.Login(); err != nil {
		return nil, err
	}

	// 設定 http requqest
	var client = http.Client{}
	req, err := http.NewRequest(http.MethodGet, ApiUrl+"/illust/detail", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("illust_id", id)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+i.Token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	output, _ := ioutil.ReadAll(resp.Body)

	// 開始解析 json
	var o Response
	if err := json.Unmarshal(output, &o); err != nil {
		return nil, err
	}
	// 開始賦值及抓圖處理
	return i.parseGetArtwork(&o)
}

func (i *instance) parseGetArtwork(input *Response) (*artwork.Artwork, error) {
	var o = &artwork.Artwork{}

	// 處理標籤順便處理分級
	o.Tags = make([]string, len(input.Illust.Tags))
	o.IsR18 = false
	for k, v := range input.Illust.Tags {
		o.Tags[k] = v.Name
		if strings.Contains(v.Name, "R-18") {
			o.IsR18 = true
		}
	}
	// 標題與作者資訊
	o.Title = input.Illust.Title
	o.ArtworkId = strconv.Itoa(int(input.Illust.Id))
	o.Author = input.Illust.User.Name
	o.AuthorId = strconv.Itoa(int(input.Illust.User.Id))
	o.Site = string(nekomaid.Pixiv)
	o.IsAnimated = false
	o.PreviewUrl = fmt.Sprintf(nekomaid.PreviewUrlPattern, o.Site, o.AuthorId, o.ArtworkId)
	o.Status = ""
	o.Description = input.Illust.Caption

	// referer
	var referer = fmt.Sprintf(RefererUrl, o.ArtworkId)
	// 處理圖片，以及 o.PreviewUrl/ o.Thumb
	o.Images = make([]artwork.Image, 0)
	if input.Illust.PageCount > 1 {
		for k, v := range input.Illust.MetaPages {
			img, thumb, err := i.getImageUpload(o.ArtworkId, v.ImageUrls.Original, k, referer)
			if err != nil {
				return o, err
			}
			o.Images = append(o.Images, img)
			if k == 0 {
				o.Thumb = thumb
			}
		}
	} else {
		img, thumb, err := i.getImageUpload(o.ArtworkId, input.Illust.MetaSinglePage.OriginalImageUrl, 0, referer)
		if err != nil {
			return o, err
		}
		o.Thumb = thumb
		o.Images = append(o.Images, img)
	}

	return o, nil
}

func (i *instance) getImageUpload(id string, u string, idx int, referer string) (artwork.Image, string, error) {
	var o = artwork.Image{}
	var client = http.Client{}
	var thumb = ""
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return o, thumb, err
	}
	req.Header.Add("Referer", referer)
	resp, err := client.Do(req)
	if err != nil {
		return o, thumb, err
	}
	defer resp.Body.Close()
	o, thumb, err = nekomaid.UploadImage(id, resp, idx)
	if err != nil {
		return o, thumb, err
	}
	return o, thumb, nil
}
