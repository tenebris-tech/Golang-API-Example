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
const ProductVersion = "0.0.2"

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

	a := api.New()
	a.ListenPort = 8080
	a.HTTPTimeout = 30
	a.HTTPIdleTimeout = 60
	a.MaxConcurrent = 100
	a.DownFile = string(os.PathSeparator) + "down.txt"

	fmt.Printf("%s %s starting HTTP server on port %d\n", ProductName, ProductVersion, a.ListenPort)
	err := a.Start()
	if err != nil {
		fmt.Println("Error starting HTTP server: " + err.Error())
	}
}

// AppCleanup provides a graceful exit point
func AppCleanup() {

	// Log exit
	fmt.Println("HTTP server stopping")

	// Exit
	os.Exit(0)
}
