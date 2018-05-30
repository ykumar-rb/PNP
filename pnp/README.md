# PnP

## Prerequisites

1. go lang: version > 1.9
2. Set env variable GOPATH
3. Set environment variable "SDP_NETWORK_INTERFACE=<Primary_interface>" e.g. `ens33`
3. Run `$ go get "github.com/micro/go-micro"`
4. Run `$ go get "github.com/micro/go-grpc"`
5. Run `$ go get "github.com/micro/cli"`
6. Running instance of Consul.
7. Generate the server certificate and key file: 
   
   7.1. Go to `'../PnP/util/'` folder, and run the `GenerateTLSCertificate.go <interface_name>`. Enter the interface name as command line argument. This generates the `server.crt` & `server.key` files in `../PnP/certs` folder.
    
   7.2. Transfer these files to Client machine in the folder: `'../PnP/certs'`.

Note: To run PnP server and client you should be a root user

## Running PnP Server

`$ go run server.go --registry_address=<consul_ip> --server_name=<pnp_server_name> --package_file=<path/of/packageInfo.json>  --cert_file "../certs/server.crt" --key_file "../certs/server.key"`

e.g.: 
`$ go run server.go --registry_address=172.16.128.132 --server_name "NewPnPService" --package_file "/../config/packageInfo.json" --cert_file "../certificate-manager/certs/server.crt" --key_file "../certificate-manager/certs/server.key"`

`packageInfo.json` recides in config directory.

## Running PnP client

`$ go run client.go --registry_address=<consul_ip> --pnp_server=<pnp server name registered to consul> --pnp_op_type=<operation_type> --certificate_manager=<certificate manager server name registered to consul>`

e.g.: 
`$ go run client.go --registry_address="192.168.50.129" --pnp_server="NewPnPService" --pnp_op_type="installPackages" --certificate_manager="CertificateManagerSevice"`
