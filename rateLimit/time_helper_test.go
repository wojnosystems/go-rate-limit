package rateLimit

import (
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

// simulateNows will return each time provided, in order, whenever the returned function is called.
// Should the returned function be called more than the number of times provided, returns time.Zero
func simulateNows(times ...time.Time) func() time.Time {
	i := 0
	return func() time.Time {
		if i >= len(times) {
			return zeroTime
		}
		now := times[i]
		i++
		return now
	}
}

func Test_simulateNows(t *testing.T) {
	cases := map[string]struct {
		input []time.Time
	}{
		"no times return zero-time": {
			input: []time.Time{},
		},
		"some times returned": {
			input: []time.Time{
				time.Date(2021, 11, 28, 21, 35, 15, 0, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 15, 1, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 15, 2, time.UTC),
			},
		},
		"another call to prove i is not shared among calls": {
			input: []time.Time{
				time.Date(2021, 11, 28, 21, 35, 15, 0, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 15, 1, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 15, 2, time.UTC),
			},
		},
	}

	for caseName, td := range cases {
		t.Run(caseName, func(t *testing.T) {
			g := NewWithT(t)
			function := simulateNows(td.input...)
			actual := make([]time.Time, 0, 0)
			for {
				item := function()
				if item == zeroTime {
					break
				}
				actual = append(actual, item)
			}
			g.Expect(actual).Should(ConsistOf(td.input))
		})
	}
}

// startAtAndAddDurations simulates times like simulateNows, but allows you to start with a time,
// then add durations to each subsequent returned value, pushing the times further into the future.
func startAtAndAddDurations(startTime time.Time, addToEach ...time.Duration) func() time.Time {
	times := make([]time.Time, len(addToEach)+1)
	times[0] = startTime
	for i, each := range addToEach {
		times[i+1] = times[i].Add(each)
	}
	return simulateNows(times...)
}

func Test_startAtAndAddDurations(t *testing.T) {
	cases := map[string]struct {
		start    time.Time
		adds     []time.Duration
		expected []time.Time
	}{
		"no times return zero-time": {
			start: time.Date(2021, 11, 28, 21, 35, 15, 0, time.UTC),
			adds: []time.Duration{
				1 * time.Second,
				4 * time.Second,
				1 * time.Second,
			},
			expected: []time.Time{
				time.Date(2021, 11, 28, 21, 35, 15, 0, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 16, 0, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 20, 0, time.UTC),
				time.Date(2021, 11, 28, 21, 35, 21, 0, time.UTC),
			},
		},
	}

	for caseName, td := range cases {
		t.Run(caseName, func(t *testing.T) {
			g := NewWithT(t)
			function := startAtAndAddDurations(td.start, td.adds...)
			actual := make([]time.Time, 0, 0)
			for {
				item := function()
				if item == zeroTime {
					break
				}
				actual = append(actual, item)
			}
			g.Expect(actual).Should(ConsistOf(td.expected))
		})
	}
}
