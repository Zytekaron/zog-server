package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/server/middleware"
	"github.com/zytekaron/zog-server/src/types"
)

type LogHandler struct {
	Config *config.Config
	Logs   database.Controller[*types.Log]
	Tokens database.Controller[*types.Token]
}

func New(cfg *config.Config, lc database.Controller[*types.Log], tc database.Controller[*types.Token]) *LogHandler {
	return &LogHandler{
		Config: cfg,
		Logs:   lc,
		Tokens: tc,
	}
}

func (l *LogHandler) Register(r *gin.Engine) {
	readAuth := middleware.Auth(l.Tokens, "read")
	writeAuth := middleware.Auth(l.Tokens, "write")

	readRateLimit := middleware.RateLimit("read")
	writeRateLimit := middleware.RateLimit("write")

	r.GET("/", readAuth, readRateLimit, l.Find())
	r.GET("/:id", readAuth, readRateLimit, l.Get())
	r.GET("/count", readAuth, readRateLimit, l.Count())

	r.POST("/", writeAuth, writeRateLimit, l.Insert())
	r.PATCH("/:id", writeAuth, writeRateLimit, l.Update())
	r.DELETE("/:id", writeAuth, writeRateLimit, l.Delete())
}
