package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	KEY_NOT_FOUND_ERROR        = errors.New("Key doesn't exist")
	TYPE_CAST_ERROR            = errors.New("Type cast failed")
	WRONG_ARGUMENT_COUNT_ERROR = errors.New("Not enough arguments")
	COMMAND_NOT_FOUND_ERROR    = errors.New("Command not found")
	LIST_NOT_FOUND_ERROR       = errors.New("Value not found in list")
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
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}

	var commands = map[string]func([]string) (string, error){
		"GET":    s.Get,
		"SET":    s.Set,
		"LPUSH":  s.AppendList,
		"LREM":   s.RemoveFromList,
		"LRANGE": s.RangeList,
		"HGET":   s.GetHash,
		"HSET":   s.SetHash,
		"EXPIR":  s.ChangeExpiration,
	}

	comm := args[0]

	cmd, found := commands[comm]
	if !found {
		return "", COMMAND_NOT_FOUND_ERROR
	}

	return cmd(args[1:])
}

func getValue[T Value](db DB, key string) (T, error) {
	tmpVal, found := db.Get(key)
	var zero T
	if !found {
		return zero, KEY_NOT_FOUND_ERROR
	}
	val, ok := tmpVal.(T)
	if !ok {
		return zero, TYPE_CAST_ERROR
	}
	return val, nil
}

func (s *DBServer) Get(args []string) (string, error) {
	if len(args) < 1 {
		return "", WRONG_ARGUMENT_COUNT_ERROR
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
		return "", WRONG_ARGUMENT_COUNT_ERROR
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
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}
	key, value := args[0], args[1]
	if _, ok := s.db.data[key]; !ok {
		s.db.Set(key, &ListValue{}, time.Now().Add(time.Duration(30)*time.Minute))
	}
	val, err := getValue[*ListValue](s.db, key)
	if err != nil {
		return "", err
	}
	val.Data = append(val.Data, value)
	return "", nil
}

func (s *DBServer) RemoveFromList(args []string) (string, error) {
	if len(args) < 2 {
		return "", WRONG_ARGUMENT_COUNT_ERROR
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
		return "", LIST_NOT_FOUND_ERROR
	}
	val.Data = append(val.Data[:idx], val.Data[idx+1:]...)
	return "", nil
}

func (s *DBServer) RangeList(args []string) (string, error) {
	if len(args) < 3 {
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}
	key, begin, end := args[0], args[1], args[2]

	val, err := getValue[*ListValue](s.db, key)
	if err != nil {
		return "", err
	}

	start, err := strconv.Atoi(begin)
	if err != nil {
		return "", err
	}
	stop, err := strconv.Atoi(end)
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
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}
	key, hash, value := args[0], args[1], args[2]

	if _, ok := s.db.data[key]; !ok {
		s.db.Set(key, &HashValue{}, time.Now().Add(time.Duration(30)*time.Minute))
	}
	val, err := getValue[*HashValue](s.db, key)
	if err != nil {
		return "", err
	}
	val.Data[hash] = value
	return "", nil
}

func (s *DBServer) GetHash(args []string) (string, error) {
	if len(args) < 2 {
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}
	key, hash := args[0], args[1]

	val, err := getValue[*HashValue](s.db, key)
	if err != nil {
		return "", err
	}
	return val.Data[hash], nil
}

func (s *DBServer) ChangeExpiration(args []string) (string, error) {
	if len(args) < 2 {
		return "", WRONG_ARGUMENT_COUNT_ERROR
	}
	key, exp := args[0], args[1]

	expiration, err := strconv.Atoi(exp)
	if err != nil {
		return "", err
	}

	entry, found := s.db.data[key]
	if !found {
		return "", KEY_NOT_FOUND_ERROR
	}

	entry.expiration = time.Now().Add(time.Duration(expiration) * time.Minute)

	return "", nil
}
