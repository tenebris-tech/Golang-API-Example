//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
// See LICENSE for further information.
//

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"Golang-API-Example/api"
)

const ProductName = "golang-api-example"
const ProductVersion = "0.0.3"

func main() {

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

	// Create the API object and set parameters
	a := api.New()

	a.Listen = "127.0.0.1:8080"
	a.HTTPTimeout = 30
	a.HTTPIdleTimeout = 60
	a.MaxConcurrent = 100
	a.DownFile = string(os.PathSeparator) + "down.txt"
	a.TLS = false // disabled
	a.TLSCertFile = "cert.pem"
	a.TLSKeyFile = "key.pem"

	fmt.Printf("%s %s starting API server on %s\n", ProductName, ProductVersion, a.Listen)

 	// Start the API
	// If the application needs to do other work, Start() could be launched as a goroutine
	err := a.Start()
	if err != nil {
		fmt.Println("Error starting API server: " + err.Error())
	}
}

// AppCleanup provides a graceful exit point
func AppCleanup() {

	// Log exit
	fmt.Println("API server stopping")

	// Exit
	os.Exit(0)
}
