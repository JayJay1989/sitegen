package project

import (
	"github.com/gosimple/slug"
	"github.com/refinedmods/sitegen/release"
	"github.com/refinedmods/sitegen/wiki"
)

type Project struct {
	Name                  string                  `json:"name"`
	Slug                  string                  `json:"-"`
	ReleaseGroups         []*release.ReleaseGroup `json:"releaseGroups"`
	ReleaseGroupsReversed []*release.ReleaseGroup `json:"-"`
	Templates             map[string]string       `json:"templates"`
	LatestStableRelease   *release.Release
	WikiPath              string `json:"wikiPath"`
	Wikis                 []*wiki.Wiki
	WikiSidebars          []*wiki.Sidebar
	WikisByGroup          map[string][]*wiki.Wiki
}

func (p *Project) Load() error {
	p.WikisByGroup = make(map[string][]*wiki.Wiki)
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

	for _, group := range p.ReleaseGroupsReversed {
		if group.StableRelease != nil {
			p.LatestStableRelease = group.StableRelease
			break
		}
	}

	if p.WikiPath != "" {
		wikis, wikisByGroup, err := wiki.LoadWikis(p.WikiPath, p.Slug)
		if err != nil {
			return err
		}

		p.Wikis = wikis
		p.WikisByGroup = wikisByGroup
	}

	return nil
}
