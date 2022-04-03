package logs

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"log"
	"net/http"
	"time"
)

func (l *LogHandler) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		entry, err := l.Logs.Get(c, id)
		if err == database.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, types.NewErrorNil("document does not exist"))
			return
		}
		if err != nil {
			log.Println("error getting document:", err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred fetching the document"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccess(entry))
	}
}
