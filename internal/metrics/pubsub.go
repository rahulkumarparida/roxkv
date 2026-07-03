package metrics

import (
	"sort"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)


// GetTopics()
// GetAllSubscribers()
// GetAllPublishers()
// GetSubscribers(topic string)
// GetPublishers(topic string)
// GetPublishStats(topic string)
// GetInactiveTopics()
// BroadcastToTopic(topic string)
// BroadcastEverywhere()
// HistoryOfTopic(topic string)
// DeleteTopic(topic string)
// CreateTopic(topic string)
// RemoveUser(portaddress string)



func GetAllTopicsOnline(client *utils.NewClient) []string{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	
	topics := pubsub.GetTopics(client)
	logger.InfoLog("System asked for all the topics name")
	return topics

}

func GetAllSubscribers(client *utils.NewClient) []*utils.NewClient{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	clients := pubsub.GetAllMembers(client,"sub")
	return clients
}


func GetAllPublisher(client *utils.NewClient) []*utils.NewClient{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	clients := pubsub.GetAllMembers(client,"pub")
	return clients
}


func GetPublishers(client *utils.NewClient,topic string) []*utils.NewClient{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}
	channel := pubsub.GetChannel(client,topic)
	

	clients := pubsub.GetClients(client,channel,"pub")

	return clients

}

func GetSubscribers(client *utils.NewClient,topic string) []*utils.NewClient{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	channel := pubsub.GetChannel(client,topic)
	clients := pubsub.GetClients(client,channel,"sub")
	return clients
}

func GetTopicStatistics(client *utils.NewClient,topic string) *pubsub.SubrChannel{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}

	channelData := pubsub.GetChannel(client,topic)	
	logger.InfoLog("System asked for statistics of a topic "+topic)
	return channelData
}

func GetInactiveTopics(client *utils.NewClient)  []*pubsub.SubrChannel{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}
	pubsub.Helper.Mu.Lock()
	channels := make([]*pubsub.SubrChannel,len(pubsub.Helper.Channels))
	copy(channels,pubsub.Helper.Channels)
	pubsub.Helper.Mu.Lock()


	sort.Slice(channels, func(i, j int) bool {
		return channels[i].CreatedAt.Before(channels[j].CreatedAt)
	})
	logger.InfoLog("System asked for inactive topics")
	return channels
}

func BroadcastToTopic(client *utils.NewClient,topic string,msg string) bool{
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.Broker(client,topic,msg)
	logger.InfoLog("System broadcasted a message to "+topic)
	return  true
}

func BroadcastEverywhere(client *utils.NewClient,msg string)  any{
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.PublishToAllTopics(client,msg)
	logger.InfoLog("System broadcasted a message every where")
	return true

}




func HistoryOfTopic(client *utils.NewClient,topic string)  []pubsub.TopicHistory{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}
	
	channel := pubsub.GetChannel(client,topic)
	logger.InfoLog("System requested for thr history of "+topic)
	return channel.History
}


func DeleteTopic(client *utils.NewClient,topic string) bool{
	if client.Role != string(utils.RoleSystem) {
		return false
	}

	pubsub.CloseChannel(client,topic)
	logger.InfoLog("System deleted a channel "+topic)
	return true
}

func CreateTopic(client *utils.NewClient,topic string) *pubsub.SubrChannel{
	if client.Role != string(utils.RoleSystem) {
		return nil
	}
	channel := pubsub.CreateTopic(client,topic)
	logger.InfoLog("System cerated a topic "+channel.Topic)
	return channel
}

func RemoveUser(client *utils.NewClient,portaddres string,reason string) any{
		if client.Role != string(utils.RoleSystem) {
		return nil
	}
	users := GetAllSubscribers(client)

	for _, user := range users {
		if user.ID == portaddres {
			user.Conn.Write([]byte(reason))
			user.Conn.Close()
			break
		}
	}

	logger.InfoLog("System removed the user"+portaddres+" beacause "+reason)
	return true
}