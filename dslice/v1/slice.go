package dslice

func In[T comparable](t T, ts ...T) bool {
	for _, oneOfT := range ts {
		if t == oneOfT {
			return true
		}
	}

	return false
}

type Slice[T any] struct {
	ts []T
}

func NewSlice[T any](ts ...T) Slice[T] {
	return Slice[T]{
		ts: append(
			make([]T, 0, cap(ts)),
			ts...,
		),
	}
}

func WrapSlice[T any](ts []T) Slice[T] {
	return Slice[T]{
		ts: ts,
	}
}

func (s *Slice[T]) Append(t T) {
	s.ts = append(s.ts, t)
}

func (s *Slice[T]) Ts() []T {
	return s.ts
}

func (s *Slice[T]) TsCopy() []T {
	return append(
		make([]T, 0, len(s.ts)),
		s.ts...,
	)
}

/*
TsF - iterates over all Ts and returns []T filtered using 'f'. 'ok' in 'f' is used to indicate if 't' should be included
in the resulting []T.
*/
func (s *Slice[T]) TsF(f func(i int, t T) (ok bool)) []T {
	var ts = make([]T, 0)
	for i, t := range s.ts {
		if f(i, t) {
			ts = append(ts, t)
		}
	}

	return ts
}

/*
T - iterates over all Ts and returns first T filtered using 'f'. 'ok' in 'f' is used to indicate if 't' is wanted and
should be returned. Resulting bool indicates if wanted T was found or not.
*/
func (s *Slice[T]) T(f func(i int, t T) (ok bool)) (T, bool) {
	for i, t := range s.ts {
		if f(i, t) {
			return t, true
		}
	}

	var t T
	return t, false
}

/*
DeleteT - iterates over all Ts and deletes first T filtered using 'f'.'ok' in'f' is used to indicate if't' is wanted
and should be deleted. Resulting bool indicates if wanted T was deleted or not.
*/
func (s *Slice[T]) DeleteT(f func(i int, t T) (ok bool)) bool {
	for i, t := range s.ts {
		if f(i, t) {
			return s.DeleteI(i)
		}
	}

	return false
}

func (s *Slice[T]) DeleteI(i int) bool {
	if i < 0 {
		return false
	}

	maxI := len(s.ts) - 1
	if maxI < 0 {
		return false
	}
	if i > maxI {
		return false
	}

	s.ts = append(s.ts[:i], s.ts[i+1:]...)

	return true
}

func (s *Slice[T]) Len() int {
	return len(s.ts)
}

func (s *Slice[T]) Range(f func(i int, t T)) {
	s.RangeNext(func(i int, t T) (next bool) {
		f(i, t)

		return true
	})
}

// RangeNext - ranges through all Ts using'f' until'next' in'f' is true or the end of Ts is reached.
func (s *Slice[T]) RangeNext(f func(i int, t T) (next bool)) {
	for i, t := range s.ts {
		if f(i, t) {
			continue
		}

		return
	}
}
