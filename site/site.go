package site

import (
	"encoding/json"
	"github.com/refinedmods/sitegen/project"
	"io/ioutil"
)

type Site struct {
	Projects        []*project.Project `json:"projects"`
	Name            string             `json:"name"`
	Templates       map[string]string  `json:"templates"`
	CopyDirectories []string           `json:"copyDirectories"`
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
	for _, proj := range c.Projects {
		err := proj.Load()
		if err != nil {
			return err
		}
	}

	return nil
}
