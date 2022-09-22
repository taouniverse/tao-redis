// Copyright 2022 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"encoding/json"
	"github.com/taouniverse/tao"
)

// ConfigKey for this repo
const ConfigKey = "redis"

// Config implements tao.Config
type Config struct {
	Addrs            []string `json:"addrs"`
	MasterName       string   `json:"master_name,omitempty"`
	MaxPoolSize      int      `json:"max_pool_size"`
	MinPoolSize      int      `json:"min_pool_size"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	SentinelPassword string   `json:"sentinel_password,omitempty"`
	DB               int      `json:"db"`
	RunAfters        []string `json:"run_after,omitempty"`
}

var defaultRedis = &Config{
	Addrs:       []string{"localhost:6379"},
	MaxPoolSize: 50,
	MinPoolSize: 5,
	RunAfters:   []string{},
}

// Name of Config
func (r *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (r *Config) ValidSelf() {
	if r.Addrs == nil || len(r.Addrs) == 0 {
		r.Addrs = defaultRedis.Addrs
	}
	if r.MaxPoolSize == 0 {
		r.MaxPoolSize = defaultRedis.MaxPoolSize
	}
	if r.MinPoolSize == 0 {
		r.MinPoolSize = defaultRedis.MinPoolSize
	}
	if r.RunAfters == nil {
		r.RunAfters = defaultRedis.RunAfters
	}
}

// ToTask transform itself to Task
func (r *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			// print some info
			marshal, err := json.Marshal(Rdb.PoolStats())
			if err != nil {
				return param, err
			}
			tao.Debugf("redis pool stats: %s\n", string(marshal))
			return param, nil
		}, tao.SetClose(func() error {
			err := Rdb.Close()
			if err != nil {
				return tao.NewErrorWrapped("redis: rdb close.", err)
			}
			return nil
		}))
}

// RunAfter defines pre task names
func (r *Config) RunAfter() []string {
	return r.RunAfters
}
