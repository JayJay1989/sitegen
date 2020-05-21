package main

import (
	"flag"
	"github.com/refinedmods/sitegen/render"
	"github.com/refinedmods/sitegen/site"
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func init() {
	log.SetOutput(os.Stdout)
}

var output = flag.String("output", "../output/", "the output directory")
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
		lastBuild := time.Now()

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
						diff := time.Now().Sub(lastBuild)
						if diff.Seconds() > 1 {
							log.Info("Noticed change, rebuilding")
							err = build(site)

							lastBuild = time.Now()
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

		err = watcher.Add("./")
		if err != nil {
			log.WithError(err).Fatal("Could not add directory to watcher")
		}
		<-done
	}
}

func build(site *site.Site) error {
	renderer := render.NewRenderer(*output, site.Templates["layout"], site)

	for _, proj := range site.Projects {
		renderer.AddProjectIndex(proj.Templates["index"], proj, site)
	}

	log.Info("Rendering all files")

	err := renderer.RenderAll()
	if err != nil {
		return err
	}

	log.Info("Done")

	return nil
}
