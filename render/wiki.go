package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
	"github.com/refinedmods/sitegen/wiki"
	"html/template"
)

func (w Wiki) Project() *project.Project {
	return w.project
}

func (w Wiki) Title() string {
	return w.Wiki.Name
}

func (w Wiki) Site() *site.Site {
	return w.site
}

func (w Wiki) NavItem() string {
	return "wiki"
}

type Wiki struct {
	Wiki    *wiki.Wiki
	site    *site.Site
	project *project.Project
	Body    template.HTML
}

func (r *Renderer) AddWiki(inputFile string, wiki *wiki.Wiki, site *site.Site, project *project.Project) {
	r.AddFile(inputFile, project.Slug+"/wiki/"+wiki.Slug+".html", &Wiki{
		Wiki:    wiki,
		site:    site,
		project: project,
		Body:    template.HTML(wiki.Body),
	})
}
