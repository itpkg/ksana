package job

import (
	"sort"
	"time"
)

type Configuration struct {
	Workers  int            `toml:"workers"`
	TimeoutI int            `toml:"timeout"`
	QueuesM  map[string]int `toml:"queues"`
}

func (p *Configuration) Timeout() time.Duration {
	return time.Duration(p.TimeoutI) * time.Second
}

func (p *Configuration) Queues() []string {
	n := map[int][]string{}
	var a []int
	for k, v := range p.QueuesM {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	names := make([]string, 0)
	for _, k := range a {
		for _, s := range n[k] {
			names = append(names, s)
		}
	}

	return names
}
