package main

import (
	"encoding/json"
	"fmt"
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

	fmt.Println("Running ", message.MethodName, message.MethodArgs)

	r := RunMethodByName(message.MethodName, message.MethodArgs)
	return r
}

func RunMethodByName(name string, params map[string]string) string {
	var q Queue
	args := []reflect.Value{}

	for _, v := range params {
		args = append(args, reflect.ValueOf(v))
	}

	fmt.Println(args)

	r := reflect.ValueOf(&q).MethodByName(name).Call(args)
	return r[0].String()
}
