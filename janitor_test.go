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
	"sync/atomic"
	"testing"
	"time"
)

func Test_janitor(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	j := newJanitor(ctx, time.Millisecond)
	doneFlag := make(chan struct{})
	j.done = doneFlag
	num := int64(0)

	j.run(func(_ context.Context) {
		atomic.AddInt64(&num, 1)
	})

	time.Sleep(5 * time.Millisecond)
	cancelFunc()

	select {
	case <-doneFlag:
		t.Log("done")
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}

	if atomic.LoadInt64(&num) < 1 {
		t.Fatalf("failed to run cleanup function, num: %d", num)
	}
}
