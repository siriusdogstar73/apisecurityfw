package socket

import (
	"context"

	"github.com/jjeffcaii/reactor-go"
	"github.com/rsocket/rsocket-go/core"
	"github.com/rsocket/rsocket-go/core/framing"
	"github.com/rsocket/rsocket-go/internal/fragmentation"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"go.uber.org/atomic"
)

type requestChannelSubscriber struct {
	sid       uint32
	n         uint32
	dc        *DuplexConnection
	requested atomic.Bool
	rcv       flux.Processor
}

func (r *requestChannelSubscriber) OnNext(item payload.Payload) {
	if !r.requested.CAS(false, true) {
		r.dc.sendPayload(r.sid, item, core.FlagNext)
		return
	}
	d := item.Data()
	m, _ := item.Metadata()
	size := framing.CalcPayloadFrameSize(d, m) + 4
	if !r.dc.shouldSplit(size) {
		metadata, _ := item.Metadata()
		r.dc.sendFrame(framing.NewWriteableRequestChannelFrame(r.sid, r.n, item.Data(), metadata, core.FlagNext))
		return
	}
	r.dc.doSplitSkip(4, d, m, func(index int, result fragmentation.SplitResult) {
		var f core.WriteableFrame
		if index == 0 {
			f = framing.NewWriteableRequestChannelFrame(r.sid, r.n, result.Data, result.Metadata, result.Flag|core.FlagNext)
		} else {
			f = framing.NewWriteablePayloadFrame(r.sid, result.Data, result.Metadata, result.Flag|core.FlagNext)
		}
		r.dc.sendFrame(f)
	})
}

func (r *requestChannelSubscriber) OnError(err error) {
	r.dc.writeError(r.sid, err)
}

func (r *requestChannelSubscriber) OnComplete() {
	complete := framing.NewWriteablePayloadFrame(r.sid, nil, nil, core.FlagComplete)
	done := make(chan struct{})
	complete.HandleDone(func() {
		close(done)
	})
	if r.dc.sendFrame(complete) {
		<-done
	}
}

func (r *requestChannelSubscriber) OnSubscribe(ctx context.Context, s rx.Subscription) {
	select {
	case <-ctx.Done():
		r.OnError(reactor.ErrSubscribeCancelled)
	default:
		cb := requestChannelCallback{
			rcv: r.rcv,
			snd: s,
		}
		r.dc.register(r.sid, cb)
		s.Request(1)
	}
}

type respondChannelSubscriber struct {
	sid        uint32
	n          uint32
	dc         *DuplexConnection
	rcv        flux.Processor
	subscribed chan<- struct{}
	calls      *atomic.Int32
}

func (r *respondChannelSubscriber) OnNext(next payload.Payload) {
	r.dc.sendPayload(r.sid, next, core.FlagNext)
}

func (r *respondChannelSubscriber) OnError(err error) {
	if r.calls.Inc() == 2 {
		r.dc.unregister(r.sid)
	}
	r.dc.writeError(r.sid, err)
}

func (r *respondChannelSubscriber) OnComplete() {
	if r.calls.Inc() == 2 {
		r.dc.unregister(r.sid)
	}
	complete := framing.NewWriteablePayloadFrame(r.sid, nil, nil, core.FlagComplete)
	done := make(chan struct{})
	complete.HandleDone(func() {
		close(done)
	})
	if r.dc.sendFrame(complete) {
		<-done
	}
}

func (r *respondChannelSubscriber) OnSubscribe(ctx context.Context, s rx.Subscription) {
	select {
	case <-ctx.Done():
		r.OnError(reactor.ErrSubscribeCancelled)
	default:
		cb := respondChannelCallback{
			rcv: r.rcv,
			snd: s,
		}
		r.dc.register(r.sid, cb)
		close(r.subscribed)
		s.Request(ToIntRequestN(r.n))
	}
}
