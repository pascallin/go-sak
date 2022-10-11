package lib

import (
	"log"
	"net"

	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = MachineID
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Fatal("sonyflake not created")
	}
}

func GenerateID() (uint64, error) {
	return sf.NextID()
}

func MachineID() (uint16, error) {
	ip := getOutboundIP()
	return ipToMachineID(ip), nil
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func ipToMachineID(ip net.IP) uint16 {
	return uint16(ip[2])<<8 + uint16(ip[3])
}
