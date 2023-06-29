//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"crypto/tls"
	"errors"
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
	}
}

// Start starts the API
func (c *Config) Start() error {

	// Update the downFile for access by the handler
	downFile = c.DownFile

	// Instantiate the HTTP router
	router := c.newRouter()

	// Add catch all and not found handler
	router.PathPrefix("/").Handler(Logger(http.HandlerFunc(Handler404)))
	router.NotFoundHandler = Logger(http.HandlerFunc(Handler404))

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
