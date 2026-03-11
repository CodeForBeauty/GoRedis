package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var db *DB = MakeDB()

func ProcessCommand(command string) (string, error) {
	args := strings.Split(command, " ")

	if len(args) < 2 {
		return "", errors.New("Not enough arguments for command")
	}

	key := args[1]
	var secondArg string
	if len(args) > 2 {
		secondArg = args[2]
	}
	var thirdArg string
	if len(args) > 3 {
		thirdArg = args[3]
	}

	switch args[0] {
	case "GET":
		return Get(key)
	case "SET":
		Set(key, secondArg, 10000)
	case "LAPP":
		AppendList(key, secondArg)
	case "LREM":
		RemoveFromList(key, secondArg)
	case "LRAN":
		start, err := strconv.Atoi(secondArg)
		if err != nil {
			return "", err
		}
		var end int
		end, err = strconv.Atoi(thirdArg)
		if err != nil {
			return "", err
		}
		l, err := RangeList(key, start, end)
		if err != nil {
			return "", err
		}
		tmp := "L " + strings.Join(l, ",")
		return tmp, nil
	case "HSET":
		SetHash(key, secondArg, thirdArg)
	case "HGET":
		return GetHash(key, secondArg)
	}

	return "", nil
}

func getValue[T Value](key string) (T, error) {
	tmpVal, found := db.Get(key)
	var zero T
	if !found {
		return zero, errors.New("Key doesn't exist")
	}
	if tmpVal.GetType() != zero.GetType() {
		return zero, errors.New("Wrong entry type")
	}
	val, ok := tmpVal.(T)
	if !ok {
		return zero, errors.New("Failed to cast Value to ListValue")
	}
	return val, nil
}

func Get(key string) (string, error) {
	val, err := getValue[*StringValue](key)
	if err != nil {
		return "", err
	}
	return val.Data, nil
}

func Set(key string, value string, expiration int) error {
	tmp := &StringValue{Data: value}
	db.Set(key, tmp, time.Now().Add(time.Duration(expiration)*time.Millisecond))
	return nil
}

func AppendList(key string, value string) error {
	val, err := getValue[*ListValue](key)
	if err != nil {
		return err
	}
	val.Data = append(val.Data, value)
	return nil
}

func RemoveFromList(key string, value string) error {
	val, err := getValue[*ListValue](key)
	if err != nil {
		return err
	}
	var idx int = -1
	for i := range val.Data {
		if val.Data[i] == value {
			idx = i
			break
		}
	}
	val.Data = append(val.Data[idx:], val.Data[:idx+1]...)
	return nil
}

func RangeList(key string, start int, end int) ([]string, error) {
	val, err := getValue[*ListValue](key)
	if err != nil {
		return nil, err
	}
	return val.Data[start:end], nil
}

func SetHash(key string, hash string, value string) error {
	val, err := getValue[*HashValue](key)
	if err != nil {
		return err
	}
	val.Data[hash] = value
	return nil
}

func GetHash(key string, hash string) (string, error) {
	val, err := getValue[*HashValue](key)
	if err != nil {
		return "", err
	}
	return val.Data[hash], nil
}
