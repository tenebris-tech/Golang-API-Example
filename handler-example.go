//
// Copyright (c) 2021 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type ExampleResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	ID     string `json:"id"`
}

// ExampleHandler accepts an 'id' variable and echos it back
func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	var resp ExampleResponse

	// Get parameter
	vars := mux.Vars(r)
	id := vars["id"]

	// Create example response
	resp.Status = "ok"
	resp.Code = 200
	resp.ID = id

	// Set reply header
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Send reply
	w.WriteHeader(resp.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("JSON encode error in ExampleHandler: %s\n", err.Error())
	}
}
