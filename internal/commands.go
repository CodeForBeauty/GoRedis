package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type DBServer struct {
	db DB
}

func MakeServer() *DBServer {
	return &DBServer{db: *MakeDB()}
}

func (s *DBServer) ProcessCommand(command string) (string, error) {
	args := strings.Split(command, " ")

	if len(args) < 2 {
		return "", errors.New("Not enough arguments for a command")
	}

	var commands = map[string]func([]string) (string, error){
		"GET":    s.Get,
		"SET":    s.Set,
		"LPUSH":  s.AppendList,
		"LREM":   s.RemoveFromList,
		"LRANGE": s.RangeList,
		"HGET":   s.GetHash,
		"HSET":   s.SetHash,
	}

	comm := args[0]

	cmd, found := commands[comm]
	if !found {
		return "", errors.New("Command not found")
	}

	return cmd(args[1:])
}

func getValue[T Value](db DB, key string) (T, error) {
	tmpVal, found := db.Get(key)
	var zero T
	if !found {
		return zero, errors.New("Key doesn't exist")
	}
	val, ok := tmpVal.(T)
	if !ok {
		return zero, errors.New("Type cast failed")
	}
	return val, nil
}

func (s *DBServer) Get(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("Not enough arguments")
	}
	key := args[0]
	val, err := getValue[*StringValue](s.db, key)
	if err != nil {
		return "", err
	}
	return val.Data, nil
}

func (s *DBServer) Set(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("Not enough arguments")
	}
	key, value, exp := args[0], args[1], args[2]
	expiration, err := strconv.Atoi(exp)
	if err != nil {
		return "", err
	}
	tmp := &StringValue{Data: value}
	s.db.Set(key, tmp, time.Now().Add(time.Duration(expiration)*time.Second))
	return "", nil
}

func (s *DBServer) AppendList(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("Not enough arguments")
	}
	key, value := args[0], args[1]
	val, err := getValue[*ListValue](s.db, key)
	if err != nil {
		return "", err
	}
	val.Data = append(val.Data, value)
	return "", nil
}

func (s *DBServer) RemoveFromList(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("Not enough arguments")
	}
	key, value := args[0], args[1]
	val, err := getValue[*ListValue](s.db, key)
	if err != nil {
		return "", err
	}
	var idx int = -1
	for i := range val.Data {
		if val.Data[i] == value {
			idx = i
			break
		}
	}
	if idx == -1 {
		return "", errors.New("Value not found in list")
	}
	val.Data = append(val.Data[:idx], val.Data[idx+1:]...)
	return "", nil
}

func (s *DBServer) RangeList(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("Not enough arguments")
	}
	key, begin, end := args[0], args[1], args[2]

	start, err := strconv.Atoi(begin)
	if err != nil {
		return "", err
	}
	stop, err := strconv.Atoi(end)
	if err != nil {
		return "", err
	}

	val, err := getValue[*ListValue](s.db, key)
	if err != nil {
		return "", err
	}

	if start < 0 || start >= len(val.Data) || stop < start || stop >= len(val.Data) {
		return "", errors.New("Out of bounds index")
	}

	tmpSlice := val.Data[start:stop]
	return "$ " + strings.Join(tmpSlice, ","), nil
}

func (s *DBServer) SetHash(args []string) (string, error) {
	if len(args) < 3 {
		return "", errors.New("Not enough arguments")
	}
	key, hash, value := args[0], args[1], args[2]

	val, err := getValue[*HashValue](s.db, key)
	if err != nil {
		return "", err
	}
	val.Data[hash] = value
	return "", nil
}

func (s *DBServer) GetHash(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("Not enough arguments")
	}
	key, hash := args[0], args[1]

	val, err := getValue[*HashValue](s.db, key)
	if err != nil {
		return "", err
	}
	return val.Data[hash], nil
}
