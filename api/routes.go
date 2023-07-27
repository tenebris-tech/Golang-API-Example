//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (c *Config) getRoutes() Routes {
	return Routes{
		Route{
			"health",
			"GET",
			"/health",
			c.HealthHandler,
		},
		Route{
			"ip",
			"GET",
			"/example",
			c.ExampleHandler,
		},
		Route{
			"ip",
			"GET",
			"/example/{id}",
			c.ExampleHandler,
		},
	}
}
