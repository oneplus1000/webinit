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

	viewinfos map[string]ViewInfo
	views     map[string]*template.Template

	methodInfos map[string]MethodInfo
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

	if me.viewinfos == nil {
		me.viewinfos = make(map[string]ViewInfo)
	}

	if _, ok := me.viewinfos[vname]; ok {
		log.Panicf("dup view name %s\n", vname)
		return
	}

	me.viewinfos[vname] = ViewInfo{
		VName:         vname,
		StartTmplName: startTmplName,
		TmplFiles:     tmplfiles,
	}
}

func (me *WebInit) View(vname string) (*template.Template, error) {
	if view, ok := me.views[vname]; ok {
		return view, nil
	}
	return nil, errors.New("no view found  ( vname = " + vname + ")")
}

func (me *WebInit) ViewInfo(vname string) (*ViewInfo, error) {

	if vinfo, ok := me.viewinfos[vname]; ok {
		return &vinfo, nil
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
	if me.methodInfos == nil {
		me.methodInfos = make(map[string]MethodInfo)
	}
	for cname, c := range me.ctrls {
		c.Init(me)
		methods := c.Methods()
		for mname, minfo := range methods {
			//pattern := fmt.Sprintf("%s/%s", cname, mname)
			patterns := me.Patterns(cname, mname)
			for _, pattern := range patterns {
				http.HandleFunc(pattern, me.GlobalHandleFunc)
				err := me.addMethodInfo(pattern, minfo)
				if err != nil {
					log.Panicf("error %s\n", err.Error())
					return
				}
				log.Printf("regit controller %s\n", pattern)
			}

		}
	}
}

func (me *WebInit) Patterns(ctrlname string, methodname string) []string {
	var patterns []string
	patterns = append(patterns, fmt.Sprintf("/%s/%s", ctrlname, methodname))
	if ctrlname == "home" && methodname == "index" {
		patterns = append(patterns, "/")
	} else if methodname == "index" {
		patterns = append(patterns, "/%s", ctrlname)
	}
	return patterns
}

func (me *WebInit) addMethodInfo(pattern string, minfo MethodInfo) error {
	if _, ok := me.methodInfos[pattern]; ok {
		return errors.New("dup pattern")
	}
	me.methodInfos[pattern] = minfo
	return nil
}

func (me *WebInit) GlobalHandleFunc(w http.ResponseWriter, r *http.Request) {
	pattern := r.RequestURI
	if minfo, ok := me.methodInfos[pattern]; ok {
		if me.setupinfo.HotReloadView {
			me.views = nil //reset
			me.bindViews() //re compile all templ
		}
		minfo.Handler(w, r)
		return
	}
	fmt.Fprintf(w, "page not found")
}
