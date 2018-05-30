package main

import (
	"os"
	"log"
	"github.com/RiverbedTechnology/sdp-ztp/ZTP/sdp-ztp/util"
	"os/exec"
	"strings"
)

const (
	port  = "8090"
)
var interfaceName string

func startZTPService() string {
	log.Println("Configuring files and downloading artifacts for ZTP services")
	status := configureZtp()
	if status {
		log.Println("Starting ZTP services with new configs")
		go startDnsmasqServer()
		startMatchbox()
		return "running"
	}
	return "failed"
}

func startDnsmasqServer() {

	cmd := exec.Command("/bin/sh", "-c", "service dnsmasq restart") //nolint
	opcmd, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Dnsmasq run cmd error '%s'", err)
	}
	log.Println("opt of Dnsmasq: ", string(opcmd))
}

func startMatchbox() {
	ipAddrWithCIDR := util.GetIPv4ForInterfaceName(interfaceName).String()
	ipAddr := util.GetIPFromIPwithCIDR(ipAddrWithCIDR)
	cmd := exec.Command("/bin/sh", "-c", "matchbox -address "+ipAddr+":"+port+" > /dev/null 2>&1 &") //nolint
	opcmd, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Matchbox run cmd error '%s'", err)
	}
	log.Println("opt of matchbox: ", string(opcmd))

}

func configureZtp() bool{
	ipAddrWithCidr := util.GetIPv4ForInterfaceName(interfaceName).String()
	cidr := util.GetCIDRFromIPwithCIDR(util.GetIPv4ForInterfaceName(interfaceName).String())
	ipAddr := util.GetIPFromIPwithCIDR(ipAddrWithCidr)
	sub := util.SubnetCalculator(ipAddr,cidr)
	config :=  util.Config {
		NetworkID:		sub.GetNetworkPortion(),
		NetMask:		sub.GetSubnetMask(),
		//IPRange:		strings.Replace(strings.Replace(sub.GetIPAddressRange(), "0 ", "1 ", -1), ".255", ".254", -1),
		IPRange:		strings.Replace(strings.Replace(sub.GetIPAddressRange(), "0,", "1,", -1), ".255", ".254", -1),
		BroadcastIP:	sub.GetBroadcastAddress(),
		MatchboxPort:	port,
		IP:				ipAddr,
		Interface:		interfaceName,
	}

	err := config.GenerateTemplates()
	if err != nil {
		log.Fatal("Error while generating config file", err)
		return false
	}

	return true
}

func main() {
	interfaceName = os.Getenv("SDP_NETWORK_INTERFACE")
	if interfaceName == "" {
		log.Fatalf("Provide \"SDP_NETWORK_INTERFACE\" environment variable")
	}
	status := startZTPService()
	if(status == "failed") {
		log.Fatalln("Failed to start ZTP services")
	} else if (status == "running") {
		log.Println("ZTP services running successfully")
	}
}
