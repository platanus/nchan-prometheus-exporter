package nchanClient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// NchanClient allows you to fetch Nchan metrics from the stub_status page.
type NchanClient struct {
	apiEndpoint string
	httpClient  *http.Client
}

// StubStats represents Nchan stub_status metrics.
type StubStats struct {
	Redis             StubRedis
	Interprocess      StubInterprocess
	Messages          StubMessages
	Channels          int64
	Subscribers       int64
	SharedMemoryUsed  int64
	SharedMemoryLimit int64
}

// StubRedis represents redis related metrics.
type StubRedis struct {
	PendingCommands  int64
	ConnectedServers int64
}

// StubInterprocess represents interprocess related metrics.
type StubInterprocess struct {
	AlertsInTransit     int64
	QueuedAlerts        int64
	TotalAlertsReceived int64
	TotalSendDelay      int64
	TotalReceiveDelay   int64
}

// StubMessages represents messages related metrics.
type StubMessages struct {
	TotalPublished int64
	Stored         int64
}

// NewNchanClient creates an NchanClient.
func NewNchanClient(httpClient *http.Client, apiEndpoint string) (*NchanClient, error) {
	client := &NchanClient{
		apiEndpoint: apiEndpoint,
		httpClient:  httpClient,
	}

	if _, err := client.GetStubStats(); err != nil {
		return nil, fmt.Errorf("Failed to create NchanClient: %v", err)
	}

	return client, nil
}

// GetStubStats fetches the stub_status metrics.
func (client *NchanClient) GetStubStats() (*StubStats, error) {
	resp, err := client.httpClient.Get(client.apiEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get %v: %v", client.apiEndpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %v response, got %v", http.StatusOK, resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %v", err)
	}

	var stats StubStats
	err = parseStubStats(body, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body %q: %v", string(body), err)
	}

	return &stats, nil
}

