// Copyright 2020 MatrixOrigin.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package prophet

import (
	"context"
	"time"

	"github.com/matrixorigin/matrixcube/components/prophet/util"
)

// StartMonitor calls systimeErrHandler if system time jump backward.
func StartMonitor(ctx context.Context, now func() time.Time, systimeErrHandler func()) {
	util.GetLogger().Info("start system time monitor")
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		last := now().UnixNano()
		select {
		case <-tick.C:
			if now().UnixNano() < last {
				util.GetLogger().Errorf("system time jump backward, last %+v", last)
				systimeErrHandler()
			}
		case <-ctx.Done():
			return
		}
	}
}
