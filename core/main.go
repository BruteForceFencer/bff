package main

import (
	"flag"
	"fmt"
	"github.com/BruteForceFencer/core/config"
	"github.com/BruteForceFencer/core/controlserver"
	"github.com/BruteForceFencer/core/dashboard"
	"github.com/BruteForceFencer/core/globals"
	"github.com/BruteForceFencer/core/hitcounter"
	"github.com/BruteForceFencer/core/version"
	"os"
	"os/signal"
	"runtime"
)

func configure() {
	// Setup multithreading
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse flags
	configFilename := flag.String("c", "config.json", "the name of the configuration file")
	displayVersion := flag.Bool("version", false, "display the version number")
	flag.Parse()

	// Display version number
	if *displayVersion {
		version.PrintVersion()
		os.Exit(0)
	}

	// Read the configuration
	var errs []error
	globals.Configuration, errs = config.ReadConfig(*configFilename)
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, "configuration error:", err)
		}

		os.Exit(1)
	}
}

func initialize() {
	// Create the hit counter
	globals.HitCounter = hitcounter.NewHitCounter(
		globals.Configuration.Directions,
		globals.Configuration.Logger,
	)

	// Create the server
	globals.Server = controlserver.New()
	globals.Server.HandleFunc = routeRequest

	for _, source := range globals.Configuration.AcceptedSources {
		globals.Server.AcceptedSources[source] = true
	}

	// Create the dashboard
	if globals.Configuration.DashboardAddress != "" {
		globals.Dashboard = dashboard.New(
			globals.Configuration,
			globals.HitCounter,
		)
	}
}

func start() {
	go globals.Server.ListenAndServe(
		globals.Configuration.ListenType,
		globals.Configuration.ListenAddress,
	)

	if globals.Dashboard != nil {
		go globals.Dashboard.ListenAndServe()
	}
}

func routeRequest(req *controlserver.Request) bool {
	if req.Type == controlserver.HitRequest {
		return globals.HitCounter.HandleRequest(req.Direction, req.Value)
	} else if req.Type == controlserver.CommandRequest && req.Direction == "die" {
		globals.Server.Close()
		os.Exit(0)
	}

	return true
}

func main() {
	configure()
	initialize()
	start()

	fmt.Fprintln(os.Stderr, "The server is running.")

	// Capture interrupt signal so that the server closes properly
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	<-interrupts

	globals.Server.Close()
}
