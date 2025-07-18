package zeroconf

import "net"

// ServerOption fills the option struct to configure server behavior.
type ServerOption func(*serverOpts)

type serverOpts struct {
	listenOn IPType
	ifaces   []net.Interface
}

// SelectServerIPTraffic selects the type of IP packets (IPv4, IPv6, or both) this
// server instance listens for.
func SelectServerIPTraffic(t IPType) ServerOption {
	return func(o *serverOpts) {
		o.listenOn = t
	}
}

// SelectServerIfaces selects the interfaces for the server to use
func SelectServerIfaces(ifaces []net.Interface) ServerOption {
	return func(o *serverOpts) {
		o.ifaces = ifaces
	}
}
