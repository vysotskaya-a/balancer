package proxy

import (
	"net/url"
	"sync/atomic"
)

type Balancer struct {
	targets []*url.URL
	counter uint64
}

func NewBalancer(targets []string) *Balancer {
	urls := make([]*url.URL, len(targets))
	for i, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			panic("Невалидный URL в конфиге: " + t)
		}
		urls[i] = u
	}

	return &Balancer{
		targets: urls,
	}
}

func (b *Balancer) NextTarget() *url.URL {
	idx := atomic.AddUint64(&b.counter, 1)
	return b.targets[int(idx)%len(b.targets)]
}
