package handlers

import (
	"fmt"
	"reflect"
	"runtime"
)

type RouteRegistrarMock struct {
	handlerPatternToFunctionsMap map[string]string
}

func (r *RouteRegistrarMock) RegisterRoute(method string, group string, path string, handler HTTPHandler,
	mw ...Middleware) {

	if r.handlerPatternToFunctionsMap == nil {
		r.handlerPatternToFunctionsMap = make(map[string]string)
	}

	r.handlerPatternToFunctionsMap[fmt.Sprintf("%s %s%s", method, group, path)] = r.getFunctionName(handler)
}

func (r *RouteRegistrarMock) IsRouteRegistered(pattern string, handler HTTPHandler) bool {
	return r.handlerPatternToFunctionsMap[pattern] == r.getFunctionName(handler)
}

func (r *RouteRegistrarMock) getFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
