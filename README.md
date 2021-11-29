# Overview

An example Rate Limiter library used to control the rate that events occur, but these can also be used as thresholds that should replenish over time, such as error rates in circuit breakers.

# How to use

`go get -u github.com/wojnosystems/go-rate-limit`

```gopackage main

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
