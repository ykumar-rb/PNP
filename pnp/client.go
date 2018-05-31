package main

import (
	"net"
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
	"log"
	"strings"
)

func main() {
	var pnpServer string
	var interfaceName string
	var pnpCertificateService string
	service := grpc.NewService(
		micro.Flags(
			cli.StringFlag{
				Name : "pnp_server",
				Value: "PnPServer",
				Usage: "PnP server name registered to registry",
			},
			cli.StringFlag{
				Name: "pnp_interface",
				Value: "ens33",
				Usage: "Client interface used to communicate with PnP Server",
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
			interfaceName = c.String("pnp_interface")

			pnp_interf, err := net.InterfaceByName(interfaceName)
			if err != nil {
				log.Fatalf("Unable to load interface specified, Error: %v", err)
			}

			if ! strings.Contains(pnp_interf.Flags.String(), "up") {
				log.Fatalf("Specified network interface %v is DOWN, specify running network interface", interfaceName)
			}
			pnpCertificateService = c.String("certificate_manager")
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

	color.Println("Initializing package management...")
	invoke.InstallMgmt(pnpClient)

	/*switch pnpOpType {
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
	}*/
}