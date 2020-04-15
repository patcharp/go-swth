package server

import (
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	echoprometheus "github.com/globocom/echo-prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

const (
	Production  = true
	Development = false
	MetricsPath = "/metrics"
	NameSpace   = "echo"
)

type Config struct {
	host string
	port string
	prod bool
	e    *echo.Echo
}

func New(host string, port string, prod bool) Config {
	cfg := Config{
		host: host,
		port: port,
		prod: prod,
		e:    echo.New(),
	}
	cfg.e.Use(middleware.Recover())
	cfg.e.Use(middleware.RemoveTrailingSlash())
	// Cr. https://echo.labstack.com/middleware/request-id
	cfg.e.Use(middleware.RequestID())
	// Cr. https://echo.labstack.com/middleware/secure
	cfg.e.Use(middleware.Secure())
	cfg.e.HTTPErrorHandler = cfg.serverErrorHandler
	listenAddr := fmt.Sprintf("%s:%s", cfg.host, cfg.port)
	cfg.e.Server = &http.Server{
		Addr:    listenAddr,
		Handler: cfg.e,
	}
	if !cfg.prod {
		cfg.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))
		cfg.e.Use(cfg.logger())
	}
	return cfg
}

func (cfg *Config) Ctx() *echo.Echo {
	return cfg.e
}

func (cfg *Config) Run() error {
	if err := gracehttp.Serve(cfg.e.Server); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) EnableMetrics(metricsPath string, nameSpace string) error {
	cfg.e.Use(echoprometheus.MetricsMiddlewareWithConfig(echoprometheus.Config{
		Namespace: nameSpace,
		Buckets: []float64{
			0.0005, // 0.5ms
			0.001,  // 1ms
			0.005,  // 5ms
			0.01,   // 10ms
			0.05,   // 50ms
			0.1,    // 100ms
			0.5,    // 500ms
			1,      // 1s
			2,      // 2s
			5,      // 5s
			10,     // 10s
		},
	}))
	cfg.e.GET(metricsPath, echo.WrapHandler(promhttp.Handler()))
	return nil
}

func (cfg *Config) serverErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
		if he.Internal != nil {
			err = fmt.Errorf("%v, %v", err, he.Internal)
		}
	} else if cfg.e.Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{"error": msg}
	}
	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			cfg.e.Logger.Error(err)
		}
	}
}

func (cfg *Config) logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}
			reqSize := req.Header.Get(echo.HeaderContentLength)
			if reqSize == "" {
				reqSize = "0"
			}
			log.Infof("%s %s [%v] %s %-7s %s %3d %s %s %13v %s %s",
				id,
				c.RealIP(),
				stop.Format(time.RFC3339),
				req.Host,
				req.Method,
				req.RequestURI,
				res.Status,
				reqSize,
				strconv.FormatInt(res.Size, 10),
				stop.Sub(start).String(),
				req.Referer(),
				req.UserAgent(),
			)
			return err
		}
	}
}

func (cfg *Config) SkipperExceptPath(c echo.Context, path map[string]string) bool {
	if method, ok := path[c.Path()]; ok && method == c.Request().Method {
		return true
	}
	return false
}
