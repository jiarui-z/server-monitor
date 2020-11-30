package terminal

import (
	"container/ring"
)

type StatRing struct {
	r *ring.Ring
}

func NewStatRing(n int) *StatRing {
	return &StatRing{r: ring.New(n)}
}

func (p *StatRing) Push(n uint64) {
	p.r.Value = n
	p.r = p.r.Next()
}

// Data, convert underline data to float64
func (p *StatRing) ToFloat() []float64 {
	var l []float64
	p.r.Do(func(x interface{}) {
		if v, ok := x.(uint64); ok {
			l = append(l, float64(v))
		} else {
			l = append(l, 0.0)
		}
	})
	return l
}

// NormalizedData return normalized data between [0,1]
func (p *StatRing) NormalizedData() []float64 {
	max := p.max()
	if max == 0 {
		return make([]float64, p.r.Len(), p.r.Len())
	}

	var l []float64
	p.r.Do(func(x interface{}) {
		var pct float64
		if v, ok := x.(uint64); ok {
			pct = float64(v) / float64(max)
		}
		l = append(l, pct)
	})
	return l
}

func (p *StatRing) max() uint64 {
	var max uint64
	// find max
	p.r.Do(func(x interface{}) {
		if v, ok := x.(uint64); ok && v > max {
			max = v
		}
	})
	return max
}
