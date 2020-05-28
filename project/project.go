package project

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"github.com/refinedmods/sitegen/release"
	"github.com/refinedmods/sitegen/wiki"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		fileList := make([]string, 0)
		err := filepath.Walk(p.WikiPath, func(path string, f os.FileInfo, err error) error {
			fileList = append(fileList, path)
			return err
		})

		if err != nil {
			return err
		}

		wikiIndex := make(map[string]*wiki.Wiki)

		for _, file := range fileList {
			if strings.HasSuffix(file, ".md") {
				wikiPage := new(wiki.Wiki)
				wikiPage.Name = filepath.Base(strings.ReplaceAll(file, ".md", ""))
				wikiPage.Slug = slug.Make(wikiPage.Name)

				data, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				wikiPage.Body = string(data)

				metaFile := strings.ReplaceAll(file, ".md", ".json")
				if _, err := os.Stat(metaFile); err == nil {
					data, err := ioutil.ReadFile(metaFile)
					if err != nil {
						return err
					}

					err = json.Unmarshal(data, &wikiPage.Meta)
					if err != nil {
						return err
					}
				}

				wikiPage.Parse()

				p.Wikis = append(p.Wikis, wikiPage)

				wikiIndex[wikiPage.Name] = wikiPage

				if wikiPage.Meta.Group != "" {
					p.WikisByGroup[wikiPage.Meta.Group] = append(p.WikisByGroup[wikiPage.Meta.Group], wikiPage)
				}
			}
		}

		for _, wikiPage := range p.Wikis {
			err = wikiPage.ParseReferences(p.Slug, wikiIndex)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
