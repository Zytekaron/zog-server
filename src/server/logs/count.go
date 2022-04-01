package logs

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/types"
	"log"
	"net/http"
	"time"
)

func (l *LogHandler) Count() func(*gin.Context) {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		count, err := l.Logs.Count(c)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred counting documents"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccess(count))
	}
}
