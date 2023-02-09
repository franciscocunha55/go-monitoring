package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	var prometheusPort = flag.Int("prometheus.port", 9150, " port to expose metrics")
	flag.Parse()
	fmt.Println(prometheusPort)

	reg := prometheus.NewRegistry()
	// register a collector that will export metrics about a process
	reg.MustRegister(collectors.NewGoCollector())

	mux := http.NewServeMux()
	prometheusHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", prometheusHandler)

	// start http server

	port := fmt.Sprintf(":%d", *prometheusPort)
	fmt.Println(port)
	log.Printf("starting nginx exporter on %q/metrics", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("cannot start nginx exporter: %s", err)
	}

}
