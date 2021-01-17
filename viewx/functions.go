package viewx

import (
	"github.com/russross/blackfriday/v2"
	"html/template"
	"time"
)

func markdown(s string) string {
	return string(blackfriday.Run([]byte(s)))
}

func humanDate(t time.Time) string { return t.Format("02 Jan 2006 at 15:04") }

// Initialize a template.FuncMap object and store it in a global variable. This is // essentially a string-keyed map which acts as a lookup between the names of our // custom template functions and the functions themselves.

var functions = template.FuncMap{
	"humanDate": humanDate,
	"markdown":  markdown,
	"htmlSafe": func(html string) template.HTML {
		return template.HTML(html)
	},
}
