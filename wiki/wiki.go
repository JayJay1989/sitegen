package wiki

import "html/template"

type WikisByName map[string]*Wiki

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
	Icon string     `json:"icon"`
	Tags []*WikiTag `json:"tags"`
}

type WikiTag struct {
	Type    string `json:"type"`
	Release string `json:"release"`
}
