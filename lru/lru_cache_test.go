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

package lru

import (
	"context"
	"testing"

	cacheError "github.com/chenmingyong0423/go-generics-cache/error"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	cache := NewCache[string, int](10)
	assert.Equal(t, 10, cache.maxEntries)
	assert.NotNil(t, cache.cache)
}

func TestCache_Set(t *testing.T) {
	testCases := []struct {
		name  string
		cache func(t *testing.T) *Cache[string, int]
		key   string
		value int

		wantKeys  []string
		wantError error
	}{
		{
			name: "set a new key",
			cache: func(_ *testing.T) *Cache[string, int] {
				return NewCache[string, int](1)
			},
			key:   "1",
			value: 1,
			wantKeys: []string{
				"1",
			},
		},
		{
			name: "set a existing key",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:   "1",
			value: 1,
			wantKeys: []string{
				"1",
			},
		},
		{
			name: "set a new key with a full cache",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:   "2",
			value: 2,
			wantKeys: []string{
				"2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache(t)
			err := cache.Set(context.Background(), tc.key, tc.value)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantKeys, cache.Keys())
		})
	}
}

func TestCache_Get(t *testing.T) {
	testCases := []struct {
		name      string
		cache     func(t *testing.T) *Cache[string, int]
		key       string
		wantValue int
		wantError error
	}{
		{
			name: "get a existing key",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:       "1",
			wantValue: 1,
		},
		{
			name: "get a non-existing key",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:       "2",
			wantValue: 0,
			wantError: cacheError.ErrNoKey,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache(t)
			got, err := cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantValue, got)
		})
	}
}

func TestCache_Delete(t *testing.T) {
	testCases := []struct {
		name      string
		cache     func(t *testing.T) *Cache[string, int]
		key       string
		wantKeys  []string
		wantError error
	}{
		{
			name: "delete a existing key",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:      "1",
			wantKeys: []string{},
		},
		{
			name: "delete a non-existing key",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](1)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				return cache
			},
			key:       "2",
			wantKeys:  []string{"1"},
			wantError: cacheError.ErrNoKey,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache(t)
			err := cache.Delete(context.Background(), tc.key)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantKeys, cache.Keys())
		})
	}
}

func TestCache_Keys(t *testing.T) {
	testCases := []struct {
		name  string
		cache func(t *testing.T) *Cache[string, int]
		want  []string
	}{
		{
			name: "get keys from empty cache",
			cache: func(t *testing.T) *Cache[string, int] {
				return NewCache[string, int](1)
			},
			want: []string{},
		},
		{
			name: "get keys from non-empty cache",
			cache: func(t *testing.T) *Cache[string, int] {
				cache := NewCache[string, int](4)
				err := cache.Set(context.Background(), "1", 1)
				assert.NoError(t, err)
				err = cache.Set(context.Background(), "3", 3)
				assert.NoError(t, err)
				err = cache.Set(context.Background(), "4", 4)
				assert.NoError(t, err)
				err = cache.Set(context.Background(), "2", 2)
				assert.NoError(t, err)
				return cache
			},
			want: []string{
				"1",
				"3",
				"4",
				"2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cache := tc.cache(t)
			got := cache.Keys()
			assert.Equal(t, tc.want, got)
		})
	}
}
