package main

import (
	"os"
	"crypto/tls"
	"crypto/x509"
	"github.com/micro/go-grpc"
	"github.com/ZTP/pnp/invoke-service"
	"github.com/micro/go-micro"
	"github.com/micro/cli"
	"github.com/micro/go-micro/transport"
	"github.com/ZTP/pnp/util/color"
	proto "github.com/ZTP/pnp/pnp-proto"
	certproto "github.com/ZTP/certificate-manager/proto/certificate"
	invokeCertManager "github.com/ZTP/certificate-manager/invoke-service"
)

func main() {
	var pnpServer string
	var pnpOpType string
	var pnpCertificateService string
	interfaceName := os.Getenv("SDP_NETWORK_INTERFACE")
	if interfaceName == "" {
		color.Fatalf("Provide \"SDP_NETWORK_INTERFACE\" environment variable")
	}
	service := grpc.NewService(
		micro.Flags(
			cli.StringFlag{
				Name : "pnp_server",
				Value: "PnPServer",
				Usage: "PnP server name registered to registry",
			},
			cli.StringFlag{
				Name : "pnp_op_type",
				Usage: "Specifies pnp operation type, supported values are" +
					"installPackages, deploySDPMaster, deploySDPSatellite",
			},
			cli.StringFlag{
				Name : "certificate_manager",
				Value: "CertificateManagerSevice",
				Usage: "Certificate-manager server name registered to registry",
			},
		),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			pnpServer = c.String("pnp_server")
			pnpOpType = c.String("pnp_op_type")
			pnpCertificateService = c.String("certificate_manager")
			if pnpOpType == "" {
				color.Fatalf("PnP operation type not specified, supported values are" +
					"installPackages, deploySDPMaster, deploySDPSatellite")
			}
		}),
	)
	pnpClient := proto.PnPServiceClient(pnpServer, service.Client())
	pnpCertClient := certproto.CertificateServiceClient(pnpCertificateService, service.Client())
	caCert := invokeCertManager.GetCertificate(pnpCertClient, interfaceName)
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	service.Init(
		micro.Transport(transport.NewTransport(transport.Secure(true))),
		grpc.WithTLS(tlsConfig),
	)

	switch pnpOpType {
	case "installPackages":
		{
			color.Println("Initializing package installation..")
			invoke.InstallPackages(pnpClient)
		}
	default:
		{
			color.Println("PnP operation type not specified, supported values are " +
			"installPackages")
		}
	}
}