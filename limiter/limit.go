package limiter

import "sync"

type limiter struct {
	Current int
	Max     int
	mu      *sync.Mutex
}

func NewLimit(limit int) limiter {
	mu := new(sync.Mutex)
	return limiter{
		Current: 0,
		Max:     limit,
		mu:      mu,
	}
}

func (l *limiter) Inc() {
	l.mu.Lock()
	l.Current++
	l.mu.Unlock()
}

func (l *limiter) Dec() {
	l.mu.Lock()
	l.Current--
	l.mu.Unlock()
}
