package dmath

import "time"

func Is50x50() bool {
	return time.Now().UnixNano()%2 == 0
}
