package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Server struct {
	*gin.Engine

	ConfigProvider ConfigProvider
	server         *http.Server
	addr           string
}

func New(config ConfigProvider) *Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	// TODO
	// engine.Use(middleware.Recovery())

	return &Server{
		ConfigProvider: config,
		Engine:         engine,
	}
}

func Run(opts ...fx.Option) {
	app := fx.New(
		fx.NopLogger, // 关闭 Fx 自身的日志
		// constructors
		fx.Provide(
			// TODO
			New,
		),
		// TODO
		fx.Options(opts...),
	)

	app.Run()
}
