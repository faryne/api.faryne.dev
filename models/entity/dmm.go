package entity

type DmmVideosList struct {
	Videos []DmmVideo `json:"videos"` // 影片列表
}

type DmmVideo struct {
	DmmVideoHeader
	DmmVideoBody
}

type DmmVideoHeader struct {
	No    string `json:"no"`    // 品番
	Title string `json:"title"` // 標題
	Url   string `json:"url"`   // 網址
	Thumb string `json:"thumb"` // 縮圖
}

type DmmVideoBody struct {
	VodDate     string          `json:"vod_date"`    // VOD 上架日期
	PublishDate string          `json:"pulish_date"` // 實體片上架日期
	Duration    int             `json:"duration"`    // 片長，單位：分鐘
	Directors   []string        `json:"directors"`   // 監督
	Series      []string        `json:"series"`      // 系列
	Makers      []string        `json:"makers"`      // 片商
	Labels      []string        `json:"labels"`      // 品牌
	Tags        []string        `json:"tags"`        // 標籤
	Actresses   []string        `json:"actresses"`   // 出演女優
	Images      []DmmVideoImage `json:"images"`      // 預覽圖
}

type DmmVideoImage struct {
	Thumb   string `json:"thumb"`   // 縮圖
	Preview string `json:"preview"` // 預覽圖
}

type DMMActress struct {
	Name         string   `json:"name",default:""` // 姓名
	Kana         string   `json:"kana",default:""` // 平假名
	Photo        string   `json:"photo"`           // 圖片
	Height       int      `json:"height"`          // 身高
	Bust         int      `json:"bust"`            // 胸圍
	Waist        int      `json:"waist"`           // 腰圍
	Hips         int      `json:"hips"`            // 臀圍
	Cup          string   `json:"cup"`             // 罩杯
	Horoscope    string   `json:"horoscope"`       // 星座
	Blood        string   `json:"blood"`           // 血型
	BornCity     string   `json:"born_city"`       // 出身地
	BirthYear    int      `json:"birth_year"`      // 出生年
	BirthMonth   int      `json:"birth_month"`     // 出生月
	BirthDay     int      `json:"birth_day"`       // 出生日
	FullBirthday string   `json:"full_birthday"`   // 完整生日
	Interests    []string `json:"interests"`       // 興趣
}
