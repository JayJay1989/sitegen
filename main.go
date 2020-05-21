package main

import (
	"fmt"
	"github.com/refinedmods/sitegen/render"
	"github.com/refinedmods/sitegen/site"
)

func main() {
	site, err := site.NewSite("site.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d projects\n", len(site.Projects))

	renderer := render.NewRenderer("../output/", site.Templates["layout"], site)

	for _, proj := range site.Projects {
		renderer.AddProjectIndex(proj.Templates["index"], proj, site)
	}

	err = renderer.RenderAll()
	if err != nil {
		panic(err)
	}
}
