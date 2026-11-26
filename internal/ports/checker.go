
package ports

import (
	"fmt"
	"net"
)

func IsAvailable(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, fmt.Errorf("invalid port: %d", port)
	}

	address := fmt.Sprintf("127.0.0.1:%d", port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false, nil
	}

	defer listener.Close()
	return true, nil
}