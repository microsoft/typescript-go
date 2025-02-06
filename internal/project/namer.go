package project

import "sync"

type namer struct {
	mu       sync.Mutex
	counters map[string]int
}

func (n *namer) next(name string) string {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.counters == nil {
		n.counters = make(map[string]int)
	}
	n.counters[name]++
	return name + string(n.counters[name]) + "*"
}
