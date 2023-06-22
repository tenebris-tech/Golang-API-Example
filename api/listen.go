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
// using netutil.LimitListener. If maxConcurrent is 0, bypass limit.
func listen(srv *http.Server, maxConcurrent int) error {

	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	if maxConcurrent == 0 {
		// Note - for TLS see ServeTLS
		return srv.Serve(listener)
	}

	limited := netutil.LimitListener(listener, maxConcurrent)
	// Note - for TLS see ServeTLS
	return srv.Serve(limited)
}
