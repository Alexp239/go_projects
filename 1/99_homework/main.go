package main

import (
	"sort"
	"strconv"
)

func ReturnInt() int {
	return 1
}

func ReturnFloat() float32 {
	return 1.1
}

func ReturnIntArray() [3]int {
	a := [...]int{1, 3, 4}
	return a
}

func ReturnIntSlice() []int {
	sl := []int{1, 2, 3}
	return sl
}

func IntSliceToString(sl []int) string {
	s := ""
	for _, val := range sl {
		s += strconv.Itoa(val)
	}
	return s
}

func MergeSlices(sl1 []float32, sl2 []int32) []int {
	var sl []int
	for _, val := range sl1 {
		sl = append(sl, int(val))
	}

	for _, val := range sl2 {
		sl = append(sl, int(val))
	}
	return sl
}

func GetMapValuesSortedByKey(mp map[int]string) []string {
	var res []string
	var keys []int
	for key := range mp {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, key := range keys {
		res = append(res, mp[key])
	}
	return res
}
