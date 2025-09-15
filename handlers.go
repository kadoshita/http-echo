// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

			// Create log attributes with basic request info
			logAttrs := []slog.Attr{
				slog.String("host", r.Host),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("protocol", r.Proto),
				slog.Int("status_code", status),
				slog.Int("response_length", length),
				slog.String("user_agent", r.UserAgent()),
				slog.Int64("duration_ms", dur.Milliseconds()),
			}

			// Add all request headers
			for name, values := range r.Header {
				// Join multiple values with comma if they exist
				headerValue := ""
				if len(values) > 0 {
					if len(values) == 1 {
						headerValue = values[0]
					} else {
						headerValue = fmt.Sprintf("[%s]", strings.Join(values, ", "))
					}
				}
				logAttrs = append(logAttrs, slog.String("header_"+strings.ToLower(strings.ReplaceAll(name, "-", "_")), headerValue))
			}

			slog.LogAttrs(context.Background(), slog.LevelInfo, "request", logAttrs...)
		}(time.Now())

		h(&mrw, r)
	}
}
