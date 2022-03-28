package main

import (
	"io/ioutil"
	"log"

	"github.com/alekstet/nsq_check_disk_memory/conf"
	"github.com/nsqio/go-nsq"
	"github.com/olebedev/config"
)

func ReadConfig() (*conf.WriteNSQ, string) {
	nsq_data := new(conf.WriteNSQ)
	nsq_data.Config_nsq = nsq.NewConfig()

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		log.Fatal(err)
	}

	nsq_prod, err := cfg.String("to_nsq.producer")
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Prod, err = nsq.NewProducer(nsq_prod, nsq_data.Config_nsq)
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Topic, err = cfg.String("to_nsq.topic")
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Test_message_period, err = cfg.Int("to_nsq.test_message_period")
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Polling_period, err = cfg.Int("to_nsq.polling_period")
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Memory_mes, err = cfg.Int("from_nsq.memory_mes")
	if err != nil {
		log.Fatal(err)
	}

	nsq_data.Disk_mes, err = cfg.Int("from_nsq.disk_mes")
	if err != nil {
		log.Fatal(err)
	}

	nsqlookupd_address, err := cfg.String("from_nsq.nsqlookupd_address")
	if err != nil {
		log.Fatal(err)
	}
	return nsq_data, nsqlookupd_address
}
