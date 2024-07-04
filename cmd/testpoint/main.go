package main

import (
	"flag"
	"fmt"
	"os"
	"testpoint/cmd/testpoint/send"
)

var version = "undefined"

var commands = map[string]func(){
	"send": send.Command,
	"compare": func() {
		fmt.Println("not implemented")
	},
	"version": func() {
		fmt.Println("testpoint version", version)
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

Use "testpoint help <command>" for more information about a command.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	command := os.Args[1]

	if command == "help" {
		if len(os.Args) < 3 {
			usage()
			return
		}
		runCommand(os.Args[2])
		return
	}

	runCommand(command)
}

func runCommand(command string) {
	if f, ok := commands[command]; ok {
		f()
	} else {
		fmt.Fprintf(flag.CommandLine.Output(), "testpoint %s: unknown command\nRun 'testpoint help' for usage.\n", command)
		os.Exit(2)
	}
}
