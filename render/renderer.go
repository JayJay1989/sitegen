package render

import (
	"bufio"
	"github.com/otiai10/copy"
	"github.com/refinedmods/sitegen/project"
	"github.com/refinedmods/sitegen/site"
	"html/template"
	"os"
	"path/filepath"
)

type renderFile struct {
	inputFile  string
	input      RenderData
	outputFile string
	project    *project.Project
	title      string
}

type Renderer struct {
	files            []*renderFile
	layoutFile       string
	releaseBadgeFile string
	site             *site.Site
	outputLocation   string
}

func NewRenderer(outputLocation string, layoutFile string, releaseBadgeFile string, site *site.Site) *Renderer {
	return &Renderer{outputLocation: outputLocation, site: site, layoutFile: layoutFile, releaseBadgeFile: releaseBadgeFile}
}

func (r *Renderer) AddFile(inputFile string, outputFile string, input RenderData) {
	r.files = append(r.files, &renderFile{
		inputFile:  inputFile,
		input:      input,
		outputFile: outputFile,
	})
}

func (r *Renderer) RenderAll() error {
	for _, f := range r.files {
		tpl, err := template.ParseFiles(r.layoutFile, f.inputFile, r.releaseBadgeFile)
		if err != nil {
			return err
		}

		err = ensureDir(r.outputLocation + f.outputFile)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(r.outputLocation+f.outputFile, os.O_TRUNC|os.O_CREATE, 0666)
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
