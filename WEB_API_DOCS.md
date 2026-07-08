# Web API SSE Endpoints

The server exposes a single Server-Sent Events endpoint in [server/WebServer.go](server/WebServer.go):

- GET /api/events/{name}

The `{name}` path parameter selects one of the supported payloads. The server emits a new JSON event every second over the SSE stream.

Each event body is a JSON object that corresponds to one of the handlers in [server/WebServer.go](server/WebServer.go).

## 1) Monitor endpoint

- Path: `/api/events/monitor`
- Handler: `HandleResponseData("monitor")`
- Returns: a JSON object representing machine usage, storage usage, network usage, and nested runtime metadata.

### Actual payload shape

```json
{
  "cpuusage": "12.5",
  "ramusage": {
    "totalRam": 16777216000,
    "freeRam": 8388608000,
    "usedPercentge": 50.1
  },
  "diskusage": {
    "total": 476,
    "free": 211,
    "avaliable": 211,
    "err": null
  },
  "networkstats": {
    "downloadspeed": 124.8,
    "uploadspeed": 32.4,
    "Name": "wlan0"
  },
  "machineinfo": {
    "goVersion": "go1.22.0",
    "os": "linux",
    "user": "my-hostname",
    "arch": "amd64",
    "cpus": "8",
    "goroutines": "42",
    "allocatedMemMB": 123,
    "totalAllocatedMB": 456,
    "systemMemMB": 789,
    "connectedclient": 2,
    "serverUptime": 123.45
  }
}
```

### Nested structures

#### `ramusage`

```json
{
  "totalRam": 16777216000,
  "freeRam": 8388608000,
  "usedPercentge": 50.1
}
```

#### `diskusage`

```json
{
  "total": 476,
  "free": 211,
  "avaliable": 211,
  "err": null
}
```

#### `networkstats`

```json
{
  "downloadspeed": 124.8,
  "uploadspeed": 32.4,
  "Name": "wlan0"
}
```

#### `machineinfo`

```json
{
  "goVersion": "go1.22.0",
  "os": "linux",
  "user": "my-hostname",
  "arch": "amd64",
  "cpus": "8",
  "goroutines": "42",
  "allocatedMemMB": 123,
  "totalAllocatedMB": 456,
  "systemMemMB": 789,
  "connectedclient": 2,
  "serverUptime": 123.45
}
```

### Notes

- This endpoint uses `MonitorStatsData()` in [server/WebServer.go](server/WebServer.go).
- `cpuusage` is a string value returned by `metrics.GetCPUUsage()`.
- `ramusage.usedPercentge` and `diskusage.avaliable` are emitted with those exact JSON field names.
- `networkstats.Name` is capitalized because the field has no explicit JSON tag.
- Runtime metadata such as `connectedclient`, `serverUptime`, `cpus`, and `arch` now live inside `machineinfo`.

---

## 2) Storage endpoint

- Path: `/api/events/storage`
- Handler: `HandleResponseData("storage")`
- Returns: a JSON object describing snapshot storage information.

### Actual payload shape

```json
{
  "totalSnapShotSize": 204800,
  "totalSnapshots": 5,
  "latestSnapshot": "2026-07-08T12:34:56Z",
  "lastCreated": "2026-07-08T12:34:56Z",
  "nextSnap": 3
}
```

### Notes

- The values are produced by `SnapshotData()` in [server/WebServer.go](server/WebServer.go).
- `latestSnapshot` and `lastCreated` are timestamps.
- `nextSnap` is a Go `time.Duration` value and is emitted as a numeric duration value.

---

## 3) Database endpoint

- Path: `/api/events/database`
- Handler: `HandleResponseData("database")`
- Returns: a JSON object containing database/storage metrics plus pub/sub topic information.

### Actual payload shape

```json
{
  "savedDataSize": 409600,
  "savedSnapShotSize": 102400,
  "totalKeysSize": 245760,
  "lastSnapShotTime": "2026-07-08T12:34:56Z",
  "ttlMetrics": {
    "activeTTLKeys": 3,
    "permanentKeys": 12,
    "totalExpiredKeys": 1,
    "expiresToday": 2,
    "nextExpiringKeys": [
      {
        "key": "user:42",
        "expiresIn": 120000000000
      }
    ]
  },
  "pubsubTopics": [
    {
      "topics": ["chat", "events"],
      "total": 2
    }
  ]
}
```

### Nested structures

#### `ttlMetrics`

```json
{
  "activeTTLKeys": 3,
  "permanentKeys": 12,
  "totalExpiredKeys": 1,
  "expiresToday": 2,
  "nextExpiringKeys": [
    {
      "key": "user:42",
      "expiresIn": 120000000000
    }
  ]
}
```

#### `pubsubTopics`

```json
[
  {
    "topics": ["chat", "events"],
    "total": 2
  }
]
```

### Notes

- This endpoint uses `DataBaseData()`.
- `ttlMetrics` is populated from `metrics.GetTTLMetrics()`.
- `pubsubTopics` is a list containing one `TopicList` object returned by `GetTopicList()`.

---

## 4) Activity endpoint

- Path: `/api/events/activity`
- Handler: `HandleResponseData("activity")`
- Returns: a JSON array of recent activity/history strings.

### Actual payload shape

```json
[
  "2026-07-08 12:34:56 INFO key set",
  "2026-07-08 12:34:54 INFO client connected",
  "2026-07-08 12:34:50 INFO snapshot created"
]
```

### Notes

- The handler reads the log content from the logger file and splits it into lines.
- The returned array is reversed and trimmed to the latest entries.
- If no data is available, the server returns:

```json
[
  "Not Data Found",
  "Chcek after a few seconds"
]
```

---

## Invalid endpoint

If a path name is not one of the supported values, the handler returns a plain string payload:

```json
"No such endpoint found"
```
