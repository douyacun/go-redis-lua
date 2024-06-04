# go-redis-lua
go redis lua常用脚本封装


### 令牌桶算法

测试每3秒5次限流，更多测试实例移步[token_bucket_limiter_test.go](token_bucket_limiter_test.go)
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

```shell
    token_bucket_limiter_test.go:66: 第0秒
    token_bucket_limiter_test.go:71: i: 3 allowed: true remain: 4 err <nil>
    token_bucket_limiter_test.go:71: i: 8 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 9 allowed: true remain: 3 err <nil>
    token_bucket_limiter_test.go:71: i: 7 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 5 allowed: true remain: 5 err <nil>
    token_bucket_limiter_test.go:71: i: 6 allowed: true remain: 1 err <nil>
    token_bucket_limiter_test.go:71: i: 2 allowed: true remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 1 allowed: true remain: 2 err <nil>
    token_bucket_limiter_test.go:71: i: 10 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 4 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:66: 第3秒
    token_bucket_limiter_test.go:71: i: 14 allowed: true remain: 4 err <nil>
    token_bucket_limiter_test.go:71: i: 13 allowed: true remain: 3 err <nil>
    token_bucket_limiter_test.go:71: i: 11 allowed: true remain: 2 err <nil>
    token_bucket_limiter_test.go:71: i: 17 allowed: true remain: 1 err <nil>
    token_bucket_limiter_test.go:71: i: 19 allowed: true remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 18 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 15 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 16 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 12 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 20 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:66: 第6秒
    token_bucket_limiter_test.go:71: i: 29 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 28 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 27 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 24 allowed: true remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 23 allowed: true remain: 1 err <nil>
    token_bucket_limiter_test.go:71: i: 22 allowed: true remain: 2 err <nil>
    token_bucket_limiter_test.go:71: i: 25 allowed: true remain: 3 err <nil>
    token_bucket_limiter_test.go:71: i: 21 allowed: true remain: 4 err <nil>
    token_bucket_limiter_test.go:71: i: 30 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 26 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:66: 第9秒
    token_bucket_limiter_test.go:71: i: 34 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 33 allowed: true remain: 4 err <nil>
    token_bucket_limiter_test.go:71: i: 32 allowed: true remain: 3 err <nil>
    token_bucket_limiter_test.go:71: i: 35 allowed: true remain: 1 err <nil>
    token_bucket_limiter_test.go:71: i: 38 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 37 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 36 allowed: true remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 39 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 40 allowed: false remain: 0 err <nil>
    token_bucket_limiter_test.go:71: i: 31 allowed: true remain: 2 err <nil>
```


