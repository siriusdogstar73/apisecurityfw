package flux

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/jjeffcaii/reactor-go"
	"github.com/jjeffcaii/reactor-go/hooks"
	"github.com/jjeffcaii/reactor-go/scheduler"
)

type fluxDelayElement struct {
	source reactor.RawPublisher
	delay  time.Duration
	sc     scheduler.Scheduler
}

type delayElementSubscriber struct {
	delay  time.Duration
	sc     scheduler.Scheduler
	s      reactor.Subscription
	actual reactor.Subscriber
	stat   int32
}

func newFluxDelayElement(source reactor.RawPublisher, delay time.Duration, sc scheduler.Scheduler) *fluxDelayElement {
	return &fluxDelayElement{
		source: source,
		delay:  delay,
		sc:     sc,
	}
}

func (d *delayElementSubscriber) Request(n int) {
	d.s.Request(n)
}

func (d *delayElementSubscriber) Cancel() {
	// TODO: support cancel
	d.s.Cancel()
}

func (d *delayElementSubscriber) OnComplete() {
	if atomic.CompareAndSwapInt32(&d.stat, 0, statComplete) {
		d.actual.OnComplete()
	}
}

func (d *delayElementSubscriber) OnError(err error) {
	if atomic.CompareAndSwapInt32(&d.stat, 0, statError) {
		d.actual.OnError(err)
		return
	}
	hooks.Global().OnErrorDrop(err)
}

func (d *delayElementSubscriber) OnNext(v Any) {
	if atomic.LoadInt32(&d.stat) == statComplete {
		hooks.Global().OnNextDrop(v)
		return
	}
	<-time.After(d.delay)
	d.actual.OnNext(v)
}

func (d *delayElementSubscriber) OnSubscribe(ctx context.Context, s reactor.Subscription) {
	d.s = s
	d.actual.OnSubscribe(ctx, d)
	s.Request(reactor.RequestInfinite)
}

func (f *fluxDelayElement) SubscribeWith(ctx context.Context, actual reactor.Subscriber) {
	f.source.SubscribeWith(ctx, &delayElementSubscriber{
		delay:  f.delay,
		actual: actual,
		sc:     f.sc,
	})
}
