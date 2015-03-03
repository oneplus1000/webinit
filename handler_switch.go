package webinit

import (
	"errors"
	"net/http"
)

var ERROR_HANDLER_NOT_FOUND = errors.New("not found handler for this pattern")

type HandlerSwitch struct {
	handlers map[string]http.HandlerFunc
}

func (me *HandlerSwitch) SetHandler(pattern string, handlerFunc http.HandlerFunc) error {

	if me.handlers == nil {
		me.handlers = make(map[string]http.HandlerFunc)
	}

	if _, ok := me.handlers[pattern]; ok {
		return errors.New("dup pattern")
	}

	me.handlers[pattern] = handlerFunc

	return nil
}

func (me *HandlerSwitch) Handler(pattern string) (http.HandlerFunc, error) {
	if val, ok := me.handlers[pattern]; ok {
		return val, nil
	}
	return nil, ERROR_HANDLER_NOT_FOUND
}
