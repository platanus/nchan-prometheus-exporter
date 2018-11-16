package nchanClient

import "testing"

const validStabStats = "total published messages: 123\nstored messages: 54353\nshared memory used: 12K\nchannels: 34\nsubscribers: 5434535\nredis pending commands: 48\nredis connected servers: 65\ntotal interprocess alerts received: 43\ninterprocess alerts in transit: 654\ninterprocess queued alerts: 765\ntotal interprocess send delay: 534\ntotal interprocess receive delay: 46\nnchan version: 1.1.5\n"

func TestParseStubStatsValidInput(t *testing.T) {
	var tests = []struct {
		input          []byte
		expectedResult StubStats
		expectedError  bool
	}{
		{
			input: []byte(validStabStats),
			expectedResult: StubStats{
				Redis: StubRedis{
					PendingCommands:  48,
					ConnectedServers: 65,
				},
				Channels:         34,
				Subscribers:      5434535,
				SharedMemoryUsed: 12,
				Interprocess: StubInterprocess{
					AlertsInTransit:     654,
					QueuedAlerts:        765,
					TotalAlertsReceived: 43,
					TotalSendDelay:      534,
					TotalReceiveDelay:   46,
				},
				Messages: StubMessages{
					TotalPublished: 123,
					Stored:         54353,
				},
			},
			expectedError: false,
		},
		{
			input:         []byte("invalid-stats"),
			expectedError: true,
		},
	}

	for _, test := range tests {
		var result StubStats

		err := parseStubStats(test.input, &result)

		if err != nil && !test.expectedError {
			t.Errorf("parseStubStats() returned error for valid input %q: %v", string(test.input), err)
		}

		if !test.expectedError && test.expectedResult != result {
			t.Errorf("parseStubStats() result %v != expected %v for input %q", result, test.expectedResult, test.input)
		}
	}
}
