package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akshatnathani/portain/internal/ports"
	"github.com/akshatnathani/portain/internal/process"
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

	case "free":
		processes, err := process.FindByPort(port)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if len(processes) == 0 {
			fmt.Printf("No process found listening on port %d\n", port)
			return
		}

		if len(processes) > 1 {
			fmt.Printf("Multiple processes found on port %d. Cannot safely terminate automatically.\n", port)
			for _, p := range processes {
				fmt.Printf("PID: %d | Name: %s\n", p.PID, p.Name)
			}
			os.Exit(1)
		}

		p := processes[0]

		fmt.Printf("Process found:\n")
		fmt.Printf("PID:     %d\n", p.PID)
		fmt.Printf("Name:    %s\n", p.Name)
		fmt.Printf("Command: %s\n", p.Command)

		fmt.Print("\nAre you sure you want to terminate this process? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil && len(answer) == 0 {
			fmt.Fprintln(os.Stderr, "Could not read confirmation:", err)
			os.Exit(1)
		}

		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("Operation cancelled.")
			return
		}

		result, err := process.Terminate(p.PID, port, 5*time.Second)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Termination failed:", err)
			os.Exit(1)
		}

		if result.Terminated && result.PortFreed {
			fmt.Printf("Successfully terminated process %d. Port %d is now available.\n", result.PID, port)
		} else {
			fmt.Printf("Process termination could not be fully verified. PID: %d, Port freed: %t\n", result.PID, result.PortFreed)
			os.Exit(1)
		}
	case "find":
		found := false

		for candidate := port; candidate <= 65535; candidate++ {
			address := fmt.Sprintf("127.0.0.1:%d", candidate)

			listener, err := net.Listen("tcp", address)
			if err != nil {
				continue
			}

			listener.Close()

			fmt.Printf("Next available port: %d\n", candidate)
			found = true
			break
		}

		if !found {
			fmt.Printf("No available port found from %d onwards\n", port)
		}

	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("Portain - Linux port management utility")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  portain check <port>")
	fmt.Println("  portain inspect <port>")
	fmt.Println("  portain find <starting-port>")
	fmt.Println("  portain free <port>")
}
