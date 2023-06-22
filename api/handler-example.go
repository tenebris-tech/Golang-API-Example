//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ExampleHandler accepts an 'id' variable and echos it back
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	var resp Response

	// Get parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Create example response
	resp.Status = "ok"
	resp.Code = 200
	resp.Details = fmt.Sprintf("received ID %s", id)

	// Set reply header
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Send reply
	w.WriteHeader(resp.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("JSON encode error in ExampleHandler: %s\n", err.Error())
	}
}
