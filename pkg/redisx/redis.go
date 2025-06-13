// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

package redisx

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Addrs        []string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
}

type RedisClient struct {
	client redis.UniversalClient
}

var Redis *RedisClient

// 支持单机版、主从模式和集群模式
func InitRedis(config RedisConfig) *RedisClient {
	var client redis.UniversalClient

	options := &redis.Options{
		Addr:         config.Addrs[0],
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout * time.Second,
		ReadTimeout:  config.ReadTimeout * time.Second,
		WriteTimeout: config.WriteTimeout * time.Second,
		MaxRetries:   config.MaxRetries,
	}

	if len(config.Addrs) == 1 {
		// 单机版
		client = redis.NewClient(options)
	} else {
		// 主从模式或集群模式
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        config.Addrs,
			Password:     config.Password,
			DB:           config.DB,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			DialTimeout:  config.DialTimeout * time.Second,
			ReadTimeout:  config.ReadTimeout * time.Second,
			WriteTimeout: config.WriteTimeout * time.Second,
			MaxRetries:   config.MaxRetries,
		})
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redis connect failed: %v\n", err)
	}

	Redis = &RedisClient{client: client}

	return Redis
}

// Set 设置键值对
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get 获取键的值
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

// Incr 原子递增
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr 原子递减
func (r *RedisClient) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// HSet 设置哈希字段，支持两种方式：
// 1. 传入字段值对：HSet(ctx, "user:1001", "name", "test", "age", "20")
// 2. 传入map：HSet(ctx, "user:1001", map[string]interface{}{"name": "test", "age": 20})
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	if len(values) == 0 {
		return nil
	}

	// 如果第一个参数是map，则转换为字段值对
	if len(values) == 1 {
		if m, ok := values[0].(map[string]interface{}); ok {
			var pairs []interface{}
			for k, v := range m {
				pairs = append(pairs, k, v)
			}
			return r.client.HSet(ctx, key, pairs...).Err()
		}
	}

	return r.client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段的值
// 示例：value, err := HGet(ctx, "user:1001", "name")
func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

// HGetList 获取哈希表的所有字段和值
// 示例：values, err := HGetList(ctx, "user:1001")
func (r *RedisClient) HGetList(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err == redis.Nil {
		return make(map[string]string), nil
	}
	return result, err
}

// LPush 向列表左侧推入元素
func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPop 从列表右侧弹出元素
func (r *RedisClient) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// Close 关闭客户端连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Lock 尝试获取一个分布式锁
func (r *RedisClient) Lock(ctx context.Context, lockKey string, lockValue string, lockExpireTime time.Duration) (bool, error) {
	// lockExpireTime 锁的过期时间, 必须设置，防止死锁
	if lockExpireTime <= 0 {
		lockExpireTime = 10 * time.Second
	}

	// SetNX 设置一个键值对，如果键不存在，则设置键值对，并返回true，否则返回false
	if flag, err := r.client.SetNX(ctx, lockKey, lockValue, lockExpireTime).Result(); err != nil {
		return false, err
	} else {
		return flag, nil
	}
}

// Unlock 释放分布式锁
func (r *RedisClient) Unlock(ctx context.Context, lockKey string, currentProcessLockValue string) error {
	// 获取锁的值, 检查是否是当前线程持有的锁
	val, err := r.client.Get(ctx, lockKey).Result()
	if err != nil {
		return err
	}

	if val == currentProcessLockValue {
		// 删除锁
		_, err = r.client.Del(ctx, lockKey).Result()
	}

	return err
}

func (r *RedisClient) AddHook(hook redis.Hook) {
	r.client.AddHook(hook)
}
