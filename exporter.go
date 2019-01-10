package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/platanus/nchan-prometheus-exporter/collector"
	"github.com/platanus/nchan-prometheus-exporter/nchanClient"
	"github.com/platanus/nchan-prometheus-exporter/nginxClient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable value for %s must be a boolean", key)
	}
	return b
}

var (
	// Set during go build
	version   string
	gitCommit string

	// Defaults values
	defaultListenAddress  = getEnv("LISTEN_ADDRESS", ":9113")
	defaultMetricsPath    = getEnv("TELEMETRY_PATH", "/metrics")
	defaultScrapeURI      = getEnv("SCRAPE_URI", "http://127.0.0.1:8080/nchan_stub_status")
	defaultNginx          = getEnvBool("NGINX", false)
	defaultNginxScrapeURI = getEnv("NGINX_SCRAPE_URI", "http://127.0.0.1:8080/stub_status")
	defaultSslVerify      = getEnvBool("SSL_VERIFY", true)

	// Command-line flags
	listenAddr = flag.String("web.listen-address", defaultListenAddress,
		"An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable.")
	metricsPath = flag.String("web.telemetry-path", defaultMetricsPath,
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	scrapeURI = flag.String("scrape-uri", defaultScrapeURI,
		`A URI for scraping Nchan metrics.
	The nchan_stub_status page must be available through the URI. The default value can be overwritten by SCRAPE_URI environment variable.`)
	nginx = flag.Bool("nginx", defaultNginx,
		"Start the exporter with NGINX metrics support. The default value can be overwritten by NGINX environment variable.")
	nginxScrapeURI = flag.String("nginx.scrape-uri", defaultScrapeURI,
		`A URI for scraping NGINX metrics.
	The stub_status page must be available through the URI. The default value can be overwritten by NGINX_SCRAPE_URI environment variable.`)
	sslVerify = flag.Bool("ssl-verify", defaultSslVerify,
		"Perform SSL certificate verification. The default value can be overwritten by SSL_VERIFY environment variable.")
)

func main() {
	flag.Parse()

	log.Printf("Starting Nchan Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	registry := prometheus.NewRegistry()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !*sslVerify},
	}

	nchanClient, err := nchanClient.NewNchanClient(&http.Client{Transport: tr}, *scrapeURI)
	if err != nil {
		log.Fatalf("Could not create Nchan Client: %v", err)
	}

	registry.MustRegister(collector.NewNchanCollector(nchanClient, "nchan"))

	if *nginx {
		nginxClient, err := nginxClient.NewNginxClient(&http.Client{Transport: tr}, *nginxScrapeURI)
		if err != nil {
			log.Fatalf("Could not create Nginx Client: %v", err)
		}

		registry.MustRegister(collector.NewNginxCollector(nginxClient, "nginx"))
	}

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nchan Exporter</title></head>
			<body>
			<h1>Nchan Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
