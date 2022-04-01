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

var buckets = map[string]*rl.Bucket{}
var bucketsMux sync.Mutex

func init() {
	// daemon to purge expired ratelimit data
	go func() {
		for {
			time.Sleep(10 * time.Minute)

			bucketsMux.Lock()
			for id, bucket := range buckets {
				if bucket.RemainingTime() == 0 {
					delete(buckets, id)
				}
			}
			bucketsMux.Unlock()
		}
	}()
}

// RateLimit middleware provides rate limiting to a route.
// Auth must be called before RateLimit to ensure the API token has been loaded.
func RateLimit(mode string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		tokenInterface, ok := ctx.Get("token")
		if !ok {
			log.Println("erroneous passthrough: RateLimit middleware called before calling Auth")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, types.NewErrorNil("Internal Server Error"))
			return
		}
		token := tokenInterface.(*types.Token)

		switch mode {
		case "read":
			bucket := getBucket(token.ID, token.ReadLimit, token.ReadReset)
			if !bucket.CanDraw(1) {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, types.NewErrorNil("Too Many Requests"))
			}
		case "write":
			bucket := getBucket(token.ID, token.WriteLimit, token.WriteReset)
			if !bucket.CanDraw(1) {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests, types.NewErrorNil("Too Many Requests"))
			}
		default:
			log.Println("erroneous RateLimit mode passed:", mode)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, types.NewErrorNil("Internal Server Error"))
		}
	}
}

func getBucket(id string, limit int, reset types.Duration) *rl.Bucket {
	bucketsMux.Lock()
	defer bucketsMux.Unlock()

	bucket, ok := buckets[id]
	if !ok {
		bucket = rl.NewBucket(limit, time.Duration(reset))
		buckets[id] = bucket
	}

	return bucket
}
