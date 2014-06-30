package main

import (
	"encoding/json"
	"reflect"
)

type Queue struct {
	QueueProducer
	QueueConsumer
}

type Message struct {
	MethodName string
	MethodArgs map[string]string
}

func QueueEncode(name string, params map[string]string) string {
	m := Message{
		MethodName: name,
		MethodArgs: params,
	}

	b, _ := json.Marshal(m)
	return string(b)
}

func QueueDecode(msg string) string {
	var message Message
	bmsg := []byte(msg)

	err := json.Unmarshal(bmsg, &message)
	if err != nil {
	}

	r := RunMethodByName(message.MethodName, message.MethodArgs)
	return r
}

func RunMethodByName(name string, params map[string]string) string {
	var q Queue
	args := []reflect.Value{}

	for i, _ := range params {
		args = append(args, reflect.ValueOf(params[i]))
	}

	r := reflect.ValueOf(&q).MethodByName(name).Call(args)
	return r[0].String()
}
