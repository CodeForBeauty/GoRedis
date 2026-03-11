# Redis clone in GO

A simple Redis clone made with Go.

## Features

Multitype data storage

Lazy expiration

Memory cap with LRU eviction

Multithreaded connection manager

## Data types

String

List

Hashmap

Hashset

## How to use

The server listens for tcp connections on port 8080. After connecting to the server following commands can be ran.

```
GET [key] // Get string value
SET [key] [value] [expire seconds] // Set string value
LPUSH [key] [value] // push to list, if doesn't exist create one
LREM [key] [value] // remove first matching value from list
LRANGE [key] [start] [end] // get elements from list in range
HGET [key] [hash] // get value of the hash
HSET [key] [hash] [value] // sets value of the hash, creates hashmap if doesn't exist
EXPIR [key] [expire seconds] // sets expiration time from current time
SPUSH [key] [value] // push to set, of doesn't exist create one
SREM [key] [value] // remove value from set
SHAS [key] [value] // check if set has value
```

Each command returns some response.

Usage examples with expected results.

```
SET user1 50 500
OK

GET user1
50

LPUSH responses 10
OK

LREM responses 10
OK

LRANGE responses 0 3
$ 10,5,15

HSET comments John Hi
OK

HGET comments John
Hi

EXPIR comments 50
OK

SPUSH users John
OK

SREM users John
OK

SHAS users John
0
```

## Prerequisites

go 1.25.0+

git

## How to run

Open terminal

Clone repo

```bash
git clone https://github.com/CodeForBeauty/GoRedis
```

Go into cloned directory

```bash
cd GoRedis
```

Run server with go

```bash
go run ./cmd/main.go
```

### Tests

To run unit tests run

```bash
go test -v ./tests
```

## Limitations

No persistens

White line characters aren't supported

## Conclusion

This was a great learning experience. I've learned to design message parser, distribute tasks across threads and to manage memory usage.