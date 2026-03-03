package amqp09

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	defaultHeartbeat         = 10 * time.Second
	defaultLocale            = "en_US"
	defaultConnectionTimeout = 30 * time.Second
)

// ShuffleDialer tries all the resolved IP addresses in a random order until succeeds or all the IP connection fails
func ShuffleDialer(network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host:port: %w", err)
	}
	ips, err := net.LookupIP(host) // Fresh DNS call
	if err != nil {
		return nil, fmt.Errorf("failed to lookup IPs for host: %w", err)
	}

	// Randomize to spread the load across nodes
	rand.Shuffle(len(ips), func(i, j int) { ips[i], ips[j] = ips[j], ips[i] })

	errs := make([]error, len(ips))
	for i, ip := range ips {
		conn, err := net.DialTimeout(network, net.JoinHostPort(ip.String(), port), defaultConnectionTimeout)
		if err == nil {
			return conn, nil
		}
		errs[i] = fmt.Errorf("%s: %w", ip.String(), err)
	}
	return nil, fmt.Errorf("all node IPs unreachable: %w", errors.Join(errs...))
}
