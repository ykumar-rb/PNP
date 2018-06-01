package main

import (
	"time"
	"crypto/tls"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/ZTP/pnp/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro/transport"
	"github.com/ZTP/pnp/common/color"
	handler "github.com/ZTP/pnp/handlers"
	proto "github.com/ZTP/pnp/pnp-proto"
)

func main() {
	service := grpc.NewService(
		micro.Name("PnPServer"),
		micro.Flags(
			cli.StringFlag{
				Name : "cert_file",
				Value: "./certs/server.crt",
				Usage: "Path of server certificate file",
			},
			cli.StringFlag{
				Name : "key_file",
				Value: "./certs/server.key",
				Usage: "Path of server key file",
			},
		),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*5),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			config.CertFile = c.String("cert_file")
			config.KeyFile = c.String("key_file")
		}),
	)

	cert, err := tls.LoadX509KeyPair(config.CertFile , config.KeyFile)
	if err != nil {
		color.Fatalf("%v",err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	service.Init(
		micro.Transport(transport.NewTransport(transport.Secure(true))),
		grpc.WithTLS(tlsConfig),
	)
	proto.RegisterPnPHandler(service.Server(), new(handler.PnPService))

	if err := service.Run(); err != nil {
		color.Fatalf("%v",err)
	}
}
