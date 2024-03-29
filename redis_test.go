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
	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
	"testing"
	"time"
)

func TestTao(t *testing.T) {
	err := tao.SetConfigPath("./test.yaml")
	assert.Nil(t, err)

	timeout, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	ping := Rdb.Ping(timeout)
	assert.Nil(t, ping.Err())

	err = tao.Run(nil, nil)
	assert.Nil(t, err)
}
