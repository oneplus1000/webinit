package webinit

import (
	"log"
	"net/http"
	//"html/template"
)

type BaseController struct {
	winit *WebInit
}

func (me *BaseController) Init(winit *WebInit) {
	me.winit = winit
}

func (me *BaseController) Render(w http.ResponseWriter, tmplname string, data interface{}) {

	tmpl, err := me.winit.GetTmpl(tmplname)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, tmplname, data)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
