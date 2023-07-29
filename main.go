//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Golang-API-Example/api"
)

const ProductName = "golang-api-example"
const ProductVersion = "0.0.7"

func main() {

	// Create the API object
	// This is done first so that the API object can be used to gracefully
	// shut down the server when a signal is received
	a := api.New()

	// Setup signal catching
	signals := make(chan os.Signal, 1)

	// Catch signals
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// method invoked upon receiving signal
	go func() {
		for {
			s := <-signals
			log.Printf("Received signal: %v", s)
			log.Printf("Asking the API server to stop")
			err := a.Stop()
			if err != nil {
				log.Printf("Error stopping the API server: %s", err.Error())
				log.Printf("Forcing exit")
				cleanup()
			}
		}
	}()

	// Set API parameters
	a.Listen = "127.0.0.1:8080"
	a.HTTPTimeout = 30
	a.HTTPIdleTimeout = 60
	a.MaxConcurrent = 100
	a.DownFile = string(os.PathSeparator) + "down.txt"
	a.Debug = true
	// a.TLS = true
	// a.TLSCertFile = "cert.pem"
	// a.TLSKeyFile = "key.pem"

	log.Printf("%s %s starting API server on %s", ProductName, ProductVersion, a.Listen)

	// Start the API
	// If the application needs to do other work, Start() could be launched as a goroutine to run in the background
	err := a.Start()
	if err != nil {
		// Server returns an error even if it shut down gracefully
		// If the error is not a server closed error, return it
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Error starting the API server: %s", err.Error())
		} else {
			log.Printf("API server shut down gracefully")
		}
	}
	cleanup()
}

// cleanup is the graceful exit point
func cleanup() {

	// Perform any cleanup here

	// Exit
	log.Printf("%s %s exiting", ProductName, ProductVersion)
	os.Exit(0)
}
