package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/olebedev/config"
)

type Write_nsq struct {
	prod                *nsq.Producer
	config_nsq          *nsq.Config
	topic               string
	memory_mes          int
	disk_mes            int
	test_message_period int
	polling_period      int
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

func (nsq_data *Write_nsq) to_nsq(address, port, topic, channel string) {
	fmt.Println(address, port, topic, channel)
	mes := []string{address, port, topic, channel}
	payload, err := json.Marshal(mes)
	if err != nil {
		log.Panicf("Error with Marshal: %v", err)
	}
	err = nsq_data.prod.Publish(nsq_data.topic, payload)
	if err != nil {
		log.Panicf("Error with publish to NSQ: %v", err)
	}
}

func (nsq_data *Write_nsq) memory_checker(nsqlookupd_address string) {
	client := http.Client{}
	nsqd_arr := [][]string{}

	url := "http://" + nsqlookupd_address + "/nodes"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	var m Nsqlookupd
	json.NewDecoder(resp.Body).Decode(&m)

	for _, j := range m.Producers {
		info := []string{j.BroadcastAddress, strconv.Itoa(j.HTTPPort)}
		nsqd_arr = append(nsqd_arr, info)
	}

	for _, j := range nsqd_arr {
		url := "http://" + j[0] + ":" + j[1] + "/stats?format=json"
		resp, err := client.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		var n Nsqd
		json.NewDecoder(resp.Body).Decode(&n)
		for _, j1 := range n.Topics {
			for _, j2 := range j1.Channels {
				if j2.Depth > int64(nsq_data.memory_mes) || j2.BackendDepth > int64(nsq_data.disk_mes) {
					nsq_data.to_nsq(j[0], j[1], j1.TopicName, j2.ChannelName)
				}
			}
		}
	}
}

func main() {
	nsq_data := new(Write_nsq)
	nsq_data.config_nsq = nsq.NewConfig()

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	nsq_prod, err := cfg.String("to_nsq.producer")
	if err != nil {
		panic(err)
	}

	nsq_data.prod, err = nsq.NewProducer(nsq_prod, nsq_data.config_nsq)
	if err != nil {
		panic(err)
	}

	nsq_data.topic, err = cfg.String("to_nsq.topic")
	if err != nil {
		panic(err)
	}

	nsq_data.test_message_period, err = cfg.Int("to_nsq.test_message_period")
	if err != nil {
		panic(err)
	}

	nsq_data.polling_period, err = cfg.Int("to_nsq.polling_period")
	if err != nil {
		panic(err)
	}

	nsq_data.memory_mes, err = cfg.Int("from_nsq.memory_mes")
	if err != nil {
		panic(err)
	}

	nsq_data.disk_mes, err = cfg.Int("from_nsq.disk_mes")
	if err != nil {
		panic(err)
	}

	nsqlookupd_address, err := cfg.String("from_nsq.nsqlookupd_address")
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(nsq_data.polling_period) * time.Second)
	ticker1 := time.NewTicker(time.Duration(nsq_data.test_message_period) * time.Second)

	for {
		go func() {
			for range ticker.C {
				nsq_data.memory_checker(nsqlookupd_address)
			}
		}()

		go func() {
			for range ticker1.C {
				nsq_data.to_nsq("I", "am", "still", "working")
			}
		}()
	}
}
