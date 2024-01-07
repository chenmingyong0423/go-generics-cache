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

package types

import "context"

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
