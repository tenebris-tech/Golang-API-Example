//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"golang.org/x/net/netutil"
	"net"
	"net/http"
)

// listen is a replacement for ListenAndServe that implements a concurrent session limit
// using netutil.LimitListener. If maxConcurrent is 0, no limit is imposed.
func (c *Config) listen(srv *http.Server) error {

	// Get listen address, default to ":http"
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	// Create listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// If maxConcurrent is 0, bypass limit
	if c.MaxConcurrent == 0 {
		return c.serve(srv, listener)
	}

	// Create server with limited listener
	limited := netutil.LimitListener(listener, c.MaxConcurrent)
	return c.serve(srv, limited)
}

// Start server using the specified listener (limited or not) and TLS if configured
func (c *Config) serve(srv *http.Server, l net.Listener) error {
	if c.TLS {
		// This will use the previously configured TLS information
		return srv.ServeTLS(l, "", "")
	} else {
		return srv.Serve(l)
	}
}
