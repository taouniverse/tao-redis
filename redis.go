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
	"github.com/go-redis/redis/v8"
	"github.com/taouniverse/tao"
	"time"
)

/**
import _ "github.com/taouniverse/tao-redis"
*/

// R config of redis
var R = new(Config)

func init() {
	err := tao.Register(ConfigKey, R, setup)
	if err != nil {
		panic(err.Error())
	}
}

// Rdb to describe redis db client
// 1. If the MasterName option is specified, a sentinel-backed FailoverClient is returned.
// 2. if the number of Addrs is two or more, a ClusterClient is returned.
// 3. Otherwise, a single-node Client is returned.
var Rdb redis.UniversalClient

// setup with redis config
// execute when init tao universe
func setup() (err error) {
	// setup
	Rdb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            R.Addrs,
		DB:               R.DB,
		Username:         R.Username,
		Password:         R.Password,
		SentinelPassword: R.SentinelPassword,
		PoolSize:         R.MaxPoolSize,
		MinIdleConns:     R.MinPoolSize,
		MasterName:       R.MasterName,
	})
	// ping pong
	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err = Rdb.Ping(timeout).Err()
	if err != nil {
		return tao.NewErrorWrapped("redis: rdb ping error", err)
	}
	return
}
