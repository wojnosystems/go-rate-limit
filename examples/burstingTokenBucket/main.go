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
			Burst: rateLimit.NewTokenBucket(rateLimit.TokenBucketOpts{
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
