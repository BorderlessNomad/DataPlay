package main

import (
	"reflect"
)

type Queue struct {
	QueueProducer
	QueueConsumer
}

func RunMethod(name string, params ...interface{}) []reflect.Value {
	var q Queue
	args := []reflect.Value{}

	for i, _ := range params {
		args = append(args, reflect.ValueOf(params[i]))
	}

	r := reflect.ValueOf(&q).MethodByName(name).Call(args)
	return r
}
