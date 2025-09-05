package kubernetes

import (
	"fmt"
	"reflect"
)

// Paginate 分页
func Paginate(items interface{}, page, pageSize int) (interface{}, error) {
	slice := reflect.ValueOf(items)

	if slice.Kind() != reflect.Slice {
		return nil, fmt.Errorf("paginate: items is not a slice")
	}

	startIndex := (page - 1) * pageSize
	endIndex := page * pageSize

	if startIndex > slice.Len() {
		startIndex = slice.Len()
	}
	if endIndex > slice.Len() {
		endIndex = slice.Len()
	}

	paginatedSlice := slice.Slice(startIndex, endIndex).Interface()

	slicePtr := reflect.New(reflect.TypeOf(items))
	slicePtr.Elem().Set(reflect.ValueOf(paginatedSlice))

	return slicePtr.Interface(), nil
}
