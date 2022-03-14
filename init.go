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
	"encoding/json"
	"github.com/taouniverse/tao"
)

/**
import _ "github.com/taouniverse/tao-redis"
*/

// R config of redis
var R *Config

func init() {
	err := tao.Register(ConfigKey, func() error {
		// 1. transfer config bytes to object
		bytes, err := tao.GetConfigBytes(ConfigKey)
		if err != nil {
			R = R.Default().(*Config)
		} else {
			err = json.Unmarshal(bytes, &R)
			if err != nil {
				return err
			}
		}

		// redis config
		R.ValidSelf()

		// 2. set object to tao
		err = tao.SetConfig(ConfigKey, R)
		if err != nil {
			return err
		}

		// redis setup
		return setup()
	})
	if err != nil {
		panic(err.Error())
	}
}
