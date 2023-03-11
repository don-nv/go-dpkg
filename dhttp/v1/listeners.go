package dhttp

import (
	"fmt"
	"net"
)

const networkTCP = "tcp"

func listenTCP(address string) (net.Listener, error) {
	return listen(networkTCP, address)
}

func listen(network, address string) (net.Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("listening on %q (%q): %w", address, network, err)
	}

	return listener, nil
}
