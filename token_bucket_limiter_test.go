package go_redis_lua

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

// TestTokenBucketLimit_Allowed_1s3 测试每1秒3次限流
func TestTokenBucketLimit_Allowed_1s3(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	limiter := NewTokenBucketLimit(&TokenBucketLimitOption{
		RDB:      rdb,
		Count:    3,
		Duration: time.Second,
		Burst:    3,
	})

	wg := sync.WaitGroup{}
	ctx := context.Background()
	for i := 1; i <= 30; i++ {
		i := i
		wg.Add(1)
		if i%10 == 0 {
			t.Logf("第%d秒", i/10)
		}
		go func() {
			allowed, remain, err := limiter.Allowed(ctx, 1)
			t.Logf("i: %d allowed: %v remain: %d err %v", i, allowed, remain, err)
			wg.Done()
		}()
		if i%10 == 0 {
			time.Sleep(1 * time.Second)
		}
	}
	wg.Wait()
}

// TestTokenBucketLimit_Allowed_3s5 测试每3秒5次限流
func TestTokenBucketLimit_Allowed_3s5(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	limiter := NewTokenBucketLimit(&TokenBucketLimitOption{
		RDB:      rdb,
		Count:    5,
		Duration: time.Second * 3,
		Burst:    10,
	})

	wg := sync.WaitGroup{}
	ctx := context.Background()
	ch := make(chan struct{}, 10)
	for i := 1; i <= 100; i++ {
		i := i
		wg.Add(1)
		if i%10 == 0 {
			t.Logf("第%d秒", (i/10-1)*3)
		}
		go func() {
			ch <- struct{}{}
			allowed, remain, err := limiter.Allowed(ctx, 1)
			t.Logf("i: %d allowed: %v remain: %d err %v", i, allowed, remain, err)
			wg.Done()
			<-ch
		}()
		if i%10 == 0 {
			time.Sleep(3 * time.Second)
		}
	}
	wg.Wait()
}

// TestTokenBucketLimit_Allowed_1m5 测试每分钟5次限流
func TestTokenBucketLimit_Allowed_1m5(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	limiter := NewTokenBucketLimit(&TokenBucketLimitOption{
		RDB:      rdb,
		Count:    5,
		Duration: time.Minute,
		Burst:    5,
	})
	wg := sync.WaitGroup{}
	ctx := context.Background()
	ch := make(chan struct{}, 10)

	n := time.Duration(5)

	t.Logf("第0秒")
	for i := 1; i <= 10; i++ {
		i := i
		wg.Add(1)
		go func() {
			ch <- struct{}{}
			allowed, remain, err := limiter.Allowed(ctx, 1)
			t.Logf("i: %d allowed: %v remain: %d err %v", i, allowed, remain, err)
			wg.Done()
			<-ch
		}()
	}
	wg.Wait()
	time.Sleep(time.Second * n)

	for i := 1; i <= 10; i++ {
		t.Logf("第%d秒\n", n*time.Duration(i))
		allowed, remain, err := limiter.Allowed(ctx, 1)
		t.Logf("i: %d allowed: %v remain: %d err %v", i, allowed, remain, err)
		time.Sleep(time.Second * n)
	}
}
