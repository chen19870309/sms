package utils

import (
	"bytes"
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

var show = template.Must(template.New("show").Parse(`<!DOCTYPE html>
	<html>
		<head>
			<title>{{.Title}}</title>
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link rel="stylesheet" href="{{.Css}}" >
		</head>
		<body>
			<article>{{.Text}}</article>
		</body>
	</html>`))

func MakeStaticPage(markd string) string {
	unsafe := blackfriday.MarkdownCommon([]byte(markd))
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
}

func TemplateMarkDown(markd string) string {
	var buff bytes.Buffer
	params := map[string]interface{}{
		"Title": "测试",
		"Css":   "markdown.css",
		"Text":  template.HTML(MakeStaticPage(markd)),
	}
	show.Execute(&buff, params)
	return buff.String()
}
