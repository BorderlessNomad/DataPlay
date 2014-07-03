package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

/* Custom config (only use if you want to change defaults) */
var (
	uri          = flag.String("uri", "amqp://playgen:aDam3ntiUm@109.231.121.13:5672/", "AMQP URI")
	exchangeName = flag.String("exchange", "playgen-dev", "Durable (non-auto-deleted) AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")

	requestQueue  = flag.String("requestqueue", "dataplay-request", "Ephemeral AMQP Request queue name")
	responseQueue = flag.String("responsequeue", "dataplay-response", "Ephemeral AMQP Response queue name")

	requestKey  = flag.String("requestkey", "api-request", "AMQP Request routing key")
	responseKey = flag.String("responsekey", "api-response", "AMQP Response routing key")

	requestTag  = flag.String("reqtag", "consumer-request", "AMQP consumer request tag (should not be blank)")
	responseTag = flag.String("restag", "consumer-response", "AMQP consumer response tag (should not be blank)")

	body     = flag.String("body", "foobar", "Body of message")
	reliable = flag.Bool("reliable", true, "Wait for the publisher confirmation before exiting")

	lifetime = flag.Duration("lifetime", 0*time.Second, "lifetime of process before shutdown (0s=infinite, 60s=1minute, 60m=1hour ..)")
)

func init() {
	flag.Parse()
}

type Queue struct {
	QueueProducer
	QueueConsumer
	QueueResponder
}

type Message struct {
	MethodName string
	MethodArgs map[string]string
}

type funcs map[string]func(map[string]string) (ret string)

func (q *Queue) Encode(name string, params map[string]string) string {
	m := Message{
		MethodName: name,
		MethodArgs: params,
	}

	b, _ := json.Marshal(m)
	return string(b)
}

func (q *Queue) Decode(msg []byte) string {
	var message Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
	}

	r := RunMethodByName(message.MethodName, message.MethodArgs)

	return r
}

func (self *funcs) registerCallback(name string, function func(map[string]string) (ret string)) {
	(*self)[name] = function
}

func (self *funcs) callFunction(name string, args map[string]string) (ret string, err error) {
	function := (*self)[name]

	if function != nil {
		return function(args), nil
	}

	return "", fmt.Errorf("function does not exist")
}

func RunMethodByName(name string, params map[string]string) string {
	result, err := myfuncs.callFunction(name, params)
	if err != nil {
		fmt.Println(err)
	}

	return result
}
