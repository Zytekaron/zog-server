package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/server/middleware"
	"github.com/zytekaron/zog-server/src/types"
	"time"
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
	auth := middleware.Auth(l.Tokens, 1, time.Second)

	readRL := middleware.RateLimitToken(25, time.Second)
	r.GET("/:id", auth, readRL, l.Get())
	r.GET("/count", auth, readRL, l.Count())

	batchRL := middleware.RateLimitToken(5, time.Second)
	r.GET("/", auth, batchRL, l.Find())

	writeRL := middleware.RateLimitToken(10, time.Second)
	r.POST("/", auth, writeRL, l.Insert())
	r.PATCH("/:id", auth, writeRL, l.Update())
	r.DELETE("/:id", auth, writeRL, l.Delete())
}
