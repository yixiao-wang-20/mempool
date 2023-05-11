package common

import (
	"errors"
	"net"
)

func NetToPortInt(addr net.Addr) (ip string, port int, err error) {
	switch addr := addr.(type) {
	case *net.UDPAddr:
		ip = addr.IP.String()
		port = addr.Port
	case *net.TCPAddr:
		ip = addr.IP.String()
		port = addr.Port
	default:
		err = errors.New("unknown IP type")
	}

	return
}