func parseStubStats(data []byte, stats *StubStats) error {
	dataStr := string(data)

	parts := strings.Split(dataStr, "\n")
	if len(parts) != 15 {
		return fmt.Errorf("invalid input %q", dataStr)
	}

	totalPublishedMessagesParts := strings.Split(strings.TrimSpace(parts[0]), " ")
	if len(totalPublishedMessagesParts) != 4 {
		return fmt.Errorf("invalid input for total published messages %q", parts[0])
	}

	totalPublishedMessages, err := strconv.ParseInt(totalPublishedMessagesParts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for total published messages %q: %v", totalPublishedMessagesParts[3], err)
	}
	stats.Messages.TotalPublished = totalPublishedMessages

	storedMessagesParts := strings.Split(strings.TrimSpace(parts[1]), " ")
	if len(storedMessagesParts) != 3 {
		return fmt.Errorf("invalid input for stored messages %q", parts[1])
	}

	storedMessages, err := strconv.ParseInt(storedMessagesParts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for stored messages %q: %v", storedMessagesParts[2], err)
	}
	stats.Messages.Stored = storedMessages

	sharedMemoryUsedParts := strings.Split(strings.TrimSpace(parts[2]), " ")
	if len(sharedMemoryUsedParts) != 4 {
		return fmt.Errorf("invalid input for shared memory used %q", parts[2])
	}

	sharedMemoryUsed, err := strconv.ParseInt(strings.Replace(sharedMemoryUsedParts[3], "K", "", 1), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for shared memory used %q: %v", sharedMemoryUsedParts[3], err)
	}
	stats.SharedMemoryUsed = sharedMemoryUsed

	sharedMemoryLimitParts := strings.Split(strings.TrimSpace(parts[3]), " ")
	if len(sharedMemoryLimitParts) != 4 {
		return fmt.Errorf("invalid input for shared memory limit %q", parts[3])
	}

	sharedMemoryLimit, err := strconv.ParseInt(strings.Replace(sharedMemoryLimitParts[3], "K", "", 1), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for shared memory limit %q: %v", sharedMemoryLimitParts[3], err)
	}
	stats.SharedMemoryLimit = sharedMemoryLimit

	channelsParts := strings.Split(strings.TrimSpace(parts[4]), " ")
	if len(channelsParts) != 2 {
		return fmt.Errorf("invalid input for channels %q", parts[4])
	}

	channels, err := strconv.ParseInt(channelsParts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for channels %q: %v", channelsParts[1], err)
	}
	stats.Channels = channels

	subscribersParts := strings.Split(strings.TrimSpace(parts[5]), " ")
	if len(subscribersParts) != 2 {
		return fmt.Errorf("invalid input for subscribers %q", parts[5])
	}

	subscribers, err := strconv.ParseInt(subscribersParts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for subscribers %q: %v", subscribersParts[1], err)
	}
	stats.Subscribers = subscribers

	redisPendingCommandsParts := strings.Split(strings.TrimSpace(parts[6]), " ")
	if len(redisPendingCommandsParts) != 4 {
		return fmt.Errorf("invalid input for redis pending commands %q", parts[6])
	}

	redisPendingCommands, err := strconv.ParseInt(redisPendingCommandsParts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for redis pending commands %q: %v", redisPendingCommandsParts[3], err)
	}
	stats.Redis.PendingCommands = redisPendingCommands

	redisConnectedServersParts := strings.Split(strings.TrimSpace(parts[7]), " ")
	if len(redisConnectedServersParts) != 4 {
		return fmt.Errorf("invalid input for redis connected servers %q", parts[7])
	}

	redisConnectedServers, err := strconv.ParseInt(redisConnectedServersParts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for redis connected servers %q: %v", redisConnectedServersParts[3], err)
	}
	stats.Redis.ConnectedServers = redisConnectedServers

	totalInterprocessAlertsReceivedParts := strings.Split(strings.TrimSpace(parts[8]), " ")
	if len(totalInterprocessAlertsReceivedParts) != 5 {
		return fmt.Errorf("invalid input for interprocess total alerts received %q", parts[8])
	}

	totalInterprocessAlertsReceived, err := strconv.ParseInt(totalInterprocessAlertsReceivedParts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for interprocess total alerts received %q: %v", totalInterprocessAlertsReceivedParts[4], err)
	}
	stats.Interprocess.TotalAlertsReceived = totalInterprocessAlertsReceived

	interprocessAlertsInTransitParts := strings.Split(strings.TrimSpace(parts[9]), " ")
	if len(interprocessAlertsInTransitParts) != 5 {
		return fmt.Errorf("invalid input for interprocess alerts in transit %q", parts[9])
	}

	interprocessAlertsInTransit, err := strconv.ParseInt(interprocessAlertsInTransitParts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for interprocess alerts in transit %q: %v", interprocessAlertsInTransitParts[4], err)
	}
	stats.Interprocess.AlertsInTransit = interprocessAlertsInTransit

	interprocessQueuedAlertsParts := strings.Split(strings.TrimSpace(parts[10]), " ")
	if len(interprocessQueuedAlertsParts) != 4 {
		return fmt.Errorf("invalid input for interprocess queued alerts %q", parts[10])
	}

	interprocessQueuedAlerts, err := strconv.ParseInt(interprocessQueuedAlertsParts[3], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for interprocess queued alerts %q: %v", interprocessQueuedAlertsParts[3], err)
	}
	stats.Interprocess.QueuedAlerts = interprocessQueuedAlerts

	totalInterprocessSendDelayParts := strings.Split(strings.TrimSpace(parts[11]), " ")
	if len(totalInterprocessSendDelayParts) != 5 {
		return fmt.Errorf("invalid input for interprocess total send delay %q", parts[11])
	}

	totalInterprocessSendDelay, err := strconv.ParseInt(totalInterprocessSendDelayParts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for interprocess total send delay %q: %v", totalInterprocessSendDelayParts[4], err)
	}
	stats.Interprocess.TotalSendDelay = totalInterprocessSendDelay

	totalInterprocessReceiveDelayParts := strings.Split(strings.TrimSpace(parts[12]), " ")
	if len(totalInterprocessReceiveDelayParts) != 5 {
		return fmt.Errorf("invalid input for interprocess total receive delay %q", parts[12])
	}

	totalInterprocessReceiveDelay, err := strconv.ParseInt(totalInterprocessReceiveDelayParts[4], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid input for interprocess total receive delay %q: %v", totalInterprocessReceiveDelayParts[4], err)
	}
	stats.Interprocess.TotalReceiveDelay = totalInterprocessReceiveDelay

	return nil
}
