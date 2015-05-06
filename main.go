package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/BruteForceFencer/bff/config"
	"github.com/BruteForceFencer/bff/controlserver"
	"github.com/BruteForceFencer/bff/dashboard"
	"github.com/BruteForceFencer/bff/globals"
	"github.com/BruteForceFencer/bff/hitcounter"
	"github.com/BruteForceFencer/bff/store"
	"github.com/BruteForceFencer/bff/version"
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
	var err error
	globals.Configuration, err = config.ReadConfig(*configFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "configuration error:", err)
		os.Exit(1)
	}
}

func initialize() {
	// Create the hit counter
	directions := make([]hitcounter.Direction, 0)
	for _, dir := range globals.Configuration.Directions {
		shardMap := store.NewShardMap(int64(dir.MaxTracked))
		shardMap.Type = dir.Typ

		directions = append(directions, hitcounter.Direction{
			Store:       shardMap,
			Name:        dir.Name,
			CleanUpTime: dir.CleanUpTime,
			MaxHits:     dir.MaxHits,
			WindowSize:  dir.WindowSize,
		})
	}
	globals.HitCounter = hitcounter.NewHitCounter(directions)

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
	return globals.HitCounter.HandleRequest(req.Direction, req.Value)
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
