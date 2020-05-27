package render

import (
	"encoding/json"
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
	"strings"
)

func (r *Renderer) AddUpdate(site *site.Site, project *project.Project) error {
	data := make(map[string]interface{})

	data["website"] = site.Url + "/" + project.Slug

	promos := make(map[string]string)
	for _, group := range project.ReleaseGroups {
		if strings.HasPrefix(group.Name, "Minecraft ") {
			mcVersion := strings.ReplaceAll(group.Name, "Minecraft ", "")

			promos[mcVersion+"-latest"] = group.LatestRelease.Version
			promos[mcVersion+"-recommended"] = group.StableRelease.Version

			releases := make(map[string]string)
			for _, release := range group.Releases {
				releases[release.Version] = release.Changelog
			}

			data[mcVersion] = releases
		}
	}
	data["promos"] = promos

	result, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.AddRawFile(project.Slug+"/update.json", string(result))

	return nil
}
