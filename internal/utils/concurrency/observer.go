package concurrency

type Observer[M any] struct {
	out chan M
}

func (o *Observer[M]) Observe(subject chan M) {
	go channelRelay(subject, o.out)
}

func NewObserver[M any](out chan M) (f *Observer[M]) {
	return &Observer[M]{
		out: out,
	}
}
