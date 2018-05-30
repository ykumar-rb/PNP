### Certificate-Manager

Certificate manager is used to supply the server certificates to clients. This certificate is used by client to initiate TLS-gRPC communication between pnp-client and pnp-server.

On request from client with its MACid, certificate manager discovers Onboarder api & makes request to Onboarder to fetch the Client details. Once done, it encrypts the server certificate with the MAC id of client and sends it to client.

#### Prerequisites

1. `go run GenerateTLSCertificate.go <primary_interface_name>`  (e.g. `go run GenerateTLSCertificate.go ens33`)
2. Start the certificate-manager service: `go run certificate.go --registry_address=<consul_ip> --server_name="CertificateManagerSevice" --onboarder_service_name="ClientOnboardService"`
