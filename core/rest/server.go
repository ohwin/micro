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

type RouterFunc func(engine *gin.Engine)

func New(routers ...RouterFunc) *Server {
	initialize.Init()
	env := gin.DebugMode
	if config.App.Env == gin.ReleaseMode {
		env = gin.ReleaseMode
	}
	gin.SetMode(env)
	engine := gin.Default()
	for _, router := range routers {
		router(engine)
	}
	return &Server{
		Ctx:    context.Background(),
		Engine: engine,
	}
}

func (s *Server) Run() {
	err := s.Engine.Run(config.Addr())
	if err != nil {
		panic(err)
		return
	}
}
