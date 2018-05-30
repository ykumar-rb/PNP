package templates

//template for profiles
var ProfilesTmlp = `
{
  "id": "Install ubuntu Linux and Reboot",
  "boot": {
    "kernel": "/assets/coreos/client/linux gfxpayload=800x600x16,800x600 netcfg/enable=true netcfg/hostname=ztp-node-client netcfg/get_hostname=ztp-node-client --- auto=true url=http://{{.IP}}:{{.MatchboxPort}}/assets/coreos/client/preseed.cfg",
    "initrd": ["/assets/coreos/client/initrd.gz"],
    "args": [
    ]
  },
  "ignition_id": "ubuntu-install-reboot.yaml"
}
`
