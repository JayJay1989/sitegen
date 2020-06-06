package project

import (
	"errors"
	"github.com/gosimple/slug"
	"github.com/refinedmods/sitegen/release"
	"github.com/refinedmods/sitegen/wiki"
)

type Project struct {
	Name                  string                      `json:"name"`
	Slug                  string                      `json:"-"`
	ReleaseGroups         []*release.ReleaseGroup     `json:"releaseGroups"`
	ReleaseGroupsReversed []*release.ReleaseGroup     `json:"-"`
	ReleasesByVersion     map[string]*release.Release `json:"-"`
	Templates             map[string]string           `json:"templates"`
	LatestStableRelease   *release.Release
	WikiPath              string          `json:"wikiPath"`
	WikiSidebars          []*wiki.Sidebar `json:"wikiSidebars"`
	Wikis                 []*wiki.Wiki
	WikisByName           wiki.WikisByName
}

func (p *Project) Load() error {
	p.WikisByName = make(map[string]*wiki.Wiki)
	p.ReleasesByVersion = make(map[string]*release.Release)
	p.Slug = slug.Make(p.Name)

	for _, group := range p.ReleaseGroups {
		err := group.Load()
		if err != nil {
			return err
		}

		for _, release := range group.Releases {
			p.ReleasesByVersion[release.Version] = release
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
		wikis, byName, err := wiki.LoadWikis(p.WikiPath, p.Slug, p.WikiSidebars)
		if err != nil {
			return err
		}

		p.Wikis = wikis
		p.WikisByName = byName

		for _, w := range p.Wikis {
			for _, tag := range w.Meta.Tags {
				if p.ReleasesByVersion[tag.Release] == nil {
					return errors.New("Version " + tag.Release + " not found on wiki " + w.Name)
				}
			}
		}
	}

	return nil
}
