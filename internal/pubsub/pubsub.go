package pubsub

import (
	"slices"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

// Single Channel and will contain many to many relationship with Subs and Pubs
type TopicHistory struct {
	Message   string    `json:"message"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

type SubrChannel struct {
	Subscribers   []*utils.NewClient `json:"subscribers"`
	Publisher     []*utils.NewClient `json:"publisher"`
	Topic         string             `json:"topic"`
	PTopic 		  map[string]string			 `json:"ptopic"`
	SubscribeChan chan string        `json:"-"`
	Wg            sync.WaitGroup     `json:"-"`
	Mu            sync.RWMutex       `json:"-"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	LastPublisher *utils.NewClient   `json:"lastPublisher"`
	PublishCount  int                `json:"publishCount"`
	History       []TopicHistory     `json:"history"`
	TotalSize     int64              `json:"totalSize"`
}

// Collects all the channel for the broker to decide the message to send to, along with all the names
type AllChannels struct {
	Channels     map[string]*SubrChannel `json:"channels"`
	ChannelNames []string       `json:"channelNames"`
	Mu           sync.RWMutex   `json:"-"`
}

var Helper AllChannels




func CreateTopic(client *utils.NewClient, topic string) *SubrChannel {
	publisher := []*utils.NewClient{client}

	channel := &SubrChannel{
		Subscribers:   []*utils.NewClient{},
		Publisher:     publisher,
		Topic:         topic,
		PTopic: 	   map[string]string{},	
		SubscribeChan: make(chan string, 100),
		Wg:            sync.WaitGroup{},
		Mu:            sync.RWMutex{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastPublisher: nil,
		History:       []TopicHistory{},
		PublishCount:  0,
		TotalSize:     0,
	}

	return channel
}

func FindChannel(collection *AllChannels, topic string) (*SubrChannel, bool) {
	if collection.Channels == nil {
		collection.Channels = make(map[string]*SubrChannel)
	}
	

	channel := collection.Channels[topic]

	if channel != nil {
			return channel, true
	}
	

	return nil, false

}

func GetChannel(client *utils.NewClient, topic string) *SubrChannel {
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()

	channel, exist := FindChannel(&Helper, topic)

	if exist {

		return channel
	}

	channel = CreateTopic(client, topic)
	Helper.ChannelNames = append(Helper.ChannelNames, channel.Topic)
	Helper.Channels[channel.Topic] = channel

	return channel
}

func HandleSubscribers(client *utils.NewClient, topic string) *SubrChannel {
	channel := GetChannel(client, topic)

	channel.Mu.Lock()
	channel.UpdatedAt = time.Now()
	channel.Subscribers = append(channel.Subscribers, client)
	channel.Mu.Unlock()

	return channel
}

func DeliverMessage(sub *utils.NewClient, msg string, wg *sync.WaitGroup) {
	defer wg.Done()

	sub.Conn.Write([]byte(msg))

	sub.Mu.Lock()
	sub.LastUsed = time.Now()
	sub.Mu.Unlock()
}

func Broker(client *utils.NewClient, topic string, msg string) {
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()

	channel, exist := FindChannel(&Helper, topic)

	if !exist && channel == nil {
		client.Conn.Write([]byte("Channel on the topic does not exist yet\n"))
		return
	}

	historymsg := TopicHistory{
		Message:   msg,
		Size:      int64(len(msg)),
		CreatedAt: time.Now(),
	}

	channel.UpdatedAt = time.Now()
	channel.LastPublisher = client
	channel.PublishCount += 1
	channel.TotalSize += int64(len(msg))
	channel.History = append(channel.History, historymsg)

	subs := make([]*utils.NewClient, len(channel.Subscribers))
	pubs := make([]*utils.NewClient, len(channel.Publisher))
	copy(subs, channel.Subscribers)
	copy(pubs, channel.Publisher)


	channel.Publisher = append(channel.Publisher, client)

	
	for _, sub := range subs {
		if sub == client {
			continue
		}else{
			logger.InfoLog("Pubsub: delivering message to subscriber " + sub.Conn.RemoteAddr().String())
			if client.Initiator == utils.RedisSource {
				continue
			}

			channel.Wg.Add(1)
			go DeliverMessage(sub, "roxkv> " +msg+"\n", &channel.Wg)
		}	

	}
	channel.Wg.Wait()



}

func HandleUnsubscribes(client *utils.NewClient, topic string) {
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()

	channel, exist := FindChannel(&Helper, topic)

	if !exist && channel == nil {
		client.Conn.Write([]byte("Channel on the topic does not exist yet\n"))
		return
	}

	channel.UpdatedAt = time.Now()
	for idx, sub := range channel.Subscribers {
		if sub == client {
			channel.Subscribers = slices.Delete(channel.Subscribers, idx, idx+1)
			break
		}
	}

	for idx, pub := range channel.Publisher {
		if pub == client {
			channel.Publisher = slices.Delete(channel.Publisher, idx, idx+1)
			break
		}
	}

	client.Conn.Write([]byte("Sucessfully Unsubscribed to " + topic + ".\n"))
}

func GetTopics(client *utils.NewClient) []string {
	Helper.Mu.Lock()
	topics := make([]string, len(Helper.ChannelNames))
	copy(topics, Helper.ChannelNames)
	Helper.Mu.Unlock()

	return topics
}

func CloseChannel(client *utils.NewClient, topic string) {
	Helper.Mu.Lock()
	topics := make([]string, len(Helper.ChannelNames))
	copy(topics, Helper.ChannelNames)

	delete(Helper.Channels,topic)	

	for idx, channel := range topics {
		if channel == topic {
			Helper.ChannelNames = slices.Delete(Helper.ChannelNames, idx, idx+1)

			break
		}
	}

	
	
	Helper.Mu.Unlock()

	GetTopics(client)

}

func PublishToAllTopics(client *utils.NewClient, message string) string {

	if client.Role != string(utils.RoleSystem) && client.Role != string(utils.RoleAdmin) {
		logger.ErrorLog("Pubsub: client " + client.Role + " is not allowed to publish everywhere")
		return ""
	}
	logger.InfoLog(" " + client.Role + " : Broadcasted a message across all topics")
	topics := GetTopics(client)

	for _, topic := range topics {

		Broker(client, topic, message)

	}

	return "Published Sucessfully"
}

func GetClients(client *utils.NewClient, topic *SubrChannel, category string) []*utils.NewClient {

	var clients []*utils.NewClient

	switch category {
	case "sub":
		clients = append(clients, topic.Subscribers...)
	case "pub":
		clients = append(clients, topic.Publisher...)
	default:
		return clients
	}
	return clients
}

func GetAllMembers(client *utils.NewClient, subOrPub string) []*utils.NewClient {

	if client.Role != string(utils.RoleSystem) && client.Role != string(utils.RoleAdmin) {
		return []*utils.NewClient{}
	}
	logger.InfoLog(" " + client.Role + " : Invoked all the members of all the topics")

	Helper.Mu.Lock()
	topicsChannel := Helper.Channels
	Helper.Mu.Unlock()
	var clients []*utils.NewClient
	for _, channel := range topicsChannel {
		clients = GetClients(client, channel, subOrPub)
	}

	return clients
}
