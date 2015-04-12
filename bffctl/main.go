// Command bffctl is a controller for BFF on Unix based systems.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/BruteForceFencer/bff/core/controlserver"
)

const usage = `BFF Controller
Copyright (C) James Hall 2015.

Usage: bffctl <command>

Commands:
  start     Launches the bff daemon.
  kill      Kills the bff daemon.`

var (
	installDir string
	lastConfig string
)

func main() {
	setupDirectories()

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "start":
		if err := start(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("The daemon has been launched.")
		}
	case "kill":
		if err := kill(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("The daemon has been killed.")
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		printUsage()
	}
}

func setupDirectories() {
	var execPath string
	if os.Args[0][0] == '/' {
		execPath = os.Args[0]
	} else {
		wd, _ := os.Getwd()
		execPath = path.Join(wd, os.Args[0])
	}

	installDir = path.Dir(path.Dir(execPath))
	lastConfig = path.Join(installDir, "var", "last_config.json")
}

func printUsage() {
	fmt.Println(usage)
}

func start() error {
	conn, err := connect()
	if err == nil {
		conn.Close()
		return fmt.Errorf("An instance of BFF is already running.")
	}

	cmd := exec.Command(
		path.Join(installDir, "util", "daemonize"),
		path.Join(installDir, "util", "bffcore"),
		"-c",
		path.Join(installDir, "config.json"),
	)

	if err := cmd.Run(); err != nil {
		return err
	}

	if err := copyLastConfig(); err != nil {
		return fmt.Errorf(
			"Unable to write to %s, daemon will have to be killed manually.",
			lastConfig,
		)
	}

	return nil
}

func kill() error {
	conn, err := connect()
	if err != nil {
		return fmt.Errorf("BFF is not currently running.")
	}
	defer conn.Close()

	req := controlserver.Request{
		Type:      controlserver.CommandRequest,
		Direction: "die",
	}

	enc := json.NewEncoder(conn)
	if err := enc.Encode(req); err != nil {
		return fmt.Errorf("Unable to kill daemon.")
	}

	return nil
}

func copyLastConfig() error {
	dest, err := os.Create(lastConfig)
	if err != nil {
		return err
	}
	defer dest.Close()

	src, err := os.Open(path.Join(installDir, "config.json"))
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	return err
}
