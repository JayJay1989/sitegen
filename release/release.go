package release

import "time"

type Release struct {
	Version   string            `json:"version"`
	Changelog string            `json:"changelog"`
	Type      string            `json:"type"`
	Slug      string            `json:"-"`
	Date      *time.Time        `json:"date"`
	Downloads map[string]string `json:"downloads"`
	Links     map[string]string `json:"links"`
}
