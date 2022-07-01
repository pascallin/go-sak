package network

import (
	"net"
	// "errors"
)

func Dial(protocol, addr string) error {
	return DialWithDialer(&net.Dialer{}, protocol, addr)
}

func DialWithDialer(dialer *net.Dialer, protocol, addr string) error {
	// if protocol != "rtmp" {
	// 	return errors.New("Unknown protocol")
	// }

	_, err := dialer.Dial("tcp", addr)
	if err != nil {
		return err
	}

	return nil
}
