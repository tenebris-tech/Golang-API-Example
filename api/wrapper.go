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

// Wrapper implements a custom logger by wrapping the handler
// It also sets headers on every request to prevent caching
func Wrapper(destinationHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the start time and source IP
		startTime := time.Now()
		src := getIP(r)

		// Set headers to prevent caching
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Service the request
		sw := statusWriter{ResponseWriter: w}
		destinationHandler.ServeHTTP(&sw, r)

		// Get duration of request
		duration := time.Since(startTime)

		// Remove parameters from URI to avoid logging confidential information
		uri := strings.Split(r.RequestURI, "?")[0]

		// Add code here to send the log event somewhere other than stdout
		fmt.Printf("%s %s %s %d %f\n", src, r.Method, uri, sw.status, duration.Seconds())
	})
}

// JSONEncodeError is called if a JSON encode error occurs
func JSONEncodeError(hName string, err error) {
	// Add code here to send the log event somewhere other than stdout
	fmt.Printf("JSON encode error in %s handler: %s", hName, err.Error())
}
