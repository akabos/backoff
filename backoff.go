package backoff

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	DefaultInitial = time.Second
	DefaultFinal = time.Second * 10
	DefaultScale = 1.5
)

type Backoff struct {
	Initial time.Duration
	Final   time.Duration
	Scale   float64
	Jitter  float64
	current time.Duration
	mux     sync.Mutex
	once    sync.Once
}

func (b *Backoff) init() {
	b.once.Do(func() {
		if b.Initial == 0 {
			b.Initial = DefaultInitial
		}
		if b.Final == 0 {
			b.Final = DefaultFinal
		}
		if b.Scale == 0 {
			b.Scale = DefaultScale
		}
	})
}

func (b *Backoff) Next() time.Duration {
	b.init()
	b.mux.Lock()
	if b.current == 0 {
		b.current = b.Initial
	} else {
		b.current = time.Duration(math.Floor(float64(b.current) * b.Scale))
	}
	if b.current > b.Final {
		b.current = b.Final
	}
	var diff time.Duration
	if b.Jitter > 0 {
		diff = time.Duration(float64(b.current) * (rand.Float64() - 0.5) * b.Jitter)
	}
	b.mux.Unlock()
	return b.current + diff
}

func (b *Backoff) Reset() {
	b.mux.Lock()
	b.current = 0
	b.mux.Unlock()
}
