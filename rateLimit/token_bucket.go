package rateLimit

import (
	"math"
	"time"
)

var (
	zeroTime = time.Time{}
)

func defaultNow() time.Time {
	return time.Now()
}

type TokenBucketOpts struct {
	// Capacity the maximum number of tokensAvailable the bucket may contain
	Capacity uint64
	// TokensAddedPerSecond how many tokensAvailable are replenished every second, tokensAvailable are added discretely but this may
	// also be larger than 1 to add tokensAvailable more quickly than each second
	TokensAddedPerSecond float64
	// InitialTokens in the bucket. Prime the bucket
	InitialTokens uint64
}

// TokenBucket is a rate-limiter using a token bucket scheme to approximate rates
// Call NewTokenBucket to create a new instance with custom initialization.
// Empty TokenBuckets contain no capacity and no tokensAvailable as well as no refill.
// Instance is _not_ thread-safe.
type TokenBucket struct {
	opts TokenBucketOpts

	tokensAvailable uint64
	// remainder is the left-over tokens generated thus far which would have been discarded because we round-down
	// This ensures we don't lose fractionally generated tokens.
	remainder   float64
	lastUpdated time.Time

	// nowFactory allows us to simulate time
	nowFactory    nowFactory
	isInitialized bool
}

func NewTokenBucket(opts TokenBucketOpts) *TokenBucket {
	return &TokenBucket{
		opts:            opts,
		tokensAvailable: opts.InitialTokens,
	}
}

// Allowed returns true only if tokenCost tokensAvailable are available. If the tokenCost is not available,
// does not deduct the tokensAvailable and returns false
func (b *TokenBucket) Allowed(tokenCost uint64) bool {
	return b.allowed(tokenCost)
}

func (b *TokenBucket) allowed(cost uint64) bool {
	b.initializeIfNeeded()
	b.tokensAvailable, b.lastUpdated, b.remainder = replenishTokens(
		b.tokensAvailable,
		b.remainder,
		b.lastUpdated,
		b.nowFactory(),
		b.opts.TokensAddedPerSecond,
		b.opts.Capacity)

	if b.tokensAvailable < cost {
		return false
	}
	b.tokensAvailable -= cost
	return true
}

// initializeIfNeeded will prepare the TokenBucket for use if NewTokenBucket was not called
func (b *TokenBucket) initializeIfNeeded() {
	if b.isInitialized {
		return
	}
	b.isInitialized = true

	if b.nowFactory == nil {
		b.nowFactory = defaultNow
	}

	b.lastUpdated = b.nowFactory()
}

// replenishTokens adds tokensAvailable if any are available and there's capacity.
// Intended to be called each time allowed is called.
// returns the new number of tokensAvailable and the new lastUpdated time
func replenishTokens(tokens uint64,
	remainder float64,
	lastUpdated time.Time,
	now time.Time,
	tokensAddedPerSecond float64,
	maxTokens uint64) (
	updatedTokens uint64,
	updatedLastUpdated time.Time,
	updatedRemainder float64) {

	amountToAdd := tokensAddedPerSecond*now.Sub(lastUpdated).Seconds() + remainder
	tokensToAdd := uint64(math.Floor(amountToAdd))
	if tokensToAdd < 1 {
		// not enough time has passed to add a single token, return without updating the last updated time
		return tokens, lastUpdated, remainder
	}

	updatedRemainder = amountToAdd - float64(tokensToAdd)

	remainingCapacity := maxTokens - tokens
	if tokensToAdd > remainingCapacity {
		tokensToAdd = remainingCapacity
	}
	tokens += tokensToAdd
	return tokens, now, updatedRemainder
}

func (b *TokenBucket) Tokens() uint64 {
	return b.tokensAvailable
}
