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
	"strconv"

	"github.com/matrixorigin/matrixcube/components/prophet/config"
	"github.com/matrixorigin/matrixcube/components/prophet/core"
	"github.com/matrixorigin/matrixcube/components/prophet/pb/metapb"
)

const (
	unknown   = "unknown"
	labelType = "label"
)

type containerStatistics struct {
	opt             *config.PersistOptions
	Up              int
	Disconnect      int
	Unhealthy       int
	Down            int
	Offline         int
	Tombstone       int
	LowSpace        int
	StorageSize     uint64
	StorageCapacity uint64
	ResourceCount   int
	LeaderCount     int
	LabelCounter    map[string]int
}

func newContainerStatistics(opt *config.PersistOptions) *containerStatistics {
	return &containerStatistics{
		opt:          opt,
		LabelCounter: make(map[string]int),
	}
}

func (s *containerStatistics) Observe(container *core.CachedContainer, stats *ContainersStats) {
	for _, k := range s.opt.GetLocationLabels() {
		v := container.GetLabelValue(k)
		if v == "" {
			v = unknown
		}
		key := fmt.Sprintf("%s:%s", k, v)
		// exclude tombstone
		if container.GetState() != metapb.ContainerState_Tombstone {
			s.LabelCounter[key]++
		}
	}
	containerAddress := container.Meta.Addr()
	id := strconv.FormatUint(container.Meta.ID(), 10)
	// Container state.
	switch container.GetState() {
	case metapb.ContainerState_UP:
		if container.DownTime() >= s.opt.GetMaxContainerDownTime() {
			s.Down++
		} else if container.IsUnhealthy() {
			s.Unhealthy++
		} else if container.IsDisconnected() {
			s.Disconnect++
		} else {
			s.Up++
		}
	case metapb.ContainerState_Offline:
		s.Offline++
	case metapb.ContainerState_Tombstone:
		s.Tombstone++
		s.resetContainerStatistics(containerAddress, id)
		return
	}
	if container.IsLowSpace(s.opt.GetLowSpaceRatio(), s.opt.GetReplicationConfig().Groups) {
		s.LowSpace++
	}

	// Container stats.
	s.StorageSize += container.StorageSize()
	s.StorageCapacity += container.GetCapacity()

	var resScore, leaderScore, resSize, resCount, leaderSize, leaderCount float64
	for _, group := range s.opt.GetReplicationConfig().Groups {
		s.ResourceCount += container.GetResourceCount(group)
		s.LeaderCount += container.GetLeaderCount(group)

		resScore += container.ResourceScore(group, s.opt.GetResourceScoreFormulaVersion(), s.opt.GetHighSpaceRatio(), s.opt.GetLowSpaceRatio(), 0, 0)
		leaderScore += container.LeaderScore(group, s.opt.GetLeaderSchedulePolicy(), 0)
		resSize += float64(container.GetResourceSize(group))
		resCount += float64(container.GetResourceCount(group))
		leaderSize += float64(container.GetLeaderSize(group))
		leaderCount += float64(container.GetLeaderCount(group))
	}

	containerStatusGauge.WithLabelValues(containerAddress, id, "container_available").Set(float64(container.GetAvailable()))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_used").Set(float64(container.GetUsedSize()))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_capacity").Set(float64(container.GetCapacity()))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_available_avg").Set(float64(container.GetAvgAvailable()))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_available_deviation").Set(float64(container.GetAvailableDeviation()))
	containerStatusGauge.WithLabelValues(containerAddress, id, "resource_score").Set(resScore)
	containerStatusGauge.WithLabelValues(containerAddress, id, "leader_score").Set(leaderScore)
	containerStatusGauge.WithLabelValues(containerAddress, id, "resource_size").Set(resSize)
	containerStatusGauge.WithLabelValues(containerAddress, id, "resource_count").Set(resCount)
	containerStatusGauge.WithLabelValues(containerAddress, id, "leader_size").Set(leaderSize)
	containerStatusGauge.WithLabelValues(containerAddress, id, "leader_count").Set(leaderCount)

	// Container flows.
	containerFlowStats := stats.GetRollingContainerStats(container.Meta.ID())
	if containerFlowStats == nil {
		return
	}
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_write_rate_bytes").Set(containerFlowStats.GetLoad(ContainerWriteBytes))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_read_rate_bytes").Set(containerFlowStats.GetLoad(ContainerReadBytes))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_write_rate_keys").Set(containerFlowStats.GetLoad(ContainerWriteKeys))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_read_rate_keys").Set(containerFlowStats.GetLoad(ContainerReadKeys))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_cpu_usage").Set(containerFlowStats.GetLoad(ContainerCPUUsage))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_disk_read_rate").Set(containerFlowStats.GetLoad(ContainerDiskReadRate))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_disk_write_rate").Set(containerFlowStats.GetLoad(ContainerDiskWriteRate))

	containerStatusGauge.WithLabelValues(containerAddress, id, "container_write_rate_bytes_instant").Set(containerFlowStats.GetInstantLoad(ContainerWriteBytes))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_read_rate_bytes_instant").Set(containerFlowStats.GetInstantLoad(ContainerReadBytes))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_write_rate_keys_instant").Set(containerFlowStats.GetInstantLoad(ContainerWriteKeys))
	containerStatusGauge.WithLabelValues(containerAddress, id, "container_read_rate_keys_instant").Set(containerFlowStats.GetInstantLoad(ContainerReadKeys))
}

