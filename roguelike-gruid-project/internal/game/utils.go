package game

import (
	"reflect"
)

func GetReflectType(component interface{}) reflect.Type {
	return reflect.TypeOf(component)
}
