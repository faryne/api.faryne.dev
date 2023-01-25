package nico

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/faryne/api-server/models/artwork"
	"github.com/faryne/api-server/service/nekomaid"
	"io/ioutil"
	"net/http"
	"net/url"
)

type instance struct {
	SimpleBaseEndpoint string `json:"simple_base_endpoint"`
	DetailBaseEndpoint string `json:"detail_base_endpoint"`
}

func New() nekomaid.RetrieverInterface {
	var i = instance{
		SimpleBaseEndpoint: "https://seiga.nicovideo.jp/api/illust/info",
		DetailBaseEndpoint: "https://sp.seiga.nicovideo.jp/ajax/seiga/%s",
	}
	return &i
}

// SimpleResponse 作品簡要資訊，年齡分級要從這邊取
// https://seiga.nicovideo.jp/api/illust/info?id=im11115799
type SimpleResponse struct {
	Image struct {
		Id           int64  `xml:"id"`
		UserId       int64  `xml:"user_id"`
		Title        string `xml:"title"`
		Description  string `xml:"description"`
		Summary      string `xml:"summary"`
		PublicStatus int64  `xml:"public_status"`
		AdultLevel   int64  `xml:"adult_level"`
	} `xml:"image"`
}

// DetailResponse 作品完整資訊，標籤還是作者資訊要從這邊取
// https://sp.seiga.nicovideo.jp/ajax/seiga/im6385426
type DetailResponse struct {
	TargetImage struct {
		Id       string `json:"id"`
		UserId   string `json:"user_id"`
		Nickname string `json:"nickname"`
		Width    string `json:"width"`
		Height   string `json:"height"`
		ImageUrl string `json:"image_url"`
		Tags     struct {
			Items []struct {
				Name string `json:"name"`
			} `json:"tag"`
		} `json:"tag_list"`
	} `json:"target_image"`
}

func (*instance) Login() error {
	// 由於已有相關的 api 可以取得圖檔，所以不實作 Login 方法
	return nil
}

func (i *instance) Get(id string) (*artwork.Artwork, error) {
	// 使用這兩個 url 取得作品資訊
	// https://seiga.nicovideo.jp/api/illust/info?id=im11115799
	// https://sp.seiga.nicovideo.jp/ajax/seiga/im6385426

	// 先取出摘要資訊
	var u = url.Values{}
	u.Add("id", id)
	resp1, err1 := http.Get(i.SimpleBaseEndpoint + "?" + u.Encode())
	if err1 != nil {
		return nil, err1
	}
	defer resp1.Body.Close()
	o1, err1 := ioutil.ReadAll(resp1.Body)
	if err1 != nil {
		return nil, err1
	}
	var output1 SimpleResponse
	if err := xml.Unmarshal(o1, &output1); err != nil {
		return nil, err
	}

	// 再取出詳細資料
	resp2, err2 := http.Get(fmt.Sprintf(i.DetailBaseEndpoint, id))
	if err2 != nil {
		return nil, err2
	}
	defer resp2.Body.Close()
	o2, err2 := ioutil.ReadAll(resp2.Body)
	if err2 != nil {
		return nil, err2
	}
	var output2 DetailResponse
	if err := json.Unmarshal(o2, &output2); err != nil {
		return nil, err
	}

	return i.parseGetArtwork(&output1, &output2)
}

func (i *instance) parseGetArtwork(simpleResponse *SimpleResponse, detailResponse *DetailResponse) (*artwork.Artwork, error) {
	var o = &artwork.Artwork{}

	// 處理 tags
	o.Tags = make([]string, len(detailResponse.TargetImage.Tags.Items))
	for k, v := range detailResponse.TargetImage.Tags.Items {
		o.Tags[k] = v.Name
	}
	// 處理分級
	o.IsR18 = false
	if simpleResponse.Image.AdultLevel > 1 {
		o.IsR18 = true
	}
	o.Title = simpleResponse.Image.Title
	o.ArtworkId = fmt.Sprintf("im%d", simpleResponse.Image.Id)
	o.Author = detailResponse.TargetImage.Nickname
	o.AuthorId = detailResponse.TargetImage.UserId
	o.Site = string(nekomaid.Nico)
	o.PreviewUrl = fmt.Sprintf(nekomaid.PreviewUrlPattern, o.Site, o.AuthorId, o.ArtworkId)
	o.IsAnimated = false
	o.Status = ""
	o.Description = simpleResponse.Image.Description

	// 處理圖片
	o.Images = make([]artwork.Image, 1)
	// 直接拿網址取圖
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, detailResponse.TargetImage.ImageUrl, nil)
	if err != nil {
		return o, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return o, err
	}
	defer resp.Body.Close()
	img, thumb, err := nekomaid.UploadImage(o.ArtworkId, resp, 0)
	if err != nil {
		return o, err
	}
	o.Thumb = thumb
	o.Images[0] = img

	return o, nil
}
