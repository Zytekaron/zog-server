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

func (l *LogHandler) Delete() func(*gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := l.Logs.Delete(c, id)
		if err == database.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, types.NewErrorNil("document does not exist"))
			return
		}
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil("an error occurred fetching the document"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccessNil())
	}
}
