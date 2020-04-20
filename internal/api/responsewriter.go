package api

import (
	"net/http"
)

type trackedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (trw *trackedResponseWriter) WriteHeader(code int) {
	trw.statusCode = code
	trw.ResponseWriter.WriteHeader(code)
}
