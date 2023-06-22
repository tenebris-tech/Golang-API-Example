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

// Replacement for ListenAndServe that implements concurrent session limit
// using netutil.LimitListener. If maxConcurrent is 0, bypass the limit.
func (c *Config) listen(srv *http.Server) error {

	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if c.MaxConcurrent == 0 {
		return c.serve(srv, listener)
	}

	limited := netutil.LimitListener(listener, c.MaxConcurrent)
	return c.serve(srv, limited)
}

// Start server
func (c *Config) serve(srv *http.Server, l net.Listener) error {
	if c.TLS {
		return srv.ServeTLS(l, c.TLSCertFile, c.TLSKeyFile)
	} else {
		return srv.Serve(l)
	}
}
