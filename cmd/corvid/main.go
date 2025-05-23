package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/CartConnoisseur/corvid/srv"
	"github.com/godbus/dbus/v5"
)

func main() {
	if len(os.Args) < 2 {
		server()
	}

	switch os.Args[1] {
	case "server":
		server()
	case "dismiss":
		if len(os.Args) < 3 {
			fmt.Println("dismiss command requires positional argument 'id'")
			os.Exit(1)
		}

		id, err := strconv.ParseInt(os.Args[2], 0, 32)
		if err != nil || id <= 0 {
			fmt.Printf("invalid value for positional argument 'id' (must be u64): %s\n", os.Args[2])
			os.Exit(1)
		}

		call("Dismiss", uint32(id))
	case "dismiss-all":
		call("DismissAll")
	case "help":
		fmt.Printf("subcommands:\n")
		fmt.Printf("  server (default) - run the server\n")
		fmt.Printf("  dismiss <id>     - dismiss specific notification\n")
		fmt.Printf("  dismiss-all      - dismiss all notifications\n")
		fmt.Printf("\n")
		fmt.Printf("environment vars:\n")
		fmt.Printf("  CORVID_DEFAULT_EXPIRATION - default notification expiration in\n")
		fmt.Printf("                              ms. -1 = never (default: 5000)\n")
		fmt.Printf("  CORVID_SORT_DIRECTION     - notification sort direction, either\n")
		fmt.Printf("                              NEWEST_FIRST or OLDEST_FIRST\n")
		fmt.Printf("                              (default: NEWEST_FIRST)\n")
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		fmt.Printf("run corvid help for command list\n")
		os.Exit(1)
	}
}

func server() {
	defaultExpiration := getEnvInt("CORVID_DEFAULT_EXPIRATION", 5000)
	sortDirection := 1

	switch os.Getenv("CORVID_SORT_DIRECTION") {
	case "NEWEST_FIRST":
		sortDirection = 1
	case "OLDEST_FIRST":
		sortDirection = -1
	}

	srv.Start(defaultExpiration, sortDirection)
	select {}
}

func getEnvInt(key string, fallback int) int {
	str := os.Getenv(key)

	if len(str) == 0 {
		return fallback
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}

	return value
}

func call(name string, args ...interface{}) error {
	const CORVID_DBUS_OBJECT = "/sh/cxl/Corvid"
	const CORVID_DBUS_NAME = "sh.cxl.Corvid"

	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	call := conn.Object(CORVID_DBUS_NAME, CORVID_DBUS_OBJECT).Call(CORVID_DBUS_NAME+"."+name, 0, args...)
	if call.Err != nil {
		return call.Err
	}

	return nil
}
