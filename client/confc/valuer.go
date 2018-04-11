package confc

import (
	"sync"
)

// Valuer provides value updates from a provider.
type Valuer interface {
	ChanForValue(name string) <-chan string
}

// defaultValuer is a default Valuer point for the created values.
var defaultValuer Valuer = newProxyValuer()

// SetValuer sets the default valuer.
// All previously made ChanForValue()-made chans are chained and proxied
// via the passed Valuer.
// This function should only be called once.
func SetValuer(v Valuer) {
	defaultValuer.(*proxyValuer).Proxy(v)
}

type proxyValuer struct {
	mu    sync.RWMutex
	realV Valuer

	rcMu     sync.RWMutex
	regChans map[string][]chan string
}

func newProxyValuer() *proxyValuer {
	return &proxyValuer{
		regChans: map[string][]chan string{},
	}
}

func (p *proxyValuer) ChanForValue(name string) <-chan string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.realV != nil {
		return p.realV.ChanForValue(name)
	}
	p.rcMu.Lock()
	defer p.rcMu.Unlock()

	ret := make(chan string, 4)
	p.regChans[name] = append(p.regChans[name], ret)

	return ret
}

func (p *proxyValuer) Proxy(v Valuer) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.rcMu.Lock()
	defer p.rcMu.Unlock()

	if p.regChans == nil || p.realV != nil {
		panic("SetValuer called twice!")
	}

	p.realV = v

	// proxy values from real chan to old fakes
	for key := range p.regChans {
		cc := p.regChans[key]

		go func(key string) {
			uCh := v.ChanForValue(key)
			for v := range uCh {
				for _, c := range cc {
					c <- v
				}
			}
		}(key)
	}
	// make sure SetValuer() isn't called twice
	p.regChans = nil
}
