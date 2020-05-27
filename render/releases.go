package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/release"
	"github.com/refinedmods/sitegen/site"
	"math"
	"sort"
	"strconv"
)

const perPage = 25

type ReleaseList []*release.Release

func (p ReleaseList) Len() int {
	return len(p)
}

func (p ReleaseList) Less(i, j int) bool {
	return p[i].Date.After(*p[j].Date)
}

func (p ReleaseList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (r Releases) Project() *project.Project {
	return r.project
}

func (r Releases) Title() string {
	return "Releases"
}

func (r Releases) Site() *site.Site {
	return r.site
}

func (p Releases) NavItem() string {
	return "releases"
}

type Releases struct {
	project    *project.Project
	site       *site.Site
	Releases   ReleaseList
	Page       int
	TotalPages int
}

func (r *Renderer) AddReleases(inputFile string, project *project.Project, site *site.Site) {
	var releases ReleaseList
	for _, group := range project.ReleaseGroups {
		for _, rel := range group.Releases {
			releases = append(releases, rel)
		}
	}

	sort.Sort(releases)

	totalPages := int(math.Ceil(float64(len(releases)) / float64(perPage)))
	for page := 0; page < totalPages; page++ {
		r.addPaginatedRelease(inputFile, project, site, releases, page, totalPages, false)
	}

	r.addPaginatedRelease(inputFile, project, site, releases, 0, totalPages, true)
}

func (r *Renderer) addPaginatedRelease(inputFile string, project *project.Project, site *site.Site, releases ReleaseList, page int, pages int, index bool) {
	endRange := int(math.Min(float64(len(releases)), float64(page*perPage+perPage)))

	pageName := strconv.Itoa(page+1) + ".html"
	if index {
		pageName = "index.html"
	}

	r.AddFile(inputFile, project.Slug+"/releases/"+pageName, &Releases{
		project:    project,
		site:       site,
		Releases:   releases[page*perPage : endRange],
		Page:       page + 1,
		TotalPages: pages,
	})
}
