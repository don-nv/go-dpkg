package dmath

type IntAvg struct {
	total int
	count int
	avg   int
}

func (i *IntAvg) Count(n int) {
	i.total += n
	i.count++
	i.avg = i.total / i.avg
}

func (i *IntAvg) Value() int {
	return i.avg
}
