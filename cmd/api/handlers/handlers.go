package handlers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

type HTTPHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type RouteRegistrar interface {
	RegisterRoute(method string, group string, path string, handler HTTPHandler)
}

// TestRegistrar - for the testing purpose only
type TestRegistrar struct {
	handlerPatternToFunctionsMap map[string]string
}

func (r *TestRegistrar) RegisterRoute(method string, group string, path string, handler HTTPHandler) {
	if r.handlerPatternToFunctionsMap == nil {
		r.handlerPatternToFunctionsMap = make(map[string]string)
	}

	r.handlerPatternToFunctionsMap[fmt.Sprintf("%s %s%s", method, group, path)] = r.getFunctionName(handler)
}

func (r *TestRegistrar) IsRouteRegistered(pattern string, handler HTTPHandler) bool {
	return r.handlerPatternToFunctionsMap[pattern] == r.getFunctionName(handler)
}

func (r *TestRegistrar) getFunctionName(i any) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
