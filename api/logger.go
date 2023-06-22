//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// Logger implements a custom logger by wrapping the handler
//
//goland:noinspection ALL
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the start time and source IP
		startTime := time.Now()
		src := getIP(r)

		// Service the request
		sw := statusWriter{ResponseWriter: w}
		inner.ServeHTTP(&sw, r)

		// Get duration of request
		duration := time.Since(startTime)

		// Remove parameters from URI to avoid logging confidential information
		uri := strings.Split(r.RequestURI, "?")[0]

		// Add code here to send the log event somewhere other than stdout
		fmt.Printf("%s %s %s %d %f\n", src, r.Method, uri, sw.status, duration.Seconds())
	})
}
