package cache

import (
	"context"
	"github.com/zyedidia/generic/cache"
	"github.com/zytekaron/zog-server/src/database"
	"github.com/zytekaron/zog-server/src/types"
	"github.com/zytekaron/zog-server/src/types/find"
	"github.com/zytekaron/zog-server/src/types/updates"
	"sync"
)

type TokenCache struct {
	supplier database.Controller[*types.Token]
	cache    *cache.Cache[string, *types.Token]
	cacheMux sync.Mutex
}

func NewTokenCache(supplier database.Controller[*types.Token], cacheSize int) *TokenCache {
	return &TokenCache{
		supplier: supplier,
		cache:    cache.New[string, *types.Token](cacheSize),
	}
}

func (t *TokenCache) Insert(ctx context.Context, token *types.Token) error {
	if _, ok := t.syncGet(token.ID); ok {
		return database.ErrDuplicateKey
	}

	err := t.supplier.Insert(ctx, token)
	if err != nil {
		return err
	}

	t.syncPut(token)
	return nil
}

func (t *TokenCache) Get(ctx context.Context, id string) (*types.Token, error) {
	if token, ok := t.syncGet(id); ok {
		return token, nil
	}

	token, err := t.supplier.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	t.syncPut(token)
	return token, nil
}

func (t *TokenCache) Update(ctx context.Context, id string, updates updates.Updates[*types.Token]) error {
	err := t.supplier.Update(ctx, id, updates)
	if err != nil {
		return err
	}

	if token, ok := t.syncGet(id); ok {
		updates.Apply(token)
	}

	return nil
}

func (t *TokenCache) Delete(ctx context.Context, id string) error {
	err := t.supplier.Delete(ctx, id)
	if err != nil {
		return err
	}

	t.syncRemove(id)
	return nil
}

func (t *TokenCache) Count(ctx context.Context) (int64, error) {
	return t.supplier.Count(ctx)
}

func (t *TokenCache) Find(ctx context.Context, query find.Query[*types.Token], options find.Options[*types.Token]) (database.Iterator[*types.Token], error) {
	return t.supplier.Find(ctx, query, options)
}

func (t *TokenCache) SetSize(size int) {
	t.cacheMux.Lock()
	t.cache.Resize(size)
	t.cacheMux.Unlock()
}

func (t *TokenCache) Evict(id string) {
	t.syncRemove(id)
}

func (t *TokenCache) Clear() {
	t.cacheMux.Lock()
	t.cache = cache.New[string, *types.Token](t.cache.Size())
	t.cacheMux.Unlock()
}

func (t *TokenCache) syncGet(id string) (*types.Token, bool) {
	t.cacheMux.Lock()
	token, ok := t.cache.Get(id)
	t.cacheMux.Unlock()
	if !ok {
		return nil, false
	}
	return token, true
}

func (t *TokenCache) syncPut(token *types.Token) {
	t.cacheMux.Lock()
	t.cache.Put(token.ID, token)
	t.cacheMux.Unlock()
}

func (t *TokenCache) syncRemove(id string) {
	t.cacheMux.Lock()
	t.cache.Remove(id)
	t.cacheMux.Unlock()
}
