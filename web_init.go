package webinit

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type WebInit struct {
	funcmap      template.FuncMap
	cahceTmpls   map[string]*template.Template
	setupinfo    *SetupInfo
	ctrls        map[string]IController
	jsbundlemap  map[string]([]string)
	cssbundlemap map[string]([]string)
}

func (me *WebInit) Setup(setupinfo *SetupInfo) {
	me.setupinfo = setupinfo
}

func (me *WebInit) GetTmpl(name string) (*template.Template, error) {
	for key, val := range me.cahceTmpls {
		if key == name {
			return val, nil
		}
	}
	return nil, errors.New("no template found  ( name = " + name + ")")
}

func (me *WebInit) RegitFunc(fnname string, fn interface{}) {
	if me.funcmap == nil {
		me.funcmap = make(template.FuncMap)
	}
	me.funcmap[fnname] = fn
}

func (me *WebInit) RegitTmpl(name string, view string) {
	if me.cahceTmpls == nil {
		me.cahceTmpls = make(map[string]*template.Template)
	}

	b, err := ioutil.ReadFile(me.setupinfo.RootFolder + "/views/" + view + ".html")
	if err != nil {
		log.Panicf("%s\n", err.Error())
		return
	}
	me.cahceTmpls[name], err = template.New(name).Funcs(me.funcmap).Parse(string(b))
	if err != nil {
		log.Panicf("%s\n", err.Error())
		return
	}
}

func (me *WebInit) RegistCtrl(ctrl IController, names ...string) {
	if me.ctrls == nil {
		me.ctrls = make(map[string]IController)
	}
	for _, name := range names {
		me.ctrls[name] = ctrl
	}
}

func (me *WebInit) RegitJsBundle(name string, jsfiles []string) {
	if me.jsbundlemap == nil {
		me.jsbundlemap = make(map[string]([]string))
	}
	me.jsbundlemap[name] = jsfiles
}

func (me *WebInit) RegitCssBundle(name string, cssfiles []string) {
	if me.cssbundlemap == nil {
		me.cssbundlemap = make(map[string]([]string))
	}
	me.cssbundlemap[name] = cssfiles
}

func (me *WebInit) ListenAndServe() {
	me.bindCtrls()
	http.ListenAndServe(me.setupinfo.Addr, nil)
}

func (me *WebInit) bindCtrls() {

	//static file
	http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s%s", me.setupinfo.RootFolder, r.URL.Path[1:]))
	})

	//all controller
	for cname, c := range me.ctrls {
		c.Init(me)
		methods := c.Methods()
		for mname, m := range methods {
			pattern := fmt.Sprintf("%s/%s", cname, mname)
			http.HandleFunc(pattern, m)
			fmt.Printf("regit controller %s\n", pattern)
		}
	}
}
