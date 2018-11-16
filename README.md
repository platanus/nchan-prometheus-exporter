[![Build Status](https://travis-ci.org/platanus/nchan-prometheus-exporter.svg?branch=master)](https://travis-ci.org/platanus/nchan-prometheus-exporter) [![](https://images.microbadger.com/badges/version/platanus/nchan-prometheus-exporter.svg)](https://microbadger.com/images/platanus/nchan-prometheus-exporter "Get your own version badge on microbadger.com") [![Go Report Card](https://goreportcard.com/badge/github.com/platanus/nchan-prometheus-exporter)](https://goreportcard.com/report/github.com/platanus/nchan-prometheus-exporter)

# Nchan Prometheus Exporter

Nchan Prometheus exporter makes it possible to monitor Nchan using Prometheus.

## Overview

[Nchan](http://nchan.io) is a scalable, flexible pub/sub server for the modern web, built as a module for the Nginx web server. It provides metrics via de [nchan_stub_status page](https://nchan.io/#nchan_stub_status-stats). Nchan Prometheus exporter fetches the metrics from a single Nchan, converts the metrics into appropriate Prometheus metrics types and finally exposes them via an HTTP server to be collected by [Prometheus](https://prometheus.io/).

## Getting Started

In this section, we show how to quickly run Nchan Prometheus Exporter for Nchan.

### Prerequisites

We assume that you have already installed Prometheus and Nchan. Additionally, you need to:
* Expose the built-in metrics in Nchan:
    * For NCHAN, expose the [nchan_stub_status page](https://nchan.io/#nchan_stub_status-stats) at `/nchan_stub_status` on port `8080`
    * (optional): For NGINX, expose the [stub_status page](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html#stub_status) at `/nginx_stub_status` on port `8080`.
* Configure Prometheus to scrape metrics from the server with the exporter. Note that the default scrape port of the exporter is `9113` and the default metrics path -- `/metrics`.

### Running the Exporter in a Docker Container

To start the exporter we use the [docker run](https://docs.docker.com/engine/reference/run/) command.

* To export Nchan metrics, run:
    ```
    $ docker run -p 9113:9113 platanus/nchan-prometheus-exporter:0.1.0 -scrape-uri http://<nchan>:8080/nchan_stub_status
    ```

* (optional): To additionaly export Nginx metrics, run:
    ```
    $ docker run -p 9113:9113 platanus/nchan-prometheus-exporter:0.1.0 -scrape-uri http://<nchan>:8080/nchan_stub_status -nginx -nginx.scrape-uri http://<nchan>:8080/nginx_stub_status
    ```

> where `<nchan>` is the IP address/DNS name, through which Nchan is available.

### Running the Exporter Binary

* To export Nchan metrics, run:
    ```
    $ nchan-prometheus-exporter -nchan.scrape-uri http://<nchan>:8080/nchan_stub_status
    ```

* (optional): To additionaly export Nginx metrics, run:
    ```
    $ nchan-prometheus-exporter -nchan.scrape-uri http://<nchan>:8080/nchan_stub_status -nginx -nginx.scrape-uri http://<nchan>:8080/nginx_stub_status
    ```

> where `<nchan>` is the IP address/DNS name, through which Nchan is available.

**Note**. The `nchan-prometheus-exporter` is not a daemon. To run the exporter as a system service (daemon), configure the init system of your Linux server (such as systemd or Upstart) accordingly. Alternatively, you can run the exporter in a Docker container.

## Usage

### Command-line Arguments

```
Usage of ./nchan-prometheus-exporter:
  -scrape-uri string
        A URI for scraping Nchan metrics.
        The nchan_stub_status page must be available through the URI. The default value can be overwritten by SCRAPE_URI environment variable. (default "http://127.0.0.1:8080/nchan_stub_status")
  -nginx bool
        Start the exporter with NGINX metrics support. The default value can be overwritten by NGINX environment variable.
  -nginx.scrape-uri string
        A URI for scraping NGINX metrics.
        For NGINX, the stub_status page must be available through the URI. The default value can be overwritten by NGINX_SCRAPE_URI environment variable. (default "http://127.0.0.1:8080/nginx_stub_status")
  -ssl-verify
        Perform SSL certificate verification. The default value can be overwritten by SSL_VERIFY environment variable.
  -web.listen-address string
        An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable. (default ":9113")
  -web.telemetry-path string
        A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable. (default "/metrics")
```

### Exported Metrics

* For Nchan, all nchan_stub_status metrics are exported. Connect to the `/metrics` page of the running exporter to see the complete list of metrics along with their descriptions.

* For NGINX, all stub_status metrics are exported. Connect to the `/metrics` page of the running exporter to see the complete list of metrics along with their descriptions.

### Troubleshooting

The exporter logs errors to the standard output. When using Docker, if the exporter doesn’t work as expected, check its logs using [docker logs](https://docs.docker.com/engine/reference/commandline/logs/) command.

## Releases

For each release, we publish the corresponding Docker image at `platanus/nchan-prometheus-exporter` [DockerHub repo](https://hub.docker.com/r/platanus/nchan-prometheus-exporter/) and the binaries on the GitHub [releases page](https://github.com/platanus/nchan-prometheus-exporter/releases).

## Building the Exporter

You can build the exporter using the provided Makefile. Before building the exporter, make sure the following software is installed on your machine:
* make
* git
* Docker for building the container image
* Go for building the binary

### Building the Docker Image

To build the Docker image with the exporter, run:
```
$ make container
```

Note: go is not required, as the exporter binary is built in a Docker container. See the [Dockerfile](Dockerfile).

### Building the Binary

To build the binary, run:
```
$ make
```

Note: the binary is built for the OS/arch of your machine. To build binaries for other platforms, see the [Makefile](Makefile).

The binary is built with the name `nchan-prometheus-exporter`.

## Credits

Thank you [contributors](https://github.com/platanus/nchan-prometheus-exporter/graphs/contributors)!

<img src="http://platan.us/gravatar_with_text.png" alt="Platanus" width="250"/>

Nchan Prometheus Exporter is maintained by [platanus](http://platan.us).

## License

Cordova Plugin Flavors is © 2017 platanus, spa. It is free software and may be redistributed under the terms specified in the LICENSE file.
