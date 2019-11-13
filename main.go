package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "ts_gen"
	addr      = ":9556"
	wait      = 3 * time.Second
)

var (
	seasonalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "seasonal_total",
			Help:      "A counter that is incremented in a seasonal fashion",
		})
	randomGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "random",
			Help:      "A Gauge with random numbers",
		})
	trendSeasonalCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        "seasonal_trend_total",
			Help:        "A counter that is increasing in a trend and seasonal fashion",
			ConstLabels: nil,
		})
)

func init() {
	prometheus.MustRegister(seasonalCounter)
	prometheus.MustRegister(randomGauge)
	prometheus.MustRegister(trendSeasonalCounter)
}

func seasonalCounterIncrease() bool {
	if time.Now().Second() <= 10 {
		return true
	}
	return false
}

func main() {
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		for {
			if seasonalCounterIncrease() {
				seasonalCounter.Add(100)
			}
			time.Sleep(wait)
		}
	}()

	go func() {
		i := 10.0
		for {

			if seasonalCounterIncrease() {
				for j := 0; j < 5; j++ {
					trendSeasonalCounter.Add(i + 1500 + rand.Float64() + 1000)
					time.Sleep(wait)
				}
			} else {
				trendSeasonalCounter.Add(i + 150 + rand.Float64() + 60)
			}
			i += float64(rand.Intn(50))
			time.Sleep(wait)
		}
	}()

	go func() {
		for {
			randomGauge.Set(rand.Float64() - float64(rand.Intn(50)) + float64(rand.Intn(50))*rand.Float64())
			time.Sleep(wait)
		}
	}()

	log.Printf("Serving on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
