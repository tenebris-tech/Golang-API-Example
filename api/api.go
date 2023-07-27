//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Listen          string
	HTTPTimeout     int
	HTTPIdleTimeout int
	MaxConcurrent   int
	DownFile        string
	StrictSlash     bool
	TLS             bool
	TLSCertFile     string
	TLSKeyFile      string
	Debug           bool
}

// Store the down file in a global variable
var downFile = ""

// New returns a new Config struct with default values
func New() Config {

	// Return default configuration
	return Config{
		Listen:          "127.0.0.1:8080",
		HTTPTimeout:     60,
		HTTPIdleTimeout: 60,
		MaxConcurrent:   100,
		DownFile:        "",
		StrictSlash:     true,
		TLS:             false,
		TLSCertFile:     "",
		TLSKeyFile:      "",
		Debug:           false,
	}
}

// Start starts the API
func (c *Config) Start() error {

	// Update the downFile for access by the handler
	downFile = c.DownFile

	// Instantiate the HTTP router
	router := c.newRouter()

	// Add catch all and not found handler
	router.PathPrefix("/").Handler(c.Wrapper(http.HandlerFunc(c.Handler404)))
	router.NotFoundHandler = c.Wrapper(http.HandlerFunc(c.Handler404))
	router.MethodNotAllowedHandler = c.Wrapper(http.HandlerFunc(c.Handler405))

	// Create server
	s := &http.Server{
		Addr:              c.Listen,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(c.HTTPTimeout) * time.Second,
		ReadTimeout:       time.Duration(c.HTTPTimeout) * time.Second,
		WriteTimeout:      time.Duration(c.HTTPTimeout) * time.Second,
		IdleTimeout:       time.Duration(c.HTTPIdleTimeout) * time.Second,
	}

	// Add TLS configuration if option is enabled
	if c.TLS {
		if c.TLSCertFile == "" || c.TLSKeyFile == "" {
			return errors.New("TLS cert or key file not specified")
		}

		// Load the cert and key
		cert, err := tls.LoadX509KeyPair(c.TLSCertFile, c.TLSKeyFile)
		if err != nil {
			return err
		}

		// Create the TLS configuration
		tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}}
		tlsConfig.MinVersion = tls.VersionTLS12

		// Add to the HTTP server config
		s.TLSConfig = &tlsConfig
	}

	// Start our customized server
	err := c.listen(s)
	if err != nil {
		return err
	}
	return nil
}

// Create gorilla/mux router and load routes from route.go
func (c *Config) newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(c.StrictSlash)

	// Get routes from routes.go
	routes := c.getRoutes()

	// Iterate through routes
	for _, route := range routes {

		// Register the route, including our wrapper
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(c.Wrapper(route.HandlerFunc))
	}
	return router
}

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

// serve requests using the specified listener (limited or not) and TLS if configured
func (c *Config) serve(srv *http.Server, l net.Listener) error {
	if c.TLS {
		// This will use the previously configured TLS information
		return srv.ServeTLS(l, "", "")
	} else {
		return srv.Serve(l)
	}
}
