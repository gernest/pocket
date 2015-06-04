package pocket

import (
	"container/heap"
	"errors"
	"sync"
	"time"

	"github.com/kr/pretty"
)

type schedule struct {
	start    time.Time
	duration time.Duration
	index    int
	event    interface{}
}

type Scheduler struct {
	table schedules
	mux   sync.RWMutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}
func (s *Scheduler) Add(sh schedule) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	heap.Push(&s.table, sh)
}
func (s *Scheduler) OnAir() (*schedule, error) {
	heap.Init(&s.table)
	pretty.Println(s.table)
	return nil, nil
}

type schedules []schedule

func (s schedules) Len() int { return len(s) }
func (s schedules) Less(i, j int) bool {
	return s[i].start.Add(s[i].duration).Before(s[j].start)
}
func (s schedules) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
	a := &s[i]
	b := &s[i]
	a.index = j
	b.index = i
}

func (s *schedules) Push(x interface{}) {
	*s = append(*s, x.(schedule))
}

func (s *schedules) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func (s schedules) Now() (*schedule, error) {
	start := time.Now()
	return findSchedule(start, s)
}

func findSchedule(starts time.Time, h schedules, index ...int) (*schedule, error) {
	last := h[0]
	first := h[len(h)-1]
	if len(index) > 0 {
		last = h[index[0]]
		first = h[index[1]]
	}
	d := direction(&last, &first, starts)
	if d != nil {
		next := nextLap(d, starts)
		if next > len(h) {
			return nil, errors.New("pocket: index out of range")
		}
		if inRange(d.index, next) {
			var rst *schedule
			for k, _ := range make([]struct{}, 3) {
				x := d.index + k
				if x < len(h) {
					v := h[x]
					if v.start.Equal(starts) {
						rst = &v
					}
				}
			}
			if rst != nil {
				return rst, nil
			}
			return nil, errors.New("pocket: no schedule found")

		}
		return findSchedule(starts, h, d.index, next)
	}
	return nil, errors.New("pocket: nothing found")
}

func direction(a, b *schedule, starts time.Time) *schedule {
	var d1, d2 time.Duration
	t1 := a.start.Add(a.duration)
	t2 := b.start.Add(b.duration)
	switch {
	case t1.After(starts):
		d1 = t1.Sub(starts)
	case t1.Before(starts):
		d1 = starts.Sub(t1)
	case t2.After(starts):
		d2 = t2.Sub(starts)
	case t2.Before(starts):
		d2 = starts.Sub(t2)
	}
	switch {
	case d1 > d2:
		return b
	case d1 < d2:
		return a
	case a.start.Equal(starts):
		return a
	case b.start.Equal(starts):
		return b
	}
	return nil
}

func nextLap(s *schedule, starts time.Time) int {
	if s.start.Equal(starts) {
		return s.index
	}
	xd := s.start.Add(s.duration)
	if xd.Before(starts) {
		delta := starts.Sub(xd)
		n := delta / s.duration
		return int(n) + s.index
	}
	delta := s.start.Sub(starts)
	n := delta / s.duration
	return s.index - int(n)
}

func inRange(i, j int) bool {
	if i > j {
		return i-j < 2
	}
	return j-i < 2
}
