package wiki

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type Wiki struct {
	Name string
	Slug string
	Body string
	Meta WikiMetadata
}

type WikiMetadata struct {
	Icon string `json:"icon"`
}

func (w *Wiki) Parse() {
	value := w.Body
	value = string(bytes.Replace([]byte(value), []byte("\r"), nil, -1))
	value = string(blackfriday.Run([]byte(value), blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink)))
	value = strings.ReplaceAll(value, "<table>", `<table class="table">`)
	value = strings.ReplaceAll(value, "<h2", `<h2 class="h4"`)
	value = strings.ReplaceAll(value, "<h3", `<h3 class="h5"`)
	w.Body = value

}

func (w *Wiki) ParseReferences(projectSlug string, index map[string]*Wiki) error {
	r, err := regexp.Compile("\\[\\[(.+?)(\\|(.+?))?]]")
	if err != nil {
		return err
	}

	w.Body = string(r.ReplaceAllFunc([]byte(w.Body), func(bytes []byte) []byte {
		reference := r.ReplaceAllString(string(bytes), `$1`)

		if strings.HasPrefix(reference, "#") {
			log.WithField("value", reference).WithField("wiki", w.Name).Warn("Variables are unsupported")
			return bytes
		}
		if strings.HasPrefix(reference, "@") {
			log.WithField("value", reference).WithField("wiki", w.Name).Warn("Includes are unsupported")
			return bytes
		}

		referencedWiki := index[reference]
		if referencedWiki == nil {
			log.WithField("value", reference).WithField("wiki", w.Name).Error("could not find reference")
			return bytes
		}

		format := r.ReplaceAllString(string(bytes), `$3`)

		if format == "" {
			return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html">%s</a>`, projectSlug, referencedWiki.Slug, referencedWiki.Name))
		} else {
			return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html">%s</a>`, projectSlug, referencedWiki.Slug, format))
		}
	}))

	return nil
}
