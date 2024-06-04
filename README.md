# go-redis-lua
go redis lua常用脚本封装


### 令牌桶算法

TestTokenBucketLimit_Allowed_3s5 测试每3秒5次限流，更多测试实例请看[token_bucket_limiter_test.go](token_bucket_limiter_test.go)
```go 
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
```


