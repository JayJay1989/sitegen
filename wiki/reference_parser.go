package wiki

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func parseReferenceLinks(body, name, projectSlug string, index wikiIndex) (string, error) {
	r, err := regexp.Compile("\\[\\[(.+?)(\\|(.+?))?]]")
	if err != nil {
		return "", err
	}

	result := string(r.ReplaceAllFunc([]byte(body), func(bytes []byte) []byte {
		if err != nil {
			return bytes
		}

		reference := r.ReplaceAllString(string(bytes), `$1`)

		if strings.HasPrefix(reference, "#") {
			log.WithField("value", reference).WithField("wiki", name).Warn("Variables are unsupported")
			return bytes
		}
		if strings.HasPrefix(reference, "@") {
			log.WithField("value", reference).WithField("wiki", name).Warn("Includes are unsupported")
			return bytes
		}

		referencedWiki := index[reference]
		if referencedWiki == nil {
			err = errors.New(fmt.Sprintf("could not find reference to %s in wiki %s", reference, name))
			return bytes
		}

		format := r.ReplaceAllString(string(bytes), `$3`)

		tooltipData := ""
		if referencedWiki.Meta.Icon != "" {
			tooltipData = fmt.Sprintf(`data-tooltip-icon="%s"`, referencedWiki.Meta.Icon)
		}

		if format == "" {
			return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html" %s>%s</a>`, projectSlug, referencedWiki.Slug, tooltipData, referencedWiki.Name))
		} else {
			return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html" %s>%s</a>`, projectSlug, referencedWiki.Slug, tooltipData, format))
		}
	}))

	return result, err
}
