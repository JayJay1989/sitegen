package wiki

import "html/template"

type Wiki struct {
	Name string
	Slug string
	Body string
	Meta WikiMetadata
}

type Sidebar struct {
	File     string `json:"file"`
	Name     string `json:"name"`
	Body     string
	BodyHtml template.HTML
}

type WikiMetadata struct {
	Icon  string `json:"icon"`
	Group string `json:"group"`
}
