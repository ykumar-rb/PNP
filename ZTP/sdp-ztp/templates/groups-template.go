package templates

//template for groups
var GroupsTmlp = `{
  "id": "client",
  "name": "ubuntu client",
  "profile": "ubuntu-install-reboot-client",
  "selector": {
  },
  "metadata": {
    "coreos_channel": "stable",
    "coreos_version": "ubuntu",
    "ignition_endpoint": "http://matchbox.foo:8080/ignition",
    "baseurl": "http://matchbox.foo:8080/assets/coreos"
  }
}
`
