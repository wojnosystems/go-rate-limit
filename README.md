# Overview

An example Rate Limiter library used to control the rate that events occur, but these can also be used as thresholds that should replenish over time, such as error rates in circuit breakers.

# How to use

`go get -u github.com/wojnosystems/go-rate-limit`

# Examples

## Regular Token Bucket

```go
package main

import (
	"github.com/wojnosystems/go-rate-limit/rateLimit"
	"log"
	"time"
)

const (
	actionCost = 2
)

func main() {
	limiter := rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
		Capacity:             10,
		TokensAddedPerSecond: 10,
		InitialTokens:        5,
	})

	for {
		if !limiter.Allowed(actionCost) {
			break
		}
		log.Println("allowed!")
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("done, finally not allowed")
}
```

In this example above, the rate limiter will fill at a rate of 10 tokens per second, up to a maximum of 10 tokens, but has a reserve to 5 tokens to start. When Allowed is called and sufficient tokens are available, it will print "allowed!". However, because our actionCost is 2, each call to Allowed will consume 2 tokens instead of just 1. If you specify 0 as the cost, it will consume no tokens and always allow.

Because we're consuming 2 tokens every 100ms, but we replenish 1 token every 100ms, we consume the initial tokens and 2 additionally generated tokens before outpacing the replenish rate. The application prints:

```text
2021/11/28 23:11:34 allowed!
2021/11/28 23:11:34 allowed!
2021/11/28 23:11:34 allowed!
2021/11/28 23:11:34 allowed!
2021/11/28 23:11:34 done, finally not allowed
```

## Bursting Token Bucket

This is just like the regular TokenBucket, but it can optionally burst over the limit and refill more slowly. This bucket also supports non-1 tokenCosts.

```go
package main

import (
	"github.com/wojnosystems/go-rate-limit/rateLimit"
	"log"
	"time"
)

const (
	actionCost = 2
)

func main() {
	limiter := rateLimit.NewBurstingTokenBucket(
		rateLimit.BurstingTokenBucketOpts{
			Bucket: rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
				Capacity:             10,
				TokensAddedPerSecond: 10,
				InitialTokens:        5,
			}),
			Burst:  rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
				Capacity:             5,
				TokensAddedPerSecond: 1,
				InitialTokens:        1,
			}),
		})

	for {
		if !limiter.Allowed(actionCost) {
			break
		}
		log.Println("allowed!")
		time.Sleep(100 * time.Millisecond)
	}
	log.Println("done, finally not allowed")
}
```

It will output:

```text
2021/11/28 23:53:18 allowed!
2021/11/28 23:53:18 allowed!
2021/11/28 23:53:18 allowed!
2021/11/28 23:53:18 allowed!
2021/11/28 23:53:18 allowed!
2021/11/28 23:53:19 done, finally not allowed
```

You can see that, compared to the regular TokenBucket example, we were able to burst out an additional "allowed!" thanks to our bursting bucket.

Usually, your main bucket fills a larger capacity quickly and your burst bucket fills a smaller capacity more slowly. That way, over time, you can burst, but it's smoothed out due to the slower re-generation rates.

For example, if your API allows 20 requests per second, with a burst of an additional 5 every 30 seconds, you could set up the BurstingTokenBucket:

```go
package main

import (
	"github.com/wojnosystems/go-rate-limit/rateLimit"
)

func main() {
	limiter := rateLimit.NewBurstingTokenBucket(
		rateLimit.BurstingTokenBucketOpts{
			Bucket: rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
				Capacity:             20,
				TokensAddedPerSecond: 20,
			}),
			Burst:  rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
				Capacity:             5,
				TokensAddedPerSecond: 5/30,
			}),
		})
		
    // ...
}
```

This means that every 6 seconds, a new bursting token will be available, while 120 new regular tokens will be available to use (with a maximum of 20 at any given second).
