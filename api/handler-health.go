//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// HealthHandler implements a health check for load balancers, etc.
//
//goland:noinspection GoUnusedParameter
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	var resp Response

	// Set reply header
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Check for presence of /tmp/down
	if _, err := os.Stat(downFile); err == nil {
		// exists -- send status down and 503
		resp.Status = "down"
		resp.Code = 503
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		// does not exist - send ok and 200
		resp.Status = "ok"
		resp.Code = 200
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Printf("JSON encode error in HealthHandler: %s\n", err.Error())
	}
}
