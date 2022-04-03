package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/gotil/v2/rl"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"net/http"
	"sync"
	"time"
)

func Auth(tc database.Controller[*types.Token], limit int, reset time.Duration) gin.HandlerFunc {
	bm := rl.New(limit, reset)
	var bmMux sync.Mutex
	go cleanBuckets(bm, &bmMux, 10*time.Minute)

	return func(ctx *gin.Context) {
		bmMux.Lock()
		bucket := bm.Get(ctx.ClientIP())
		bmMux.Unlock()

		if !bucket.CanDraw(1) {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, types.NewErrorNil("Too Many Requests"))
			return
		}

		header := ctx.GetHeader("Authorization")
		if header == "" {
			bmMux.Lock()
			bucket.Draw(1)
			bmMux.Unlock()

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.NewErrorNil("Unauthorized"))
			return
		}

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		token, err := tc.Get(c, makeHash(header))
		if err != nil {
			if err == database.ErrNoDocuments {
				bmMux.Lock()
				bucket.Draw(1)
				bmMux.Unlock()

				ctx.AbortWithStatusJSON(http.StatusForbidden, types.NewErrorNil("Forbidden"))
				return
			}

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred fetching the token"))
			return
		}

		ctx.Set("token", token)
		ctx.Next()
	}
}

func makeHash(str string) string {
	sha := sha512.New()
	sha.Write([]byte(str))
	sum := sha.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(sum)
}
