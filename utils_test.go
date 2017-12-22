package main

import "testing"

func TestSortMapDescending(t *testing.T) {

	list := make(map[interface{}]int, 4)
	list["a"] = 1
	list["b"] = 2
	list["c"] = 3
	list["d"] = 4

	sorted := sortMapDescending(list)
	first := sorted[0]

	if first.Key != "d" {
		t.Errorf("SortMapDescending test failed. Expected 'd', got %s", first.Key)
	}
}

func TestSortMapAscending(t *testing.T) {

	list := make(map[interface{}]int, 4)
	list["a"] = 8
	list["b"] = 7
	list["c"] = 6
	list["d"] = 5

	sorted := sortMapAscending(list)
	first := sorted[0]

	if first.Key != "d" {
		t.Errorf("SortMapAscending test failed. Expected 'd', got %s", first.Key)
	}
}
