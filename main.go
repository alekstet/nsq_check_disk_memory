package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/alekstet/nsq_check_disk_memory/models"
)

func (nsq_data *models.WriteNSQ) ToNSQ(address, port, topic, channel string) {
	fmt.Println(address, port, topic, channel)
	mes := []string{address, port, topic, channel}
	payload, err := json.Marshal(mes)
	if err != nil {
		log.Fatalf("Error with Marshal: %v\n", err)
	}
	err = nsq_data.Prod.Publish(nsq_data.Topic, payload)
	if err != nil {
		log.Fatalf("Error with publish to NSQ: %v\n", err)
	}
}

func (nsq_data *models.WriteNSQ) MemoryChecker(nsqlookupd_address string) {
	client := http.Client{}
	nsqd_arr := [][]string{}

	url := "http://" + nsqlookupd_address + "/nodes"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	var m *models.Nsqlookupd
	json.NewDecoder(resp.Body).Decode(&m)

	for _, i := range m.Producers {
		info := []string{i.BroadcastAddress, strconv.Itoa(i.HTTPPort)}
		nsqd_arr = append(nsqd_arr, info)
	}

	for _, i := range nsqd_arr {
		url := "http://" + i[0] + ":" + i[1] + "/stats?format=json"
		resp, err := client.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		var n *models.Nsqd
		json.NewDecoder(resp.Body).Decode(&n)
		for _, j := range n.Topics {

			empty_ch := []*models.ChannelStats{}

			if !reflect.DeepEqual(empty_ch, j.Channels) {
				for _, k := range j.Channels {
					if k.Depth > int64(nsq_data.Memory_mes) || k.BackendDepth > int64(nsq_data.Disk_mes) {
						nsq_data.ToNSQ(i[0], i[1], j.TopicName, k.ChannelName)
					}
				}
			}
		}
	}
}

func main() {
	nsq_data, nsqlookupd_address := conf.ReadConfig()

	log.Println("Config read Ok")

	ticker_polling := time.NewTicker((time.Duration(nsq_data.Polling_period)) * time.Second)
	ticker_message := time.NewTicker((time.Duration(nsq_data.Test_message_period)) * time.Second)

	for {
		go func() {
			for range ticker_polling.C {
				nsq_data.MemoryChecker(nsqlookupd_address)
			}
		}()

		go func() {
			for range ticker_message.C {
				nsq_data.ToNSQ("I", "am", "still", "working")
			}
		}()
	}
}
