//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"encoding/json"
	"net/http"
)

// Response provides a consistent format for all API responses
// Data is used to hold an appropriate structure
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Details string      `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// respond sends an HTTP response and logs any JSON encoding errors
// This approach allows headers to be set consistently
func respond(w http.ResponseWriter, resp Response, caller string) {

	// Set reply headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "-1")

	// Send reply
	w.WriteHeader(resp.Code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		JSONEncodeError(caller, err)
	}
}
