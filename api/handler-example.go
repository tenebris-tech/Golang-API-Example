//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ExampleHandler accepts an optional 'id' variable and echos it back
// This is an example of a handler that can receive a variable in the URL or not
// Note that two routes are defined in routes.go, one with the variable and one without
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	var resp Response

	// Get parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Create example response
	resp.Status = "ok"
	resp.Code = http.StatusOK
  
	if id == "" {
		resp.Details = "no ID received"
	} else {
		resp.Details = fmt.Sprintf("received ID %s", id)
	}

	// Send Response
	respond(w, resp, "example")
}
