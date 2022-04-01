package main

import (
	"context"
	"github.com/zytekaron/zog-server/src/config"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/database/cache"
	"github.com/zytekaron/zog-server/src/database/mongodb"
	"github.com/zytekaron/zog-server/src/server"
	"log"
	"time"
)

var configDirs = []string{
	"/etc/opt/zog.yml",
	"config.yml",
}

var cfg *config.Config
var lc database.LogController
var tc database.TokenController

func init() {
	var err error
	cfg, err = config.Load(configDirs)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch cfg.Database.Use {
	case database.MongoDB:
		client := try(mongodb.NewClient(ctx, cfg.Database.MongoDB))
		db := client.Database(cfg.Database.MongoDB.Database)

		lc = mongodb.NewLogRepository(db, cfg.Database.MongoDB)
		tc = mongodb.NewTokenRepository(db, cfg.Database.MongoDB)
		tc = cache.NewTokenCache(tc, cfg.Cache.Tokens)
	default:
		log.Fatal("invalid config option: db.use (expected one of: mongodb)")
	}
}

func main() {
	err := server.New(cfg, lc, tc).Start()
	if err != nil {
		log.Fatal(err)
	}
}

func try[T any](result T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// use createToken to insert tokens manually
//func createToken() {
//	err := fillAndSave(&types.Token{
//		OwnerID:    "zytekaron",
//		CreatedAt:  types.Time(time.Now()),
//		ExpiresAt:  types.Time(time.Unix(0, 0)),
//		Read:       true,
//		ReadLimit:  100,
//		ReadReset:  types.Duration(time.Second),
//		Write:      true,
//		WriteLimit: 10,
//		WriteReset: types.Duration(time.Second),
//	})
//	if err != nil {
//		log.Println(err)
//	}
//}
//func fillAndSave(token *types.Token) error {
//	tokenString := random.MustSecureString(32, "0123456789abcdef") // make sure you print it
//	log.Println(tokenString)
//
//	token.ID = makeHash(tokenString)
//
//	return tc.Insert(context.Background(), token)
//}
//func makeHash(str string) string {
//	sha := sha512.New()
//	sha.Write([]byte(str))
//	sum := sha.Sum(nil)
//
//	return base64.RawURLEncoding.EncodeToString(sum)
//}