func (s *containerStatistics) Collect() {
	metrics := make(map[string]float64)
	metrics["container_up_count"] = float64(s.Up)
	metrics["container_disconnected_count"] = float64(s.Disconnect)
	metrics["container_down_count"] = float64(s.Down)
	metrics["container_unhealth_count"] = float64(s.Unhealthy)
	metrics["container_offline_count"] = float64(s.Offline)
	metrics["container_tombstone_count"] = float64(s.Tombstone)
	metrics["container_low_space_count"] = float64(s.LowSpace)
	metrics["resource_count"] = float64(s.ResourceCount)
	metrics["leader_count"] = float64(s.LeaderCount)
	metrics["storage_size"] = float64(s.StorageSize)
	metrics["storage_capacity"] = float64(s.StorageCapacity)

	for typ, value := range metrics {
		clusterStatusGauge.WithLabelValues(typ).Set(value)
	}

	// Current scheduling configurations of the cluster
	configs := make(map[string]float64)
	configs["leader-schedule-limit"] = float64(s.opt.GetLeaderScheduleLimit())
	configs["resource-schedule-limit"] = float64(s.opt.GetResourceScheduleLimit())
	configs["merge-schedule-limit"] = float64(s.opt.GetMergeScheduleLimit())
	configs["replica-schedule-limit"] = float64(s.opt.GetReplicaScheduleLimit())
	configs["max-replicas"] = float64(s.opt.GetMaxReplicas())
	configs["high-space-ratio"] = s.opt.GetHighSpaceRatio()
	configs["low-space-ratio"] = s.opt.GetLowSpaceRatio()
	configs["tolerant-size-ratio"] = s.opt.GetTolerantSizeRatio()
	configs["hot-resource-schedule-limit"] = float64(s.opt.GetHotResourceScheduleLimit())
	configs["hot-resource-cache-hits-threshold"] = float64(s.opt.GetHotResourceCacheHitsThreshold())
	configs["max-pending-peer-count"] = float64(s.opt.GetMaxPendingPeerCount())
	configs["max-snapshot-count"] = float64(s.opt.GetMaxSnapshotCount())
	configs["max-merge-resource-size"] = float64(s.opt.GetMaxMergeResourceSize())
	configs["max-merge-resource-keys"] = float64(s.opt.GetMaxMergeResourceKeys())

	var enableMakeUpReplica, enableRemoveDownReplica, enableRemoveExtraReplica, enableReplaceOfflineReplica float64
	if s.opt.IsMakeUpReplicaEnabled() {
		enableMakeUpReplica = 1
	}
	if s.opt.IsRemoveDownReplicaEnabled() {
		enableRemoveDownReplica = 1
	}
	if s.opt.IsRemoveExtraReplicaEnabled() {
		enableRemoveExtraReplica = 1
	}
	if s.opt.IsReplaceOfflineReplicaEnabled() {
		enableReplaceOfflineReplica = 1
	}

	configs["enable-makeup-replica"] = enableMakeUpReplica
	configs["enable-remove-down-replica"] = enableRemoveDownReplica
	configs["enable-remove-extra-replica"] = enableRemoveExtraReplica
	configs["enable-replace-offline-replica"] = enableReplaceOfflineReplica

	for typ, value := range configs {
		configStatusGauge.WithLabelValues(typ).Set(value)
	}

	for name, value := range s.LabelCounter {
		placementStatusGauge.WithLabelValues(labelType, name).Set(float64(value))
	}
}

func (s *containerStatistics) resetContainerStatistics(containerAddress string, id string) {
	metrics := []string{
		"resource_score",
		"leader_score",
		"resource_size",
		"resource_count",
		"leader_size",
		"leader_count",
		"container_available",
		"container_used",
		"container_capacity",
		"container_write_rate_bytes",
		"container_read_rate_bytes",
		"container_write_rate_keys",
		"container_read_rate_keys",
	}
	for _, m := range metrics {
		containerStatusGauge.DeleteLabelValues(containerAddress, id, m)
	}
}

type containerStatisticsMap struct {
	opt   *config.PersistOptions
	stats *containerStatistics
}

// NewContainerStatisticsMap create a container statistics map
func NewContainerStatisticsMap(opt *config.PersistOptions) *containerStatisticsMap {
	return &containerStatisticsMap{
		opt:   opt,
		stats: newContainerStatistics(opt),
	}
}

func (m *containerStatisticsMap) Observe(container *core.CachedContainer, stats *ContainersStats) {
	m.stats.Observe(container, stats)
}

func (m *containerStatisticsMap) Collect() {
	m.stats.Collect()
}

func (m *containerStatisticsMap) Reset() {
	containerStatusGauge.Reset()
	clusterStatusGauge.Reset()
	placementStatusGauge.Reset()
}
