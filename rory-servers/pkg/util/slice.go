package util

import (
	"reflect"
)

// InArray ?
func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

// IntArrayRemoveIndex ?
func IntArrayRemoveIndex(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

// IntArrayShuffle ?
func IntArrayShuffle(slice []int) []int {
	arr := make(map[int]int)

	for idx, c := range slice {
		arr[idx] = c
	}

	result := make([]int, 0)

	for _, c := range arr {
		result = append(result, c)
	}
	return result
}
