package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/exec"
)

func runDockerCompose(args ...string) error {
	cmd := exec.Command("docker-compose", args...)
	fmt.Println(cmd)
	var stdout, stderr bytes.Buffer
	// & because the value will change
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Printf("stdout: %s\nstderr: %s\n", stdout.String(), stderr.String())

	return nil
}

func main() {
	errDocker := runDockerCompose("up", "-d")
	if errDocker != nil {
		log.Fatalf("error running docker-compose: %v", errDocker)
	}

	var prometheusPort = flag.Int("prometheus.port", 9150, " port to expose metrics")
	flag.Parse()
	fmt.Println(prometheusPort)

	reg := prometheus.NewRegistry()
	// register a GO collector that will export metrics about a process

	//reg.MustRegister(collectors.NewGoCollector())

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
