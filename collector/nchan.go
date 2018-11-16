package collector

import (
	"log"
	"sync"

	"github.com/platanus/nchan-prometheus-exporter/nchanClient"
	"github.com/prometheus/client_golang/prometheus"
)

// NchanCollector collects Nchan metrics. It implements prometheus.Collector interface.
type NchanCollector struct {
	nchanClient *nchanClient.NchanClient
	metrics     map[string]*prometheus.Desc
	mutex       sync.Mutex
}

// NewNchanCollector creates an NchanCollector.
func NewNchanCollector(nchanClient *nchanClient.NchanClient, namespace string) *NchanCollector {
	return &NchanCollector{
		nchanClient: nchanClient,
		metrics: map[string]*prometheus.Desc{
			"messages_published_total":           newGlobalMetric(namespace, "messages_published_total", "Number of messages published to all channels through this Nchan server"),
			"messages_stored":                    newGlobalMetric(namespace, "messages_stored", "Number of messages currently buffered in memory"),
			"shared_memory_used":                 newGlobalMetric(namespace, "shared_memory_used", "Total shared memory used for buffering messages, storing channel information, and other purposes"),
			"shared_memory_limit":                newGlobalMetric(namespace, "shared_memory_limit", "Total shared memory used for buffering messages, storing channel information, and other purposes"),
			"channels":                           newGlobalMetric(namespace, "channels", "Number of channels present on this Nchan server"),
			"subscribers":                        newGlobalMetric(namespace, "subscribers", "Number of subscribers to all channels on this Nchan server"),
			"redis_pending_commands":             newGlobalMetric(namespace, "redis_pending_commands", "Number of commands sent to Redis that are awaiting a reply"),
			"redis_connected_servers":            newGlobalMetric(namespace, "redis_connected_servers", "Number of redis servers to which Nchan is currently connected"),
			"interprocess_alerts_received_total": newGlobalMetric(namespace, "interprocess_alerts_received_total", "Number of interprocess communication packets transmitted between Nginx workers processes for Nchan"),
			"interprocess_alerts_in_transit":     newGlobalMetric(namespace, "interprocess_alerts_in_transit", "Number of interprocess communication packets in transit between Nginx workers"),
			"interprocess_queued_alerts":         newGlobalMetric(namespace, "interprocess_queued_alerts", "Number of interprocess communication packets waiting to be sent"),
			"interprocess_send_delay_total":      newGlobalMetric(namespace, "interprocess_send_delay_total", "Total amount of time interprocess communication packets spend being queued if delayed"),
			"interprocess_receive_delay_total":   newGlobalMetric(namespace, "interprocess_receive_delay_total", "Total amount of time interprocess communication packets spend in transit if delayed"),
		},
	}
}

// Describe sends the super-set of all possible descriptors of NGINX metrics
// to the provided channel.
func (c *NchanCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

// Collect fetches metrics from NGINX and sends them to the provided channel.
func (c *NchanCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // To protect metrics from concurrent collects
	defer c.mutex.Unlock()

	stats, err := c.nchanClient.GetStubStats()
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.metrics["messages_published_total"],
		prometheus.CounterValue, float64(stats.Messages.TotalPublished))
	ch <- prometheus.MustNewConstMetric(c.metrics["messages_stored"],
		prometheus.CounterValue, float64(stats.Messages.Stored))
	ch <- prometheus.MustNewConstMetric(c.metrics["shared_memory_used"],
		prometheus.GaugeValue, float64(stats.SharedMemoryUsed))
	ch <- prometheus.MustNewConstMetric(c.metrics["shared_memory_limit"],
		prometheus.GaugeValue, float64(stats.SharedMemoryLimit))
	ch <- prometheus.MustNewConstMetric(c.metrics["channels"],
		prometheus.GaugeValue, float64(stats.Channels))
	ch <- prometheus.MustNewConstMetric(c.metrics["subscribers"],
		prometheus.GaugeValue, float64(stats.Subscribers))
	ch <- prometheus.MustNewConstMetric(c.metrics["redis_pending_commands"],
		prometheus.GaugeValue, float64(stats.Redis.PendingCommands))
	ch <- prometheus.MustNewConstMetric(c.metrics["redis_connected_servers"],
		prometheus.GaugeValue, float64(stats.Redis.ConnectedServers))
	ch <- prometheus.MustNewConstMetric(c.metrics["interprocess_alerts_received_total"],
		prometheus.CounterValue, float64(stats.Interprocess.TotalAlertsReceived))
	ch <- prometheus.MustNewConstMetric(c.metrics["interprocess_alerts_in_transit"],
		prometheus.GaugeValue, float64(stats.Interprocess.AlertsInTransit))
	ch <- prometheus.MustNewConstMetric(c.metrics["interprocess_queued_alerts"],
		prometheus.GaugeValue, float64(stats.Interprocess.QueuedAlerts))
	ch <- prometheus.MustNewConstMetric(c.metrics["interprocess_send_delay_total"],
		prometheus.CounterValue, float64(stats.Interprocess.TotalSendDelay))
	ch <- prometheus.MustNewConstMetric(c.metrics["interprocess_receive_delay_total"],
		prometheus.CounterValue, float64(stats.Interprocess.TotalReceiveDelay))
}
