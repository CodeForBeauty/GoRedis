package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var db *DB = MakeDB()

func ProcessCommand(command string) error {
	args := strings.Split(command, " ")

	if len(args) < 2 {
		return errors.New("Not enough arguments for command")
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
		Get(key)
	case "SET":
		Set(key, secondArg, 10000)
	case "LAPP":
		AppendList(key, secondArg)
	case "LREM":
		RemoveFromList(key, secondArg)
	case "LRAN":
		start, err := strconv.Atoi(secondArg)
		if err != nil {
			return err
		}
		var end int
		end, err = strconv.Atoi(thirdArg)
		if err != nil {
			return err
		}
		RangeList(key, start, end)
	case "HSET":
		SetHash(key, secondArg, thirdArg)
	case "HGET":
		GetHash(key, secondArg)
	}

	return nil
}

func Get(key string) (string, error) {
	tmpVal, found := db.Get(key)
	if !found {
		return "", errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_STRING {
		return "", errors.New("Key is not a string")
	}
	val, ok := (*tmpVal).(*StringValue)
	if !ok {
		return "", errors.New("Failed to cast Value to ListValue")
	}
	return val.Data, nil
}

func Set(key string, value string, expiration int) error {
	tmp := &StringValue{Data: value}
	db.Set(key, tmp, time.Now().Add(time.Duration(expiration)*time.Millisecond))
	return nil
}

func AppendList(key string, value string) error {
	tmpVal, found := db.Get(key)
	if !found {
		return errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_LIST {
		return errors.New("Wrong entry type")
	}
	val, ok := (*tmpVal).(*ListValue)
	if !ok {
		return errors.New("Failed to cast Value to ListValue")
	}
	val.Data = append(val.Data, value)
	return nil
}

func RemoveFromList(key string, value string) error {
	tmpVal, found := db.Get(key)
	if !found {
		return errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_LIST {
		return errors.New("Wrong entry type")
	}
	val, ok := (*tmpVal).(*ListValue)
	if !ok {
		return errors.New("Failed to cast Value to ListValue")
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
	tmpVal, found := db.Get(key)
	if !found {
		return nil, errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_LIST {
		return nil, errors.New("Wrong entry type")
	}
	val, ok := (*tmpVal).(*ListValue)
	if !ok {
		return nil, errors.New("Failed to cast Value to ListValue")
	}
	return val.Data, nil
}

func SetHash(key string, hash string, value string) error {
	tmpVal, found := db.Get(key)
	if !found {
		return errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_HASH {
		return errors.New("Wrong entry type")
	}
	val, ok := (*tmpVal).(*HashValue)
	if !ok {
		return errors.New("Failed to cast Value to ListValue")
	}
	val.Data[hash] = value
	return nil
}

func GetHash(key string, hash string) (string, error) {
	tmpVal, found := db.Get(key)
	if !found {
		return "", errors.New("Key doesn't exist")
	}
	if (*tmpVal).GetType() != TYPE_HASH {
		return "", errors.New("Wrong entry type")
	}
	val, ok := (*tmpVal).(*HashValue)
	if !ok {
		return "", errors.New("Failed to cast Value to ListValue")
	}
	return val.Data[hash], nil
}
