package utils

import (
	"errors"
	"sync"
	"time"
)

type Promise struct {
	ch        chan interface{}
	cancel    chan error
	timer     *time.Timer
	hasResult bool
	mu        sync.Mutex
}

func NewPromise(timer *time.Timer) *Promise {
	return &Promise{
		ch:        make(chan interface{}, 1),
		cancel:    make(chan error, 1),
		timer:     timer,
		hasResult: false,
	}
}
func (p *Promise) Cancel() error {
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.hasResult {
		return errors.New("HasResult")
	}
	p.hasResult = true
	p.cancel <- errors.New("PromiseCanceled")
	return nil
}
func (p *Promise) SetSuccess(e interface{}) error {
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.hasResult {
		return errors.New("HasResult")
	}
	p.hasResult = true
	p.ch <- e
	return nil
}
func (p *Promise) Get() (interface{}, error) {
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.hasResult {
		return nil, errors.New("HasResult")
	}
	if p.timer == nil {
		select {
		case g := <-p.ch:
			p.hasResult = true
			return g, nil
		case err := <-p.cancel:
			p.hasResult = true
			return err, err
		}
	} else {
		select {
		case g := <-p.ch:
			p.hasResult = true
			return g, nil
		case err := <-p.cancel:
			p.hasResult = true
			return err, nil
		case <-p.timer.C:
			return nil, errors.New("TimeOut")
		}
	}
}
