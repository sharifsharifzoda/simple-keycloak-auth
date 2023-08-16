package main

import (
	"falconapi/config"
	"falconapi/internal/handler"
	"falconapi/internal/service"
	"falconapi/pkg/logging"
	"github.com/gin-gonic/gin"
	"net"
)

func main() {
	cfg := config.InitViper()

	router := gin.Default()

	log := logging.GetLogger()

	newService := service.NewService(cfg, log)

	newHandler := handler.NewHandler(log, router, newService, cfg)
	newHandler.InitRoutes()

	log.Fatal(router.Run(net.JoinHostPort(cfg.Listen.ListenIP, cfg.Listen.ListenPort)))
}
