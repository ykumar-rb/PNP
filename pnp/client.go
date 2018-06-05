package main

import (
	"net"
	"crypto/tls"
	"crypto/x509"
	"log"
	"strings"
	"github.com/micro/go-grpc"
	"github.com/ZTP/pnp/invoke-service"
	"github.com/micro/go-micro"
	"github.com/micro/cli"
	"github.com/micro/go-micro/transport"
	"github.com/ZTP/pnp/util/client"
	"github.com/ZTP/pnp/common/color"
	proto "github.com/ZTP/pnp/pnp-proto"
	certproto "github.com/ZTP/certificate-manager/proto/certificate"
	invokeCertManager "github.com/ZTP/certificate-manager/invoke-service"
	"github.com/ZTP/onboarder/publisher/pubsub-proto"
	"github.com/micro/go-micro/metadata"
	"github.com/ZTP/pnp/config"
	"context"
	"os"
	"time"
)

//type Sub struct {}

var pnpClient proto.PnPService
var clientInfo proto.ClientInfo

func main() {
	var pnpCertificateService string
	service := grpc.NewService(
		micro.Name("PnPClient"),
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
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*5),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			config.PnpServerName = c.String("pnp_server")
			config.ClientInterface = c.String("pnp_interface")

			pnp_interf, err := net.InterfaceByName(config.ClientInterface)
			if err != nil {
				log.Fatalf("Unable to load interface specified, Error: %v", err)
			}

			if ! strings.Contains(pnp_interf.Flags.String(), "up") {
				log.Fatalf("Specified network interface %v is DOWN, specify running network interface", config.ClientInterface)
			}
			pnpCertificateService = c.String("certificate_manager")
		}),
	)
	pnpClient = proto.PnPServiceClient(config.PnpServerName, service.Client())
	pnpCertClient := certproto.CertificateServiceClient(pnpCertificateService, service.Client())
	clientInfo = client.PopulateClientDetails(config.ClientInterface)

	caCert, clientEnvName := invokeCertManager.GetCertificate(pnpCertClient, clientInfo)
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
	invoke.InitPkgMgmt(pnpClient, clientInfo)

	if err := micro.RegisterSubscriber(clientEnvName, service.Server(), initiateClientUpdate); err != nil {
		color.Warnf("Unable to subscribe topic %v, Error: %v", clientEnvName, err)
		os.Exit(1)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initiateClientUpdate(cxt context.Context, event *publisher.Event) error {
	md, _ := metadata.FromContext(cxt)
	color.Printf("[PubSub] Received update event with metatdata %v\n.. Initiating package update", md)
	invoke.InitPkgMgmt(pnpClient, clientInfo)
	return nil
}