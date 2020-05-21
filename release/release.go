package release

type Release struct {
	Version   string            `json:"version"`
	Changelog string            `json:"changelog"`
	Type      string            `json:"type"`
	Downloads map[string]string `json:"downloads"`
	Links     map[string]string `json:"links"`
}
