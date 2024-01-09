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

package cache

import (
	"context"
	"sync"
	"time"
)

func newJanitor(ctx context.Context, interval time.Duration) *janitor {
	return &janitor{
		ctx:      ctx,
		interval: interval,
		done:     make(chan struct{}),
	}
}

type janitor struct {
	ctx      context.Context
	interval time.Duration
	done     chan struct{}
	once     sync.Once
}

func (j *janitor) stop() {
	j.once.Do(func() { close(j.done) })
}

func (j *janitor) run(cleanup func(ctx context.Context)) {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cleanup(j.ctx)
			case <-j.ctx.Done():
				j.stop()
			case <-j.done:
				cleanup(j.ctx)
				return
			}
		}
	}()
}
