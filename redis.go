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
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/taouniverse/tao"
)

/**
import _ "github.com/taouniverse/tao-redis"
*/

// R is the global config instance for tao-redis
var R = &Config{}

// Factory is the global factory instance for managing redis.UniversalClient
var Factory *tao.BaseFactory[redis.UniversalClient]

func init() {
	var err error
	Factory, err = tao.Register(ConfigKey, R, NewRedis)
	if err != nil {
		panic(err.Error())
	}
}

// NewRedis creates a new Redis client for factory pattern
func NewRedis(name string, config InstanceConfig) (redis.UniversalClient, func() error, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            config.Addrs,
		DB:               config.DB,
		Username:         config.Username,
		Password:         config.Password,
		SentinelPassword: config.SentinelPassword,
		PoolSize:         config.MaxPoolSize,
		MinIdleConns:     config.MinPoolSize,
		MasterName:       config.MasterName,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		return nil, nil, tao.NewErrorWrapped("redis: rdb ping error", err)
	}

	closer := func() error {
		return client.Close()
	}

	return client, closer, nil
}

// Rdb returns the default redis client instance
func Rdb() (redis.UniversalClient, error) {
	return Factory.Get(R.GetDefaultInstanceName())
}

// GetRdb returns the redis client instance by name
func GetRdb(name string) (redis.UniversalClient, error) {
	return Factory.Get(name)
}
