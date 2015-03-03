package webinit

import (
	"errors"
	"fmt"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
)

type WebInit struct {
	funcmap template.FuncMap

	setupinfo    *SetupInfo
	ctrls        map[string]IController
	jsbundlemap  map[string]([]string)
	cssbundlemap map[string]([]string)

	viewinfos []ViewInfo
	views     map[string]*template.Template
}

func (me *WebInit) Setup(setupinfo *SetupInfo) {
	me.setupinfo = setupinfo
}

func (me *WebInit) RegitFunc(fnname string, fn interface{}) {
	if me.funcmap == nil {
		me.funcmap = make(template.FuncMap)
	}
	me.funcmap[fnname] = fn
}

func (me *WebInit) RegitView(vname string, startTmplName string, tmplfiles []string) {
	for _, vinfo := range me.viewinfos {
		if vinfo.VName == vname {
			log.Panicf("dup view name %s\n", vname)
			return
		}
	}

	me.viewinfos = append(me.viewinfos, ViewInfo{
		VName:         vname,
		StartTmplName: startTmplName,
		TmplFiles:     tmplfiles,
	})
}

func (me *WebInit) View(vname string) (*template.Template, error) {

	for key, val := range me.views {
		if key == vname {
			return val, nil
		}
	}
	return nil, errors.New("no view found  ( vname = " + vname + ")")
}

func (me *WebInit) ViewInfo(vname string) (*ViewInfo, error) {
	for _, val := range me.viewinfos {
		if val.VName == vname {
			return &val, nil
		}
	}
	return nil, errors.New("no viewinfo found  ( vname = " + vname + ")")
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
	//me.bindTmpls()
	me.bindViews()
	err := http.ListenAndServe(me.setupinfo.Addr, nil)
	if err != nil {
		log.Panicf("%s\n", err.Error())
	}
}

func (me *WebInit) bindViews() {
	if me.views == nil {
		me.views = make(map[string]*template.Template)
	}

	for _, vinfo := range me.viewinfos {
		vname := vinfo.VName
		tmplfiles := vinfo.TmplFiles
		var filepaths []string
		for _, tmplfile := range tmplfiles {
			filepaths = append(filepaths, me.setupinfo.RootFolder+"/tmpls/"+tmplfile)
		}
		t, err := template.New(vname).Funcs(me.funcmap).ParseFiles(filepaths...)
		if err != nil {
			log.Panicf("error %s\n", err.Error())
			return
		}
		me.views[vname] = t
	}

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

func (me *WebInit) GlobalHandleFunc(w http.ResponseWriter, r *http.Request) {

}
