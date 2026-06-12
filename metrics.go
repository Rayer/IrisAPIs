package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

var (
	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_cpu_usage_percent",
		Help: "CPU usage percentage (overall).",
	})
	memTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_memory_total_bytes",
		Help: "Total memory in bytes.",
	})
	memUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_memory_used_bytes",
		Help: "Used memory in bytes.",
	})
	memAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_memory_available_bytes",
		Help: "Available memory in bytes.",
	})
	diskTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_disk_total_bytes",
		Help: "Total disk space in bytes.",
	})
	diskUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_disk_used_bytes",
		Help: "Used disk space in bytes.",
	})
	diskFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_disk_free_bytes",
		Help: "Free disk space in bytes.",
	})
	load1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_load_1m",
		Help: "Load average 1 minute.",
	})
	load5 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_load_5m",
		Help: "Load average 5 minutes.",
	})
	load15 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "iris_load_15m",
		Help: "Load average 15 minutes.",
	})
)

func init() {
	prometheus.MustRegister(cpuUsage, memTotal, memUsed, memAvailable,
		diskTotal, diskUsed, diskFree, load1, load5, load15)
}

func collectMetrics() {
	if percent, err := cpu.Percent(0, false); err == nil && len(percent) > 0 {
		cpuUsage.Set(percent[0])
	}
	if v, err := mem.VirtualMemory(); err == nil {
		memTotal.Set(float64(v.Total))
		memUsed.Set(float64(v.Used))
		memAvailable.Set(float64(v.Available))
	}
	if v, err := disk.Usage("/"); err == nil {
		diskTotal.Set(float64(v.Total))
		diskUsed.Set(float64(v.Used))
		diskFree.Set(float64(v.Free))
	}
	if v, err := load.Avg(); err == nil {
		load1.Set(v.Load1)
		load5.Set(v.Load5)
		load15.Set(v.Load15)
	}
}

func metricsHandler() http.Handler {
	return promhttp.Handler()
}
