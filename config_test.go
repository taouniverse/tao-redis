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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
)

func TestConfig(t *testing.T) {
	r := &Config{
		BaseMultiConfig: tao.BaseMultiConfig[InstanceConfig]{
			Instances: map[string]InstanceConfig{
				"default": {},
			},
		},
	}
	r.ValidSelf()

	instance := r.Instances["default"]
	assert.Equal(t, []string{"localhost:6379"}, instance.Addrs)
	assert.Equal(t, 50, instance.MaxPoolSize)

	t.Log(r.RunAfter())
	t.Log(r.ToTask())
}

func TestConfigDefaultInstanceName(t *testing.T) {
	c := &Config{}
	assert.Equal(t, "default", c.GetDefaultInstanceName())

	c.Default = "master"
	assert.Equal(t, "master", c.GetDefaultInstanceName())
}
