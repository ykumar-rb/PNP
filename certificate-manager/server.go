package main

import (
	"time"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/cli"
	"github.com/ZTP/pnp/common/color"
	certHandler "github.com/ZTP/certificate-manager/handlers"
	certproto "github.com/ZTP/certificate-manager/proto/certificate"
	"github.com/ZTP/certificate-manager/helper"
)

func main() {
	certService := grpc.NewService(
		micro.Name("CertificateManagerSevice"),
		micro.Flags(
			cli.StringFlag{
				Name : "onboarder_service_name",
				Value: "ClientOnboardService",
				Usage: "Service name of Client onboarder Rest api",
			},
		),
		micro.RegisterTTL(time.Second*15),
		micro.RegisterInterval(time.Second*5),
	)
	certService.Init(
		micro.Action(func(c *cli.Context) {
			helper.ConsulServiceName = c.String("onboarder_service_name")
		}),
	)
	certproto.RegisterCertificateHandler(certService.Server(), new(certHandler.PnPCertificateService))
	if err := certService.Run(); err != nil {
		color.Fatalf("%v",err)
	}
}
