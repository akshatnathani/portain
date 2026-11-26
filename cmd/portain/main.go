
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/akshatnathani/portain/internal/process"
	"github.com/akshatnathani/portain/internal/ports"
)

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	command := os.Args[1]

	port, err := strconv.Atoi(os.Args[2])
	if err != nil || port < 1 || port > 65535 {
		fmt.Fprintln(os.Stderr, "Error: invalid port")
		os.Exit(1)
	}

	switch command {
	case "check":
		available, err := ports.IsAvailable(port)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if available {
			fmt.Printf("Port %d is available\n", port)
		} else {
			fmt.Printf("Port %d is occupied\n", port)
		}

	case "inspect":
		processes, err := process.FindByPort(port)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(processes) == 0 {
			fmt.Printf("No listening process found on port %d\n", port)
			return
		}

		for _, p := range processes {
			fmt.Printf("Port:    %d\n", port)
			fmt.Printf("PID:     %d\n", p.PID)
			fmt.Printf("Name:    %s\n", p.Name)
			fmt.Printf("Command: %s\n\n", p.Command)
		}

	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("PortGuard - Linux port management utility")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  portguard check <port>")
	fmt.Println("  portguard inspect <port>")
}
