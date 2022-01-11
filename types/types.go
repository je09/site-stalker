package types

import "net"

type Site struct {
	Status int
	Domain string
	Address string
	Ips []net.IP
	Header string
	Body string
}
