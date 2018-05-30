package main

import (
	"github.com/micro/go-web"
	"github.com/emicklei/go-restful"
	"github.com/RiverbedTechnology/sdp-ztp/pnp/util/color"
	"github.com/RiverbedTechnology/sdp-ztp/onboarder/handlers"
)
var onboarderSvc = handlers.Onboarder{}

func main() {
	onboarderService := web.NewService(
		web.Name("ClientOnboardService"),
		web.Address(":8099"),
	)
	onboarderService.Init()
	onboarderSvc := new(handlers.Onboarder)
	wSvc := new(restful.WebService)
	wCntnr := restful.NewContainer()
	wSvc.Consumes(restful.MIME_JSON)
	wSvc.Produces(restful.MIME_JSON)
	wSvc.Path("/pnp")
	wSvc.Route(wSvc.GET("/clients").To(onboarderSvc.GetAllRegisteredClients))
	wSvc.Route(wSvc.GET("/clients/{mac}").To(onboarderSvc.GetRegisteredClientDetails))
	wSvc.Route(wSvc.POST("/clients").To(onboarderSvc.RegisterClient))
	wSvc.Route(wSvc.DELETE("/clients/{mac}").To(onboarderSvc.DeregisterClient))
	wCntnr.Add(wSvc)
	onboarderService.Handle("/", wCntnr)
	if err := onboarderService.Run(); err != nil {
		color.Fatalf("%v",err)
	}
}
