import InfoTable from '../../../components/docs/InfoTable';

export default function ToolRegistry() {
  const kvTools = [
    { name: 'get', description: 'Retrieve a value by key' },
    { name: 'set', description: 'Set a key-value pair' },
    { name: 'del', description: 'Delete a key' },
    { name: 'keys', description: 'List all keys matching a pattern' },
  ];

  const persistenceTools = [
    { name: 'save', description: 'Trigger background save' },
    { name: 'load', description: 'Load from snapshot' },
    { name: 'get_snapshot_count', description: 'Number of snapshots' },
    { name: 'get_latest_snapshot', description: 'Info on latest snapshot' },
    { name: 'get_snapshot_size', description: 'Size of snapshots' },
    { name: 'get_persistence_health', description: 'Health of persistence layer' },
  ];

  const monitoringTools = [
    { name: 'get_computer_usage', description: 'General system usage' },
    { name: 'get_cpu_usage', description: 'CPU utilization' },
    { name: 'get_ram_usage', description: 'Memory utilization' },
    { name: 'get_disk_usage', description: 'Storage utilization' },
    { name: 'get_runtime_stats', description: 'Go runtime statistics' },
    { name: 'get_server_uptime', description: 'Server uptime duration' },
    { name: 'get_connected_clients', description: 'Number of active clients' },
    { name: 'get_command_count', description: 'Total commands processed' },
  ];

  const pubsubTools = [
    { name: 'get_all_topics_online', description: 'List active topics' },
    { name: 'get_all_subscribers', description: 'List all subscribers' },
    { name: 'get_all_publishers', description: 'List all publishers' },
    { name: 'get_publishers', description: 'Publishers for specific topic' },
    { name: 'get_subscribers', description: 'Subscribers for specific topic' },
    { name: 'get_topic_statistics', description: 'Stats for specific topic' },
    { name: 'get_inactive_topics', description: 'Topics with no activity' },
    { name: 'broadcast_to_topic', description: 'Send message to topic' },
    { name: 'broadcast_everywhere', description: 'Global broadcast' },
    { name: 'history_of_topic', description: 'Message history of topic' },
    { name: 'delete_topic', description: 'Remove a topic' },
    { name: 'create_topic', description: 'Create a new topic' },
    { name: 'remove_user', description: 'Disconnect a user/subscriber' },
  ];

  const storageAnalyticsTools = [
    { name: 'get_total_keys', description: 'Count of all keys' },
    { name: 'get_largest_keys', description: 'Keys taking most memory' },
    { name: 'get_smallest_keys', description: 'Keys taking least memory' },
    { name: 'get_average_value_size', description: 'Average value size' },
    { name: 'get_namespaces', description: 'List key namespaces' },
    { name: 'get_ttl_metrics', description: 'TTL distribution' },
    { name: 'get_expired_keys', description: 'Keys pending cleanup' },
    { name: 'get_upcoming_expirations', description: 'Keys expiring soon' },
    { name: 'get_keys_without_ttl', description: 'Persistent keys' },
    { name: 'get_oldest_key', description: 'Oldest inserted key' },
    { name: 'get_newest_key', description: 'Most recently inserted key' },
    { name: 'get_most_accessed_key', description: 'Highest read frequency' },
    { name: 'get_least_accessed_key', description: 'Lowest read frequency' },
    { name: 'get_recently_modified_keys', description: 'Recently changed keys' },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Tool Registry</h1>
        <p className="text-xl text-[#9898ab]">Comprehensive list of 42+ capabilities accessible to RoxAI</p>
      </div>
      <div className="section-divider" />

      <h2 id="overview" className="text-2xl font-semibold text-white mb-4">Registration System</h2>
      <p className="text-[#9898ab] leading-relaxed mb-4">
        All tools in RoxAI are registered globally via the <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">registry.Register</code> method defined in <code className="text-[#a78bfa] bg-[#1e1e2e] px-1.5 py-0.5 rounded">agents/Master/tools_registry.go</code>.
        The definition structure includes the tool's name, description, required parameters, category, keywords, and example phrases (used for semantic routing).
      </p>

      <h2 id="kv-tools" className="text-2xl font-semibold text-white mb-4">Key-Value Operations (4 Tools)</h2>
      <InfoTable headers={['Tool Name', 'Description']} rows={kvTools.map(t => [t.name, t.description])} />

      <h2 id="persistence-tools" className="text-2xl font-semibold text-white mb-4">Persistence (6 Tools)</h2>
      <InfoTable headers={['Tool Name', 'Description']} rows={persistenceTools.map(t => [t.name, t.description])} />

      <h2 id="monitoring-tools" className="text-2xl font-semibold text-white mb-4">Monitoring & Metrics (8 Tools)</h2>
      <InfoTable headers={['Tool Name', 'Description']} rows={monitoringTools.map(t => [t.name, t.description])} />

      <h2 id="pubsub-tools" className="text-2xl font-semibold text-white mb-4">PubSub Management (13 Tools)</h2>
      <InfoTable headers={['Tool Name', 'Description']} rows={pubsubTools.map(t => [t.name, t.description])} />

      <h2 id="analytics-tools" className="text-2xl font-semibold text-white mb-4">Storage Analytics (11 Tools)</h2>
      <InfoTable headers={['Tool Name', 'Description']} rows={storageAnalyticsTools.map(t => [t.name, t.description])} />
    </div>
  );
}
