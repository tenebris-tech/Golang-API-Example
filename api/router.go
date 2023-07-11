//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"github.com/gorilla/mux"
)

// Create gorilla/mux router and load routes from route.go
func (c *Config) newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(c.StrictSlash)

	// Get routes from routes.go
	routes := getRoutes()

	// Iterate through routes
	for _, route := range routes {

		// Wrap the handler in our logger
		handler := Wrapper(route.HandlerFunc)

		// Register the route
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
