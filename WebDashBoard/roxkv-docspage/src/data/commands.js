export const commands = [
  {
    name: "SET",
    category: "Strings",
    description: "Set the string value of a key.",
    syntax: "SET key value",
    arguments: [
      { name: "key", description: "The key to store the value under." },
      { name: "value", description: "The string value to store." }
    ],
    behavior: "RoxKV evaluates the command and stores the key-value pair in its internal store. Currently ignores additional options like EX/PX in standard SET, handling TTL via the EXPIRE command.",
    returnValue: "Simple String response `+OK\\r\\n` on success, or an Error response `-ERR wrong number of arguments for 'set' command\\r\\n` if arguments are invalid.",
    examples: [
      "SET mykey \"Hello World\""
    ],
    notes: "RoxKV treats all values as string lists internally, joining them when required.",
    related: ["GET", "DEL", "EXISTS"]
  },
  {
    name: "GET",
    category: "Strings",
    description: "Get the value of a key.",
    syntax: "GET key",
    arguments: [
      { name: "key", description: "The key to retrieve." }
    ],
    behavior: "Looks up the key in the store. If the internal value is stored as an array of strings, it joins them with a space.",
    returnValue: "Bulk String response with the value, or Null value `$-1\\r\\n` if the key does not exist.",
    examples: [
      "GET mykey"
    ],
    notes: "",
    related: ["SET", "MGET"]
  },
  {
    name: "MGET",
    category: "Strings",
    description: "Get the values of all given keys.",
    syntax: "MGET key [key ...]",
    arguments: [
      { name: "key", description: "The keys to retrieve." }
    ],
    behavior: "Looks up all the specified keys in the store.",
    returnValue: "Array response containing the values. For every key that does not exist, a Null value is returned.",
    examples: [
      "MGET key1 key2 nonexisting"
    ],
    notes: "",
    related: ["GET"]
  },
  {
    name: "DEL",
    category: "Keys",
    description: "Delete a key.",
    syntax: "DEL key [key ...]",
    arguments: [
      { name: "key", description: "The keys to delete." }
    ],
    behavior: "Iterates through the provided keys and deletes them from the store.",
    returnValue: "Integer response representing the number of keys that were removed.",
    examples: [
      "DEL key1 key2"
    ],
    notes: "",
    related: ["EXISTS"]
  },
  {
    name: "EXISTS",
    category: "Keys",
    description: "Determine if a key exists.",
    syntax: "EXISTS key [key ...]",
    arguments: [
      { name: "key", description: "The keys to check." }
    ],
    behavior: "Checks the internal store for the presence of each specified key.",
    returnValue: "Integer response representing the number of keys that exist from those specified as arguments.",
    examples: [
      "EXISTS mykey1 mykey2"
    ],
    notes: "",
    related: ["DEL"]
  },
  {
    name: "KEYS",
    category: "Keys",
    description: "Find all keys matching the given pattern.",
    syntax: "KEYS pattern",
    arguments: [
      { name: "pattern", description: "The search pattern." }
    ],
    behavior: "Returns all keys matching pattern. Note: RoxKV currently implements a simple substring search rather than a full glob-style matching for keys. If a single character is provided, it returns all keys.",
    returnValue: "Array response containing the matching keys.",
    examples: [
      "KEYS *",
      "KEYS user"
    ],
    notes: "Use with caution on large databases as it may affect performance.",
    related: ["RANDOMKEY", "EXISTS"]
  },
  {
    name: "EXPIRE",
    category: "Expiration",
    description: "Set a key's time to live in seconds.",
    syntax: "EXPIRE key seconds",
    arguments: [
      { name: "key", description: "The key to set the timeout on." },
      { name: "seconds", description: "The TTL in seconds." }
    ],
    behavior: "Internally re-evaluates the key and sets the TTL metadata.",
    returnValue: "Integer response `1` if the timeout was set, `0` if the key does not exist.",
    examples: [
      "EXPIRE mykey 10"
    ],
    notes: "",
    related: ["TTL", "PERSIST"]
  },
  {
    name: "TTL",
    category: "Expiration",
    description: "Get the time to live for a key in seconds.",
    syntax: "TTL key",
    arguments: [
      { name: "key", description: "The key to check." }
    ],
    behavior: "Returns the remaining time to live of a key that has a timeout.",
    returnValue: "Integer response. Returns `-2` if the key exists but has no associated expire. Returns `-1` if the key does not exist.",
    examples: [
      "TTL mykey"
    ],
    notes: "Note that standard Redis returns -1 for no expire and -2 for non-existent, but RoxKV returns -2 for no expire and throws an error for non-existent.",
    related: ["EXPIRE", "PERSIST"]
  },
  {
    name: "PERSIST",
    category: "Expiration",
    description: "Remove the expiration from a key.",
    syntax: "PERSIST key",
    arguments: [
      { name: "key", description: "The key to remove the timeout from." }
    ],
    behavior: "Removes the TTL metadata from the key.",
    returnValue: "Integer response `1` if the timeout was removed, `0` if key does not exist or does not have an associated timeout.",
    examples: [
      "PERSIST mykey"
    ],
    notes: "",
    related: ["EXPIRE", "TTL"]
  },
  {
    name: "RENAME",
    category: "Keys",
    description: "Rename a key.",
    syntax: "RENAME key newkey",
    arguments: [
      { name: "key", description: "The current key name." },
      { name: "newkey", description: "The new key name." }
    ],
    behavior: "Retrieves the value of the old key, sets it to the new key, and deletes the old key.",
    returnValue: "Simple String response `+OK\\r\\n` on success.",
    examples: [
      "RENAME mykey myotherkey"
    ],
    notes: "",
    related: ["KEYS"]
  },
  {
    name: "STRLEN",
    category: "Strings",
    description: "Get the length of the value stored in a key.",
    syntax: "STRLEN key",
    arguments: [
      { name: "key", description: "The key whose length you want." }
    ],
    behavior: "Returns the byte length of the value associated with the key.",
    returnValue: "Integer response representing the string length, or `0` if the key does not exist.",
    examples: [
      "STRLEN mykey"
    ],
    notes: "",
    related: ["GET"]
  },
  {
    name: "INCR",
    category: "Strings",
    description: "Increment the integer value of a key by one.",
    syntax: "INCR key",
    arguments: [
      { name: "key", description: "The key to increment." }
    ],
    behavior: "Parses the value as an integer (base 10, 64-bit), increments it, and saves it back.",
    returnValue: "Integer response representing the value of key after the increment.",
    examples: [
      "INCR mycounter"
    ],
    notes: "If the value is not an integer, an error is returned.",
    related: ["DECR", "INCRBY", "DECRBY"]
  },
  {
    name: "DECR",
    category: "Strings",
    description: "Decrement the integer value of a key by one.",
    syntax: "DECR key",
    arguments: [
      { name: "key", description: "The key to decrement." }
    ],
    behavior: "Parses the value as an integer, decrements it, and saves it back.",
    returnValue: "Integer response representing the value of key after the decrement.",
    examples: [
      "DECR mycounter"
    ],
    notes: "",
    related: ["INCR", "INCRBY", "DECRBY"]
  },
  {
    name: "INCRBY",
    category: "Strings",
    description: "Increment the integer value of a key by the given amount.",
    syntax: "INCRBY key increment",
    arguments: [
      { name: "key", description: "The key to increment." },
      { name: "increment", description: "The amount to increment by." }
    ],
    behavior: "Parses both the value and the argument as integers, adds them, and saves it back.",
    returnValue: "Integer response representing the value of key after the increment.",
    examples: [
      "INCRBY mycounter 5"
    ],
    notes: "",
    related: ["INCR", "DECR", "DECRBY"]
  },
  {
    name: "DECRBY",
    category: "Strings",
    description: "Decrement the integer value of a key by the given amount.",
    syntax: "DECRBY key decrement",
    arguments: [
      { name: "key", description: "The key to decrement." },
      { name: "decrement", description: "The amount to decrement by." }
    ],
    behavior: "Parses both the value and the argument as integers, subtracts the decrement, and saves it back.",
    returnValue: "Integer response representing the value of key after the decrement.",
    examples: [
      "DECRBY mycounter 5"
    ],
    notes: "",
    related: ["INCR", "DECR", "INCRBY"]
  },
  {
    name: "PING",
    category: "Connection",
    description: "Ping the server.",
    syntax: "PING [message]",
    arguments: [
      { name: "message", description: "Optional message to echo back." }
    ],
    behavior: "Returns PONG if no argument is provided, otherwise returns a copy of the argument. Available in subscriber mode.",
    returnValue: "Simple String response `+PONG\\r\\n` or Bulk String response of the message.",
    examples: [
      "PING",
      "PING \"hello world\""
    ],
    notes: "",
    related: ["ECHO", "QUIT"]
  },
  {
    name: "ECHO",
    category: "Connection",
    description: "Echo the given string.",
    syntax: "ECHO message",
    arguments: [
      { name: "message", description: "The string to echo." }
    ],
    behavior: "Returns the argument provided.",
    returnValue: "Bulk String response.",
    examples: [
      "ECHO \"Hello World!\""
    ],
    notes: "",
    related: ["PING"]
  },
  {
    name: "QUIT",
    category: "Connection",
    description: "Close the connection.",
    syntax: "QUIT",
    arguments: [],
    behavior: "Cleans up client state, removes the client from any topics/pattern subscriptions, and closes the TCP connection. Available in subscriber mode.",
    returnValue: "Simple String response `+OK\\r\\n` (followed by connection close).",
    examples: [
      "QUIT"
    ],
    notes: "",
    related: ["RESET"]
  },
  {
    name: "RESET",
    category: "Connection",
    description: "Reset the connection state.",
    syntax: "RESET",
    arguments: [],
    behavior: "Cleans up client state, dropping any topic and pattern subscriptions, and resetting the client mode to default. Available in subscriber mode.",
    returnValue: "Simple String response `+RESET\\r\\n`.",
    examples: [
      "RESET"
    ],
    notes: "",
    related: ["QUIT"]
  },
  {
    name: "SAVE",
    category: "Persistence",
    description: "Synchronously save the dataset to disk.",
    syntax: "SAVE",
    arguments: [],
    behavior: "Triggers a synchronous backup of the internal store using the persistence module.",
    returnValue: "Simple String response `+OK\\r\\n`.",
    examples: [
      "SAVE"
    ],
    notes: "",
    related: ["LASTSAVE"]
  },
  {
    name: "LASTSAVE",
    category: "Persistence",
    description: "Get the UNIX time stamp of the last successful save to disk.",
    syntax: "LASTSAVE",
    arguments: [],
    behavior: "Retrieves the timestamp of the last persistence operation.",
    returnValue: "Integer response representing a UNIX timestamp.",
    examples: [
      "LASTSAVE"
    ],
    notes: "",
    related: ["SAVE"]
  },
  {
    name: "DBSIZE",
    category: "Server",
    description: "Return the number of keys in the selected database.",
    syntax: "DBSIZE",
    arguments: [],
    behavior: "Counts the total number of keys currently in the store.",
    returnValue: "Integer response.",
    examples: [
      "DBSIZE"
    ],
    notes: "",
    related: ["KEYS", "RANDOMKEY"]
  },
  {
    name: "RANDOMKEY",
    category: "Keys",
    description: "Return a random key from the keyspace.",
    syntax: "RANDOMKEY",
    arguments: [],
    behavior: "Selects a random key from all currently active keys in the store.",
    returnValue: "Bulk String response representing the key name.",
    examples: [
      "RANDOMKEY"
    ],
    notes: "",
    related: ["KEYS"]
  },
  {
    name: "UPTIME",
    category: "Server",
    description: "Get the server uptime.",
    syntax: "UPTIME",
    arguments: [],
    behavior: "Calculates the time since the RoxKV server started.",
    returnValue: "Array response with two values: uptime in seconds, and uptime in milliseconds.",
    examples: [
      "UPTIME"
    ],
    notes: "Note that this differs from standard Redis where uptime is a field in the INFO command.",
    related: []
  },
  {
    name: "SUBSCRIBE",
    category: "Pub/Sub",
    description: "Listen for messages published to the given channels.",
    syntax: "SUBSCRIBE channel [channel ...]",
    arguments: [
      { name: "channel", description: "The channel(s) to subscribe to." }
    ],
    behavior: "Puts the client into subscriber mode. The client will only accept subscriber mode commands.",
    returnValue: "For each channel, returns an Array response containing 'subscribe', the channel name, and the total number of subscriptions.",
    examples: [
      "SUBSCRIBE news sport"
    ],
    notes: "Once a client enters subscriber mode, commands like SET or GET are rejected.",
    related: ["UNSUBSCRIBE", "PUBLISH", "PSUBSCRIBE"]
  },
  {
    name: "UNSUBSCRIBE",
    category: "Pub/Sub",
    description: "Stop listening for messages posted to the given channels.",
    syntax: "UNSUBSCRIBE [channel [channel ...]]",
    arguments: [
      { name: "channel", description: "The channel(s) to unsubscribe from." }
    ],
    behavior: "Removes the client from the specified channels. If no channels are remaining, the client drops out of subscriber mode.",
    returnValue: "Array response containing 'unsubscribe', the channel name, and the remaining subscription count.",
    examples: [
      "UNSUBSCRIBE news"
    ],
    notes: "",
    related: ["SUBSCRIBE"]
  },
  {
    name: "PUBLISH",
    category: "Pub/Sub",
    description: "Post a message to a channel.",
    syntax: "PUBLISH channel message",
    arguments: [
      { name: "channel", description: "The target channel." },
      { name: "message", description: "The message to send." }
    ],
    behavior: "Sends the message to all exact-match subscribers, as well as pattern-match subscribers.",
    returnValue: "Integer response representing the number of clients that received the message.",
    examples: [
      "PUBLISH news \"breaking news\""
    ],
    notes: "",
    related: ["SUBSCRIBE", "PSUBSCRIBE"]
  },
  {
    name: "PSUBSCRIBE",
    category: "Pub/Sub",
    description: "Listen for messages published to channels matching the given patterns.",
    syntax: "PSUBSCRIBE pattern [pattern ...]",
    arguments: [
      { name: "pattern", description: "The pattern(s) to subscribe to." }
    ],
    behavior: "Puts the client into subscriber mode and registers a pattern. RoxKV's implementation splits topics by `.` or `:` and handles the `*` wildcard for specific parts.",
    returnValue: "Array response containing 'psubscribe', the pattern, and the active pattern subscription count.",
    examples: [
      "PSUBSCRIBE news.*", 
      "PSUBSCRIBE news.*.notification"
    ],
    notes: "",
    related: ["PUNSUBSCRIBE", "PUBLISH", "SUBSCRIBE"]
  },
  {
    name: "PUNSUBSCRIBE",
    category: "Pub/Sub",
    description: "Stop listening for messages posted to channels matching the given patterns.",
    syntax: "PUNSUBSCRIBE [pattern [pattern ...]]",
    arguments: [
      { name: "pattern", description: "The pattern(s) to unsubscribe from." }
    ],
    behavior: "Removes the pattern from the client's subscriptions.",
    returnValue: "Array response containing 'punsubscribe', the pattern, and the remaining pattern count.",
    examples: [
      "PUNSUBSCRIBE news.*",
      "PUNSUBSCRIBE news.*.notification"
    ],
    notes: "",
    related: ["PSUBSCRIBE"]
  },
  {
    name: "HISTORY",
    category: "Server",
    description: "Returns the command history.",
    syntax: "HISTORY",
    arguments: [],
    behavior: "Retrieves recent server commands using the internal monitor.",
    returnValue: "Returns strings in a stream via repeated Bulk String writes, separated by sleep.",
    examples: [
      "HISTORY"
    ],
    notes: "Invoked via the 'monitor' string in the CLI.",
    related: []
  }
];
