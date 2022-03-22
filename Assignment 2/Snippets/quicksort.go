package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

func main() {
	a := []int{6, 3, -6, -45, 7, 3, 9, 4, 35, 6, 8, 4, 24, 8, 5}
	fmt.Println("before sorting: ", a)
	quicksort(a)
	fmt.Println("after sorting: ", a)

	y := []string{"Foo", "foo", "foobar", "FizzBuzz", "abc", "cbde", "aaaa", "zzz", "fses"}
	fmt.Println("before sorting: ", y)
	quicksort(y)
	fmt.Println("after sorting: ", y)

	z := []float64{10.45, 3.141, -49, 25.24, 924.1, 4.5, 6.2, 9.5, -3.5}
	fmt.Println("before sorting: ", z)
	quicksort(z)
	fmt.Println("after sorting: ", z)

	x := []uint{6, 3, 7, 3, 9, 4, 35, 6, 8, 4, 24, 8, 5}
	fmt.Println("before sorting: ", x)
	quicksort(x)
	fmt.Println("after sorting: ", x)
}

// generics are only supported since Go 1.18, so make sure to have downloaded the recent version
// function arguments as in https://golangexample.com/generic-sort-for-slices-in-golang/
// constraints.ordered contains types that are comparable with <, <=, ==, !=, >= and >.
// see https://pkg.go.dev/golang.org/x/exp/constraints
func quicksort[E constraints.Ordered](list []E) {
	l := 0
	r := len(list) - 1
	quicksortHelper(list, l, r)
}

// this quicksort algorithm uses dual pivot, meaning in the partitioning we create two pivots instead of one
// the algorithm is semi efficient, because we do not gain very much from dual pivot, and the pivots are not chosen randomly,
// so the worst case does not improve (O(n^2)). But I wanted some practice :)
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

// partitions the list from lo to hi and returns two pivots
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
