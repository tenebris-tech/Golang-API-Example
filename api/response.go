//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

// Response provides a consistent format for all API responses
// Data is used to hold an appropriate structure
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Details string      `json:"details,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
