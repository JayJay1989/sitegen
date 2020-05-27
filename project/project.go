package project

import (
	"github.com/gosimple/slug"
	"github.com/refinedmods/sitegen/release"
)

type Project struct {
	Name                  string                  `json:"name"`
	Slug                  string                  `json:"-"`
	ReleaseGroups         []*release.ReleaseGroup `json:"releaseGroups"`
	ReleaseGroupsReversed []*release.ReleaseGroup `json:"-"`
	Templates             map[string]string       `json:"templates"`
}

func (p *Project) Load() error {
	p.Slug = slug.Make(p.Name)

	for _, group := range p.ReleaseGroups {
		err := group.Load()
		if err != nil {
			return err
		}
	}

	for i := len(p.ReleaseGroups) - 1; i >= 0; i-- {
		p.ReleaseGroupsReversed = append(p.ReleaseGroupsReversed, p.ReleaseGroups[i])
	}

	return nil
}
