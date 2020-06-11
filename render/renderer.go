package render

import (
	"bufio"
	"github.com/otiai10/copy"
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
	log "github.com/sirupsen/logrus"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type renderFile struct {
	inputFile  string
	input      RenderData
	outputFile string
}

type Renderer struct {
	files            []*renderFile
	rawFiles         map[string]string
	layoutFile       string
	releaseBadgeFile string
	site             *site.Site
	outputLocation   string
}

func NewRenderer(outputLocation string, layoutFile string, releaseBadgeFile string, site *site.Site) *Renderer {
	return &Renderer{outputLocation: outputLocation, site: site, layoutFile: layoutFile, releaseBadgeFile: releaseBadgeFile, rawFiles: make(map[string]string)}
}

func (r *Renderer) AddFile(inputFile string, outputFile string, input RenderData) {
	r.files = append(r.files, &renderFile{
		inputFile:  inputFile,
		input:      input,
		outputFile: outputFile,
	})
}

func (r *Renderer) AddRawFile(outputFile, data string) {
	r.rawFiles[outputFile] = data
}

func (r *Renderer) RenderAll() error {
	log.WithField("amount", len(r.files)+len(r.rawFiles)).Info("Rendering files")

	for _, f := range r.files {
		log.WithField("inputFile", f.inputFile).WithField("outputFile", f.outputFile).Debug("Rendering file")

		tpl, err := template.New(r.layoutFile).Funcs(template.FuncMap{
			"nl2br": func(text string) template.HTML {
				return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
			},
			"rangeList": func(startInclusive int, endInclusive int) []int {
				var result []int
				for i := startInclusive; i <= endInclusive; i++ {
					result = append(result, i)
				}
				return result
			},
			"add": func(i, change int) int {
				return i + change
			},
			"sub": func(i, change int) int {
				return i - change
			},
			"wikiIcon": func(project *project.Project, wikiName string) template.HTML {
				wiki := project.WikisByName[wikiName]
				if wiki == nil {
					panic("not found: " + wikiName)
				}

				return template.HTML(`<a href="/` + project.Slug + `/wiki/` + wiki.Slug + `.html">
                                            <img src="` + wiki.Meta.Icon + `"
                                                 alt="` + wiki.Name + `"
                                                 class="m-1 home-icon"
                                                 data-tippy-content="` + wiki.Name + `"
                                            >
                                        </a>`)
			},
		}).ParseFiles(r.layoutFile, f.inputFile, r.releaseBadgeFile)
		if err != nil {
			return err
		}

		err = ensureDir(r.outputLocation + f.outputFile)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(r.outputLocation+f.outputFile, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		w := bufio.NewWriter(file)
		err = tpl.Execute(w, f.input)
		if err != nil {
			return err
		}

		err = w.Flush()
		if err != nil {
			return err
		}
	}

	for filename, value := range r.rawFiles {
		err := ensureDir(r.outputLocation + filename)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(r.outputLocation+filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		w := bufio.NewWriter(file)
		_, err = w.WriteString(value)
		if err != nil {
			return err
		}

		err = w.Flush()
		if err != nil {
			return err
		}
	}

	for _, directory := range r.site.CopyDirectories {
		err := copy.Copy(directory, r.outputLocation+directory)
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureDir(fileName string) error {
	err := os.MkdirAll(filepath.Dir(fileName), os.ModeDir)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
