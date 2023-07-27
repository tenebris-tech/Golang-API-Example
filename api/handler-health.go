//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"net/http"
	"os"
)

// HealthHandler implements a health check for load balancers, etc.
//
//goland:noinspection GoUnusedParameter
func (c *Config) HealthHandler(w http.ResponseWriter, r *http.Request) {
	var resp Response

	// Check for presence of the file that indicates the server is down
	if _, err := os.Stat(downFile); err == nil {
		// exists -- send status down and 503
		resp.Status = "down"
		resp.Code = http.StatusServiceUnavailable
		resp.Details = "server is shutting down"
	} else {
		// does not exist - send ok and 200
		resp.Status = "ok"
		resp.Code = http.StatusOK
		resp.Details = "health check ok"
	}

	// Send Response
	c.respond(w, resp, "health")
}
