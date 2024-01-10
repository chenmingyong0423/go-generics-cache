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
	"testing"
	"time"

	cacheError "github.com/chenmingyong0423/go-generics-cache/error"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleCache(t *testing.T) {
	cache := NewSimpleCache[int, int](context.Background(), 0, 3*time.Second)
	assert.NotNil(t, cache)
	err := cache.Set(context.Background(), 1, 1, WithExpiration(time.Millisecond))
	assert.NoError(t, err)
	time.Sleep(5 * time.Millisecond)
	cache.DeleteExpired(context.Background())
}

func Test_newItem(t *testing.T) {

	testCases := []struct {
		name  string
		value int
		opts  []ItemOption
		now   time.Time

		want *Item[int]
	}{
		{
			name:  "Creates an item with only a value",
			value: 1,
			opts:  nil,
			want: &Item[int]{
				value: 1,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newItem(tt.value, tt.opts...))
		})
	}
}

func TestCache_Get(t *testing.T) {
	testCases := []struct {
		name      string
		cache     func(t *testing.T) *Cache[int, int]
		ctx       context.Context
		waitTime  time.Duration
		key       int
		wantValue int
		wantErr   error
	}{
		{
			name: "Lookup for non-existent key in empty cache",
			cache: func(t *testing.T) *Cache[int, int] {
				return NewSimpleCache[int, int](context.Background(), 0, time.Minute)
			},
			ctx:       context.Background(),
			key:       1,
			wantValue: 0,
			wantErr:   cacheError.ErrNoKey,
		},
		{
			name: "Lookup for non-existent key in non-empty cache",
			cache: func(t *testing.T) *Cache[int, int] {
				cache := NewSimpleCache[int, int](context.Background(), 0, time.Minute)
				assert.NoError(t, cache.Set(context.Background(), 1, 1))
				return cache
			},
			ctx:       context.Background(),
			key:       2,
			wantValue: 0,
			wantErr:   cacheError.ErrNoKey,
		},
		{
			name: "Lookup and match",
			cache: func(t *testing.T) *Cache[int, int] {
				cache := NewSimpleCache[int, int](context.Background(), 0, time.Minute)
				assert.NoError(t, cache.Set(context.Background(), 1, 1))
				return cache
			},
			ctx:       context.Background(),
			key:       1,
			wantValue: 1,
			wantErr:   nil,
		},
		{
			name: "Lookup the key after the key expires",
			cache: func(t *testing.T) *Cache[int, int] {
				cache := NewSimpleCache[int, int](context.Background(), 0, time.Minute)
				assert.NoError(t, cache.Set(context.Background(), 1, 1, WithExpiration(time.Second)))
				return cache
			},
			waitTime:  time.Second * 2,
			ctx:       context.Background(),
			key:       1,
			wantValue: 0,
			wantErr:   cacheError.ErrNoKey,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.cache(t)
			time.Sleep(tt.waitTime)
			got, err := cache.Get(tt.ctx, tt.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantValue, got)
		})
	}
}

func TestCache_Set(t *testing.T) {
	testCases := []struct {
		name   string
		cache  *Cache[int, int]
		ctx    context.Context
		keys   []int
		values []int

		wantKeys   []int
		wantValues []int
		wantErr    []error
	}{
		{
			name:   "first set",
			cache:  NewSimpleCache[int, int](context.Background(), 0, time.Minute),
			ctx:    context.Background(),
			keys:   []int{1},
			values: []int{1},

			wantKeys:   []int{1},
			wantValues: []int{1},
			wantErr:    []error{nil},
		},
		{
			name:   "set multiple keys",
			cache:  NewSimpleCache[int, int](context.Background(), 0, time.Minute),
			ctx:    context.Background(),
			keys:   []int{1, 2, 3},
			values: []int{1, 2, 3},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{1, 2, 3},
			wantErr:    []error{nil, nil, nil},
		},
		{
			name:   "set multiple keys with duplicates",
			cache:  NewSimpleCache[int, int](context.Background(), 0, time.Minute),
			ctx:    context.Background(),
			keys:   []int{1, 2, 3, 2},
			values: []int{1, 2, 3, 4},

			wantKeys:   []int{1, 2, 3},
			wantValues: []int{1, 4, 3},
			wantErr:    []error{nil, nil, nil, nil},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.keys {
				err := tt.cache.Set(tt.ctx, tt.keys[i], tt.values[i])
				assert.NoError(t, err)
			}
			keys := tt.cache.Keys()
			assert.ElementsMatch(t, tt.wantKeys, keys)
			for i := range tt.wantKeys {
				get, err := tt.cache.Get(tt.ctx, tt.keys[i])
				assert.Equal(t, tt.wantErr[i], err)
				assert.Equal(t, tt.wantValues[i], get)
			}
		})
	}
}

