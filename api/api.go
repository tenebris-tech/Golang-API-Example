//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	ListenPort      int
	HTTPTimeout     int
	HTTPIdleTimeout int
	MaxConcurrent   int
	DownFile        string
	StrictSlash     bool
}

// Store the down file in a global variable
var downFile = ""

func New() Config {
	// Return default configuration
	return Config{
		ListenPort:      80,
		HTTPTimeout:     60,
		HTTPIdleTimeout: 60,
		MaxConcurrent:   100,
		DownFile:        "",
		StrictSlash:     true,
	}
}

func (c *Config) Start() error {

	// Update the downFile for access by the handler
	// This will create an issue if more than one API is initialized, they can not
	// have separate down files
	downFile = c.DownFile

	// Instantiate HTTP router
	router := c.newRouter()

	// Add catch all and not found handler
	router.PathPrefix("/").Handler(Logger(http.HandlerFunc(Handler404), "404"))
	router.NotFoundHandler = Logger(http.HandlerFunc(Handler404), "404")

	// Validate port
	if c.ListenPort < 1 || c.ListenPort > 65535 {
		return errors.New(fmt.Sprintf("invalid ListenPort: %d\n", c.ListenPort))
	}
	listenPort := strconv.Itoa(c.ListenPort)

	// Create server
	s := &http.Server{
		Addr:              ":" + listenPort,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(c.HTTPTimeout) * time.Second,
		ReadTimeout:       time.Duration(c.HTTPTimeout) * time.Second,
		WriteTimeout:      time.Duration(c.HTTPTimeout) * time.Second,
		IdleTimeout:       time.Duration(c.HTTPIdleTimeout) * time.Second,
	}

	// Start our customized server
	err := listen(s, c.MaxConcurrent)
	if err != nil {
		return err
	}
	return nil
}
