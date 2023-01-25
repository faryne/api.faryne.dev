package tinami

import (
	"encoding/xml"
	"fmt"
	"github.com/faryne/api-server/config"
	"github.com/faryne/api-server/models/artwork"
	"github.com/faryne/api-server/service/nekomaid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type instance struct {
	ApiKey string `json:"api_key"`
	Base   string `json:"base"` // base endpoint, currently is https://www.tinami.com/api
}

// image 圖檔格式
type image struct {
	Url    string `xml:"url"`
	Width  int64  `xml:"width"`
	Height int64  `xml:"height"`
}

// Response Tinami 回傳的 xml 格式，只截取會用到的部分
type Response struct {
	Stat    string `xml:"stat,attr"` // rsp.stat
	Content struct {
		Id           string // 回應中的 xml 沒有 id，所以會在傳入時設定
		Type         string `xml:"type,attr"`         // rsp > content.type
		Issupport    string `xml:"issupport,attr"`    // rsp > content.issupport
		Iscollection string `xml:"iscollection,attr"` // rsp > content.iscollection
		Title        string `xml:"title"`             // rsp > content > title
		Creator      struct {
			Id        int64  `xml:"id,attr"`   // creator_id
			Name      string `xml:"name"`      // 暱稱
			Thumbnail string `xml:"thumbnail"` //
		} `xml:"creator"` // rsp > content > creator，創作者資訊
		AgeLevel    int64  `xml:"age_level"`   // 是否為限制級，1：沒踩線 2：有泳裝或內衣等灰色地帶 3：直接踩線
		Description string `xml:"description"` // 描述
		Image       image  `xml:"image"`       // 圖檔，content.type==illust 時使用
		Images      struct {
			Items []image `xml:"image"`
		} `xml:"images"` // 圖檔，在多圖時使用
		Tags struct {
			Items []string `xml:"tag"`
		} `xml:"tags"` // 標籤列表
	} `xml:"content"`
}

func New() nekomaid.RetrieverInterface {
	var i = instance{
		ApiKey: config.Config.Tinami.ApiKey,
		Base:   "https://www.tinami.com/api",
	}
	return &i
}

func (*instance) Login() error {
	// 由於已有相關的 api 可以取得圖檔，所以不實作 Login 方法
	return nil
}

func (i *instance) Get(id string) (*artwork.Artwork, error) {
	//
	var u = url.Values{}
	u.Add("api_key", i.ApiKey)
	u.Add("cont_id", id)
	u.Add("models", "1")
	resp, err := http.Get(i.Base + "/content/info?" + u.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var a Response
	if err := xml.Unmarshal(content, &a); err != nil {
		return nil, err
	}
	a.Content.Id = id
	return i.parseGetArtwork(&a)
}

func (i *instance) parseGetArtwork(response *Response) (*artwork.Artwork, error) {
	var o = &artwork.Artwork{}
	// 處理標籤
	o.Tags = response.Content.Tags.Items
	// 0：普遍級  1：可能會有泳裝等輕微暗示等  2：踩線（全裸或露點）
	o.IsR18 = false
	if response.Content.AgeLevel > 1 {
		o.IsR18 = true
	}
	o.Title = response.Content.Title
	o.ArtworkId = fmt.Sprintf("ti%s", response.Content.Id)
	o.Author = response.Content.Creator.Name
	o.AuthorId = strconv.Itoa(int(response.Content.Creator.Id))
	o.Site = string(nekomaid.Tinami)
	o.Status = ""
	o.IsAnimated = false
	o.PreviewUrl = fmt.Sprintf(nekomaid.PreviewUrlPattern, o.Site, o.AuthorId, o.ArtworkId)
	o.Description = response.Content.Description

	// 處理圖片
	o.Images = make([]artwork.Image, 0)
	// 直接拿網址取圖
	if response.Content.Type == "illust" {
		img, thumb, err := i.getImageUpload(o.ArtworkId, response.Content.Image.Url, 0)
		if err != nil {
			return o, err
		}
		o.Thumb = thumb
		o.Images = append(o.Images, img)
	} else {
		for k, v := range response.Content.Images.Items {
			img, thumb, err := i.getImageUpload(o.ArtworkId, v.Url, k)
			if k == 0 { // 只針對第一張產生縮圖
				o.Thumb = thumb
			}
			if err != nil {
				return o, err
			}
			o.Images = append(o.Images, img)
		}
	}

	return o, nil
}

func (i *instance) getImageUpload(id string, u string, idx int) (artwork.Image, string, error) {
	var o = artwork.Image{}
	var client = http.Client{}
	var thumb = ""
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return o, thumb, err
	}
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
