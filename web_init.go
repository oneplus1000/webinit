package webinit

import (
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"
	//"strings"
	"bytes"
	"io/ioutil"
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

	jstmplinfos map[string]string
}

func (w *WebInit) Setup(setupinfo *SetupInfo) {
	w.setupinfo = setupinfo
}

//Regit template func to all template
func (w *WebInit) RegitFunc(fnname string, fn interface{}) {
	if w.funcmap == nil {
		w.funcmap = make(template.FuncMap)
	}
	w.funcmap[fnname] = fn
}

//Regist view
func (w *WebInit) RegitView(vname string, startTmplName string, tmplfiles []string) {

	if w.viewinfos == nil {
		w.viewinfos = make(map[string]ViewInfo)
	}

	if _, ok := w.viewinfos[vname]; ok {
		log.Panicf("dup view name %s\n", vname)
		return
	}

	w.viewinfos[vname] = ViewInfo{
		VName:         vname,
		StartTmplName: startTmplName,
		TmplFiles:     tmplfiles,
	}
}

func (w *WebInit) View(vname string) (*template.Template, error) {
	if view, ok := w.views[vname]; ok {
		return view, nil
	}
	return nil, errors.New("no view found  ( vname = " + vname + ")")
}

func (w *WebInit) ViewInfo(vname string) (*ViewInfo, error) {

	if vinfo, ok := w.viewinfos[vname]; ok {
		return &vinfo, nil
	}

	return nil, errors.New("no viewinfo found  ( vname = " + vname + ")")
}

func (w *WebInit) RegistCtrl(ctrl IController, names ...string) {
	if w.ctrls == nil {
		w.ctrls = make(map[string]IController)
	}
	for _, name := range names {
		w.ctrls[name] = ctrl
	}
}

func (w *WebInit) RegitJsBundle(name string, jsfiles []string) {
	if w.jsbundlemap == nil {
		w.jsbundlemap = make(map[string]([]string))
	}
	w.jsbundlemap[name] = jsfiles
}

func (w *WebInit) RegitJsTmpl(name string, file string) {
	if w.jstmplinfos == nil {
		w.jstmplinfos = make(map[string]string)
	}
	w.jstmplinfos[name] = file
}

func (w *WebInit) RegitCssBundle(name string, cssfiles []string) {
	if w.cssbundlemap == nil {
		w.cssbundlemap = make(map[string]([]string))
	}
	w.cssbundlemap[name] = cssfiles
}

func (w *WebInit) ListenAndServe() {
	w.bindCtrls()
	w.bindViews()
	err := http.ListenAndServe(w.setupinfo.Addr, nil)
	if err != nil {
		log.Panicf("%s\n", err.Error())
	}
}

