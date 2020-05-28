package wiki

import "github.com/russross/blackfriday/v2"

type Wiki struct {
	Name string
	Slug string
	Body string
	Meta WikiMetadata
}

type WikiMetadata struct {
	Icon string `json:"icon"`
}

func (w *Wiki) Parse() {
	w.Body = string(blackfriday.Run([]byte(w.Body)))
}
