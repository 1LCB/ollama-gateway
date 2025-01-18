package utils

import "net/http"

type MiddlewareFunc = func(next http.HandlerFunc) http.HandlerFunc

type MiddlewareApplyer struct {
	middlewares []MiddlewareFunc
}

func NewMiddlewareApplyer(middlewares ...MiddlewareFunc) *MiddlewareApplyer {
	return &MiddlewareApplyer{
		middlewares: middlewares,
	}
}

func (x *MiddlewareApplyer) IncludeMiddleware(middleware MiddlewareFunc) {
	x.middlewares = append(x.middlewares, middleware)
}

func (x *MiddlewareApplyer) Apply(handler http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range x.middlewares {
		handler = middleware(handler)
	}
	return handler
}
