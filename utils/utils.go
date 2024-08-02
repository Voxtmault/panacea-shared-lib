package utils

import (
	"reflect"
)

func ClearObj(arrObj ...interface{}) {
	for _, obj := range arrObj {
		p := reflect.ValueOf(obj).Elem()
		p.Set(reflect.Zero(p.Type()))
	}
}
