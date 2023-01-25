package artwork

/**
{
"status":"ok, but not indexed.",
"from":"pixiv",
"is_r18":1,
"artwork_id":"92670929",
"title":"\u30d3\u30fc\u30c1\u3067\u3042\u305d\u307c06",
"author":"kanabun\u308f\u304b\u308b\u30c6\u30a3\u30c3\u30b7\u30e5",
"author_id":45722,"photos":[{"url":"http:\/\/pcdn2.ha2.tw\/92670929_6992f.jpg","original":"https:\/\/i.pximg.net\/img-original\/img\/2021\/09\/11\/20\/23\/23\/92670929_p0.jpg","width":812,"height":1416,"mime":"image\/jpeg","ext":"jpg","raw":"https:\/\/i.pximg.net\/img-original\/img\/2021\/09\/11\/20\/23\/23\/92670929_p0.jpg","size":897508,"file_id":"92670929","key_id":"6992f","filename":"92670929_6992f.jpg","index":0},{"url":"http:\/\/cdn-pixiv.maid.tw\/92670929_801db_p1.jpg","original":"https:\/\/i.pximg.net\/img-original\/img\/2021\/09\/11\/20\/23\/23\/92670929_p1.jpg","width":1003,"height":1416,"mime":"image\/jpeg","ext":"jpg","raw":"https:\/\/i.pximg.net\/img-original\/img\/2021\/09\/11\/20\/23\/23\/92670929_p1.jpg","size":1323099,"file_id":"92670929","key_id":"801db","filename":"92670929_801db_p1.jpg","index":1}],"thumb":"http:\/\/pcdn2.ha2.tw\/thumb\/92670929_6992f_thumb.jpg","tags":["R-18","\u30aa\u30ea\u30b8\u30ca\u30eb","\u65e5\u713c\u3051\u8de1","\u9a0e\u4e57\u4f4d","\u9670\u6bdb","\u4e2d\u51fa\u3057"],"is_animated":0,"preview_url":"http:\/\/neko.maid.tw\/pixiv\/45722\/92670929"}
*/

type Image struct {
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	Mime     string `json:"mime"`
	Ext      string `json:"ext"`
	Raw      string `json:"raw"`
	Size     int64  `json:"size"`
	Filename string `json:"filename"`
	Index    int64  `json:"index"`
	Url      string `json:"url"`
	Original string `json:"original"`
	FileId   string `json:"file_id"`
	KeyId    string `json:"key_id"`
}

type Artwork struct {
	Site        string   `json:"from"`
	Status      string   `json:"status"`
	AuthorId    string   `json:"author_id"`
	ArtworkId   string   `json:"artwork_id"`
	IsR18       bool     `json:"is_r18"`
	Title       string   `json:"title"`
	Author      string   `json:"author"`
	Images      []Image  `json:"photos"`
	Tags        []string `json:"tags"`
	Thumb       string   `json:"thumb"`
	IsAnimated  bool     `json:"is_animated"`
	PreviewUrl  string   `json:"preview_url"`
	Description string   `json:"description"`
}
