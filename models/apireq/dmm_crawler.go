package apireq

type DMMCrawlerRequest struct {
	Url  string `json:"url" query:"url" validate:"required,url"`
	Type string `json:"type" query:"type" validate:"required,oneof=video actress" enums:"video,actress"`
}
