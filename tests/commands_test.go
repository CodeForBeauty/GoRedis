package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/CodeForBeauty/GoRedis/internal"
)

func TestStrings(t *testing.T) {
	var dbServer = internal.MakeServer()

	key := "test"
	value := "50"

	_, err := dbServer.ProcessCommand(fmt.Sprintf("SET %s %s 5", key, value))

	if err != nil {
		t.Errorf("Error setting key %e", err)
	}

	output, err := dbServer.ProcessCommand(fmt.Sprintf("GET %s", key))

	if err != nil {
		t.Errorf("Error getting value %e", err)
	}

	if output != value {
		t.Errorf("Incorrect value saved. Expected: %s, got: %s", value, output)
	}
}

func TestLists(t *testing.T) {
	var dbServer = internal.MakeServer()

	key := "test"

	value1 := "1"
	value2 := "2"
	value3 := "3"

	_, err := dbServer.ProcessCommand(fmt.Sprintf("LPUSH %s %s", key, value1))
	if err != nil {
		t.Error("Failed to push to list")
	}
	_, err = dbServer.ProcessCommand(fmt.Sprintf("LPUSH %s %s", key, value2))
	if err != nil {
		t.Error("Failed to push to list")
	}
	_, err = dbServer.ProcessCommand(fmt.Sprintf("LPUSH %s %s", key, value3))
	if err != nil {
		t.Error("Failed to push to list")
	}

	output, err := dbServer.ProcessCommand(fmt.Sprintf("LRANGE %s %d %d", key, 0, 3))
	if err != nil {
		t.Error("Failed to get range: ", err)
	}

	if output != ("$ " + strings.Join([]string{value1, value2, value3}, ",")) {
		t.Errorf("Wrong number of elements")
	}

	_, err = dbServer.ProcessCommand(fmt.Sprintf("LREM %s %s", key, value2))
	if err != nil {
		t.Error("Failed to remove from list: ", err)
	}

	output, err = dbServer.ProcessCommand(fmt.Sprintf("LRANGE %s %d %d", key, 0, 2))
	if err != nil {
		t.Error("Failed to get range: ", err)
	}

	if output != ("$ " + strings.Join([]string{value1, value3}, ",")) {
		t.Errorf("Wrong number of elements")
	}
}

func TestHashes(t *testing.T) {
	var dbServer = internal.MakeServer()

	mainKey := "test"

	key := "test"
	value := "50"

	_, err := dbServer.ProcessCommand(fmt.Sprintf("HSET %s %s %s", mainKey, key, value))
	if err != nil {
		t.Error("Failed to set hash: ", err)
	}

	output, err := dbServer.ProcessCommand(fmt.Sprintf("HGET %s %s", mainKey, key))
	if err != nil {
		t.Error("Failed to get hash: ", err)
	}

	if output != value {
		t.Error("Wrong hash value")
	}
}

func TestExpiration(t *testing.T) {
	var dbServer = internal.MakeServer()

	key := "test"
	value := "50"

	_, err := dbServer.ProcessCommand(fmt.Sprintf("SET %s %s 1", key, value))
	if err != nil {
		t.Errorf("Error setting key %e", err)
	}

	_, err = dbServer.ProcessCommand(fmt.Sprintf("EXPIR %s 2", key))
	if err != nil {
		t.Errorf("Error setting expiration %e", err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	_, err = dbServer.ProcessCommand(fmt.Sprintf("GET %s", key))
	if err != nil {
		t.Errorf("Error getting value %e", err)
	}

	time.Sleep(time.Duration(1) * time.Second)

	output, err := dbServer.ProcessCommand(fmt.Sprintf("GET %s", key))
	if err != nil {
		t.Errorf("Error getting value %e", err)
	}

	if output != value {
		t.Error("Wrong value")
	}
}
