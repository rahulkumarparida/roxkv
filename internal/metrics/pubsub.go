package metrics

import (
	"sort"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type ClientSnapshot struct {
	ID           string    `json:"id"`
	Role         string    `json:"role"`
	LastUsed     time.Time `json:"lastInteraction"`
	ConnectedAt  time.Time `json:"connectedAt"`
	Interactions int       `json:"totalInteractions"`
}

type TopicSnapshot struct {
	Subscribers   []ClientSnapshot      `json:"subscribers"`
	Publishers    []ClientSnapshot      `json:"publishers"`
	Topic         string                `json:"topic"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	LastPublisher *ClientSnapshot       `json:"lastPublisher"`
	PublishCount  int                   `json:"publishCount"`
	History       []pubsub.TopicHistory `json:"history"`
	TotalSize     int64                 `json:"totalSize"`
}

func snapshotClient(client *utils.NewClient) ClientSnapshot {
	if client == nil {
		return ClientSnapshot{}
	}

	client.Mu.RLock()
	defer client.Mu.RUnlock()

	id, _ := client.ID.(string)

	return ClientSnapshot{
		ID:           id,
		Role:         client.Role,
		LastUsed:     client.LastUsed,
		ConnectedAt:  client.ConnectedAt,
		Interactions: client.Interactions,
	}
}

func snapshotClients(clients []*utils.NewClient) []ClientSnapshot {
	if len(clients) == 0 {
		return []ClientSnapshot{}
	}

	snapshots := make([]ClientSnapshot, 0, len(clients))
	for _, client := range clients {
		if client == nil {
			continue
		}
		snapshots = append(snapshots, snapshotClient(client))
	}

	return snapshots
}

func snapshotTopic(channel *pubsub.SubrChannel) TopicSnapshot {
	if channel == nil {
		return TopicSnapshot{}
	}

	channel.Mu.RLock()
	defer channel.Mu.RUnlock()

	history := append([]pubsub.TopicHistory(nil), channel.History...)

	var lastPublisher *ClientSnapshot
	if channel.LastPublisher != nil {
		snapshot := snapshotClient(channel.LastPublisher)
		lastPublisher = &snapshot
	}

	return TopicSnapshot{
		Subscribers:   snapshotClients(channel.Subscribers),
		Publishers:    snapshotClients(channel.Publisher),
		Topic:         channel.Topic,
		CreatedAt:     channel.CreatedAt,
		UpdatedAt:     channel.UpdatedAt,
		LastPublisher: lastPublisher,
		PublishCount:  channel.PublishCount,
		History:       history,
		TotalSize:     channel.TotalSize,
	}
}

func GetAllTopicsOnline(client *utils.NewClient) []string {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	topics := pubsub.GetTopics(client)
	logger.InfoLog("System asked for all the topics name")
	return topics
}

func GetAllSubscribers(client *utils.NewClient) []ClientSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	clients := pubsub.GetAllMembers(client, "sub")
	return snapshotClients(clients)
}

func GetAllPublisher(client *utils.NewClient) []ClientSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	clients := pubsub.GetAllMembers(client, "pub")
	return snapshotClients(clients)
}

func GetPublishers(client *utils.NewClient, topic string) []ClientSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	channel := pubsub.GetChannel(client, topic)
	clients := pubsub.GetClients(client, channel, "pub")

	return snapshotClients(clients)
}

func GetSubscribers(client *utils.NewClient, topic string) []ClientSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	channel := pubsub.GetChannel(client, topic)
	clients := pubsub.GetClients(client, channel, "sub")
	return snapshotClients(clients)
}

func GetTopicStatistics(client *utils.NewClient, topic string) TopicSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return TopicSnapshot{}
	}

	channelData := pubsub.GetChannel(client, topic)
	logger.InfoLog("System asked for statistics of a topic " + topic)
	return snapshotTopic(channelData)
}

func GetInactiveTopics(client *utils.NewClient) []TopicSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	pubsub.Helper.Mu.RLock()
	channels := make([]*pubsub.SubrChannel, len(pubsub.Helper.Channels))
	copy(channels, pubsub.Helper.Channels)
	pubsub.Helper.Mu.RUnlock()

	sort.Slice(channels, func(i, j int) bool {
		return channels[i].CreatedAt.Before(channels[j].CreatedAt)
	})

	inactiveTopics := make([]TopicSnapshot, 0, len(channels))
	for _, channel := range channels {
		inactiveTopics = append(inactiveTopics, snapshotTopic(channel))
	}

	logger.InfoLog("System asked for inactive topics")
	return inactiveTopics
}

func BroadcastToTopic(client *utils.NewClient, topic string, msg string) bool {
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.Broker(client, topic, msg)
	logger.InfoLog("System broadcasted a message to " + topic)
	return true
}

func BroadcastEverywhere(client *utils.NewClient, msg string) bool {
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.PublishToAllTopics(client, msg)
	logger.InfoLog("System broadcasted a message every where")
	return true
}

func HistoryOfTopic(client *utils.NewClient, topic string) []pubsub.TopicHistory {
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	channel := pubsub.GetChannel(client, topic)
	channel.Mu.RLock()
	history := append([]pubsub.TopicHistory(nil), channel.History...)
	channel.Mu.RUnlock()

	logger.InfoLog("System requested for thr history of " + topic)
	return history
}

func DeleteTopic(client *utils.NewClient, topic string) bool {
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.CloseChannel(client, topic)
	logger.InfoLog("System deleted a channel " + topic)
	return true
}

func CreateTopic(client *utils.NewClient, topic string) TopicSnapshot {
	if client.Role != string(utils.RoleSystem) {
		return TopicSnapshot{}
	}

	channel := pubsub.CreateTopic(client, topic)
	logger.InfoLog("System cerated a topic " + channel.Topic)
	return snapshotTopic(channel)
}

func RemoveUser(client *utils.NewClient, portaddres string, reason string) bool {
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	users := pubsub.GetAllMembers(client, "sub")

	for _, user := range users {
		if user != nil && user.ID == portaddres {
			user.Conn.Write([]byte(reason))
			user.Conn.Close()
			break
		}
	}

	logger.InfoLog("System removed the user" + portaddres + " beacause " + reason)
	return true
}
