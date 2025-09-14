// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kadoshita/http-echo/version"
)

const (
	httpHeaderAppName    string = "X-App-Name"
	httpHeaderAppVersion string = "X-App-Version"
)

// withAppHeaders adds application headers such as X-App-Version and X-App-Name.
func withAppHeaders(c int, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(httpHeaderAppName, version.Name)
		w.Header().Set(httpHeaderAppVersion, version.Version)
		w.WriteHeader(c)
		h(w, r)
	}
}

// metaResponseWriter is a response writer that saves information about the
// response for logging.
type metaResponseWriter struct {
	writer http.ResponseWriter
	status int
	length int
}

// Header implements the http.ResponseWriter interface.
func (w *metaResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// WriteHeader implements the http.ResponseWriter interface.
func (w *metaResponseWriter) WriteHeader(s int) {
	w.status = s
	w.writer.WriteHeader(s)
}

// Write implements the http.ResponseWriter interface.
func (w *metaResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.length = len(b)
	return w.writer.Write(b)
}

// httpLog logs the request and response objects using structured logging.
func httpLog(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mrw metaResponseWriter
		mrw.writer = w

		defer func(start time.Time) {
			status := mrw.status
			length := mrw.length
			end := time.Now()
			dur := end.Sub(start)
			xForwardedFor := r.Header.Get("X-Forwarded-For")
			if xForwardedFor == "" {
				xForwardedFor = "-"
			}

			slog.Info("request",
				"host", r.Host,
				"remote_addr", r.RemoteAddr,
				"method", r.Method,
				"path", r.URL.Path,
				"protocol", r.Proto,
				"status_code", status,
				"response_length", length,
				"user_agent", r.UserAgent(),
				"x_forwarded_for", xForwardedFor,
				"duration_ms", dur.Milliseconds(),
			)
		}(time.Now())

		h(&mrw, r)
	}
}
