
package process

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func findSocketInodes(port int) (map[string]bool, error) {
	files := []string{
		"/proc/net/tcp",
		"/proc/net/tcp6",
	}

	targets := make(map[string]bool)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(f)

		// Skip the header.
		if scanner.Scan() {
			for scanner.Scan() {
				fields := strings.Fields(scanner.Text())

				if len(fields) < 10 {
					continue
				}

				if fields[3] != "0A" {
					continue
				}

				localAddress := fields[1]
				parts := strings.Split(localAddress, ":")

				if len(parts) != 2 {
					continue
				}

				portNumber, err := strconv.ParseInt(parts[1], 16, 32)
				if err != nil || int(portNumber) != port {
					continue
				}

				targets[fields[9]] = true
			}
		}

		f.Close()
	}

	return targets, nil
}

func FindByPort(port int) ([]ProcessInfo, error) {
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %d", port)
	}

	inodes, err := findSocketInodes(port)
	if err != nil {
		return nil, err
	}

	if len(inodes) == 0 {
		return nil, nil
	}

	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("cannot read /proc: %w", err)
	}

	var results []ProcessInfo
	seen := make(map[int]bool)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}

		fdDir := filepath.Join("/proc", entry.Name(), "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}

		for _, fd := range fds {
			link := filepath.Join(fdDir, fd.Name())

			target, err := os.Readlink(link)
			if err != nil {
				continue
			}

			if !strings.HasPrefix(target, "socket:[") {
				continue
			}

			inode := strings.TrimSuffix(
				strings.TrimPrefix(target, "socket:["),
				"]",
			)

			if !inodes[inode] || seen[pid] {
				continue
			}

			info, err := Inspect(pid)
			if err != nil {
				continue
			}

			results = append(results, *info)
			seen[pid] = true
			break
		}
	}

	return results, nil
}
