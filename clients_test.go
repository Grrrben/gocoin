package main

import (
	"testing"
	"reflect"
	"fmt"
)

func TestInitClients(t *testing.T) {
	cls = initClients()


	typeof := fmt.Sprint(reflect.TypeOf(cls))


	if typeof != "*main.Clients" {
		t.Errorf("Wrong type, expected *main.Clients, got %s", typeof)
	}
}

func TestAddClient(t *testing.T) {
	testClient := Client{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8000,
		Name:     "client1",
	}
	cls.addClient(testClient)

	if len(cls.List) != 1 {
		t.Error("Added one client, list should have length 1.")
	}

	cl := cls.List[0]

	if cl.Name != "client1" {
		t.Error("Wrong Client in cls in test TestAddClient")
	}
}

func TestNum (t *testing.T) {
	if cls.num() != 1 {
		t.Errorf("Expected 1 client, got %d.", cls.num())
	}

	secondClient := Client{
		Protocol: "http://",
		Hostname: "127.0.0.1",
		Port:     8001,
		Name:     "client2",
	}

	cls.addClient(secondClient)

	if cls.num() != 2 {
		t.Errorf("Expected 2 clients, got %d.", cls.num())
	}
}
