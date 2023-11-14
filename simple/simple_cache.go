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

package simple

import (
	"context"
	cacheError "github.com/chenmingyong0423/go-generics-cache/error"
)

type Cache[K comparable, V any] struct {
	cache map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		cache: make(map[K]V, 0),
	}
}

func (c *Cache[K, V]) Set(ctx context.Context, key K, value V) error {
	c.cache[key] = value
	return nil
}

func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, error) {
	var (
		value V
		ok    bool
	)
	if value, ok = c.cache[key]; !ok {
		return value, cacheError.ErrNoKey
	}
	return value, nil
}

func (c *Cache[K, V]) Delete(ctx context.Context, key K) error {
	delete(c.cache, key)
	return nil
}

func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, 0)
	for key := range c.cache {
		keys = append(keys, key)
	}
	return keys
}
