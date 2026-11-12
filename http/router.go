package http

import (
	err "github.com/adagit94/golib/errors"
	str "strings"
)

const (
	Post    = "POST"
	Get     = "GET"
	Put     = "PUT"
	Patch   = "PATCH"
	Delete  = "DELETE"
	Options = "OPTIONS"
	Connect = "CONNECT"
	Head    = "HEAD"
	Trace   = "TRACE"
)

type PathParams map[string]string

type IRouter[H any] interface {
	Route(path string) IRoute[H]
	Select(path, method string) (H, PathParams)
}

// Creates router with support for dynamic segments & query params.
func CreateRouter[H any]() IRouter[H] {
	router := &router[H]{tree: make(tree[H])}
	return router
}

type routeSeg struct {
	seg string
	dyn bool
}

type dynRoute[H any] struct {
	segs           []*routeSeg
	dynSegsIndices []int
	handler        H
}

type statRoutes[H any] map[string]H
type dynRoutes[H any] map[int][]*dynRoute[H]

type routes[H any] struct {
	stat statRoutes[H]
	dyn  dynRoutes[H]
}

type tree[H any] map[string]*routes[H]

type router[H any] struct {
	tree tree[H]
}

func (router *router[H]) registerHandler(path, method string, handler H) {
	methodRoutes, methodRoutesFound := router.tree[method]

	if !methodRoutesFound {
		router.tree[method] = &routes[H]{stat: make(statRoutes[H]), dyn: make(dynRoutes[H])}
		methodRoutes = router.tree[method]
	}

	segs := str.Split(path, "/")
	segsCount := len(segs)
	dynSegsIndices := collectDynSegsIndices(segs)

	if len(dynSegsIndices) > 0 {
		dynRoutes, dynRoutesFound := methodRoutes.dyn[segsCount]

		if !dynRoutesFound {
			methodRoutes.dyn[segsCount] = make([]*dynRoute[H], 0, 1)
			dynRoutes = methodRoutes.dyn[segsCount]
		}

		dynRoutes = append(dynRoutes, &dynRoute[H]{segs: genRouteSegs(segs), dynSegsIndices: dynSegsIndices, handler: handler})
		sortRoutes(dynRoutes)
		methodRoutes.dyn[segsCount] = dynRoutes
	} else {
		methodRoutes.stat[path] = handler
	}
}

// Create new route for which respective http methods handlers can be defined.
func (router *router[H]) Route(path string) IRoute[H] {
	route := &route[H]{path: path, registerHandler: router.registerHandler}
	return route
}

// Select is mean't to be wrapped inside of specific http implementation entry point for incoming requests. It returns selected handler with parsed path parameters for a found route. It panics in case no route gets matched or method doesn't exists in register, so it's important to recover and handle such potential errors gracefully in outer code.
func (router *router[H]) Select(path, method string) (H, PathParams) {
	methodRoutes, methodRoutesFound := router.tree[method]

	if !methodRoutesFound {
		panic(&err.Err{Code: err.HttpNotFoundCode, Message: "Method not registered."})
	}

	// Faster check for a static route.
	if statRouteHandler, statRouteHandlerFound := methodRoutes.stat[path]; statRouteHandlerFound {
		return statRouteHandler, nil
	}

	segs := str.Split(path, "/")
	dynRoutes, dynRoutesFound := methodRoutes.dyn[len(segs)]

	if !dynRoutesFound {
		panic(&err.Err{Code: err.HttpNotFoundCode, Message: "Route not registered."})
	}

outter:
	for _, dynRoute := range dynRoutes {
		for i, routeSeg := range dynRoute.segs {
			if routeSeg.dyn {
				continue
			}

			if routeSeg.seg != segs[i] {
				continue outter
			}
		}

		pathParams := make(PathParams, len(dynRoute.dynSegsIndices))

		for _, segIndex := range dynRoute.dynSegsIndices {
			pathParams[dynRoute.segs[segIndex].seg] = segs[segIndex]
		}

		return dynRoute.handler, pathParams
	}

	panic(&err.Err{Code: err.HttpNotFoundCode, Message: "Route not registered."})
}

type IRoute[H any] interface {
	Post(H) IRoute[H]
	Get(H) IRoute[H]
	Put(H) IRoute[H]
	Patch(H) IRoute[H]
	Delete(H) IRoute[H]
	Options(H) IRoute[H]
	Connect(H) IRoute[H]
	Head(H) IRoute[H]
	Trace(H) IRoute[H]
}

type route[H any] struct {
	path            string
	registerHandler func(path, method string, handler H)
}

func (route *route[H]) Post(handler H) IRoute[H] {
	route.registerHandler(route.path, Post, handler)
	return route
}

func (route *route[H]) Get(handler H) IRoute[H] {
	route.registerHandler(route.path, Get, handler)
	return route
}

func (route *route[H]) Put(handler H) IRoute[H] {
	route.registerHandler(route.path, Put, handler)
	return route
}

func (route *route[H]) Patch(handler H) IRoute[H] {
	route.registerHandler(route.path, Patch, handler)
	return route
}

func (route *route[H]) Delete(handler H) IRoute[H] {
	route.registerHandler(route.path, Delete, handler)
	return route
}

func (route *route[H]) Options(handler H) IRoute[H] {
	route.registerHandler(route.path, Options, handler)
	return route
}

func (route *route[H]) Connect(handler H) IRoute[H] {
	route.registerHandler(route.path, Connect, handler)
	return route
}

func (route *route[H]) Head(handler H) IRoute[H] {
	route.registerHandler(route.path, Head, handler)
	return route
}

func (route *route[H]) Trace(handler H) IRoute[H] {
	route.registerHandler(route.path, Trace, handler)
	return route
}
