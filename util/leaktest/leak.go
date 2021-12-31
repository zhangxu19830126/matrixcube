// Copyright 2021 MatrixOrigin.
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

package leaktest

import (
	"os"
	"testing"

	ltlib "github.com/lni/goutils/leaktest"
)

func goroutineLeakCheckDisabled() bool {
	glc := os.Getenv("NO_GOROUTINE_LEAK_CHECK")
	return len(glc) > 0
}

func AfterTest(t testing.TB) func() {
	if goroutineLeakCheckDisabled() {
		return func() {}
	}
	return ltlib.AfterTest(t)
}
