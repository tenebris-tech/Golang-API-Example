//
// Copyright (c) 2021 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"health",
		"GET",
		"/health",
		HealthHandler,
	},
	Route{
		"ip",
		"GET",
		"/example/{id}",
		ExampleHandler,
	},
}
