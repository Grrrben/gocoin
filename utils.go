package main

import (
	"sort"
)

// A data structure to hold key/value pairs
type Pair struct {
	Key   interface{}
	Value int
}

// A slice of pairs that implements sort.Interface to sort by values
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// sortMap Sorts any key type basedd on the Int value.
func sortMap(presorted map[interface{}]int) PairList {
	sortedList := make(PairList, len(presorted))
	i := 0
	for k, v := range presorted {
		sortedList[i] = Pair{k, v}
		i++
	}

	// sorting it highest to lowest
	sort.Sort(sort.Reverse(sortedList))
	return sortedList
}
