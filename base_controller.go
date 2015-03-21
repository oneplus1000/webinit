package webinit

import (
	"fmt"
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

func (me *BaseController) NewMapMethodInfo() MapMethodInfo {
	return make(MapMethodInfo)
}

func (me *BaseController) BindMethodInfo(
	m *MapMethodInfo,
	name string,
	handler http.HandlerFunc,
) {
	if _, ok := (*m)[name]; ok {
		log.Panicf("dup method name %s", name)
		return
	}
	(*m)[name] = MethodInfo{
		Name:    name,
		Handler: handler,
	}
}

func (me *BaseController) Render(w http.ResponseWriter, r *http.Request, viewname string, data interface{}) {

	view, err := me.winit.View(viewname)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//create page model
	var pagemodel PageModel
	pagemodel.HttpRequest = r
	pagemodel.Data = data

	//render
	viewinfo, err := me.winit.ViewInfo(viewname)
	err = view.ExecuteTemplate(w, viewinfo.StartTmplName, pagemodel)
	if err != nil {
		log.Printf("%s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (me *BaseController) WriteHttpErr(w http.ResponseWriter, r *http.Request, errcode int, errmsg string) {
	w.WriteHeader(errcode)
	fmt.Fprintf(w, "err: %s", errmsg)
}
