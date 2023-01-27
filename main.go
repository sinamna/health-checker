package main

import (
	"fmt"
	"health_checker/api"
	"health_checker/api/middleware"
	"health_checker/modules/checker"
	"health_checker/pkg/config"
	"health_checker/pkg/repository"
)

func main() {
	fmt.Println("loading config...")
	config.LoadConfig()

	fmt.Println("setting up postgres...")
	repository.SetupPostgres()

	server := api.CreateHttpServer()
	// user signup
	server.Router.POST("/user/signup", api.SignupHandler)
	server.Router.POST("/user/login", api.LoginHandler)

	// create endpoint
	server.Router.POST("/endpoint/", middleware.AuthMiddleware(), api.EndpointRegisterHandler)

	// retrieve user endpoints
	server.Router.GET("/endpoint/", middleware.AuthMiddleware(), api.EndpointRetrieverHandler)

	// details of success requests
	server.Router.GET("/endpoint/:id", middleware.AuthMiddleware(), api.EndpointStatusHandler)

	// get alerts of endpoints
	server.Router.GET("/alerts/:id", middleware.AuthMiddleware(), api.AlertHandler)

	checker := checker.NewChecker(config.Conf.Server.WorkerNum, config.Conf.Server.Interval)
	checker.Start()

	if err := server.StartServer(fmt.Sprintf("%d", config.Conf.Server.Port)); err != nil {
		panic(err)
	}
}
