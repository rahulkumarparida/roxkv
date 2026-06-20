# RoxKV

RoxKV is a lightweight networked key-value database written in Go.

The project started as an attempt to understand how systems like Redis work internally and gradually evolved into a database server featuring custom command parsing, persistence, TTL support, background workers, logging, and TCP-based client communication.

## Features

### Storage Engine

* Store key-value pairs in memory
* Thread-safe operations using `sync.RWMutex`
* Shared memory across multiple TCP clients

### Persistence

* Save data to disk using snapshot files
* Date-based storage files
* Load snapshots from oldest to newest
* Latest values automatically overwrite stale entries

### TTL (Time To Live)

* Optional expiration support
* Background expiry worker
* Automatic cleanup of expired keys

Example:

```bash
set session abc123 --ttl 30 min
```

### Networking

* Raw TCP server
* Multi-client support
* Concurrent client handling using goroutines

Connect using:

```bash
nc localhost 6969
```

### Logging & History

* Success logs
* Error logs
* Informational logs

History commands:

```bash
history
history --success
history --error
history --info
```

### Custom Command Parser

Supports:

```bash
set name Rahul
set game "Battlefield 4"
set quote 'Never Stop Learning'
set desc `Testing on Build`
```

The parser was implemented manually to support quoted and unquoted values while preserving spaces.

## Supported Commands

```bash
get <key>

set <key> <value>

set <key> <value> --ttl <time> <sec|min|hr>

del <key>

keys

save

load

history

history --success

history --info

history --error
```

## Architecture

```text
TCP Clients
      │
      ▼
 TCP Server
      │
      ▼
 Command Parser
      │
      ▼
 Storage Engine
      │
 ┌────┼────┐
 ▼    ▼    ▼
TTL  Logs Persistence
 │
 ▼
Worker
```

## Project Goals

The purpose of RoxKV is educational:

* Learn storage engine design
* Understand persistence strategies
* Explore networking with TCP
* Build custom parsers
* Work with concurrency and goroutines
* Understand how key-value databases operate internally

## Future Plans

* Pub/Sub messaging
* Message broker functionality
* RESP protocol support
* Improved persistence engine
* Graceful shutdown handling
* Automated snapshots

## Build Linux Only

#### Download only the roxkv executables from the repo , No need to clone the entire repo

```bash
cd <to the downloads>

sudo cp roxkv /usr/local/bin

roxkv

```

## Author

Rahul Kumar Parida

GitHub:
https://github.com/rahulkumarparida
