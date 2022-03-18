package apiresp

type (
	TelegraphPage struct {
		Path        string        `json:"path"`
		URL         string        `json:"url"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		AuthorName  string        `json:"author_name,omitempty"`
		AuthorURL   string        `json:"author_url,omitempty"`
		ImageURL    string        `json:"image_url,omitempty"`
		Content     []interface{} `json:"content,omitempty"`
		Views       int           `json:"views"`
		CanEdit     bool          `json:"can_edit,omitempty"`
	}

	TelegraphPagesList struct {
		TotalCount int             `json:"total_count"`
		Pages      []TelegraphPage `json:"pages"`
	}
)
