package tests

import (
	"testing"
	"time"

	"github.com/CodeForBeauty/GoRedis/internal"
)

func TestSetGetStringValue(t *testing.T) {
	var db = internal.MakeDB(128)

	key := "test"
	value := "Test"

	db.Set(key, &internal.StringValue{Data: value}, time.Now().Add(time.Duration(5)*time.Second))

	val, found := db.Get(key)

	if !found {
		t.Error("Key not added")
	}

	strVal, ok := val.(*internal.StringValue)

	if !ok {
		t.Error("Type cast filed")
	}

	if strVal.Data != value {
		t.Error("Incorrect value")
	}
}

func TestKeyExpiration(t *testing.T) {
	var db = internal.MakeDB(128)

	key := "test"
	value := "Test"

	duration := time.Duration(1) * time.Second

	db.Set(key, &internal.StringValue{Data: value}, time.Now().Add(duration))

	time.Sleep(duration)

	_, found := db.Get(key)

	if found {
		t.Error("Key not removed added")
	}
}

func TestLRUEviction(t *testing.T) {
	var db = internal.MakeDB(18)

	testValue := "Test"

	firstKey := "test"

	db.Set(firstKey, &internal.StringValue{Data: testValue}, time.Now().Add(time.Duration(5)*time.Second))

	db.Set("key1", &internal.StringValue{Data: testValue}, time.Now().Add(time.Duration(5)*time.Second))
	db.Set("key2", &internal.StringValue{Data: testValue}, time.Now().Add(time.Duration(5)*time.Second))
	db.Set("key3", &internal.StringValue{Data: testValue}, time.Now().Add(time.Duration(5)*time.Second))

	db.Set("key4", &internal.StringValue{Data: testValue}, time.Now().Add(time.Duration(5)*time.Second))

	if db.Len() != 4 {
		t.Error("Key not removed or added")
	}

	_, found := db.Get(firstKey)

	if found {
		t.Error("Wrong key removed")
	}
}
