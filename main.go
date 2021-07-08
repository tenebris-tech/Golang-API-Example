//
// Copyright (c) 2021 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const ProductName = "golang-api-example"
const ProductVersion = "0.0.1"

type Config struct {
	ListenPort      int
	HTTPTimeout     int
	HTTPIdleTimeout int
	MaxConcurrent   int
	DownFile        string
}

// Instantiate global config structure
var config Config

func main() {
	var err error

	// Create configuration
	config.ListenPort = 8080
	config.HTTPTimeout = 60
	config.HTTPIdleTimeout = 60
	config.MaxConcurrent = 100
	config.DownFile = string(os.PathSeparator) + "down.txt"

	// Get listen port
	if config.ListenPort < 1 || config.ListenPort > 65535 {
		fmt.Printf("Invalid ListenPort: %d\n", config.ListenPort)
		os.Exit(1)
	}
	listenPort := strconv.Itoa(config.ListenPort)

	// Setup signal catching
	signals := make(chan os.Signal, 1)

	// Catch signals
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// method invoked upon seeing signal
	go func() {
		for {
			s := <-signals
			fmt.Printf("Received signal: %v\n", s)
			AppCleanup()
		}
	}()

	// Instantiate HTTP router
	router := newRouter()

	// Create server
	s := &http.Server{
		Addr:              ":" + listenPort,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(config.HTTPTimeout) * time.Second,
		ReadTimeout:       time.Duration(config.HTTPTimeout) * time.Second,
		WriteTimeout:      time.Duration(config.HTTPTimeout) * time.Second,
		IdleTimeout:       time.Duration(config.HTTPIdleTimeout) * time.Second,
	}

	fmt.Printf("%s %s starting HTTP server on port %s\n", ProductName, ProductVersion, listenPort)

	// Start our customized server
	err = listen(s, config.MaxConcurrent)
	if err != nil {
		fmt.Printf("Failure starting HTTP server: %s\n", err.Error())
		os.Exit(1)
	}
}

// Create gorilla/mux router and load routes from route.go
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

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
		return srv.Serve(listener)
	}

	limited := netutil.LimitListener(listener, maxConcurrent)
	return srv.Serve(limited)
}

// AppCleanup provides a graceful exit point
func AppCleanup() {

	// Log exit
	fmt.Println("HTTP server stopping")

	// Exit
	os.Exit(0)
}
