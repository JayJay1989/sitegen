package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
)

type RenderData interface {
	Project() *project.Project
	Title() string
	Site() *site.Site
}
