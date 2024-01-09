// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	cacheError "github.com/chenmingyong0423/go-generics-cache/error"
	"github.com/chenmingyong0423/go-generics-cache/simple"
)

var _ ICache[int, any] = (*simple.Cache[int, any])(nil)

// ICache defines an interface for a key-value cache.
type ICache[K comparable, V any] interface {

	// Set stores the given key-value pair in the cache.
	Set(ctx context.Context, key K, value V) error

	// Get retrieves the value associated with the given key from the cache.
	Get(ctx context.Context, key K) (V, error)

	// Delete removes the value associated with the given key from the cache.
	Delete(ctx context.Context, key K) error

	Keys() []K
}

type Cache[K comparable, V any] struct {
	cache ICache[K, *Item[V]]
	mutex sync.RWMutex

	janitor *janitor
}

func NewSimpleCache[K comparable, V any](ctx context.Context, size int, interval time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		cache:   simple.NewCache[K, *Item[V]](size),
		janitor: newJanitor(ctx, interval),
	}
	cache.janitor.run(cache.DeleteExpired)
	return cache
}

type ItemOption func(*itemOptions)

type itemOptions struct {
	expiration time.Time
}

func WithExpiration(exp time.Duration) ItemOption {
	return func(o *itemOptions) {
		o.expiration = time.Now().Add(exp)
	}
}

type Item[V any] struct {
	value      V
	expiration time.Time
}

func newItem[V any](value V, opts ...ItemOption) *Item[V] {
	var item = &itemOptions{}
	for _, opt := range opts {
		opt(item)
	}
	return &Item[V]{
		value:      value,
		expiration: item.expiration,
	}
}

func (i *Item[V]) Expired() bool {
	return !i.expiration.IsZero() && i.expiration.Before(time.Now())
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (v V, err error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	item, err := c.cache.Get(ctx, key)
	if err != nil {
		return
	}
	if item.Expired() {
		return v, cacheError.ErrNoKey
	}
	return item.value, nil
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V, opts ...ItemOption) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item := newItem[V](value, opts...)
	return c.cache.Set(ctx, key, item)
}

func (c *Cache[K, V]) SetNX(ctx context.Context, key K, value V, opts ...ItemOption) (b bool, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err = c.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cacheError.ErrNoKey) {
			item := newItem[V](value, opts...)
			return true, c.cache.Set(ctx, key, item)
		}
		return false, err
	}
	return false, nil
}

func (c *Cache[K, V]) Delete(ctx context.Context, key K) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.cache.Delete(ctx, key)
}

func (c *Cache[K, V]) Keys() []K {
	return c.cache.Keys()
}

func (c *Cache[K, V]) DeleteExpired(ctx context.Context) {
	c.mutex.RLock()
	keys := c.Keys()
	c.mutex.RUnlock()
	i := 0
	for _, key := range keys {
		if i > 10000 {
			return
		}
		c.mutex.Lock()
		if item, err := c.cache.Get(ctx, key); err == nil && item.Expired() {
			_ = c.cache.Delete(ctx, key)
		}
		c.mutex.Unlock()
		i++
	}
}
