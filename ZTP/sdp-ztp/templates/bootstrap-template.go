package templates

//template for bootstrap
var BootstrapTmlp = `#!/bin/bash

cd /root/PnP
cat resolv.conf > /etc/resolvconf/resolv.conf.d/base
/etc/init.d/networking restart
apt-get install -y curl
chmod +x client
export SDP_NETWORK_INTERFACE={{.Interface}}
./client --registry_address="{{.IP}}" --pnp_server="NewPnPService" --certificate_manager="CertificateManagerService"

`
