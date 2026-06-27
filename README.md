# RoxKV

RoxKV is a lightweight networked in-memory key-value database and message broker written in Go.

The project started as an attempt to understand how systems like Redis work internally and gradually evolved into a database server featuring a custom storage engine, persistence, TTL expiration, TCP networking, background workers, and a Publish/Subscribe messaging system.

> **Note:** RoxKV is an educational project built from scratch to explore systems programming concepts and database internals.

---

# Features

## Storage Engine

* In-memory key-value storage
* Thread-safe operations using `sync.RWMutex`
* Shared storage across multiple TCP clients
* Automatic key overwrite on duplicate inserts

---

## Persistence

* Snapshot-based persistence
* Daily snapshot files
* Load snapshots from oldest to newest
* Latest values automatically overwrite stale data
* Persistent keys survive server restarts

---

## TTL (Time To Live)

* Optional key expiration
* Background expiry worker
* Automatic deletion of expired keys
* Supports:

```bash
set session abc123 --ttl 30 sec
set token xyz --ttl 15 min
set cache data --ttl 1 hr
```

---

## TCP Networking

* Raw TCP server
* Concurrent client handling
* Shared database across all clients
* Configurable maximum client connections
* Automatic inactive client cleanup

Connect using:

```bash
nc localhost 6969
```

---

## Publish / Subscribe

RoxKV includes a lightweight Pub/Sub broker.

Supported operations:

```bash
Subscribe <topic>

Publish <topic> <message>

Unsubscribe <topic>

Topics

CloseChannel <topic>
```

Features:

* Multiple subscribers per topic
* Instant message delivery
* Dynamic topic creation
* Topic removal
* Concurrent message broadcasting

---

## Logging & History

Three logging levels are supported.

```bash
history

history --success

history --info

history --error
```

Useful for debugging, monitoring, and auditing server activity.

---

## Custom Command Parser

The parser was implemented manually without relying on shell parsing.

Supports:

```bash
set game Battlefield4

set game "Battlefield 4"

set quote 'Never Stop Learning'

set desc `Testing on Build`
```

Quoted values preserve spaces while unquoted values continue to work normally.

---

# Supported Commands

## Database

```bash
set <key> <value>

set <key> <value> --ttl <time> <sec|min|hr>

get <key>

del <key>

keys

save

load
```

---

## History

```bash
history

history --success

history --info

history --error
```

---

## Pub/Sub

```bash
Subscribe <topic>

Publish <topic> <message>

Unsubscribe <topic>

Topics

CloseChannel <topic>
```

---

# Architecture

```text
                    TCP Clients
                         │
         ┌───────────────┼───────────────┐
         │               │               │
      Client A        Client B        Client C
         │               │               │
         └───────────────┼───────────────┘
                         │
                  TCP Server
                         │
                  Command Parser
                         │
         ┌───────────────┼───────────────┐
         │               │               │
   Storage Engine   Pub/Sub Broker   Logger
         │
 ┌───────┼──────────────┐
 │       │              │
TTL   Persistence   Workers
```

---

# Technologies Used

* Go
* Goroutines
* Channels
* sync.RWMutex
* Context
* Raw TCP Networking
* File-based Persistence

---

# Project Goals

The purpose of RoxKV is educational.

The project is used to understand:

* Storage engine implementation
* TCP networking
* Concurrent programming
* Background workers
* Mutexes and synchronization
* Channels
* Publish/Subscribe systems
* Database persistence
* Message brokers
* Custom protocol design

---

# Roadmap

Planned improvements:

* RESP (Redis Serialization Protocol) support
* Redis CLI compatibility
* Publisher permissions & access control
* Better persistence engine
* Graceful shutdown
* Automatic periodic snapshots
* Configuration file support
* Benchmarking suite
* Unit and integration tests

---

# Build (Linux)

Download the latest executable from the Releases section.

Install:

```bash
sudo cp roxkv /usr/local/bin/
```

Run:

```bash
roxkv
```

---

# Author

**Rahul Kumar Parida**

GitHub:

https://github.com/rahulkumarparida

---

*"Built to understand how databases, message brokers, and network servers work under the hood."*
