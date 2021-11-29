package rateLimit

type Limiter interface {
	Allowed(tokenCost uint64) bool
}

type TokenLimiter interface {
	Limiter
	Tokens() uint64
}

type BurstingTokenBucketOpts struct {
	Bucket TokenLimiter
	Burst  TokenLimiter
}

// BurstingTokenBucket is just like TokenBucket, but will Allow requests to temporarily exceed the token cost within a
// refillable limit. Use NewBurstingBucket. The default will panic. Is _not_ thread-safe.
type BurstingTokenBucket struct {
	bucket TokenLimiter
	burst  TokenLimiter
}

func NewBurstingTokenBucket(opts BurstingTokenBucketOpts) *BurstingTokenBucket {
	return &BurstingTokenBucket{
		bucket: opts.Bucket,
		burst:  opts.Burst,
	}
}

// Allowed returns true only if tokenCost tokensAvailable are available in both the regular bucket and bursting bucket.
// If the tokenCost is not available, does not deduct the tokensAvailable and returns false
func (b *BurstingTokenBucket) Allowed(tokenCost uint64) bool {
	if b.bucket.Allowed(tokenCost) {
		return true
	}
	remainingBucketTokens := b.bucket.Tokens()
	if !b.burst.Allowed(tokenCost - remainingBucketTokens) {
		return false
	}
	_ = b.bucket.Allowed(remainingBucketTokens)
	return true
}
