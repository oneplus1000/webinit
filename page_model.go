package webinit

import (
	"net/http"
)

type PageModel struct {
	HttpRequest *http.Request
	Data        interface{} //data from controller
}
