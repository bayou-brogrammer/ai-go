package main

import (
	"reflect"
)

func GetReflectType(component interface{}) reflect.Type {
	return reflect.TypeOf(component)
}
