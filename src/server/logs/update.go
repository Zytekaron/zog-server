package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/updates"
	"log"
	"net/http"
	"time"
)

func (l *LogHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var updateData *updates.Log
		err := json.NewDecoder(ctx.Request.Body).Decode(&updateData)
		if err != nil {
			err = fmt.Errorf("error parsing the request body: %w", err)
			ctx.JSON(http.StatusBadRequest, types.NewErrorNil(err.Error()))
			return
		}

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = l.Logs.Update(c, id, updateData)
		if err == database.ErrNoDocuments {
			ctx.JSON(http.StatusNotFound, types.NewErrorNil("document does not exist"))
			return
		}
		if err != nil {
			log.Println("error updating document:", err)
			ctx.JSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred updating the document"))
			return
		}

		ctx.JSON(http.StatusOK, types.NewSuccessNil())
	}
}
