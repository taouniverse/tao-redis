# github.com/taouniverse/tao-redis

[![Go Report Card](https://goreportcard.com/badge/github.com/taouniverse/tao-redis)](https://goreportcard.com/report/github.com/taouniverse/tao-redis)
[![GoDoc](https://pkg.go.dev/badge/github.com/taouniverse/tao-redis?status.svg)](https://pkg.go.dev/github.com/taouniverse/tao-redis?tab=doc)

Tao Universe 组件单元（Unit），基于泛型工厂模式封装 **Redis** 缓存数据库。

## 安装

```bash
go get github.com/taouniverse/tao-redis
```

## 使用

### 导入

```go
import _ "github.com/taouniverse/tao-redis"
```

### 配置

```yaml
# 单实例配置
redis:
  addrs:
    - localhost:6379
  db: 0
  max_pool_size: 50
  min_pool_size: 5

# Sentinel 哨兵模式
redis:
  addrs:
    - sentinel-1:26379
    - sentinel-2:26379
  master_name: mymaster
  sentinel_password: sent_pass
  password: redis_pass
  db: 0

# 多实例配置
redis:
  default_instance: cache
  cache:
    addrs: ["localhost:6379"]
    db: 0
    max_pool_size: 50
  session:
    addrs: ["localhost:6380"]
    db: 1
    max_pool_size: 20
```

### 配置字段说明

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `addrs` | []string | `["localhost:6379"]` | Redis 地址列表 |
| `password` | string | - | 密码 |
| `db` | int | `0` | 数据库索引 |
| `max_pool_size` | int | `50` | 连接池最大连接数 |
| `min_pool_size` | int | `5` | 连接池最小连接数 |
| `max_retries` | int | `3` | 最大重试次数 |
| `dial_timeout` | duration | `5s` | 连接超时 |
| `read_timeout` | duration | `3s` | 读取超时 |
| `write_timeout` | duration | `3s` | 写入超时 |
| `master_name` | string | - | Sentinel 主节点名称 |
| `sentinel_password` | string | - | Sentinel 密码 |

## 工厂模式 API

| API | 说明 |
|-----|------|
| `redis.M` | 配置实例 `*Config` |
| `redis.Factory` | `*tao.BaseFactory[redis.UniversalClient]` 工厂实例 |
| `redis.Rdb()` | 获取默认 Redis 客户端 `(UniversalClient, error)` |
| `redis.GetRdb(name)` | 获取指定名称的客户端 `(UniversalClient, error)` |

## 使用示例

### 获取客户端并执行操作

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/taouniverse/tao-redis"
)

func main() {
    // 获取默认实例
    client, err := redis.Rdb()
    if err != nil {
        log.Fatal(err)
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Ping 测试
    err = client.Ping(ctx).Err()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Redis 连接成功")
}
```

### 基本操作

```go
client, _ := redis.Rdb()
ctx := context.Background()

// 设置键值
err := client.Set(ctx, "key", "value", time.Hour).Err()

// 获取值
val, err := client.Get(ctx, "key").Result()

// 删除键
err := client.Del(ctx, "key").Err()

// 设置哈希
err := client.HSet(ctx, "user:1", "name", "tao").Err()

// 获取哈希字段
name, err := client.HGet(ctx, "user:1", "name").Result()
```

### 多实例使用

```go
// 获取缓存实例
cache, _ := redis.GetRdb("cache")

// 获取会话实例
session, _ := redis.GetRdb("session")

// 缓存操作
cache.Set(ctx, "data", "value", time.Minute)

// 会话操作
session.Set(ctx, "session_id", "user_data", time.Hour)
```

## 单元测试

### 快速测试（无需 Docker）

```bash
# 仅运行配置相关测试
go test -v -run "TestConfig" ./...
```

### 完整集成测试（需要 Docker）

```bash
# 启动 Redis 并运行单实例测试
make test

# 启动 Redis 并运行多实例测试
make test-multi

# 启动 Redis 并运行所有测试
make test-all

# 生成覆盖率报告
make coverage

# 停止 Redis 服务
make down
```

### 手动测试

```bash
# 1. 启动 Redis
docker-compose up -d

# 2. 运行单实例测试
go test -v ./...

# 3. 运行多实例测试
TAO_TEST_MULTI_INSTANCE=true go test -v ./...

# 4. 停止 Redis
docker-compose down
```

## 开发指南

| 文件 | 说明 |
|------|------|
| `config.go` | InstanceConfig 字段 + ValidSelf 默认值 |
| `redis.go` | NewRedis 构造器 + 工厂注册 |
