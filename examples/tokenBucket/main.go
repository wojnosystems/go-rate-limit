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
