// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fifo

import (
	"container/list"
	"context"

	cacheError "github.com/chenmingyong0423/go-generics-cache/error"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewCache[K comparable, V any](cap int) *Cache[K, V] {
	return &Cache[K, V]{
		maxEntries:       cap,
		cache:            make(map[K]*list.Element, cap),
		linkedDoublyList: list.New(),
	}
}

type Cache[K comparable, V any] struct {
	maxEntries       int
	cache            map[K]*list.Element
	linkedDoublyList *list.List
}

func (c *Cache[K, V]) Set(_ context.Context, key K, value V) error {
	if e, ok := c.cache[key]; ok {
		// 元素存在
		c.linkedDoublyList.MoveToBack(e)
		e.Value.(*entry[K, V]).value = value
		return nil
	}
	// 元素不存在
	if c.linkedDoublyList.Len() >= c.maxEntries {
		e := c.linkedDoublyList.Front()
		c.linkedDoublyList.Remove(e)
		delete(c.cache, e.Value.(*entry[K, V]).key)
	}
	e := &entry[K, V]{
		key:   key,
		value: value,
	}
	c.cache[key] = c.linkedDoublyList.PushBack(e)
	return nil
}

func (c *Cache[K, V]) Get(_ context.Context, key K) (v V, err error) {
	if e, ok := c.cache[key]; ok {
		return e.Value.(*entry[K, V]).value, nil
	}
	return v, cacheError.ErrNoKey
}

func (c *Cache[K, V]) Delete(_ context.Context, key K) error {
	if e, ok := c.cache[key]; ok {
		c.linkedDoublyList.Remove(e)
		delete(c.cache, key)
		return nil
	}
	return cacheError.ErrNoKey
}

func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, 0)
	for e := c.linkedDoublyList.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry[K, V]).key)
	}
	return keys
}
