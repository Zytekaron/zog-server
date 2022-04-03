package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/find"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (l *LogHandler) Find() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var query *find.LogQuery
		err := json.NewDecoder(ctx.Request.Body).Decode(&query)
		if err != nil {
			err = fmt.Errorf("error parsing the request body: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil(err.Error()))
			return
		}

		limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "100"), 10, 64)
		if err != nil || limit < 1 || limit > 100 {
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil("limit must be a valid integer from 1 to 100"))
			return
		}
		offset, err := strconv.ParseInt(ctx.DefaultQuery("offset", "0"), 10, 64)
		if err != nil || offset < 0 {
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil("offset must be a valid non-negative integer"))
			return
		}

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		options := find.NewLogOptions().
			SortCreatedAt(1).
			Limit(limit).
			Skip(limit * offset)

		iter, err := l.Logs.Find(c, query, options)
		if err != nil {
			log.Println("error finding documents:", err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred fetching documents"))
			return
		}

		entries, err := database.SliceBuf(iter, limit)
		if iter.Err() != nil {
			log.Println("error decoding found documents:", err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred fetching documents"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccess(entries))
	}
}
