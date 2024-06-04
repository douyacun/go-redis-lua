package go_redis_lua

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenBucketLimitOption struct {
	RDB      *redis.Client
	Count    int64         // 速率，请求数
	Duration time.Duration // 时长，比如：每秒3次/每分钟10次
	Burst    int64         // 令牌桶的大小
	Prefix   string        // redis key前缀
}

type TokenBucketLimiter struct {
	script *redis.Script
	option *TokenBucketLimitOption
}

func NewTokenBucketLimit(option *TokenBucketLimitOption) *TokenBucketLimiter {
	luaScript := `
local duration = tonumber(ARGV[1]) -- 多少微妙产生一条token
local capacity = tonumber(ARGV[2]) -- 容量
local now = tonumber(ARGV[3]) -- 当前时间
local requested = tonumber(ARGV[4]) -- 请求量

-- 上一次剩余令牌
local last_tokens = tonumber(redis.call("get", KEYS[1])) 
if last_tokens == nil then
	-- 第一次请求
    last_tokens = capacity 
end

-- 上一次请求的时间
local last_time = tonumber(redis.call("get", KEYS[2]))
if last_time == nil then
    last_time = 0
end

-- 时间差
local delta = math.max(0, now - last_time)

-- 剩余令牌 + 时间差 / 单个token消耗的时间
local new_tokens = math.min(capacity, last_tokens + delta/duration)

-- 返回 1:是否允许 2:剩余令牌
local res = {0, new_tokens} 
if new_tokens >= requested then
    res[2] = new_tokens - requested
	res[1] = 1
end

redis.call("set", KEYS[1], res[2])
redis.call("set", KEYS[2], now)

return res
`
	return &TokenBucketLimiter{
		script: redis.NewScript(luaScript),
		option: option,
	}
}

func (t *TokenBucketLimiter) Allowed(ctx context.Context, requested int64) (allowed bool, remainingTokens int64, err error) {
	// 计算每个token需要多少微妙
	duration := t.option.Duration.Microseconds() / t.option.Count
	//fmt.Printf("duration: %d\n", duration)
	now := time.Now().UnixMicro()
	// 调用 Lua 脚本
	result, err := t.script.Run(ctx, t.option.RDB, []string{t.option.Prefix + "tokens", t.option.Prefix + "last_time"}, duration, t.option.Burst, now, requested).Result()
	if err != nil {
		return false, 0, fmt.Errorf("TokenBucketLimiter Allowed error running script: %v", err)
	}
	// 解析返回结果
	resultArray := result.([]interface{})

	allowed = resultArray[0].(int64) == 1
	remainingTokens = resultArray[1].(int64)

	return allowed, remainingTokens, nil
}
