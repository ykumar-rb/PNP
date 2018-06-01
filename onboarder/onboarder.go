package main

import (
	"github.com/micro/go-web"
	"github.com/emicklei/go-restful"
	"github.com/ZTP/pnp/common/color"
	"github.com/ZTP/onboarder/handlers"
	"github.com/go-redis/redis"
	"os"
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
	clientEnv := NewEnv()
	wSvc.Route(wSvc.GET("/clients").To(onboarderSvc.GetAllRegisteredClients))
	wSvc.Route(wSvc.GET("/clients/{mac}").To(onboarderSvc.GetRegisteredClientDetails))
	wSvc.Route(wSvc.POST("/clients").To(onboarderSvc.RegisterClient))
	wSvc.Route(wSvc.DELETE("/clients/{mac}").To(onboarderSvc.DeregisterClient))
	wSvc.Route(wSvc.POST("/environment").To(clientEnv.CreateEnvironment))
	wSvc.Route(wSvc.PUT("/environment").To(clientEnv.UpdateEnvironment))
	wCntnr.Add(wSvc)
	onboarderService.Handle("/", wCntnr)
	if err := onboarderService.Run(); err != nil {
		color.Fatalf("%v",err)
	}
}

func NewEnv () *handlers.InstallEnv {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		color.Fatalf("Provide \"REDIS_ADDR\" environment variable")
	}
	return &handlers.InstallEnv{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}
