package pocket

import (
	"testing"
	"time"
)

var (
	duration = time.Millisecond
	now      = time.Now()
	data     = []string{
		"hello world",
		"hello africa",
		"hello Tanzania",
	}
)

func TestScheduler(t *testing.T) {
	sh := NewScheduler()
	for k, v := range data {
		shd := schedule{
			start:    now.Add(duration * time.Duration(k)),
			duration: duration,
			event:    v,
		}
		sh.Add(shd)
	}
	sh.Add(schedule{
		start:    now.Add(-time.Second),
		duration: duration,
		event:    "yoyo",
	})
	//w := sync.WaitGroup{}
	//w.Add(1)
	//go func(wg sync.WaitGroup, s *Scheduler) {
	//	ch := time.Tick(duration)
	//END:
	//	for {
	//		select {
	//		case <-ch:
	//			log.Println(s.OnAir())
	//		case <-time.After(4 * time.Millisecond):
	//			break END
	//		}
	//	}
	//	wg.Done()
	//}(w, sh)
	//w.Wait()
	sh.OnAir()
}
