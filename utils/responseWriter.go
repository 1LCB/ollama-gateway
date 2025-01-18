package utils

import "net/http"

type CustomResponseWriter struct {
	StatusCode int
	http.ResponseWriter
}

func (x *CustomResponseWriter) WriteHeader(statusCode int)  {
	x.StatusCode = statusCode
	x.ResponseWriter.WriteHeader(statusCode)
}