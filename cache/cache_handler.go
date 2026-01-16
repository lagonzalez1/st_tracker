package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

type CacheHandler struct {
	valkey valkey.Client
}

func New(valkey valkey.Client) *CacheHandler {
	return &CacheHandler{valkey: valkey}
}

const (
	dataTTL        = 10 * time.Minute      // cache TTL for the data
	lockTTL        = 10 * time.Second      // should cover worst-case DB fetch time
	waiterMax      = 2 * time.Second       // how long non-holders will poll
	initialBackoff = 50 * time.Millisecond // waiter polling backoff
	maxBackoff     = 200 * time.Millisecond
)

var unlockScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`

func (h *CacheHandler) LockKey(key string) string {
	return key + ":lock"
}

/*
Check the cache
Return byte, bool
*/
func (h *CacheHandler) CheckCache(ctx context.Context, key string) ([]byte, bool) {
	res := h.valkey.Do(ctx, h.valkey.B().Get().Key(key).Build())
	value, err := res.ToString()
	if err != nil {
		return nil, false
	}
	if err == nil {
		return []byte(value), true
	}
	return nil, false
}

/*
Create the lock using Nx()
Returns True if you are the lock
Returns False if you are not the lock
*/
func (h *CacheHandler) TryAcquireLock(ctx context.Context, lockKey string, token string) (bool, error) {
	res := h.valkey.B().Set().Key(lockKey).Value(token).Nx().Ex(10 * time.Second).Build()
	_, err := h.valkey.Do(ctx, res).ToString()
	if err == nil {
		return true, nil
	}
	return false, nil
}

func (h *CacheHandler) SetCacheByString(ctx context.Context, key string, payload string) (bool, error) {
	err := h.valkey.Do(ctx, h.valkey.B().Set().Key(key).Value(payload).Ex(10*time.Hour).Build()).Error()
	if err != nil {
		return false, fmt.Errorf("unable to set cache with key %s", key)
	}
	return true, nil
}

func (h *CacheHandler) SetCache(ctx context.Context, key string, payload any) (bool, error) {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}
	err = h.valkey.Do(ctx, h.valkey.B().Set().Key(key).Value(string(jsonString)).Ex(10*time.Hour).Build()).Error()
	if err != nil {
		return false, fmt.Errorf("unable to set cache with key %s", key)
	}
	return true, nil

}

func (h *CacheHandler) SafeUnlock(ctx context.Context, lockKey string, token string) {
	_ = h.valkey.Do(ctx, h.valkey.B().Eval().Script(unlockScript).Numkeys(1).Key(lockKey).Arg(token).Build()).Error()
}

func (h *CacheHandler) ClearCache(ctx context.Context, key string) {
	_ = h.valkey.Do(ctx, h.valkey.B().Del().Key(key).Build()).Error()
}

func (h *CacheHandler) WaitForCacheUpdate(ctx context.Context, key string) bool {
	deadline := time.Now().Add(waiterMax)
	sleep := initialBackoff
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return false
		case <-time.After(sleep):
			_, hit := h.CheckCache(ctx, key)
			if hit {
				return true
			}
			if sleep < maxBackoff {
				sleep *= 2
			}
		}
	}
	return false
}
