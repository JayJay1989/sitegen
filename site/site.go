package site

import (
	"encoding/json"
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/wiki"
	"io/ioutil"
)

type Site struct {
	Projects        []*project.Project `json:"projects"`
	Name            string             `json:"name"`
	Templates       map[string]string  `json:"templates"`
	CopyDirectories []string           `json:"copyDirectories"`
	Url             string             `json:"url"`
}

func NewSite(filename string) (*Site, error) {
	site := new(Site)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, site)
	if err != nil {
		return nil, err
	}

	err = site.load()
	if err != nil {
		return nil, err
	}

	return site, nil
}

func (c *Site) load() error {
	projectNameToProjectSlug := make(map[string]string)
	projectNameToWikiIndex := make(map[string]wiki.WikisByName)

	for _, proj := range c.Projects {
		proj.Init()

		projectNameToProjectSlug[proj.Name] = proj.Slug
		projectNameToWikiIndex[proj.Name] = make(wiki.WikisByName)

		err := proj.Load(projectNameToProjectSlug, projectNameToWikiIndex)
		if err != nil {
			return err
		}
	}

	for _, proj := range c.Projects {
		err := proj.PostLoad(projectNameToProjectSlug, projectNameToWikiIndex)
		if err != nil {
			return err
		}
	}

	return nil
}
