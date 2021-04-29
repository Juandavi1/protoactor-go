package actor

import (
	"sync/atomic"
	"time"
)

type ShouldThrottle func() Valve

type Valve int32

const (
	Open Valve = iota
	Closing
	Closed
)

// NewThrottle
// This has no guarantees that the throttle opens exactly after the period, since it is reset asynchronously
// Throughput has been prioritized over exact re-opening
// throttledCallBack, This will be called with the number of events what was throttled after the period
func NewThrottle(maxEventsInPeriod int32, period time.Duration, throttledCallBack func(int32)) ShouldThrottle {

	var currentEvents = int32(0)

	startTimer := func(duration time.Duration, back func(int32)) {
		go func() {
			time.Sleep(duration)
			timesCalled := atomic.SwapInt32(&currentEvents, 0)
			if timesCalled > maxEventsInPeriod {
				throttledCallBack(timesCalled - maxEventsInPeriod)
			}
		}()
	}

	return func() Valve {

		tries := atomic.AddInt32(&currentEvents, 1)
		if tries == 1 {
			startTimer(period, throttledCallBack)
		}

		if tries == maxEventsInPeriod {
			return Closing
		} else if tries > maxEventsInPeriod {
			return Closed
		} else {
			return Open
		}
	}
}
