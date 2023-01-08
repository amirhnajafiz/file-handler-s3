package metric

import "github.com/prometheus/client_golang/prometheus"

const (
	Namespace = "fhs"
	Subsystem = "system"
)

// Metrics has all the client metrics.
type Metrics struct {
	Failed   prometheus.Counter
	Requests prometheus.Counter
	BitRate  prometheus.Histogram
}

func newCounter(counterOpts prometheus.CounterOpts) prometheus.Counter {
	ev := prometheus.NewCounter(counterOpts)

	if err := prometheus.Register(ev); err != nil {
		panic(err)
	}

	return ev
}

func newHistogram(histogramOpts prometheus.HistogramOpts) prometheus.Histogram {
	ev := prometheus.NewHistogram(histogramOpts)

	if err := prometheus.Register(ev); err != nil {
		panic(err)
	}

	return ev
}

func NewMetrics() Metrics {
	return Metrics{
		Failed: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "connection_errors_total",
			Help:        "total number of connection errors",
			ConstLabels: nil,
		}),
		Requests: newCounter(prometheus.CounterOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "total_number_of_requests",
			Help:        "total number of requests",
			ConstLabels: nil,
		}),
		BitRate: newHistogram(prometheus.HistogramOpts{
			Namespace:   Namespace,
			Subsystem:   Subsystem,
			Name:        "stream_bit_rate",
			Help:        "bit rate of our stream server",
			ConstLabels: nil,
			Buckets:     prometheus.DefBuckets,
		}),
	}
}
