package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/refinedmods/sitegen/render"
	"github.com/refinedmods/sitegen/site"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
}

var output = flag.String("output", "./output/", "the output directory")
var config = flag.String("config", "site.json", "the configuration file")

func main() {
	flag.Parse()

	site, err := site.NewSite(*config)
	if err != nil {
		log.WithError(err).Fatal("Could not load site config")
	}

	var cmd = flag.Arg(0)

	switch cmd {
	case "build":
		err = build(site)
		if err != nil {
			log.WithError(err).Fatal("Could not build")
		}
	case "watch":
		err = build(site)
		if err != nil {
			log.WithError(err).Error("Build error")
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.WithError(err).Fatal("Could not create watcher")
		}
		defer watcher.Close()

		done := make(chan bool)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write {
						err = build(site)
						if err != nil {
							log.WithError(err).Error("Build error")
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.WithError(err).Error("File watcher error occurred")
				}
			}
		}()

		err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() && !strings.Contains(path, "node_modules") {
				log.WithField("dir", path).Info("Watching directory")
				return watcher.Add(path)
			}
			return nil
		})

		if err != nil {
			log.WithError(err).Fatal("Could not add directory to watcher")
		}
		<-done
	}
}

func build(site *site.Site) error {
	renderer := render.NewRenderer(*output, site.Templates["layout"], site.Templates["releaseBadge"], site)

	for _, proj := range site.Projects {
		renderer.AddProjectIndex(proj.Templates["index"], proj, site)

		for _, group := range proj.ReleaseGroups {
			for _, release := range group.Releases {
				renderer.AddRelease(site.Templates["release"], proj, site, release)
			}
		}

		renderer.AddReleases(site.Templates["releases"], proj, site)

		for _, wiki := range proj.Wikis {
			renderer.AddWiki(site.Templates["wiki"], wiki, site, proj)
		}

		renderer.AddWikiIndex(site.Templates["wikiIndex"], site, proj)

		err := renderer.AddUpdate(site, proj)
		if err != nil {
			log.WithError(err).Fatal("Could not render update")
		}
	}

	err := renderer.RenderAll()
	if err != nil {
		return err
	}

	log.Info("Done")

	return nil
}
