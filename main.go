package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/percona/rds_exporter/basic"
	"github.com/percona/rds_exporter/client"
	"github.com/percona/rds_exporter/config"
	"github.com/percona/rds_exporter/enhanced"
	"github.com/percona/rds_exporter/enhanced2"
	"github.com/percona/rds_exporter/sessions"
)

var (
	listenAddressF       = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9042").String()
	basicMetricsPathF    = kingpin.Flag("web.basic-telemetry-path", "Path under which to expose exporter's basic metrics.").Default("/basic").String()
	enhancedMetricsPathF = kingpin.Flag("web.enhanced-telemetry-path", "Path under which to expose exporter's enhanced metrics.").Default("/enhanced").String()
	configFileF          = kingpin.Flag("config.file", "Path to configuration file.").Default("config.yml").String()
	logTraceF            = kingpin.Flag("log.trace", "Enable verbose tracing of AWS requests (will log credentials).").Default("false").Bool()
)

func main() {
	log.AddFlags(kingpin.CommandLine)
	log.Infoln("Starting RDS exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())
	kingpin.Parse()

	cfg, err := config.Load(*configFileF)
	if err != nil {
		log.Fatalf("Can't read configuration file: %s", err)
	}

	client := client.New()
	sessions, err := sessions.New(cfg.Instances, client.HTTP(), *logTraceF)
	if err != nil {
		log.Fatalf("Can't create sessions: %s", err)
	}

	// Basic Metrics
	{
		// Create new Exporter with provided settings and session pool.
		exporter := basic.New(cfg, sessions)
		registry := prometheus.NewRegistry()
		registry.MustRegister(exporter)
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

		http.Handle(*basicMetricsPathF, handler)
	}

	// Enhanced Metrics
	{
		// Create new Exporter with provided settings and session pool.
		exporter := enhanced.New(cfg, sessions)
		registry := prometheus.NewRegistry()
		registry.MustRegister(exporter)
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

		http.Handle(*enhancedMetricsPathF, handler)
	}

	// Enhanced 2 Metrics
	{
		// Create new Exporter with provided settings and session pool.
		exporter := enhanced2.NewCollector(sessions)
		registry := prometheus.NewRegistry()
		registry.MustRegister(exporter)
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorLog:      log.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})

		http.Handle("/enhanced2", handler)
	}

	// Inform user we are ready.
	log.Infoln("RDS exporter started on", *listenAddressF)

	// Start serving for clients
	log.Fatal(http.ListenAndServe(*listenAddressF, nil))
}
