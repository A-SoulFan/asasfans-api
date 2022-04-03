package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	logger *zap.Logger
	router *gin.Engine

	httpServer http.Server
}

type InitRouters func(r *gin.Engine)

func NewRouter(c *Config, logger *zap.Logger, init InitRouters) *gin.Engine {
	gin.SetMode(c.Mode)
	// 初始化 gin
	r := gin.New()

	// panic之后自动恢复
	r.Use(gin.Recovery())
	r.Use(ginZap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginZap.RecoveryWithZap(logger, true))

	if c.MaxMultipartMemory != 0 {
		// 最大上传文件大小 mb
		r.MaxMultipartMemory = int64(c.MaxMultipartMemory) << 20
	}

	init(r)

	return r
}

func NewServer(c *Config, logger *zap.Logger, router *gin.Engine) (*Server, error) {
	var s = &Server{
		config: c,
		router: router,
		logger: logger.With(zap.String("type", "http.Server")),
	}

	return s, nil
}

func (s *Server) Start() error {
	if s.config.Port == 0 {
		s.config.Port = 8080
	}

	// 默认 8080 端口
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	s.httpServer = http.Server{Addr: addr, Handler: s.router}
	s.logger.Info("http server starting ...", zap.String("addr", addr))

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("start http server err", zap.Error(err))
			return
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping http server")
	// 平滑关闭,等待5秒钟处理
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}

func Provide() fx.Option {
	return fx.Provide(NewConfig, NewRouter, NewServer)
}
