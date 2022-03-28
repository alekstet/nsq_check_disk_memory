package main

import "github.com/nsqio/go-nsq"

type WriteNSQ struct {
	Prod                *nsq.Producer
	Config_nsq          *nsq.Config
	Topic               string
	Memory_mes          int
	Disk_mes            int
	Test_message_period int
	Polling_period      int
}

type Nsqlookupd struct {
	Producers []Producer `json:"producers"`
}

type Producer struct {
	RemoteAddress    string   `json:"producers"`
	Hostname         string   `json:"hostname"`
	BroadcastAddress string   `json:"broadcast_address"`
	TCPPort          int      `json:"tcp_port"`
	HTTPPort         int      `json:"http_port"`
	Version          string   `json:"version"`
	Tombstones       []bool   `json:"tombstones"`
	Topics           []string `json:"topics"`
}

type Nsqd struct {
	Version    string          `json:"version"`
	Health     string          `json:"health"`
	Start_time int64           `json:"start_time"`
	Topics     []TopicStats    `json:"topics"`
	Memory     MemStats        `json:"memory"`
	Producers  []ProducerStats `json:"producers"`
}

type TopicStats struct {
	TopicName            string              `json:"topic_name"`
	Channels             []ChannelStats      `json:"channels"`
	Depth                int64               `json:"depth"`
	BackendDepth         int64               `json:"backend_depth"`
	MessageCount         uint64              `json:"message_count"`
	MessageBytes         uint64              `json:"message_bytes"`
	Paused               bool                `json:"paused"`
	E2eProcessingLatency []ProcessingLatency `json:"e2e_processing_latency"`
}

type ChannelStats struct {
	ChannelName          string            `json:"channel_name"`
	Depth                int64             `json:"depth"`
	BackendDepth         int64             `json:"backend_depth"`
	InFlightCount        int               `json:"in_flight_count"`
	DeferredCount        int               `json:"deferred_count"`
	MessageCount         uint64            `json:"message_count"`
	RequeueCount         uint64            `json:"requeue_count"`
	TimeoutCount         uint64            `json:"timeout_count"`
	ClientCount          int               `json:"client_count"`
	Clients              []ClientStats     `json:"clients"`
	Paused               bool              `json:"paused"`
	E2eProcessingLatency ProcessingLatency `json:"e2e_processing_latency"`
}

type ClientStats struct {
	ClientID                      string `json:"client_id"`
	Hostname                      string `json:"hostname"`
	Version                       string `json:"version"`
	RemoteAddress                 string `json:"remote_address"`
	State                         int    `json:"state"`
	ReadyCount                    int    `json:"ready_count"`
	InFlightCount                 int    `json:"in_flight_count"`
	MessageCount                  int64  `json:"message_count"`
	FinishCount                   int64  `json:"finish_count"`
	RequeueCount                  int64  `json:"requeue_count"`
	ConnectTs                     int64  `json:"connect_ts"`
	SampleRate                    int32  `json:"sample_rate"`
	Deflate                       bool   `json:"deflate"`
	Snappy                        bool   `json:"snappy"`
	UserAgent                     string `json:"user_agent"`
	TLS                           bool   `json:"tls"`
	CipherSuite                   string `json:"tls_cipher_suite"`
	TLSVersion                    string `json:"tls_version"`
	TLSNegotiatedProtocol         string `json:"tls_negotiated_protocol"`
	TLSNegotiatedProtocolIsMutual bool   `json:"tls_negotiated_protocol_is_mutual"`
}

type MemStats struct {
	HeapObjects       uint64 `json:"heap_objects"`
	HeapIdleBytes     uint64 `json:"heap_idle_bytes"`
	HeapInUseBytes    uint64 `json:"heap_in_use_bytes"`
	HeapReleasedBytes uint64 `json:"heap_released_bytes"`
	GCPauseUsec100    uint64 `json:"gc_pause_usec_100"`
	GCPauseUsec99     uint64 `json:"gc_pause_usec_99"`
	GCPauseUsec95     uint64 `json:"gc_pause_usec_95"`
	NextGCBytes       uint64 `json:"next_gc_bytes"`
	GCTotalRuns       uint32 `json:"gc_total_runs"`
}

type ProducerStats struct {
	ClientID                      string       `json:"client_id"`
	Hostname                      string       `json:"hostname"`
	Version                       string       `json:"version"`
	RemoteAddress                 string       `json:"remote_address"`
	State                         int          `json:"state"`
	ReadyCount                    int          `json:"ready_count"`
	InFlightCount                 int          `json:"in_flight_count"`
	MessageCount                  int64        `json:"message_count"`
	FinishCount                   int64        `json:"finish_count"`
	RequeueCount                  int64        `json:"requeue_count"`
	ConnectTs                     int64        `json:"connect_ts"`
	SampleRate                    int32        `json:"sample_rate"`
	Deflate                       bool         `json:"deflate"`
	Snappy                        bool         `json:"snappy"`
	UserAgent                     string       `json:"user_agent"`
	PubCounts                     []Pub_Counts `json:"pub_counts"`
	TLS                           bool         `json:"tls"`
	CipherSuite                   string       `json:"tls_cipher_suite"`
	TLSVersion                    string       `json:"tls_version"`
	TLSNegotiatedProtocol         string       `json:"tls_negotiated_protocol"`
	TLSNegotiatedProtocolIsMutual bool         `json:"tls_negotiated_protocol_is_mutual"`
}

type ProcessingLatency struct {
	Count       int                  `json:"count"`
	Percentiles []map[string]float64 `json:"percentiles"`
}

type Pub_Counts struct {
	Count       string `json:"topic"`
	Percentiles int    `json:"count"`
}
