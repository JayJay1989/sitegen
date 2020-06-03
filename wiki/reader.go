package wiki

import (
	"encoding/json"
	"github.com/gosimple/slug"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadWikis(path, projectSlug string, sidebars []*Sidebar) ([]*Wiki, map[string][]*Wiki, error) {
	var result []*Wiki
	byGroup := make(map[string][]*Wiki)

	fileList := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return err
	})

	if err != nil {
		return nil, nil, err
	}

	index := make(wikiIndex)

	for _, file := range fileList {
		if strings.HasSuffix(file, ".md") {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				return nil, nil, err
			}

			page := new(Wiki)
			page.Name = filepath.Base(strings.ReplaceAll(file, ".md", ""))
			page.Slug = slug.Make(page.Name)
			page.Body = parseBody(string(data))

			metaFile := strings.ReplaceAll(file, ".md", ".json")
			if _, err := os.Stat(metaFile); err == nil {
				data, err := ioutil.ReadFile(metaFile)
				if err != nil {
					return nil, nil, err
				}

				err = json.Unmarshal(data, &page.Meta)
				if err != nil {
					return nil, nil, err
				}
			}

			result = append(result, page)
			index[page.Name] = page

			if page.Meta.Group != "" {
				byGroup[page.Meta.Group] = append(byGroup[page.Meta.Group], page)
			}
		}
	}

	for _, sidebar := range sidebars {
		data, err := ioutil.ReadFile(sidebar.File)
		if err != nil {
			return nil, nil, err
		}

		sidebar.Body = parseBody(string(data))
	}

	for _, page := range result {
		result, err := parseReferenceLinks(page.Body, page.Name, projectSlug, index)
		if err != nil {
			return nil, nil, err
		}

		page.Body = result
	}

	for _, sidebar := range sidebars {
		result, err := parseReferenceLinks(sidebar.Body, sidebar.Name, projectSlug, index)
		if err != nil {
			return nil, nil, err
		}

		sidebar.BodyHtml = template.HTML(result)
	}

	return result, byGroup, nil
}