//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"net/http"
)

// Handler404 returns a 404 error
//
//goland:noinspection GoUnusedParameter
func Handler404(w http.ResponseWriter, r *http.Request) {
	var resp Response

	// Create error response
	resp.Status = "error"
	resp.Code = http.StatusNotFound
	resp.Details = "object does not exist"

	// Send Response
	respond(w, resp, "404")
}
