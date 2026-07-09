package master

import (
	"github.com/rahulkumarparida/roxkv/agents/registry"
)

// init registers all Master-level tools with the global registry.
// Each tool is wrapped with routing metadata (keywords, synonyms,
// examples, priority, mutation flag) so the router can score and
// filter them without hardcoded switch statements.
func init() {
	// ================================================================
	// KV Tools
	// ================================================================

	registry.Register(registry.ToolMetadata{
		Name:        "get",
		Tool:        GetKeyTool(),
		Category:    "KV",
		Description: "Read one stored value by key",
		Keywords:    []string{"get", "read", "fetch", "key", "value", "lookup"},
		Synonyms:    []string{"retrieve", "find", "show", "check", "query"},
		Examples: []string{
			"get key name",
			"what is the value of",
			"read key",
			"fetch the value for",
			"show me the value of key",
		},
		Priority:   8,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "set",
		Tool:        SetKeyTool(),
		Category:    "KV",
		Description: "Create or overwrite one key-value pair",
		Keywords:    []string{"set", "write", "store", "save", "update", "create", "put"},
		Synonyms:    []string{"add", "insert", "assign", "define"},
		Examples: []string{
			"set key name to value",
			"store this value",
			"save key",
			"write key value pair",
			"create a new key",
		},
		Priority:   9,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "del",
		Tool:        DeleteKeyTool(),
		Category:    "KV",
		Description: "Delete one key and its value from memory",
		Keywords:    []string{"delete", "del", "remove", "drop", "erase"},
		Synonyms:    []string{"destroy", "purge", "wipe", "clear"},
		Examples: []string{
			"delete key",
			"remove key",
			"drop the key",
			"erase key from store",
		},
		Priority:   10,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "keys",
		Tool:        KeysTool(),
		Category:    "KV",
		Description: "List all keys currently in memory",
		Keywords:    []string{"keys", "list", "all", "enumerate", "database"},
		Synonyms:    []string{"show", "display", "view", "browse", "print"},
		Examples: []string{
			"show all keys",
			"what keys exist",
			"list stored keys",
			"display stored keys",
			"enumerate all keys",
		},
		Priority:   8,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "save",
		Tool:        SaveTool(),
		Category:    "Persistence",
		Description: "Persist in-memory database state to disk",
		Keywords:    []string{"save", "persist", "dump", "snapshot", "backup", "disk"},
		Synonyms:    []string{"write", "export", "flush"},
		Examples: []string{
			"save the database",
			"persist current state",
			"write data to disk",
			"create a backup",
			"dump to disk",
		},
		Priority:   9,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "load",
		Tool:        LoadTool(),
		Category:    "Persistence",
		Description: "Restore previously saved data from disk into memory",
		Keywords:    []string{"load", "restore", "recover", "import", "reload"},
		Synonyms:    []string{"read", "resume", "rehydrate"},
		Examples: []string{
			"load saved data",
			"restore from disk",
			"recover database",
			"import persisted state",
		},
		Priority:   9,
		IsMutation: true,
	})

	// ================================================================
	// Monitoring Tools
	// ================================================================

	registry.Register(registry.ToolMetadata{
		Name:        "get_computer_usage",
		Tool:        GetComputerUsageTool(),
		Category:    "Monitoring",
		Description: "High-level host summary including OS, CPU, memory, runtime",
		Keywords:    []string{"computer", "machine", "host", "system", "info", "hardware"},
		Synonyms:    []string{"overview", "summary", "about", "specs"},
		Examples: []string{
			"system info",
			"tell me about this computer",
			"machine details",
			"host overview",
			"what machine is this",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_cpu_usage",
		Tool:        GetCPUUsageTool(),
		Category:    "Monitoring",
		Description: "Current CPU utilization percentage",
		Keywords:    []string{"cpu", "processor", "utilization", "load", "core"},
		Synonyms:    []string{"processing", "compute"},
		Examples: []string{
			"cpu usage",
			"how much cpu",
			"processor load",
			"cpu utilization",
			"show cpu",
		},
		Priority:   7,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_ram_usage",
		Tool:        GetRAMUsageTool(),
		Category:    "Monitoring",
		Description: "Current RAM usage details",
		Keywords:    []string{"ram", "memory", "mem", "heap"},
		Synonyms:    []string{"allocation", "free", "used"},
		Examples: []string{
			"ram usage",
			"memory usage",
			"how much memory",
			"show ram",
			"free memory",
		},
		Priority:   7,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_disk_usage",
		Tool:        GetDiskUsageTool(),
		Category:    "Monitoring",
		Description: "Disk capacity and availability for a filesystem path",
		Keywords:    []string{"disk", "storage", "filesystem", "space", "mount"},
		Synonyms:    []string{"capacity", "volume", "partition", "drive"},
		Examples: []string{
			"disk usage",
			"disk space",
			"how much disk left",
			"storage available",
			"filesystem usage",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_runtime_stats",
		Tool:        GetRuntimeStatsTool(),
		Category:    "Monitoring",
		Description: "Go runtime statistics for the RoxKV process",
		Keywords:    []string{"runtime", "goroutine", "goroutines", "heap", "gc", "allocation", "stats"},
		Synonyms:    []string{"internals", "process", "profiling"},
		Examples: []string{
			"runtime stats",
			"goroutine count",
			"go runtime info",
			"process statistics",
			"heap usage",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_server_uptime",
		Tool:        GetServerUptimeTool(),
		Category:    "Monitoring",
		Description: "How long the RoxKV server has been running",
		Keywords:    []string{"uptime", "running", "started", "duration"},
		Synonyms:    []string{"alive", "online", "elapsed"},
		Examples: []string{
			"server uptime",
			"how long running",
			"uptime",
			"when did server start",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_connected_clients",
		Tool:        GetConnectedClientsTool(),
		Category:    "Monitoring",
		Description: "Connected client metadata including IDs and activity",
		Keywords:    []string{"clients", "connected", "connections", "sessions", "users"},
		Synonyms:    []string{"active", "online", "who"},
		Examples: []string{
			"connected clients",
			"who is connected",
			"active sessions",
			"show clients",
			"list users",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_command_count",
		Tool:        GetCommandCountTool(),
		Category:    "Monitoring",
		Description: "Total commands processed since startup",
		Keywords:    []string{"command", "commands", "count", "requests", "processed"},
		Synonyms:    []string{"activity", "volume", "throughput", "operations"},
		Examples: []string{
			"command count",
			"how many commands",
			"total requests",
			"activity volume",
		},
		Priority:   5,
		IsMutation: false,
	})

	// ================================================================
	// PubSub Tools
	// ================================================================

	registry.Register(registry.ToolMetadata{
		Name:        "get_all_topics_online",
		Tool:        GetAllTopicsOnlineTool(),
		Category:    "PubSub",
		Description: "List all currently available pubsub topics",
		Keywords:    []string{"topics", "topic", "pubsub", "channels"},
		Synonyms:    []string{"list", "show", "enumerate", "available"},
		Examples: []string{
			"list topics",
			"show all topics",
			"what topics exist",
			"available channels",
		},
		Priority:   7,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_all_subscribers",
		Tool:        GetAllSubscribersTool(),
		Category:    "PubSub",
		Description: "Subscriber information across all topics",
		Keywords:    []string{"subscribers", "subscriber", "subscriptions"},
		Synonyms:    []string{"listeners", "consumers", "watchers"},
		Examples: []string{
			"all subscribers",
			"who subscribes",
			"show subscribers",
			"list all subscribers",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_all_publishers",
		Tool:        GetAllPublishersTool(),
		Category:    "PubSub",
		Description: "Publisher information across all topics",
		Keywords:    []string{"publishers", "publisher", "producers"},
		Synonyms:    []string{"senders", "emitters", "writers"},
		Examples: []string{
			"all publishers",
			"who publishes",
			"show publishers",
			"list all publishers",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_publishers",
		Tool:        GetPublishersTool(),
		Category:    "PubSub",
		Description: "Publishers for one specific topic",
		Keywords:    []string{"publishers", "publisher", "topic"},
		Synonyms:    []string{"senders", "emitters"},
		Examples: []string{
			"publishers on topic",
			"who publishes to",
			"show publishers for topic",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_subscribers",
		Tool:        GetSubscribersTool(),
		Category:    "PubSub",
		Description: "Subscribers for one specific topic",
		Keywords:    []string{"subscribers", "subscriber", "topic"},
		Synonyms:    []string{"listeners", "consumers"},
		Examples: []string{
			"subscribers on topic",
			"who subscribes to",
			"show subscribers for topic",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_topic_statistics",
		Tool:        GetTopicStatisticsTool(),
		Category:    "PubSub",
		Description: "Full snapshot of a topic including publishers, subscribers, history",
		Keywords:    []string{"statistics", "stats", "topic", "details", "snapshot"},
		Synonyms:    []string{"info", "summary", "analysis", "inspect"},
		Examples: []string{
			"topic statistics",
			"details about topic",
			"topic info",
			"inspect topic",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_inactive_topics",
		Tool:        GetInactiveTopicsTool(),
		Category:    "PubSub",
		Description: "Identify stale or idle topics",
		Keywords:    []string{"inactive", "stale", "idle", "old", "unused", "topics"},
		Synonyms:    []string{"dormant", "dead", "quiet"},
		Examples: []string{
			"inactive topics",
			"stale topics",
			"idle topics",
			"unused channels",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "broadcast_to_topic",
		Tool:        BroadcastToTopicTool(),
		Category:    "PubSub",
		Description: "Publish one message to a single topic",
		Keywords:    []string{"broadcast", "publish", "send", "message", "topic"},
		Synonyms:    []string{"emit", "push", "post", "notify"},
		Examples: []string{
			"send message to topic",
			"broadcast to topic",
			"publish message",
			"push to channel",
		},
		Priority:   8,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "broadcast_everywhere",
		Tool:        BroadcastEverywhereTool(),
		Category:    "PubSub",
		Description: "Publish one message across all topics",
		Keywords:    []string{"broadcast", "everywhere", "all", "global", "message"},
		Synonyms:    []string{"announce", "blast", "flood"},
		Examples: []string{
			"broadcast everywhere",
			"send to all topics",
			"global broadcast",
			"announce to everyone",
			"broadcast to all topics",
		},
		Priority:   8,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "history_of_topic",
		Tool:        HistoryOfTopicTool(),
		Category:    "PubSub",
		Description: "Recorded publish history for one topic",
		Keywords:    []string{"history", "messages", "log", "topic", "past"},
		Synonyms:    []string{"archive", "record", "timeline", "previous"},
		Examples: []string{
			"topic history",
			"message history",
			"what was sent on topic",
			"show history",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "delete_topic",
		Tool:        DeleteTopicTool(),
		Category:    "PubSub",
		Description: "Delete one pubsub topic",
		Keywords:    []string{"delete", "remove", "topic", "destroy"},
		Synonyms:    []string{"drop", "purge", "close"},
		Examples: []string{
			"delete topic",
			"remove topic",
			"destroy topic",
			"close channel",
		},
		Priority:   10,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "create_topic",
		Tool:        CreateTopicTool(),
		Category:    "PubSub",
		Description: "Create one new pubsub topic",
		Keywords:    []string{"create", "new", "topic", "add", "open"},
		Synonyms:    []string{"make", "start", "initialize", "register"},
		Examples: []string{
			"create topic",
			"new topic",
			"add topic",
			"open channel",
			"start a topic",
		},
		Priority:   9,
		IsMutation: true,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "remove_user",
		Tool:        RemoveUserTool(),
		Category:    "PubSub",
		Description: "Disconnect one user from the pubsub layer",
		Keywords:    []string{"remove", "disconnect", "kick", "user", "client", "ban"},
		Synonyms:    []string{"expel", "eject", "boot", "terminate"},
		Examples: []string{
			"remove user",
			"disconnect user",
			"kick client",
			"ban user",
		},
		Priority:   10,
		IsMutation: true,
	})

	// ================================================================
	// Storage Analytics Tools
	// ================================================================

	registry.Register(registry.ToolMetadata{
		Name:        "get_total_keys",
		Tool:        GetTotalKeysTool(),
		Category:    "Storage",
		Description: "Count active keys in memory",
		Keywords:    []string{"total", "count", "keys", "number", "how many"},
		Synonyms:    []string{"amount", "quantity", "size"},
		Examples: []string{
			"total keys",
			"how many keys",
			"key count",
			"number of keys",
		},
		Priority:   7,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_largest_keys",
		Tool:        GetLargestKeysTool(),
		Category:    "Storage",
		Description: "Rank largest keys by stored value size",
		Keywords:    []string{"largest", "biggest", "heaviest", "keys", "size", "top"},
		Synonyms:    []string{"huge", "big", "massive", "fat"},
		Examples: []string{
			"largest keys",
			"biggest keys",
			"keys consuming most space",
			"top keys by size",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_smallest_keys",
		Tool:        GetSmallestKeysTool(),
		Category:    "Storage",
		Description: "Rank smallest keys by stored value size",
		Keywords:    []string{"smallest", "tiniest", "lightest", "keys", "size"},
		Synonyms:    []string{"small", "tiny", "minimal", "compact"},
		Examples: []string{
			"smallest keys",
			"tiniest keys",
			"lightest keys by size",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_average_value_size",
		Tool:        GetAverageValueSizeTool(),
		Category:    "Storage",
		Description: "Average size of stored values",
		Keywords:    []string{"average", "mean", "size", "value", "typical"},
		Synonyms:    []string{"avg", "median"},
		Examples: []string{
			"average value size",
			"typical value size",
			"mean size of values",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_namespaces",
		Tool:        GetNamespacesTool(),
		Category:    "Storage",
		Description: "Summarize key distribution by namespace",
		Keywords:    []string{"namespace", "namespaces", "group", "distribution", "partition"},
		Synonyms:    []string{"bucket", "segment", "category"},
		Examples: []string{
			"show namespaces",
			"key distribution",
			"how are keys grouped",
			"list namespaces",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_ttl_metrics",
		Tool:        GetTTLMetricsTool(),
		Category:    "Storage",
		Description: "Aggregate TTL metrics",
		Keywords:    []string{"ttl", "expiration", "expiry", "lifetime", "metrics"},
		Synonyms:    []string{"timeout", "expire", "duration"},
		Examples: []string{
			"ttl metrics",
			"expiration stats",
			"ttl health",
			"key lifetime",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_expired_keys",
		Tool:        GetExpiredKeysTool(),
		Category:    "Storage",
		Description: "List keys whose TTL has already expired",
		Keywords:    []string{"expired", "stale", "ttl", "dead", "keys"},
		Synonyms:    []string{"overdue", "past", "lapsed"},
		Examples: []string{
			"expired keys",
			"which keys expired",
			"show expired",
			"stale keys",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_upcoming_expirations",
		Tool:        GetUpcomingExpirationsTool(),
		Category:    "Storage",
		Description: "List keys about to expire soon",
		Keywords:    []string{"upcoming", "expiring", "soon", "ttl", "next"},
		Synonyms:    []string{"imminent", "near", "approaching"},
		Examples: []string{
			"upcoming expirations",
			"keys expiring soon",
			"what will expire next",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_keys_without_ttl",
		Tool:        GetKeysWithoutTTLTool(),
		Category:    "Storage",
		Description: "List keys with no TTL assigned",
		Keywords:    []string{"permanent", "no ttl", "without", "ttl", "forever", "unbounded"},
		Synonyms:    []string{"persistent", "immortal", "eternal"},
		Examples: []string{
			"keys without ttl",
			"permanent keys",
			"which keys never expire",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_snapshot_count",
		Tool:        GetSnapshotCountTool(),
		Category:    "Persistence",
		Description: "Count snapshot files on disk",
		Keywords:    []string{"snapshot", "snapshots", "count", "backup", "backups"},
		Synonyms:    []string{"checkpoints", "saves"},
		Examples: []string{
			"how many snapshots",
			"snapshot count",
			"number of backups",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_latest_snapshot",
		Tool:        GetLatestSnapshotTool(),
		Category:    "Persistence",
		Description: "Most recent snapshot metadata",
		Keywords:    []string{"latest", "recent", "last", "snapshot", "newest"},
		Synonyms:    []string{"current", "most recent"},
		Examples: []string{
			"latest snapshot",
			"last backup",
			"most recent snapshot",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_snapshot_size",
		Tool:        GetSnapshotSizeTool(),
		Category:    "Persistence",
		Description: "Total disk space consumed by snapshots",
		Keywords:    []string{"snapshot", "size", "disk", "space", "storage"},
		Synonyms:    []string{"footprint", "weight"},
		Examples: []string{
			"snapshot size",
			"how much space do snapshots use",
			"snapshot disk usage",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_persistence_health",
		Tool:        GetPersistenceHealthTool(),
		Category:    "Persistence",
		Description: "Persistence health summary",
		Keywords:    []string{"persistence", "health", "snapshot", "status", "check"},
		Synonyms:    []string{"diagnostics", "wellness", "integrity"},
		Examples: []string{
			"persistence health",
			"is persistence healthy",
			"check persistence",
			"backup status",
		},
		Priority:   6,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_oldest_key",
		Tool:        GetOldestKeyTool(),
		Category:    "Storage",
		Description: "Metadata for the oldest active key",
		Keywords:    []string{"oldest", "first", "earliest", "key", "created"},
		Synonyms:    []string{"original", "initial", "ancient"},
		Examples: []string{
			"oldest key",
			"first key created",
			"earliest key",
			"which key is oldest",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_newest_key",
		Tool:        GetNewestKeyTool(),
		Category:    "Storage",
		Description: "Metadata for the newest active key",
		Keywords:    []string{"newest", "latest", "recent", "key", "last created"},
		Synonyms:    []string{"fresh", "new", "most recent"},
		Examples: []string{
			"newest key",
			"last key created",
			"most recent key",
			"which key is newest",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_most_accessed_key",
		Tool:        GetMostAccessedKeyTool(),
		Category:    "Storage",
		Description: "Metadata for the most frequently accessed key",
		Keywords:    []string{"most accessed", "hottest", "popular", "frequent", "accessed", "hot"},
		Synonyms:    []string{"busiest", "top", "heaviest traffic"},
		Examples: []string{
			"most accessed key",
			"hottest key",
			"most popular key",
			"most used key",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_least_accessed_key",
		Tool:        GetLeastAccessedKeyTool(),
		Category:    "Storage",
		Description: "Metadata for the least frequently accessed key",
		Keywords:    []string{"least accessed", "coldest", "unused", "infrequent", "cold"},
		Synonyms:    []string{"quietest", "rarely", "bottom"},
		Examples: []string{
			"least accessed key",
			"coldest key",
			"least used key",
			"rarely accessed key",
		},
		Priority:   5,
		IsMutation: false,
	})

	registry.Register(registry.ToolMetadata{
		Name:        "get_recently_modified_keys",
		Tool:        GetRecentlyModifiedKeysTool(),
		Category:    "Storage",
		Description: "List keys modified most recently",
		Keywords:    []string{"recently", "modified", "changed", "updated", "recent"},
		Synonyms:    []string{"latest changes", "mutations", "edits"},
		Examples: []string{
			"recently modified keys",
			"what changed recently",
			"latest key changes",
			"recent updates",
		},
		Priority:   6,
		IsMutation: false,
	})
}
