package ros

import (
	gotime "time"
)

// Time defines a struct that represents ros time.
type Time struct {
	temporal
}

// NewTime creates and returns a new instance of ros time.
func NewTime(sec uint32, nsec uint32) Time {
	sec, nsec = normalizeTemporal(int64(sec), int64(nsec))
	return Time{temporal{sec, nsec}}
}

// Now returns the current ros time.
func Now() Time {
	var t Time
	t.FromNSec(uint64(gotime.Now().UnixNano()))
	return t
}

// Diff returns the duration difference between
// time t and time from.
func (t *Time) Diff(from Time) Duration {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(from.Sec),
		int64(t.NSec)-int64(from.NSec))
	return Duration{temporal{sec, nsec}}
}

// Add adds d duration to time t.
func (t *Time) Add(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)+int64(d.Sec),
		int64(t.NSec)+int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

// Sub subtracts d duration from time t.
func (t *Time) Sub(d Duration) Time {
	sec, nsec := normalizeTemporal(int64(t.Sec)-int64(d.Sec),
		int64(t.NSec)-int64(d.NSec))
	return Time{temporal{sec, nsec}}
}

// Cmp compares time t and time other.
// Return integer representing the comparison.
//  1 : t > other
//  0 : t == other
// -1 : t < other
func (t *Time) Cmp(other Time) int {
	return cmpUint64(t.ToNSec(), other.ToNSec())
}