type errorCache[K comparable, V any] struct {
}

func (e errorCache[K, V]) Set(ctx context.Context, key K, value V) error {
	return nil
}

func (e errorCache[K, V]) Get(ctx context.Context, key K) (v V, err error) {
	return v, errors.New("get error")
}

func (e errorCache[K, V]) Delete(ctx context.Context, key K) error {
	return nil
}

func (e errorCache[K, V]) Keys() []K {
	return nil
}

func TestCache_SetNX(t *testing.T) {
	testCases := []struct {
		name   string
		cache  *Cache[int, int]
		ctx    context.Context
		keys   []int
		values []int

		wantBoolValues []bool
		wantErr        []error
	}{
		{
			name:           "not exist",
			cache:          NewSimpleCache[int, int](context.Background(), 0, time.Minute),
			ctx:            context.Background(),
			keys:           []int{1},
			values:         []int{1},
			wantBoolValues: []bool{true},
			wantErr:        []error{nil},
		},
		{
			name:           "exist",
			cache:          NewSimpleCache[int, int](context.Background(), 0, time.Minute),
			ctx:            context.Background(),
			keys:           []int{1, 1, 2},
			values:         []int{1, 2, 3},
			wantBoolValues: []bool{true, false, true},
			wantErr:        []error{nil, nil, nil},
		},
		{
			name:           "error",
			cache:          &Cache[int, int]{cache: &errorCache[int, *Item[int]]{}},
			ctx:            context.Background(),
			keys:           []int{1},
			values:         []int{1},
			wantBoolValues: []bool{false},
			wantErr:        []error{errors.New("get error")},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.keys {
				b, err := tt.cache.SetNX(tt.ctx, tt.keys[i], tt.values[i])
				assert.Equal(t, err, tt.wantErr[i])
				assert.Equal(t, b, tt.wantBoolValues[i])
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	testCases := []struct {
		name  string
		cache func(t *testing.T) *Cache[int, int]
		ctx   context.Context
		keys  int

		wantKeys []int
		wantErr  error
	}{
		{
			name: "Delete non-existent key from the empty cache",
			cache: func(t *testing.T) *Cache[int, int] {
				return NewSimpleCache[int, int](context.Background(), 0, time.Minute)
			},
			keys:     1,
			wantKeys: []int{},
			wantErr:  nil,
		},
		{
			name: "Delete non-existent key from the empty cache",
			cache: func(t *testing.T) *Cache[int, int] {
				cache := NewSimpleCache[int, int](context.Background(), 0, time.Minute)
				assert.NoError(t, cache.Set(context.Background(), 1, 1))
				return cache
			},
			keys:     2,
			wantKeys: []int{1},
			wantErr:  nil,
		},
		{
			name: "Delete existing keys from the cache",
			cache: func(t *testing.T) *Cache[int, int] {
				cache := NewSimpleCache[int, int](context.Background(), 0, time.Minute)
				assert.NoError(t, cache.Set(context.Background(), 1, 1))
				assert.NoError(t, cache.Set(context.Background(), 2, 2))
				return cache
			},
			keys:     1,
			wantKeys: []int{2},
			wantErr:  nil,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.cache(t)
			err := cache.Delete(tt.ctx, tt.keys)
			assert.Equal(t, tt.wantErr, err)
			keys := cache.Keys()
			assert.ElementsMatch(t, tt.wantKeys, keys)
		})
	}
}

func TestNewLruCache(t *testing.T) {
	cache := NewLruCache[int, int](context.Background(), 0, 3*time.Second)
	assert.NotNil(t, cache)
}
