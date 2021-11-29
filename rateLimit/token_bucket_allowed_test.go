package rateLimit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("TokenBucket.Allowed", func() {
	var (
		bucket *TokenBucket
	)
	When("tokenCost is 1", func() {
		var unit uint64
		BeforeEach(func() {
			unit = 1
		})
		When("does not refill", func() {
			BeforeEach(func() {
				bucket = &TokenBucket{}
			})
			When("bucket has no tokens", func() {
				It("rejects", func() {
					Expect(bucket.Allowed(unit)).Should(BeFalse())
				})
				When("time passes", func() {
					BeforeEach(func() {
						bucket.nowFactory = startAtAndAddDurations(
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							10*time.Second,
							1*time.Millisecond,
							1*time.Millisecond,
							1*time.Millisecond,
						)
					})
					It("does not refill", func() {
						Expect(bucket.Allowed(unit)).Should(BeFalse())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
					})
				})
			})
			When("bucket has tokens", func() {
				BeforeEach(func() {
					bucket.tokens = 2
				})
				It("consumes them", func() {
					Expect(bucket.Allowed(unit)).Should(BeTrue())
					Expect(bucket.Allowed(unit)).Should(BeTrue())
					Expect(bucket.Allowed(unit)).Should(BeFalse())
				})
			})
		})
		When("does refill", func() {
			BeforeEach(func() {
				// used small numbers here to make examples easy to test
				bucket = NewTokenBucket(TokenBucketOpts{
					Capacity:             2,
					TokensAddedPerSecond: 2.0,
					InitialTokens:        1,
				})
			})
			When("bucket has no tokens", func() {
				BeforeEach(func() {
					bucket.tokens = 0
				})
				It("rejects", func() {
					Expect(bucket.Allowed(unit)).Should(BeFalse())
				})
				When("enough time has passed to fill", func() {
					BeforeEach(func() {
						bucket.nowFactory = startAtAndAddDurations(
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							500*time.Millisecond,
						)
					})
					It("allows", func() {
						Expect(bucket.Allowed(unit)).Should(BeTrue())
					})
				})
			})
			When("bucket has tokens", func() {
				When("all tokens are expended", func() {
					BeforeEach(func() {
						bucket.nowFactory = startAtAndAddDurations(
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							1*time.Millisecond,
							1*time.Millisecond,
						)
					})
					It("rejects", func() {
						Expect(bucket.Allowed(unit)).Should(BeTrue())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
					})
				})
				When("max tokens is reached", func() {
					BeforeEach(func() {
						bucket.nowFactory = startAtAndAddDurations(
							time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							10*time.Second,
							1*time.Millisecond,
							1*time.Millisecond,
						)
					})
					It("does not add more than Capacity", func() {
						Expect(bucket.Allowed(unit)).Should(BeTrue())
						Expect(bucket.Allowed(unit)).Should(BeTrue())
						Expect(bucket.Allowed(unit)).Should(BeFalse())
					})
				})
			})
			When("called with fractionally-generated tokens", func() {
				BeforeEach(func() {
					bucket.opts.Capacity = 5
					bucket.opts.TokensAddedPerSecond = 10.0
					bucket.opts.InitialTokens = 0
					bucket.tokens = 0
					bucket.nowFactory = startAtAndAddDurations(
						time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
						125*time.Millisecond,
						125*time.Millisecond,
						125*time.Millisecond,
						125*time.Millisecond,
						1*time.Millisecond,
					)
				})
				It("does not lose fractions of tokens", func() {
					_ = bucket.Allowed(unit)
					_ = bucket.Allowed(unit)
					_ = bucket.Allowed(unit)
					_ = bucket.Allowed(unit)
					Expect(bucket.Allowed(unit)).Should(BeTrue())
				})
			})
		})
	})
	When("tokenCost is 10", func() {
		var unit uint64
		BeforeEach(func() {
			unit = 10
		})

		When("insufficient tokens", func() {
			BeforeEach(func() {
				bucket = NewTokenBucket(TokenBucketOpts{InitialTokens: unit / 2})
			})
			It("rejects", func() {
				Expect(bucket.Allowed(unit)).Should(BeFalse())
			})
		})
		When("sufficient tokens", func() {
			BeforeEach(func() {
				bucket = NewTokenBucket(TokenBucketOpts{InitialTokens: unit})
			})
			It("rejects", func() {
				Expect(bucket.Allowed(unit)).Should(BeTrue())
			})
		})
	})
})
