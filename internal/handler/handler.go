package handler

import (
	"falconapi/internal/model"
	"falconapi/internal/service"
	"falconapi/pkg/logging"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger  *logging.Logger
	Engine  *gin.Engine
	Service *service.Service
	Cfg     *model.Config
}

func NewHandler(logger *logging.Logger, engine *gin.Engine, service *service.Service, cfg *model.Config) *Handler {
	return &Handler{
		Logger:  logger,
		Engine:  engine,
		Service: service,
		Cfg:     cfg,
	}
}

func (h *Handler) InitRoutes() {
	v1 := h.Engine.Group("/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login-user", h.LoginUser)
		auth.POST("/generate-otp", h.GenerateOtp)
		auth.POST("/validate-otp", h.ValidateOtp)
	}

	v1.Group("/api")
	v1.Use(h.CorsMiddleware())
}
