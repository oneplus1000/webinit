package webinit

import (
//"net/http"
)

type IController interface {
	Init(winit *WebInit)
	Methods() MapMethodInfo
}
