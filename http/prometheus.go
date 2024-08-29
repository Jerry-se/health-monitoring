package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// https://pkg.go.dev/github.com/prometheus/client_golang@v1.20.2/prometheus

type PrometheusMetrics struct {
	reg                 *prometheus.Registry
	utilizationGPUGauge *prometheus.GaugeVec
	memoryTotalGauge    *prometheus.GaugeVec
	memoryUsedGauge     *prometheus.GaugeVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	pm := &PrometheusMetrics{
		reg: prometheus.NewRegistry(),
		utilizationGPUGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "utilization_gpu",
				Help: "utilization of GPU",
			},
			[]string{"job", "instance"},
		),
		memoryTotalGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "memory_total",
				Help: "total memory",
			},
			[]string{"job", "instance"},
		),
		memoryUsedGauge: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "memory_used",
				Help: "used memory",
			},
			[]string{"job", "instance"},
		),
	}
	pm.reg.MustRegister(pm.utilizationGPUGauge)
	pm.reg.MustRegister(pm.memoryTotalGauge)
	pm.reg.MustRegister(pm.memoryUsedGauge)
	return pm
}

func (pm PrometheusMetrics) Metrics(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	// pm.gauge.Set(30)
	pm.utilizationGPUGauge.WithLabelValues("test", "machine1").Set(30)
	pm.memoryTotalGauge.WithLabelValues("test", "machine1").Set(24564)
	pm.memoryUsedGauge.WithLabelValues("test", "machine1").Set(22128)
	// pm.memoryUsedGauge.With(prometheus.Labels{"job": "test", "instance": "machine1"}).Set(22128)
	gatherers := prometheus.Gatherers{
		pm.reg,
	}

	h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
