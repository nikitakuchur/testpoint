package main

import (
	"flag"
	"fmt"
	"os"
	"testpoint/cmd/testpoint/compare"
	"testpoint/cmd/testpoint/send"
)

var version = "undefined"

var commands = map[string]func(){
	"send":    send.Command,
	"compare": compare.Command,
	"version": func() {
		fmt.Println("testpoint version", version)
	},
	"help": func() {
		usage()
	},
}

func usage() {
	fmt.Println(`Testpoint is a simple CLI tool for testing REST endpoints.

Usage: 

	testpoint <command> [arguments]

The commands are:

	send		send prepared requests to specified REST endpoints
	compare		compare responses and generete a report
	version		print Testpoint version

Use "testpoint <command> --help" for more information about a command.`)
}

func runCommand(command string) {
	if f, ok := commands[command]; ok {
		f()
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "testpoint %s: unknown command\nRun 'testpoint help' for usage.\n", command)
		os.Exit(2)
	}
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	command := os.Args[1]

	// for flag parsing during command
	os.Args = os.Args[1:]

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", command)
		flag.PrintDefaults()
	}

	runCommand(command)
}
