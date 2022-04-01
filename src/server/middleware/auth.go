package middleware

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"log"
	"net/http"
	"time"
)

func Auth(tc database.TokenController, mode string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, types.NewErrorNil("Unauthorized"))
			return
		}

		hash := makeHash(authorization)

		c, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		token, err := tc.Get(c, hash)
		if err != nil {
			if err == database.ErrNoDocuments {
				ctx.AbortWithStatusJSON(http.StatusForbidden, types.NewErrorNil("Forbidden"))
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, types.NewErrorNil("an error occurred fetching the token"))
			return
		}

		ctx.Set("token", token)

		switch mode {
		case "read":
			if !token.Read {
				ctx.AbortWithStatusJSON(http.StatusForbidden, types.NewErrorNil("Forbidden"))
			}
		case "write":
			if !token.Write {
				ctx.AbortWithStatusJSON(http.StatusForbidden, types.NewErrorNil("Forbidden"))
			}
		default:
			log.Fatal("unknown endpoint mode", mode)
		}
	}
}

func makeHash(str string) string {
	sha := sha512.New()
	sha.Write([]byte(str))
	sum := sha.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(sum)
}
