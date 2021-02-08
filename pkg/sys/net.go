package sys

import (
	"net"
)

// GetIpAddress returns the string formatted ip address of the machine we are running on,
// as best can be determined.
func GetIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}
	// TODO: if the above errors, what other ways can we do this?
	return ""
}
