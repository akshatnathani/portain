
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/akshatnathani/portain/internal/process"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: port <inspect> <port>")
		os.Exit(1)
	}

	command := os.Args[1]

	port, err := strconv.Atoi(os.Args[2])
	if err != nil || port < 1 || port > 65535 {
		fmt.Println("Invalid port")
		os.Exit(1)
	}

	switch command {
	case "inspect":
		processes, err := process.FindByPort(port)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if len(processes) == 0 {
			fmt.Printf("No listening process found on port %d\n", port)
			return
		}

		for _, p := range processes {
			fmt.Printf("Port: %d\n", port)
			fmt.Printf("PID: %d\n", p.PID)
			fmt.Printf("Name: %s\n", p.Name)
			fmt.Printf("Command: %s\n", p.Command)
		}

	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
