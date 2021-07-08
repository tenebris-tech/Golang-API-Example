//
// Copyright (c) 2021 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

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

// Logger implements a custom logger
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		inner.ServeHTTP(&sw, r)
		src := getIP(r)

		// Remove parameters from URI to avoid logging confidential information
		uri := strings.Split(r.RequestURI, "?")[0]

		// Don't log health checks to reduce log noise
		if name != "health" {

			// Add code here to send the log event somewhere other than stdout
			fmt.Printf("%s %s %d %d\n", src, uri, sw.status, time.Since(start))
		}
	})
}
