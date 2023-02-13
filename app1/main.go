package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type metrics struct {
	devices prometheus.Gauge
	info    *prometheus.GaugeVec
}

// returns the pointer to the metrics struct
func NewMetrics(register prometheus.Registerer) *metrics {
	// since it is a pointer in the return -> work it the memory address
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "app1",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices.",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "app1",
			Name:      "info",
			Help:      "Information about the environment",
		}, []string{"version"}),
	}
	register.MustRegister(m.devices, m.info)
	return m
}

// Hold the connected Devices
var dvs []Device
var version string

func init() {
	version = "2.10.5"

	dvs = []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	// because m now has access to the devices struct
	m.devices.Set(float64(len(dvs)))
	m.info.With(prometheus.Labels{"version": version}).Set(1)

	dMux := http.NewServeMux()
	dMux.HandleFunc("/devices", getDevices)

	prometheusMux := http.NewServeMux()
	prometheusHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	prometheusMux.Handle("/metrics", prometheusHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", dMux))
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8081", prometheusMux))
	}()

	select {}
}

// return all the connected devices
func getDevices(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)

}
