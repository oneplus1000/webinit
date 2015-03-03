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

func (me *BaseController) BindMethodInfo(
	m *map[string]MethodInfo,
	name string,
	isSessionStart bool,
	handler http.HandlerFunc,
) {
	if _, ok := (*m)[name]; ok {
		log.Panicf("dup method name %s", name)
		return
	}
	(*m)[name] = MethodInfo{
		Name:           name,
		IsSessionStart: isSessionStart,
		Handler:        handler,
	}
}

func (me *BaseController) Render(w http.ResponseWriter, viewname string, data interface{}) {

	view, err := me.winit.View(viewname)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewinfo, err := me.winit.ViewInfo(viewname)
	err = view.ExecuteTemplate(w, viewinfo.StartTmplName, data)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
