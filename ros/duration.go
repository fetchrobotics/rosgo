package ros

import (
	"time"
)

// Duration is ros duration primitive that defines a period of time.
type Duration struct {
	temporal
}

// NewDuration creates a new instance of ros duration from seconds and nanoseconds.
func NewDuration(sec uint32, nsec uint32) Duration {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Duration{temporal{sec, nsec}}
}

// Add adds and returns duration (d+other).
func (d *Duration) Add(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)+int64(other.Sec),
		int64(d.NSec)+int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

// Sub subtracts and returns duration (d-other).
func (d *Duration) Sub(other Duration) Duration {
	sec, nsec := normalizeTemporal(int64(d.Sec)-int64(other.Sec),
		int64(d.NSec)-int64(other.NSec))
	return Duration{temporal{sec, nsec}}
}

// Cmp compares duration d with duration other.
// Return integer representing the comparison.
//  1 - d > other
//  0 - d == other
// -1 - d < other
func (d *Duration) Cmp(other Duration) int {
	return cmpUint64(d.ToNSec(), other.ToNSec())
}

// Sleep sleeps for period of time specified in d.
func (d *Duration) Sleep() error {
	if !d.IsZero() {
		time.Sleep(time.Duration(d.ToNSec()) * time.Nanosecond)
	}
	return nil
}
