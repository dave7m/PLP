package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func main() {
	x := []int{6, 3, 7, 3, 9, 4, 35, 6, 8, 4, 24, 8, 5}
	quicksort(x)
	fmt.Println(x)

	y := []string{"abc", "cbde", "aaaa", "zzz", "fses"}
	quicksort(y)
	fmt.Println(y)
}

// generics are only supported since Go 1.18, so make sure to have downloaded the recent version
// function arguments as in https://golangexample.com/generic-sort-for-slices-in-golang/
// constraints.ordered contains types that are comparable with <, <=, ==, !=, >= and >.
func quicksort[E constraints.Ordered](list []E) {
	l := 0
	r := len(list) - 1
	quicksortHelper(list, l, r)
}

// this quicksort algorithm uses dual pivot, meaning in the partitioning we create two pivots instead of one
// the algorithm is semi efficient, because we do not gain very much from dual pivot, and the pivots are not chosen randomly, so the worst case does not improve (O(n^2)).
func quicksortHelper[E constraints.Ordered](list []E, l int, r int) {
	if r-l <= 0 {
		return
	}
	if list[l] > list[r] {
		list[l], list[r] = list[r], list[l]
	}
	p, q := partition(list, l, r)
	quicksortHelper(list, l, p-1)
	quicksortHelper(list, p+1, q-1)
	quicksortHelper(list, q+1, r)
}

func partition[E constraints.Ordered](list []E, lo int, hi int) (int, int) {
	l := lo + 1
	m := lo + 1
	g := hi
	for m < g {
		if list[m] < list[lo] {
			list[l], list[m] = list[m], list[l]
			l++
			m++
		} else if list[m] >= list[hi] {
			g--
			list[m], list[g] = list[g], list[m]
		} else {
			m++
		}
	}
	l--
	list[lo], list[l] = list[l], list[lo]
	list[hi], list[m] = list[m], list[hi]
	return l, m
}
