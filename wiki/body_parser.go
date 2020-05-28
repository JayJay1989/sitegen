package wiki

import (
	"bytes"
	"github.com/russross/blackfriday/v2"
	"strings"
)

func parseBody(body string) string {
	body = string(bytes.Replace([]byte(body), []byte("\r"), nil, -1))
	body = string(blackfriday.Run([]byte(body), blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak|blackfriday.AutoHeadingIDs|blackfriday.Autolink)))
	body = strings.ReplaceAll(body, "<table>", `<table class="table">`)
	body = strings.ReplaceAll(body, "<h2", `<h2 class="h4"`)
	body = strings.ReplaceAll(body, "<h3", `<h3 class="h5"`)
	return body
}
