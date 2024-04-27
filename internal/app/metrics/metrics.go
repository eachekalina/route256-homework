package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	OrdersGiven      prometheus.Counter
	OrdersReturned   prometheus.Counter
	TimeBeforeGiven  prometheus.Histogram
	TimeBeforeReturn prometheus.Histogram
}

func NewMetrics(reg *prometheus.Registry) *Metrics {
	res := &Metrics{
		OrdersGiven: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_given",
			Help: "Number of orders given to customers.",
		}),
		OrdersReturned: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_returned",
			Help: "Number of orders returned by customers.",
		}),
		TimeBeforeGiven: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "time_before_given_sec",
			Help: "Time before order is given to customer.",
		}),
		TimeBeforeReturn: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "time_before_return_sec",
			Help: "Time before order is returned by customer.",
		}),
	}
	reg.MustRegister(res.OrdersGiven, res.OrdersReturned, res.TimeBeforeGiven, res.TimeBeforeReturn)
	return res
}
