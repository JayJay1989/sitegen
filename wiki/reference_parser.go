package wiki

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type stage int

const (
	referencesAndVariables stage = iota
	includes
)

func parseReferenceLinks(body, name, projectSlug string, stage stage, byName WikisByName) (string, error) {
	r, err := regexp.Compile("\\[\\[(.+?)(\\|(.+?))?(#(.+?))?]]")
	if err != nil {
		return "", err
	}

	result := string(r.ReplaceAllFunc([]byte(body), func(bytes []byte) []byte {
		if err != nil {
			return bytes
		}

		reference := r.ReplaceAllString(string(bytes), `$1`)

		if strings.HasPrefix(reference, "#") {
			if stage == referencesAndVariables {
				if reference == "#name" {
					return []byte(name)
				}

				err = errors.New("Unknown variable " + reference[1:] + " in wiki " + name)
			}

			return bytes
		} else if strings.HasPrefix(reference, "@") {
			if stage == includes {
				referencedWiki := byName[reference[1:]]
				if referencedWiki == nil {
					err = errors.New(fmt.Sprintf("could not find include to %s in wiki %s", reference[1:], name))

					return bytes
				}

				return []byte(referencedWiki.Body)
			}

			return bytes
		} else if stage == referencesAndVariables {
			referencedWiki := byName[reference]
			if referencedWiki == nil {
				err = errors.New(fmt.Sprintf("could not find reference to %s in wiki %s", reference, name))

				return bytes
			}

			format := r.ReplaceAllString(string(bytes), `$3`)
			heading := r.ReplaceAllString(string(bytes), `$4`)
			headingData := ""

			if heading != "" {
				headingData = "#" + heading[1:]
			}

			tooltipData := ""
			if referencedWiki.Meta.Icon != "" {
				tooltipData = fmt.Sprintf(`data-tooltip-icon="%s"`, referencedWiki.Meta.Icon)
			}

			if format == "" {
				return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html%s" %s>%s</a>`, projectSlug, referencedWiki.Slug, headingData, tooltipData, referencedWiki.Name))
			} else {
				return []byte(fmt.Sprintf(`<a href="/%s/wiki/%s.html%s" %s>%s</a>`, projectSlug, referencedWiki.Slug, headingData, tooltipData, format))
			}
		} else {
			return bytes
		}
	}))

	return result, err
}
