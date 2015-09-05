package rwjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sendo/sdlog"
)

var respformat = "{ \"data\" : %s ,  \"errmsg\" : \"%s\", \"errcode\" : %d   }"

//write json data to http response
func WriteJsonData(w http.ResponseWriter, r *http.Request, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		sdlog.Err(err, "WriteJsonData Marshal json")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, respformat, string(b), "", 0)
}

//write json error to http response
func WriteJsonErr(w http.ResponseWriter, r *http.Request, err error, errcode int) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(errcode)
	fmt.Fprintf(w, respformat, "null", msg, errcode)
}
