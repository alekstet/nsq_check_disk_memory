package main

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/olebedev/config"
)

func config_file() []int {
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	test_message_period, err := cfg.Int("to_nsq.test_message_period")
	if err != nil {
		panic(err)
	}

	polling_period, err := cfg.Int("to_nsq.polling_period")
	if err != nil {
		panic(err)
	}

	memory_mes, err := cfg.Int("from_nsq.memory_mes")
	if err != nil {
		panic(err)
	}

	disk_mes, err := cfg.Int("from_nsq.disk_mes")
	if err != nil {
		panic(err)
	}

	number_of_periods_ch, err := cfg.Int("testing.number_of_periods_ch")
	if err != nil {
		panic(err)
	}

	number_of_periods_srv, err := cfg.Int("testing.number_of_periods_srv")
	if err != nil {
		panic(err)
	}

	res := []int{test_message_period, polling_period, number_of_periods_ch, number_of_periods_srv, memory_mes, disk_mes}
	return res
}

func check_filling() int64 {
	infos := config_file()
	period := infos[1]
	number_of_periods_ch := infos[2]
	test_time := (number_of_periods_ch * period) + 1
	res_before := get_number_of_messages("10.50.0.201:4161", "memory_checker")
	need_quentity_mes := infos[4] + 2
	for i := 0; i < need_quentity_mes; i++ {
		add_mes("cip3_metrics")
	}
	go main()
	time.Sleep(time.Duration(test_time) * time.Second)
	res_after := get_number_of_messages("10.50.0.201:4161", "memory_checker")
	if res_after-res_before == int64(number_of_periods_ch) || res_after-res_before == int64(number_of_periods_ch)+1 {
		return int64(number_of_periods_ch)
	} else {
		return res_after - res_before
	}
}

func check_release() int64 {
	infos := config_file()
	period := infos[0]
	number_of_periods_srv := infos[3]
	testing_period := (number_of_periods_srv * period) + 1
	res_before := get_number_of_messages("10.50.0.201:4161", "memory_checker")
	go main()
	time.Sleep(time.Duration(testing_period) * time.Second)
	res_after := get_number_of_messages("10.50.0.201:4161", "memory_checker")
	if res_after-res_before == int64(number_of_periods_srv) || res_after-res_before == int64(number_of_periods_srv)+1 {
		return int64(number_of_periods_srv)
	} else {
		return res_after - res_before
	}
}

func Test(t *testing.T) {
	infos := config_file()
	number_of_periods_ch := infos[2]
	number_of_periods_srv := infos[3]

	empty_topic("cip3_metrics", "prometheus")

	check1 := check_filling()
	if check1 != int64(number_of_periods_ch) && check1 != int64(number_of_periods_ch)+1 {
		t.Errorf("Expected %d (or + 1) messages, got %d", int64(number_of_periods_ch), check1)
	}

	empty_topic("cip3_metrics", "prometheus")

	check2 := check_release()
	if check2 != int64(number_of_periods_srv) && check2 != int64(number_of_periods_srv)+1 {
		t.Errorf("Expected %d (or + 1) messages, got %d", int64(number_of_periods_srv), check2)
	}
}
