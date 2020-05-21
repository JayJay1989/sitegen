package render

import (
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
	"html/template"
)

type layout struct {
	Site    *site.Site
	Project *project.Project
	Body    template.HTML
}
