package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"log"
)

type Ip struct {
	ip          string
	networkSize int
	subnet_mask int
}

func GetIPv4ForInterfaceName(ifname string) (ifaceip *net.IPNet) {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == ifname {
			if addrs, err := inter.Addrs(); err == nil {
				for _, addr := range addrs {
					switch ip := addr.(type) {
					case *net.IPNet:
						if ip.IP.DefaultMask() != nil {
							return (ip)
						}
					}
				}
			}
		}
	}
	log.Fatalln("Check the interface name provided")
	return (nil)
}

func GetCIDRFromIPwithCIDR(ipCidr string) int {
	ipCidrArr := strings.Split(ipCidr, "/")
	cidr := ipCidrArr[1]
	cidrInt, err := strconv.Atoi(cidr)
	if err == nil {
		fmt.Println(err)
	}
	return cidrInt
}

func GetIPFromIPwithCIDR(ipCidr string) string {
	ipCidrArr := strings.Split(ipCidr, "/")
	ip := ipCidrArr[0]
	return ip
}

func SubnetCalculator(ip string, networkSize int) *Ip {

	s := &Ip{
		ip:          ip,
		networkSize: networkSize,
		subnet_mask: 0xFFFFFFFF << uint(32-networkSize),
	}

	return s
}

func (s *Ip) GetNetworkPortion() string {
	return s.networkCalculation("%d", ".")
}

func (s *Ip) networkCalculation(format, separator string) string {
	splits := s.GetIPAddressQuads()
	networkQuards := []string{}
	networkQuards = append(networkQuards, fmt.Sprintf(format, splits[0]&(s.subnet_mask>>24)))
	networkQuards = append(networkQuards, fmt.Sprintf(format, splits[1]&(s.subnet_mask>>16)))
	networkQuards = append(networkQuards, fmt.Sprintf(format, splits[2]&(s.subnet_mask>>8)))
	networkQuards = append(networkQuards, fmt.Sprintf(format, splits[3]&(s.subnet_mask>>0)))

	return strings.Join(networkQuards, separator)
}

func (s *Ip) GetIPAddressQuads() []int {
	splits := strings.Split(s.ip, ".")

	return convertQuardsToInt(splits)
}

func convertQuardsToInt(splits []string) []int {
	quardsInt := []int{}

	for _, quard := range splits {
		j, err := strconv.Atoi(quard)
		if err != nil {
			panic(err)
		}
		quardsInt = append(quardsInt, j)
	}

	return quardsInt
}

func (s *Ip) GetIPAddressRange() string {
	//return strings.Join([]string{s.GetNetworkPortion(), s.GetBroadcastAddress()}, " ")
	return strings.Join([]string{s.GetNetworkPortion(), s.GetBroadcastAddress()}, ",")
}

func (s *Ip) GetBroadcastAddress() string {
	networkQuads := s.GetNetworkPortionQuards()
	numberIPAddress := s.GetNumberIPAddresses()
	networkRangeQuads := []string{}
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[0]&(s.subnet_mask>>24))+(((numberIPAddress-1)>>24)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[1]&(s.subnet_mask>>16))+(((numberIPAddress-1)>>16)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[2]&(s.subnet_mask>>8))+(((numberIPAddress-1)>>8)&0xFF)))
	networkRangeQuads = append(networkRangeQuads, fmt.Sprintf("%d", (networkQuads[3]&(s.subnet_mask>>0))+(((numberIPAddress-1)>>0)&0xFF)))

	return strings.Join(networkRangeQuads, ".")
}

func (s *Ip) GetNetworkPortionQuards() []int {
	return convertQuardsToInt(strings.Split(s.networkCalculation("%d", "."), "."))
}

func (s *Ip) GetNumberIPAddresses() int {
	return 2 << uint(31-s.networkSize)
}

func (s *Ip) GetSubnetMask() string {
	return s.subnetCalculation("%d", ".")
}

func (s *Ip) subnetCalculation(format, separator string) string {
	maskQuards := []string{}
	maskQuards = append(maskQuards, fmt.Sprintf(format, (s.subnet_mask>>24)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (s.subnet_mask>>16)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (s.subnet_mask>>8)&0xFF))
	maskQuards = append(maskQuards, fmt.Sprintf(format, (s.subnet_mask>>0)&0xFF))

	return strings.Join(maskQuards, separator)
}
