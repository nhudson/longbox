package api

import (
	"context"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

const (
	requestIDHeader = "X-Request-ID"
)

func logger(log logrus.FieldLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/healthz" {
			next.ServeHTTP(w, req)
			return
		}

		trw := &trackedResponseWriter{w, http.StatusOK}

		reqID := req.Header.Get(requestIDHeader)
		if reqID == "" {
			reqID = ksuid.New().String()
		}

		// https://godoc.org/net/http#Handler
		// Except for reading the body, handlers should not modify the provided Request.
		r2 := new(http.Request)
		*r2 = *req
		// propagate headers in case the request gets used for a remote call
		r2.Header.Set(requestIDHeader, reqID)

		reqLogger := log.WithFields(logrus.Fields{
			"request_id": reqID,
			"path":       r2.RequestURI,
			"method":     r2.Method,
		})

		sourceIP := r2.Header.Get("X-Forwarded-For")
		if sourceIP != "" {
			reqLogger = reqLogger.WithField("source_ip", sourceIP)
		}

		reqLogger.Infoln("processing request...")
		defer func(begin time.Time) {
			reqLogger.WithFields(logrus.Fields{
				"process_time": time.Since(begin),
				"status_code":  trw.statusCode,
			}).Infoln("request processed")
		}(time.Now())

		ctx := context.WithValue(r2.Context(), requestIDContextKey, reqID)
		next.ServeHTTP(trw, r2.WithContext(ctx))
	})
}
