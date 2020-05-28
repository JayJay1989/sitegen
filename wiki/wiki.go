package wiki

type Wiki struct {
	Name string
	Slug string
	Body string
	Meta WikiMetadata
}

type Sidebar struct {
	File string `json:"file"`
	Body string
}

type WikiMetadata struct {
	Icon  string `json:"icon"`
	Group string `json:"group"`
}
