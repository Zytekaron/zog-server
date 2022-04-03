package middleware

import (
	"github.com/zytekaron/gotil/v2/rl"
	"sync"
	"time"
)

// periodically delete buckets that have expired
func cleanBuckets(bm *rl.BucketManager, mux *sync.Mutex, interval time.Duration) {
	for {
		time.Sleep(interval)

		mux.Lock()
		clean(bm)
		mux.Unlock()
	}
}

func clean(bm *rl.BucketManager) {
	for id, bucket := range bm.Buckets {
		if bucket.RemainingTime() == 0 {
			delete(bm.Buckets, id)
		}
	}
}
