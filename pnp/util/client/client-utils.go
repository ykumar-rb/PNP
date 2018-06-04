package client

import (
	"runtime"
	"log"
	"net"
	"github.com/ZTP/pnp/executor"
	proto "github.com/ZTP/pnp/pnp-proto"
)

func PopulateClientDetails(ifname string) (clientInfo proto.ClientInfo) {
	archType := runtime.GOARCH
	osType := runtime.GOOS

	getOSFlavorCmd := "lsb_release -a | grep Description | awk -F':' '{print $2}'"
	osFlavor, err := executor.ExecuteCommand(getOSFlavorCmd)
	if err != nil {
		log.Fatalf("Error while getting OS type: %v", err)
	}

	MACAddr := GetMACForInterfaceName(ifname)

	clientInfo = proto.ClientInfo{OsType: osType, ArchType: archType, OsFlavor: osFlavor, MACAddr: MACAddr}
	return
}

func GetMACForInterfaceName(ifname string) (string) {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == ifname {
			mac := inter.HardwareAddr.String()
			return mac
		}
	}
	return ""
}

