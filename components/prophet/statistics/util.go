// Copyright 2020 PingCAP, Inc.
// Modifications copyright (C) 2021 MatrixOrigin.
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

package statistics

import (
	"fmt"
)

const (
	// ContainerHeartBeatReportInterval is the heartbeat report interval of a container.
	ContainerHeartBeatReportInterval = 10
	// ResourceHeartBeatReportInterval is the heartbeat report interval of a resource.
	ResourceHeartBeatReportInterval = 60
)

func containerTag(id uint64) string {
	return fmt.Sprintf("container-%d", id)
}
