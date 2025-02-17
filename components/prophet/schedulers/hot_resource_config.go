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

package schedulers

import (
	"sync"
	"time"

	"github.com/matrixorigin/matrixcube/components/prophet/schedule"
	"github.com/matrixorigin/matrixcube/components/prophet/statistics"
	"github.com/matrixorigin/matrixcube/components/prophet/storage"
)

// params about hot resource.
func initHotResourceScheduleConfig() *hotResourceSchedulerConfig {
	return &hotResourceSchedulerConfig{
		MinHotByteRate:        100,
		MinHotKeyRate:         10,
		MaxZombieRounds:       3,
		ByteRateRankStepRatio: 0.05,
		KeyRateRankStepRatio:  0.05,
		CountRankStepRatio:    0.01,
		GreatDecRatio:         0.95,
		MinorDecRatio:         0.99,
		MaxPeerNum:            1000,
		SrcToleranceRatio:     1.05, // Tolerate 5% difference
		DstToleranceRatio:     1.05, // Tolerate 5% difference
	}
}

type hotResourceSchedulerConfig struct {
	sync.RWMutex
	storage storage.Storage

	MinHotByteRate  float64 `json:"min-hot-byte-rate"`
	MinHotKeyRate   float64 `json:"min-hot-key-rate"`
	MaxZombieRounds int     `json:"max-zombie-rounds"`
	MaxPeerNum      int     `json:"max-peer-number"`

	// rank step ratio decide the step when calculate rank
	// step = max current * rank step ratio
	ByteRateRankStepRatio float64 `json:"byte-rate-rank-step-ratio"`
	KeyRateRankStepRatio  float64 `json:"key-rate-rank-step-ratio"`
	CountRankStepRatio    float64 `json:"count-rank-step-ratio"`
	GreatDecRatio         float64 `json:"great-dec-ratio"`
	MinorDecRatio         float64 `json:"minor-dec-ratio"`
	SrcToleranceRatio     float64 `json:"src-tolerance-ratio"`
	DstToleranceRatio     float64 `json:"dst-tolerance-ratio"`
}

func (conf *hotResourceSchedulerConfig) EncodeConfig() ([]byte, error) {
	conf.RLock()
	defer conf.RUnlock()
	return schedule.EncodeConfig(conf)
}

func (conf *hotResourceSchedulerConfig) GetMaxZombieDuration() time.Duration {
	conf.RLock()
	defer conf.RUnlock()
	return time.Duration(conf.MaxZombieRounds) * statistics.ContainerHeartBeatReportInterval * time.Second
}

func (conf *hotResourceSchedulerConfig) GetMaxPeerNumber() int {
	conf.RLock()
	defer conf.RUnlock()
	return conf.MaxPeerNum
}

func (conf *hotResourceSchedulerConfig) GetSrcToleranceRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.SrcToleranceRatio
}

func (conf *hotResourceSchedulerConfig) SetSrcToleranceRatio(tol float64) {
	conf.Lock()
	defer conf.Unlock()
	conf.SrcToleranceRatio = tol
}

func (conf *hotResourceSchedulerConfig) GetDstToleranceRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.DstToleranceRatio
}

func (conf *hotResourceSchedulerConfig) SetDstToleranceRatio(tol float64) {
	conf.Lock()
	defer conf.Unlock()
	conf.DstToleranceRatio = tol
}

func (conf *hotResourceSchedulerConfig) GetByteRankStepRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.ByteRateRankStepRatio
}

func (conf *hotResourceSchedulerConfig) GetKeyRankStepRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.KeyRateRankStepRatio
}

func (conf *hotResourceSchedulerConfig) GetCountRankStepRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.CountRankStepRatio
}

func (conf *hotResourceSchedulerConfig) GetGreatDecRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.GreatDecRatio
}

func (conf *hotResourceSchedulerConfig) GetMinorGreatDecRatio() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.MinorDecRatio
}

func (conf *hotResourceSchedulerConfig) GetMinHotKeyRate() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.MinHotKeyRate
}

func (conf *hotResourceSchedulerConfig) GetMinHotByteRate() float64 {
	conf.RLock()
	defer conf.RUnlock()
	return conf.MinHotByteRate
}
