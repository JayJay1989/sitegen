package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
)

func (w WikiIndex) Project() *project.Project {
	return w.project
}

func (w WikiIndex) Title() string {
	return "Wiki"
}

func (w WikiIndex) Site() *site.Site {
	return w.site
}

func (w WikiIndex) NavItem() string {
	return "wiki"
}

type WikiIndex struct {
	site    *site.Site
	project *project.Project
}

func (r *Renderer) AddWikiIndex(inputFile string, site *site.Site, project *project.Project) {
	r.AddFile(inputFile, project.Slug+"/wiki/index.html", &WikiIndex{
		site:    site,
		project: project,
	})
}
