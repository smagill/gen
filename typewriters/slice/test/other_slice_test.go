// Generated by: setup
// TypeWriter: slice
// Directive: +test on Other

// See http://clipperhouse.github.io/gen for documentation

package main

import (
	"errors"
	"sort"
)

// OtherSlice is a slice of type Other. Use it where you would use []Other.
type OtherSlice []Other

// Max returns the maximum value of OtherSlice. In the case of multiple items being equally maximal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Max
func (rcv OtherSlice) Max() (result Other, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Max of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v > result {
			result = v
		}
	}
	return
}

// Min returns the minimum value of OtherSlice. In the case of multiple items being equally minimal, the first such element is returned. Returns error if no elements. See: http://clipperhouse.github.io/gen/#Min
func (rcv OtherSlice) Min() (result Other, err error) {
	l := len(rcv)
	if l == 0 {
		err = errors.New("cannot determine the Min of an empty slice")
		return
	}
	result = rcv[0]
	for _, v := range rcv {
		if v < result {
			result = v
		}
	}
	return
}

// Average sums OtherSlice over all elements and divides by len(OtherSlice). See: http://clipperhouse.github.io/gen/#Average
func (rcv OtherSlice) Average() (Other, error) {
	var result Other

	l := len(rcv)
	if l == 0 {
		return result, errors.New("cannot determine Average of zero-length OtherSlice")
	}
	for _, v := range rcv {
		result += v
	}
	result = result / Other(l)
	return result, nil
}

// Sum sums Other elements in OtherSlice. See: http://clipperhouse.github.io/gen/#Sum
func (rcv OtherSlice) Sum() (result Other) {
	for _, v := range rcv {
		result += v
	}
	return
}

// Sort returns a new ordered OtherSlice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv OtherSlice) Sort() OtherSlice {
	result := make(OtherSlice, len(rcv))
	copy(result, rcv)
	sort.Sort(result)
	return result
}

// IsSorted reports whether OtherSlice is sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv OtherSlice) IsSorted() bool {
	return sort.IsSorted(rcv)
}

// SortDesc returns a new reverse-ordered OtherSlice. See: http://clipperhouse.github.io/gen/#Sort
func (rcv OtherSlice) SortDesc() OtherSlice {
	result := make(OtherSlice, len(rcv))
	copy(result, rcv)
	sort.Sort(sort.Reverse(result))
	return result
}

// IsSortedDesc reports whether OtherSlice is reverse-sorted. See: http://clipperhouse.github.io/gen/#Sort
func (rcv OtherSlice) IsSortedDesc() bool {
	return sort.IsSorted(sort.Reverse(rcv))
}

func (rcv OtherSlice) Len() int {
	return len(rcv)
}
func (rcv OtherSlice) Less(i, j int) bool {
	return rcv[i] < rcv[j]
}
func (rcv OtherSlice) Swap(i, j int) {
	rcv[i], rcv[j] = rcv[j], rcv[i]
}
