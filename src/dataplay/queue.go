package main

import (
	"encoding/json"
	"fmt"
)

type Queue struct {
	QueueProducer
	QueueConsumer
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
	fmt.Println(*self)
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
