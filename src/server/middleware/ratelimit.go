package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/gotil/v2/rl"
	"github.com/zytekaron/zog-server/src/types"
	"log"
	"net/http"
	"sync"
	"time"
)

// RateLimitToken middleware provides rate limiting to a route.
// Auth must be called before RateLimitToken to ensure the API token has been loaded.
func RateLimitToken(limit int, reset time.Duration) gin.HandlerFunc {
	bm := rl.New(limit, reset)
	var bmMux sync.Mutex
	go cleanBuckets(bm, &bmMux, 10*time.Minute)

	return func(ctx *gin.Context) {
		t, ok := ctx.Get("token")
		if !ok {
			log.Println("erroneous passthrough: RateLimitToken middleware called before calling Auth")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, types.NewErrorNil("Internal Server Error"))
			return
		}
		token := t.(*types.Token)

		bucket := bm.Get(token.ID)
		if !bucket.Draw(1) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, types.NewErrorNil("Too Many Requests"))
			return
		}

		ctx.Next()
	}
}
