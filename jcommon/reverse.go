package jcommon

import "reflect"

func Reverse(slice interface{}) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice
	}

	ns := reflect.Indirect(reflect.New(v.Type()))
	for i := v.Len(); i > 0; i-- {
		ns = reflect.Append(ns, v.Index(i-1))
	}
	return ns.Interface()
}
