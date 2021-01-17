package viewx

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"path/filepath"
)

func View(view string, data interface{}, w http.ResponseWriter) {
	// load all files
	name := filepath.Base(fmt.Sprintf("./authentication/ui/html/views/%s.gohtml", view))

	ts, err := template.New(name).Funcs(functions).ParseFiles(fmt.Sprintf("./authentication/ui/html/views/%s.gohtml", view))
	if err != nil {

		log.Println(err.Error())
		return
	}

	ts, err = ts.ParseGlob(filepath.Join("./authentication/ui/html/layouts/*.gohtml"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	ts, err = ts.ParseGlob(filepath.Join("./authentication/ui/html/components/*.gohtml"))
	if err != nil {
		log.Println(err.Error())

		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
	}

}
