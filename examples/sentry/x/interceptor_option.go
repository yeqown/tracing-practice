package x

import "time"

type Option interface {
	Apply(*options)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option) *options {
	optsCopy := *defaultOptions
	c := &optsCopy
	for _, o := range opts {
		o.Apply(c)
	}
	return c
}

type repanicOption struct {
	Repanic bool
}

func (r *repanicOption) Apply(o *options) {
	o.Repanic = r.Repanic
}

func WithRepanicOption(b bool) Option {
	return &repanicOption{Repanic: b}
}

type waitForDeliveryOption struct {
	WaitForDelivery bool
}

func (w *waitForDeliveryOption) Apply(o *options) {
	o.WaitForDelivery = w.WaitForDelivery
}

func WithWaitForDelivery(b bool) Option {
	return &waitForDeliveryOption{WaitForDelivery: b}
}

type timeoutOption struct {
	Timeout time.Duration
}

func (t *timeoutOption) Apply(o *options) {
	o.Timeout = t.Timeout
}

func WithTimeout(t time.Duration) Option {
	return &timeoutOption{Timeout: t}
}

type reporter func(error) bool

type reportOnOption struct {
	ReportOn reporter
}

func (r *reportOnOption) Apply(o *options) {
	o.ReportOn = r.ReportOn
}

func WithReportOn(r reporter) Option {
	return &reportOnOption{ReportOn: r}
}
