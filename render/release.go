package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/release"
	"github.com/refinedmods/sitegen/site"
)

func (r Release) Project() *project.Project {
	return r.project
}

func (r Release) Title() string {
	return r.Release.Version
}

func (r Release) Site() *site.Site {
	return r.site
}

type Release struct {
	project      *project.Project
	site         *site.Site
	Release      *release.Release
	ReleaseGroup *release.ReleaseGroup
}

func (r *Renderer) AddRelease(inputFile string, project *project.Project, site *site.Site, rel *release.Release, group *release.ReleaseGroup) {
	r.AddFile(inputFile, project.Slug+"/releases/"+rel.Slug+".html", &Release{
		project:      project,
		Release:      rel,
		ReleaseGroup: group,
		site:         site,
	})
}
