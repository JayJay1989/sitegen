package project

import "github.com/refinedmods/sitegen/release"

type Project struct {
	Name          string                  `json:"name"`
	Slug          string                  `json:"slug"`
	ReleaseGroups []*release.ReleaseGroup `json:"releaseGroups"`
	Templates     map[string]string       `json:"templates"`
}

func (p *Project) Load() error {
	for _, group := range p.ReleaseGroups {
		err := group.Load()
		if err != nil {
			return err
		}
	}

	return nil
}
