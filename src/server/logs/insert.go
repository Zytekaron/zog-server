package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"log"
	"net/http"
	"time"
)

func (l *LogHandler) Insert() func(*gin.Context) {
	return func(ctx *gin.Context) {
		var entry *types.Log
		err := json.NewDecoder(ctx.Request.Body).Decode(&entry)
		if err != nil {
			err = fmt.Errorf("error parsing the request body: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil(err.Error()))
			return
		}

		errs := entry.Errors()
		if len(errs) > 0 {
			ctx.JSON(http.StatusBadRequest, types.NewError("invalid log data", errs))
			return
		}

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = l.Logs.Insert(c, entry)
		if err == database.ErrDuplicateKey {
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil("a document with the specified id already exists"))
			return
		}
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred inserting the document"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccess(entry))
	}
}
