package server

import (
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"syscall"
	"time"
)

type GinServer struct {
	host string
	port string
	prod bool
	e    *gin.Engine
}

func NewGinServer(host string, port string, prod bool) GinServer {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.RemoveExtraSlash = true
	r.NoRoute(noRouteHandler)
	r.NoMethod(methodNotAllowHandler)
	if !prod {
		r.Use(gin.LoggerWithFormatter(apacheLogFormat))
		gin.SetMode(gin.DebugMode)
	}
	return GinServer{
		host: host,
		port: port,
		prod: prod,
		e:    r,
	}
}

func (s *GinServer) Run() error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	server := endless.NewServer(addr, s.e)
	server.BeforeBegin = func(add string) {
		log.Printf("Server start at pid %d", syscall.Getpid())
	}
	return server.ListenAndServe()
}

func (s *GinServer) Ctx() *gin.Engine {
	return s.e
}

func (s *GinServer) EnableMetrics(path string, subSystem string) {
	p := ginprom.New(
		ginprom.Engine(s.e),
		ginprom.Subsystem(subSystem),
		ginprom.Path(path),
		ginprom.Ignore(path),
	)
	s.e.Use(p.Instrument())
}

func apacheLogFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] %s %s %s %d %s \"%s\" %s\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func noRouteHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, ApiResult{Error: "Request api not found"})
}

func methodNotAllowHandler(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, ApiResult{Error: "Request method not allowed"})
}
