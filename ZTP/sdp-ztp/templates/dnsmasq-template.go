package templates

//template for Dnsmasq
var DnsmasqTmlp = `
# Configuration file for dnsmasq.

interface={{.Interface}}
bind-interfaces
dhcp-range={{.IPRange}}
enable-tftp
tftp-root=/var/lib/tftpboot

dhcp-userclass=set:ipxe,iPXE
dhcp-boot=tag:#ipxe,undionly.kpxe
dhcp-boot=tag:ipxe,http://{{.IP}}:{{.MatchboxPort}}/boot.ipxe

log-queries
log-dhcp
`

