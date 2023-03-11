package dtime

import "time"

/*
Equals - compares 't1' and 't2' UTCs. If 'round' > 0, 't1' and 't2' get rounded before comparison. Half values get
rounded down.
*/
func Equals(t1, t2 time.Time, round time.Duration) bool {
	var (
		t1Ns = t1.UTC().UnixNano()
		t2Ns = t2.UTC().UnixNano()
	)

	if round < 1 {
		return t1Ns == t2Ns
	}

	// time.Time.Round() may round up half values.
	return t1Ns/int64(round) == t2Ns/int64(round)
}

// EqualsMs - is the same as Equals(), but 'round' is time.Millisecond.
func EqualsMs(t1, t2 time.Time) bool {
	return Equals(t1, t2, time.Millisecond)
}
