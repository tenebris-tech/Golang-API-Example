//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import "net/http"

//goland:noinspection GoUnusedParameter
func (c *Config) Handler401(w http.ResponseWriter, r *http.Request) {
	c.status4xx(w, http.StatusUnauthorized, "not authorized")
}

//goland:noinspection GoUnusedParameter
func (c *Config) Handler404(w http.ResponseWriter, r *http.Request) {
	c.status4xx(w, http.StatusNotFound, "object does not exist")
}

//goland:noinspection GoUnusedParameter
func (c *Config) Handler405(w http.ResponseWriter, r *http.Request) {
	c.status4xx(w, http.StatusMethodNotAllowed, "method not allowed")
}

// status4xx returns a 4xx error
func (c *Config) status4xx(w http.ResponseWriter, code int, message string) {
	var resp = Response{Details: message, Status: "error", Code: code}
	c.respond(w, resp, "4xx")
}
