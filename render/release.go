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

func (p Release) NavItem() string {
	return "releases"
}

type Release struct {
	project *project.Project
	site    *site.Site
	Release *release.Release
}

func (r *Renderer) AddRelease(inputFile string, project *project.Project, site *site.Site, rel *release.Release) {
	r.AddFile(inputFile, project.Slug+"/releases/"+rel.Slug+"/index.html", &Release{
		project: project,
		Release: rel,
		site:    site,
	})
}
