package ros

// Rate defines a struct that can be used to
// run loops at a desired frequency.
// Based on: http://docs.ros.org/diamondback/api/rostime/html/classros_1_1Rate.html
type Rate struct {
	actualCycleTime   Duration
	expectedCycleTime Duration
	start             Time
}

// NewRate creates and returns rate based on frequency specified in hertz.
func NewRate(frequency float64) Rate {
	var actualCycleTime, expectedCycleTime Duration
	expectedCycleTime.FromSec(1.0 / frequency)
	start := Now()
	return Rate{actualCycleTime, expectedCycleTime, start}
}

// CycleTime creates and returns rate based on cycle time specified in ros::duration.
func CycleTime(d Duration) Rate {
	var actualCycleTime Duration
	start := Now()
	return Rate{actualCycleTime, d, start}
}

// CycleTime returns the actual runtime of a cycle from start to sleep.
func (r *Rate) CycleTime() Duration {
	return r.actualCycleTime
}

// ExpectedCycleTime returns the expected cycle timw which is one over the frequency passed
// while create rate.
func (r *Rate) ExpectedCycleTime() Duration {
	return r.expectedCycleTime
}

// Reset resets the start time of rate to now.
func (r *Rate) Reset() {
	r.actualCycleTime = NewDuration(0, 0)
	r.start = Now()
}

// Sleep sleeps for any leftover time in a cycle.
// Calculated from the last time sleep, reset, or the constructor was called.
func (r *Rate) Sleep() error {
	end := Now()
	diff := end.Diff(r.start)
	var remaining Duration
	if r.expectedCycleTime.Cmp(diff) >= 0 {
		remaining = r.expectedCycleTime.Sub(diff)
	}
	remaining.Sleep()
	now := Now()
	r.actualCycleTime = now.Diff(r.start)
	r.start = r.start.Add(r.expectedCycleTime)
	return nil
}
