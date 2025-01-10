package ginx

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/redis/go-redis/v9"

	"github.com/virzz/vlog"

	"github.com/virzz/ginx/apikey"
)

var engine *gin.Engine

func R() *gin.Engine { return engine }

func RegisterVersion(version, commit string) {
	engine.GET("/system/version", func(c *gin.Context) {
		c.String(200, version+" "+commit)
	})
}

func New(prefix ...string) (*http.Server, error) {
	if Conf == nil {
		return nil, fmt.Errorf("HTTP Config is nil")
	}

	engine = gin.New()

	f, err := os.OpenFile(filepath.Join("logs", "gin.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		vlog.Warn("Failed to open gin log file", "err", err.Error())
		return nil, err
	}
	if Conf.Debug {
		gin.SetMode(gin.DebugMode)
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		gin.DefaultErrorWriter = io.MultiWriter(f, os.Stderr)
		engine.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
	}
	engine.Use(gin.Recovery())

	if Conf.System != "" {
		engine.POST("/system/upgrade/:token", handleUpgrade)
		engine.POST("/system/upload/:token", handleUpload)
	}

	engine.Any("/health", HealthCheckHandler)

	if Conf.Metrics {
		m := ginmetrics.GetMonitor()
		m.SetMetricPath("/metrics")
		m.SetSlowTime(10)
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(engine)
	}

	if Conf.Pprof {
		pprof.Register(engine)
	}

	if Conf.RequestID {
		engine.Use(requestid.New())
	}

	// CORS
	c := cors.DefaultConfig()
	c.AddAllowHeaders(Conf.Headers...)
	if len(Conf.Origins) > 0 {
		c.AllowAllOrigins = false
		c.AllowOrigins = Conf.Origins
	} else {
		c.AllowAllOrigins = true
	}
	engine.Use(cors.New(c))

	// Session
	if Conf.Store.Enabled {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", Conf.Store.Addr, Conf.Store.Port),
			DB:       Conf.Store.DB,
			Password: Conf.Store.Pass,
		})
		store, err := apikey.NewRedisStore(apikey.WithClient(client))
		if err != nil {
			panic(err)
		}
		engine.Use(apikey.Init(store))
	}

	// Register Before Middleware
	if len(mwBefore) > 0 {
		engine.Use(mwBefore...)
	}

	// Register Router
	var api *gin.RouterGroup
	if len(prefix) > 0 {
		api = engine.Group(prefix[0])
	} else {
		api = engine.Group("/")
	}

	for _, register := range Routers {
		register(api)
	}

	// Register After Middleware
	if len(mwAfter) > 0 {
		engine.Use(mwAfter...)
	}

	addr := fmt.Sprintf("%s:%d", Conf.Addr, Conf.Port)
	vlog.Info("HTTP Server Listening on : " + addr)
	return &http.Server{Addr: addr, Handler: engine}, nil
}
