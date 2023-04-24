package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/initialize"
)

type Server struct {
	Ctx    context.Context
	Engine *gin.Engine
}

func New() *Server {
	initialize.Init()
	gin.SetMode(config.App.Env)

	return &Server{
		Ctx:    context.Background(),
		Engine: gin.Default(),
	}
}

func (s *Server) AddRouters(routers ...func(r *gin.Engine)) *Server {
	for _, f := range routers {
		f(s.Engine)
	}
	return s
}

func (s *Server) Run() {

	////router.InitRouter(s.Engine)
	//s.Engine.GET("/ss", func(c *gin.Context) {
	//	c.JSON(200, "sd")
	//})
	err := s.Engine.Run(config.Addr())
	if err != nil {
		panic(err)
		return
	}
}
