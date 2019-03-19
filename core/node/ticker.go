package node

import (
	"sync"
	"time"
)

type Ticker struct {
	notifys []func()
	done    chan chan struct{}
	mu      sync.Mutex
}

func newTicker() *Ticker {
	ticker := new(Ticker)
	ticker.notifys = make([]func(), 0, 1)
	ticker.done = make(chan chan struct{}, 1)
	go ticker.startup()
	return ticker
}

func (t *Ticker) startup() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			for _, fn := range t.notifys {
				go fn()
			}
		case exit := <-t.done:
			ticker.Stop()
			exit <- struct{}{}
			return
		}
	}
}

func (t *Ticker) exit() {
	exit := make(chan struct{}, 1)
	t.done <- exit
	<-exit
	return
}

func (t *Ticker) Bind(fn func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.notifys = append(t.notifys, fn)
	return
}

func (t *Ticker) Close() error {
	t.exit()
	return nil
}
