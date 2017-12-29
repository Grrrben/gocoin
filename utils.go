package main

import (
	"sort"
	"regexp"
	"log"
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

// sortMapDescending Sorts any key type high to low, based on the Int value.
func sortMapDescending(presorted map[interface{}]int) PairList {
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

// sortMapAscending Sorts any key type low to high, based on the Int value.
func sortMapAscending(presorted map[interface{}]int) PairList {
	sortedList := make(PairList, len(presorted))
	i := 0
	for k, v := range presorted {
		sortedList[i] = Pair{k, v}
		i++
	}

	sort.Sort(sortedList)
	return sortedList
}

// hasValidHash checks a hash for length and regex of the hex.
// It does _not_ check for existince of a specific hash.
func hasValidHash(object hashable) bool {
	return validHash(object.getHash())
}

// validHash checks a hash for length and regex of the hex.
// It does _not_ check for existince of a wallet with this specific hash.
func validHash(hash string) bool {
	// fad5e7a92f1c43b1523614336a07f98b894bb80fee06b6763b50ab03b597d5f4
	regex, err := regexp.Compile(`[a-f0-9]{64}`)

	if err != nil {
		log.Fatal("Could not compile regex")
	}
	if regex.MatchString(hash) {
		return true
	} else {
		return false
	}
}
