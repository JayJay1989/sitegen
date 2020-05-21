package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
)

type ProjectIndex struct {
	project *project.Project
	site    *site.Site
}

func (p *ProjectIndex) Project() *project.Project {
	return p.project
}

func (p *ProjectIndex) Site() *site.Site {
	return p.site
}

func (p *ProjectIndex) Title() string {
	return p.project.Name
}

func (r *Renderer) AddProjectIndex(inputFile string, project *project.Project, site *site.Site) {
	r.AddFile(inputFile, project.Slug+"/index.html", &ProjectIndex{
		project: project,
		site:    site,
	})
}
