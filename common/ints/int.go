package ints

import (
	"strconv"
	"strings"
)

const MaxInT = int(^uint(0) >> 1)
const MinInT = ^MaxInT

func Max(ns ...int) int {
	if len(ns) == 0 {
		return 0
	}
	max := MinInT
	for _, n := range ns {
		if n > max {
			max = n
		}
	}
	return max
}

func Min(ns ...int) int {
	if len(ns) == 0 {
		return 0
	}
	min := MaxInT
	for _, n := range ns {
		if n < min {
			min = n
		}
	}
	return min
}

func Join(ns []int, sep string) string {
	as := make([]string, 0)
	for _, n := range ns {
		as = append(as, strconv.Itoa(n))
	}
	return strings.Join(as, sep)
}
