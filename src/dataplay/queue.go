package main

import (
	"reflect"
)

type Q struct{}

func RunMethod(name string, params ...interface{}) []reflect.Value {
	var q Q
	var r []reflect.Value
	p := len(params)

	switch p {
	case 0:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{})
	case 1:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{reflect.ValueOf(params[0])})
	case 2:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{reflect.ValueOf(params[0]), reflect.ValueOf(params[1])})
	case 3:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{reflect.ValueOf(params[0]), reflect.ValueOf(params[1]), reflect.ValueOf(params[2])})
	case 4:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{reflect.ValueOf(params[0]), reflect.ValueOf(params[1]), reflect.ValueOf(params[2]), reflect.ValueOf(params[3])})
	case 5:
		r = reflect.ValueOf(&q).MethodByName(name).Call([]reflect.Value{reflect.ValueOf(params[0]), reflect.ValueOf(params[1]), reflect.ValueOf(params[2]), reflect.ValueOf(params[3]), reflect.ValueOf(params[4])})
	default:
		panic("too many arguments")
	}
	return r
}
