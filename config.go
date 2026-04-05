// Copyright 2021-2026 huija
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

// InstanceConfig 单实例配置
type InstanceConfig struct {
	Addrs            []string `json:"addrs" yaml:"addrs"`
	MasterName       string   `json:"master_name,omitempty" yaml:"master_name,omitempty"`
	MaxPoolSize      int      `json:"max_pool_size" yaml:"max_pool_size"`
	MinPoolSize      int      `json:"min_pool_size" yaml:"min_pool_size"`
	Username         string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password         string   `json:"password,omitempty" yaml:"password,omitempty"`
	SentinelPassword string   `json:"sentinel_password,omitempty" yaml:"sentinel_password,omitempty"`
	DB               int      `json:"db" yaml:"db"`
}

// Config 总配置，实现 tao.MultiConfig 接口
type Config struct {
	tao.BaseMultiConfig[InstanceConfig]
	RunAfters []string `json:"run_after,omitempty" yaml:"run_after,omitempty"`
}

var defaultInstance = &InstanceConfig{
	Addrs:       []string{"localhost:6379"},
	MaxPoolSize: 50,
	MinPoolSize: 5,
}

// Name of Config
func (r *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (r *Config) ValidSelf() {
	for name, instance := range r.Instances {
		if len(instance.Addrs) == 0 {
			instance.Addrs = defaultInstance.Addrs
		}
		if instance.MaxPoolSize == 0 {
			instance.MaxPoolSize = defaultInstance.MaxPoolSize
		}
		if instance.MinPoolSize == 0 {
			instance.MinPoolSize = defaultInstance.MinPoolSize
		}
		r.Instances[name] = instance
	}
	if r.RunAfters == nil {
		r.RunAfters = []string{}
	}
}

// ToTask transform itself to Task
func (r *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			for name := range r.Instances {
				client, err := Factory.Get(name)
				if err != nil {
					return param, err
				}
				marshal, err := json.Marshal(client.PoolStats())
				if err != nil {
					return param, err
				}
				tao.Debugf("redis[%s] pool stats: %s\n", name, string(marshal))
			}
			return param, nil
		},
		tao.SetClose(func() error {
			return Factory.CloseAll()
		}))
}

// RunAfter defines pre task names
func (r *Config) RunAfter() []string {
	return r.RunAfters
}
