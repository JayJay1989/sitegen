package release

type Release struct {
	Version   string      `json:"version"`
	Changelog string      `json:"changelog"`
	Downloads []*Download `json:"downloads"`
}

type Download struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}
