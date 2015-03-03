package webinit

import (
	"net/http"
)

type MethodInfo struct {
	Name           string
	Handler        http.HandlerFunc
	IsSessionStart bool
}
