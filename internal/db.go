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
	GetSize() int
}

type Entry struct {
	key        string
	value      Value
	expiration time.Time
	next       *Entry
	prev       *Entry
}

type DB struct {
	data     map[string]*Entry
	dataSize int
	maxSize  int
	head     *Entry
	tail     *Entry
}

func MakeDB(maxMemory int) *DB {
	return &DB{data: map[string]*Entry{}, maxSize: maxMemory, head: nil, tail: nil}
}

func (d *DB) Set(key string, value Value, expiration time.Time) {
	val, found := d.data[key]
	size := value.GetSize()
	if found {
		size -= val.value.GetSize()
		if val.next != nil {
			val.next.prev = val.prev
		}
		if val.prev != nil {
			val.prev.next = val.next
		} else {
			d.tail = val.next
		}
	}
	for d.dataSize+size > d.maxSize {
		d.dataSize -= d.tail.value.GetSize()
		tail := d.tail
		d.tail = tail.next
		d.Remove(tail.key)
	}

	newEntry := &Entry{key: key, value: value, expiration: expiration}
	newEntry.prev = d.head
	d.head = newEntry
	d.dataSize += size
	if d.tail == nil {
		d.tail = newEntry
	}

	d.data[key] = newEntry
}

func (d *DB) Get(key string) (Value, bool) {
	tmp, found := d.data[key]
	if found {
		if time.Until(tmp.expiration) <= 0 {
			delete(d.data, key)
			return nil, false
		}
		if tmp.next != nil {
			tmp.next.prev = tmp.prev
		}
		if tmp.prev != nil {
			tmp.prev.next = tmp.next
		} else {
			d.tail = tmp.next
		}
		d.head = tmp
		return tmp.value, found
	}
	return nil, found
}

func (d *DB) Remove(key string) {
	delete(d.data, key)
}

func (d *DB) Len() int {
	return len(d.data)
}

type StringValue struct {
	Data string
}

func (s *StringValue) GetType() ValueType {
	return TYPE_STRING
}

func (s *StringValue) GetSize() int {
	return len(s.Data)
}

type ListValue struct {
	Data []string
}

func (l *ListValue) GetType() ValueType {
	return TYPE_LIST
}

func (l *ListValue) GetSize() int {
	var size int = 0
	for i := range l.Data {
		size += len(l.Data[i])
	}
	return size
}

type HashValue struct {
	Data map[string]string
}

func (h *HashValue) GetType() ValueType {
	return TYPE_HASH
}

func (h *HashValue) GetSize() int {
	var size int = 0
	for i := range h.Data {
		size += len(h.Data[i])
	}
	return size
}
