//
// Copyright (c) 2021 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response404 struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

// Handler404 accepts an 'id' variable and echos it back
func Handler404(w http.ResponseWriter, r *http.Request) {
	var resp Response404

	// Create error response
	resp.Status = "ok"
	resp.Code = 200

	// Set reply header
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Send reply
	w.WriteHeader(resp.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("JSON encode error in ExampleHandler: %s\n", err.Error())
	}
}
