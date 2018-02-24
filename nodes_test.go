package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInitNodes(t *testing.T) {
	cls = initNodes()

	typeof := fmt.Sprint(reflect.TypeOf(cls))

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
	cls.addNode(testNode)

	if len(cls.List) != 1 {
		t.Error("Added one node, list should have length 1.")
	}

	cl := cls.List[0]

	if cl.Name != "node_x" {
		t.Error("Wrong node in cls in test TestAddNode")
	}
}

func TestNum(t *testing.T) {
	if cls.num() != 1 {
		t.Errorf("Expected 1 node, got %d.", cls.num())
	}

	secondNode := Node{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8001,
		Name:     "node2",
	}
	cls.addNode(secondNode)

	if cls.num() != 2 {
		t.Errorf("Expected 2 nodes, got %d.", cls.num())
	}
}
