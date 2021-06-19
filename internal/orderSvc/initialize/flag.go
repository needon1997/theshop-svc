package initialize

import (
	"flag"
	"fmt"
)

var (
	ConfigPath            *string = new(string)
	DevMode               *bool   = new(bool)
	timeoutOrderQueueName *string = new(string)
)

func ParseFlag() {
	flag.StringVar(ConfigPath, "config", "./config.yaml", "the location of the configuration file")
	flag.BoolVar(DevMode, "dev", false, "dev or prod")
	flag.StringVar(timeoutOrderQueueName, "q", "", "the name of timeout order queue in aws sqs_consumer")
	flag.Parse()
	if *timeoutOrderQueueName == "" {
		fmt.Println("You must supply the name of a queue (-q QUEUE)")
		return
	}
}
