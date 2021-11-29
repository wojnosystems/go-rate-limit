package rateLimit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstingTokenBucket.Allowed", func() {
	var (
		bursting *BurstingTokenBucket
	)

	When("sufficient capacity in bucket", func() {
		BeforeEach(func() {
			bursting = NewBurstingTokenBucket(BurstingTokenBucketOpts{
				Bucket: NewTokenBucket(TokenBucketOpts{InitialTokens: 5}),
				Burst:  NewTokenBucket(TokenBucketOpts{InitialTokens: 1}),
			})
		})
		It("uses the burst last", func() {
			Expect(bursting.Allowed(5)).Should(BeTrue())
			Expect(bursting.Allowed(1)).Should(BeTrue())
			Expect(bursting.Allowed(1)).Should(BeFalse())
		})
	})
	When("insufficient capacity in bucket", func() {
		BeforeEach(func() {
			bursting = NewBurstingTokenBucket(BurstingTokenBucketOpts{
				Bucket: NewTokenBucket(TokenBucketOpts{InitialTokens: 5}),
				Burst:  NewTokenBucket(TokenBucketOpts{InitialTokens: 1}),
			})
		})
		It("can use up the burst", func() {
			Expect(bursting.Allowed(6)).Should(BeTrue())
			Expect(bursting.Allowed(1)).Should(BeFalse())
		})
		It("rejects if over both", func() {
			Expect(bursting.Allowed(7)).Should(BeFalse())
		})
		It("does not consume the bucket if rejected", func() {
			bursting.Allowed(7)
			Expect(bursting.Allowed(6)).Should(BeTrue())
		})
	})

})
