package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInitNodes(t *testing.T) {
	nodes = initNodes()

	typeof := fmt.Sprint(reflect.TypeOf(nodes))

	if typeof != "*main.Nodes" {
		t.Errorf("Wrong type, expected *main.Nodes, got %s", typeof)
	}
}

func TestAddNode(t *testing.T) {
	testNode := Node{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8000,
		Name:     "node_x",
	}
	nodes.addNode(&testNode)

	if len(nodes.List) != 1 {
		t.Error("Added one node, list should have length 1.")
	}

	cl := nodes.List[0]

	if cl.Name != "node_x" {
		t.Error("Wrong node in cls in test TestAddNode")
	}

	if len(cl.Hash) != 64 {
		t.Error("Node has no valid hash in test TestAddNode")
	}
}

func TestNum(t *testing.T) {
	if nodes.num() != 1 {
		t.Errorf("Expected 1 node, got %d.", nodes.num())
	}

	secondNode := Node{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8001,
		Name:     "node2",
	}
	nodes.addNode(&secondNode)

	if nodes.num() != 2 {
		t.Errorf("Expected 2 nodes, got %d.", nodes.num())
	}
}
