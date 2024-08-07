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
	"github.com/virzz/ginx/code"
	"github.com/virzz/ginx/rsp"
)

func CodesHandler(c *gin.Context) {
	c.JSON(200, rsp.S(code.Codes))
}

func New(prefix ...string) (*http.Server, error) {
	if Conf == nil {
		return nil, fmt.Errorf("HTTP Config is nil")
	}
	engine := gin.New()

	f, _ := os.Create(filepath.Join("logs", "gin.log"))
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
	if Conf.Store.Addr != "" {
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

	// Register Router
	var api *gin.RouterGroup
	if len(prefix) > 0 {
		api = engine.Group(prefix[0])
	} else {
		api = engine.Group("/")
	}
	api.POST("/errcode", CodesHandler)
	api.POST("/captcha", CaptchaHandler)
	for _, register := range Routers {
		register(api)
	}

	addr := fmt.Sprintf("%s:%d", Conf.Addr, Conf.Port)
	vlog.Info("HTTP Server Listening on : " + addr)
	return &http.Server{Addr: addr, Handler: engine}, nil
}
