package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://pkg.go.dev/github.com/prometheus/client_golang@v1.20.2/prometheus

type PrometheusMetrics struct {
	reg   *prometheus.Registry
	gauge prometheus.Gauge
}

func NewPrometheusMetrics() *PrometheusMetrics {
	pm := &PrometheusMetrics{
		reg: prometheus.NewRegistry(),
		gauge: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "test",
			Help: "Utilization of GPU",
		}),
	}
	pm.reg.MustRegister(pm.gauge)
	return pm
}

func (pm PrometheusMetrics) Metrics(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	// registry := prometheus.NewRegistry()

	// gauge := prometheus.NewGauge(prometheus.GaugeOpts{
	// 	Name: "test",
	// 	Help: "Utilization of GPU",
	// })
	// registry.MustRegister(gauge)

	// gauge.Set(30)
	// gatherers := prometheus.Gatherers{
	// 	registry,
	// }

	pm.gauge.Set(30)
	gatherers := prometheus.Gatherers{
		pm.reg,
	}

	h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