func (w *WebInit) bindView(vinfo *ViewInfo) (*template.Template, error) {
	vname := vinfo.VName
	tmplfiles := vinfo.TmplFiles
	var filepaths []string
	for _, tmplfile := range tmplfiles {
		filepaths = append(filepaths, w.setupinfo.RootFolder+"/tmpls/"+tmplfile)
	}

	delimLeft := w.delimLeft()
	delimRight := w.delimRight()

	//install build-in tmpl func
	funcmap := w.funcmap
	funcmap["JsBundle"] = w.JsBundle
	funcmap["JsTmpl"] = w.JsTmpl
	funcmap["JsTmplWithData"] = w.JsTmplWithData
	funcmap["CssBundle"] = w.CssBundle

	tmpl, err := template.New(vname).
		Delims(delimLeft, delimRight).
		Funcs(funcmap).
		ParseFiles(filepaths...)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (w *WebInit) reBindView(vname string) {
	if w.views == nil {
		log.Panicf("view is empty!")
		return
	}

	if vinfo, ok := w.viewinfos[vname]; ok {
		tmpl, err := w.bindView(&vinfo)
		if err != nil {
			log.Panicf("%s", err.Error())
			return
		}
		w.views[vname] = tmpl
	} else {
		log.Panicf("view %s not found!", vname)
		return
	}
}

func (w *WebInit) bindViews() {
	if w.views == nil {
		w.views = make(map[string]*template.Template)
	}

	for _, vinfo := range w.viewinfos {
		vname := vinfo.VName
		tmpl, err := w.bindView(&vinfo)
		if err != nil {
			log.Panicf("error %s\n", err.Error())
			return
		}
		w.views[vname] = tmpl
	}

}

func (w *WebInit) bindCtrls() {

	//static file
	http.HandleFunc("/public/", func(wr http.ResponseWriter, r *http.Request) {
		http.ServeFile(wr, r, fmt.Sprintf("%s%s", w.setupinfo.RootFolder, r.URL.Path[1:]))
	})

	//all controller
	if w.methodInfos == nil {
		w.methodInfos = make(map[string]MethodInfo)
	}
	for cname, c := range w.ctrls {
		c.Init(w)
		methods := c.Methods()
		for mname, minfo := range methods {
			//pattern := fmt.Sprintf("%s/%s", cname, mname)
			patterns := w.UrlPatterns(cname, mname)
			for _, pattern := range patterns {
				http.HandleFunc(pattern, w.GlobalHandleFunc)
				err := w.addMethodInfo(pattern, minfo)
				if err != nil {
					log.Panicf("error %s\n", err.Error())
					return
				}
				log.Printf("regit controller %s\n", pattern)
			}

		}
	}
}

func (w *WebInit) delimLeft() string {
	delimLeft := "{{"
	if w.setupinfo.DelimLeft != "" {
		delimLeft = w.setupinfo.DelimLeft
	}
	return delimLeft
}

func (w *WebInit) delimRight() string {
	delimRight := "}}"
	if w.setupinfo.DelimRight != "" {
		delimRight = w.setupinfo.DelimRight
	}
	return delimRight
}

func (w *WebInit) JsTmplWithData(data interface{}, name string) template.HTML {

	delimLeft := w.delimLeft()
	delimRight := w.delimRight()

	if file, ok := w.jstmplinfos[name]; ok {

		path := w.setupinfo.RootFolder + file
		filename := filepath.Base(path)
		tmpl, err := template.New(filename).Delims(delimLeft, delimRight).ParseFiles(path)
		if err != nil {
			log.Panicf("error %s\n", err.Error())
			return template.HTML("<!--[ERROR] JsTmpl  (JsTmplWithViewData) " + name + " err=" + err.Error() + "  -->")
		}

		var content bytes.Buffer
		err = tmpl.Execute(&content, data)
		if err != nil {
			log.Panicf("error %s\n", err.Error())
			return template.HTML("<!--[ERROR] JsTmpl  (JsTmplWithViewData) " + name + " err=" + err.Error() + "  -->")
		}

		var buff bytes.Buffer
		buff.WriteString("\n<script id='" + name + "' type='text/template' >\n")
		buff.WriteString(content.String())
		buff.WriteString("\n</script>\n")
		return template.HTML(buff.String())

	}

	return template.HTML("<!--[ERROR] not found JsTmpl (JsTmplWithViewData) " + name + " -->")
}

func (w *WebInit) JsTmpl(data interface{}, name string) template.HTML {

	if file, ok := w.jstmplinfos[name]; ok {
		path := w.setupinfo.RootFolder + file
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Printf("[ERROR] JsTmpl %s err=%s", name, err.Error())
			return template.HTML("<!--[ERROR] JsTmpl " + name + " err=" + err.Error() + "  -->")
		}
		var buff bytes.Buffer
		buff.WriteString("\n<script id='" + name + "' type='text/template' >\n")
		buff.Write(content)
		buff.WriteString("\n</script>\n")
		return template.HTML(buff.String())
	}
	log.Printf("[ERROR] not found JsTmpl %s ", name)
	return template.HTML("<!--[ERROR] not found JsTmpl " + name + " -->")
}

func (w *WebInit) CssBundle(data interface{}, name string) template.HTML {

	cssversion := strings.TrimSpace(w.setupinfo.CssVersion)
	if cssversion == "" {
		cssversion = "0"
	}

	if files, ok := w.cssbundlemap[name]; ok {
		var buff bytes.Buffer
		for _, file := range files {
			buff.Write([]byte(fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s?cssv=%s\" ></link>\n", file, cssversion)))
		}
		return template.HTML(buff.String())
	}
	return template.HTML("<!-- not found CssBundle " + name + " -->")
}

func (w *WebInit) JsBundle(data interface{}, name string) template.HTML { //TODO  next time make compress js

	jsversion := strings.TrimSpace(w.setupinfo.JsVersion)
	if jsversion == "" {
		jsversion = "0"
	}

	if files, ok := w.jsbundlemap[name]; ok {
		var buff bytes.Buffer
		for _, file := range files {
			buff.Write([]byte(fmt.Sprintf("<script type=\"text/javascript\" src=\"%s?jsv=%s\" ></script>\n", file, jsversion)))
		}
		return template.HTML(buff.String())
	}
	log.Printf("[ERROR] not found JsBundle %s ", name)
	return template.HTML("<!-- not found JsBundle " + name + " -->")
}

func (w *WebInit) UrlPatterns(ctrlname string, methodname string) []string {
	var patterns []string
	patterns = append(patterns, fmt.Sprintf("/%s/%s/", ctrlname, methodname))
	if ctrlname == "home" && methodname == "index" {
		patterns = append(patterns, "/")
	} else if methodname == "index" {
		patterns = append(patterns, fmt.Sprintf("/%s/", ctrlname))
	}
	return patterns
}

func (w *WebInit) addMethodInfo(pattern string, minfo MethodInfo) error {
	if _, ok := w.methodInfos[pattern]; ok {
		return errors.New("dup pattern")
	}
	w.methodInfos[pattern] = minfo
	return nil
}

func (w *WebInit) GlobalHandleFunc(wr http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Path
	if minfo, ok := w.methodInfos[pattern]; ok {
		minfo.Handler(wr, r) //Go!
		return
	}
	wr.WriteHeader(404)
	fmt.Fprintf(wr, "page not found %s", pattern)
	log.Printf("page not found %s", pattern)
}

func LogPprof() {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	log.Printf("{ Alloc: %d,TotalAlloc: %d,HeapAlloc: %d,HeapSys: %d  ", m.Alloc, m.TotalAlloc, m.HeapAlloc, m.HeapSys)
}
