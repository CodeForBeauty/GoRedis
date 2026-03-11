package internal

import "time"

type ValueType int

const (
	TYPE_STRING ValueType = iota
	TYPE_LIST
	TYPE_HASH
)

type Value interface {
	GetType() ValueType
}

type Entry struct {
	value      Value
	expiration time.Time
}

type DB struct {
	data map[string]*Entry
}

func MakeDB() *DB {
	return &DB{data: map[string]*Entry{}}
}

func (d *DB) Set(key string, value Value, expiration time.Time) {
	d.data[key] = &Entry{value: value, expiration: expiration}
}

func (d *DB) Get(key string) (Value, bool) {
	tmp, found := d.data[key]
	if found {
		return tmp.value, found
	}
	return nil, found
}

func (d *DB) Remove(key string) {
	delete(d.data, key)
}

type StringValue struct {
	Data string
}

func (s *StringValue) GetType() ValueType {
	return TYPE_STRING
}

type ListValue struct {
	Data []string
}

func (l *ListValue) GetType() ValueType {
	return TYPE_LIST
}

type HashValue struct {
	Data map[string]string
}

func (h *HashValue) GetType() ValueType {
	return TYPE_HASH
}
