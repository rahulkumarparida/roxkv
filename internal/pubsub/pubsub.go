package pubsub

import (
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)

// Single Channel and will contain many to many relationship with Subs and Pubs
type SubrChannel struct {
	Subscribers []*utils.NewClient
	Publisher []*utils.NewClient
	Topic             string
	SubscribeChan chan string
	Wg sync.WaitGroup
	Mu sync.RWMutex
	CreatedAt time.Time
	UpdatedAt time.Time
	LastPublisher *utils.NewClient
	PublishCount int
}



// Collects all the channel for the broker to decide the message to send to, along with all the names
type AllChannels struct{
	Channels []*SubrChannel
	ChannelNames []string
	Mu sync.RWMutex
}

var Helper AllChannels


func CreateTopic(client *utils.NewClient,topic string) *SubrChannel{
	publisher := []*utils.NewClient{client}

	channel :=  &SubrChannel{
		Subscribers: []*utils.NewClient{},
		Publisher: publisher,
		Topic: topic,
		SubscribeChan: make(chan string,100),
		Wg: sync.WaitGroup{},
		Mu: sync.RWMutex{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		LastPublisher: nil,
		PublishCount: 1,
	}


	return channel
}

func FindChannel(collection *AllChannels,topic string) (*SubrChannel,bool){


	for _, channel := range  collection.Channels{
		if channel.Topic==topic  {
			return channel , true
		}
	}	

	return nil,false

} 



func GetChannel(client *utils.NewClient,topic string) *SubrChannel{
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()	

	channel, exist := FindChannel(&Helper,topic)

	if exist {
		
		return channel
	}
	
	channel = CreateTopic(client,topic)
	Helper.ChannelNames = append(Helper.ChannelNames, channel.Topic)
	Helper.Channels = append(Helper.Channels, channel)
		
	return channel
}


func HandleSubscribers(client *utils.NewClient,topic string) bool{
	channel := GetChannel(client,topic)

	channel.Mu.Lock()
	channel.UpdatedAt = time.Now()
	channel.Subscribers = append(channel.Subscribers, client)
	channel.Mu.Unlock()

	return true
}

func DeliverMessage(sub *utils.NewClient, msg string,wg *sync.WaitGroup){
		defer wg.Done()
		
		sub.Conn.Write([]byte("roxkv> "+msg+"\n"))

		sub.Mu.Lock()
		sub.LastUsed = time.Now()
		sub.Mu.Unlock()
}

func Broker(client *utils.NewClient,topic string , msg string) {
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()

	channel , exist := FindChannel(&Helper,topic)

	if !exist && channel == nil {
		client.Conn.Write([]byte("Channel on the topic does not exist yet\n"))
		return
	}

	channel.UpdatedAt = time.Now()
	channel.LastPublisher = client
	channel.PublishCount += 1
	


	subs := make([]*utils.NewClient, len(channel.Subscribers))
	pubs := make([]*utils.NewClient,len(channel.Publisher))
	copy(subs,channel.Subscribers)
	copy(pubs,channel.Publisher)

	var IsPublisher bool = false
	for _, pub := range pubs {
		if client == pub{
			IsPublisher=true
		}
	}
	
	if IsPublisher {
		for _, sub := range subs {
		
		channel.Wg.Add(1)
		go DeliverMessage(sub,msg,&channel.Wg)

		}
		channel.Wg.Wait()
		
	}else{
		client.Conn.Write([]byte("Only Publisher can publish on the channel\n"))
		return

	}
	

}


func HandleUnsubscribes(client *utils.NewClient, topic string)  {
	Helper.Mu.Lock()
	defer Helper.Mu.Unlock()

	channel , exist := FindChannel(&Helper,topic)

	if !exist && channel == nil {
		client.Conn.Write([]byte("Channel on the topic does not exist yet\n"))
		return
	}

	channel.UpdatedAt = time.Now()
	for idx , sub := range channel.Subscribers {
		if sub == client {
			channel.Subscribers = slices.Delete(channel.Subscribers,idx,idx+1)	
			break
		}
	}					

	for idx , pub := range channel.Publisher {
		if pub == client {
			channel.Publisher = slices.Delete(channel.Publisher,idx,idx+1)	
			break
		}
	}

	client.Conn.Write([]byte("Sucessfully Unsubscribed to "+topic+".\n"))
}


func GetTopics(client *utils.NewClient) []string{
	Helper.Mu.Lock()
	topics := make([]string,len(Helper.ChannelNames))
	copy(topics,Helper.ChannelNames)
	Helper.Mu.Unlock()

	for idx, topic := range topics {
		client.Conn.Write([]byte(strconv.Itoa(idx)+". "+topic+"\n"))
	}
	return topics
}

func CloseChannel(client *utils.NewClient,topic string){
	Helper.Mu.Lock()
	topics := make([]string,len(Helper.ChannelNames))
	copy(topics,Helper.ChannelNames)

	for idx, channel := range topics {
		if channel == topic {
			Helper.ChannelNames = slices.Delete(Helper.ChannelNames,idx,idx+1)

			break
		}
	}

	topicsChannel := make([]*SubrChannel,len(Helper.Channels))
	copy(topicsChannel,Helper.Channels)

	for idx, channel := range topicsChannel {
		if channel.Topic == topic {
			Helper.Channels = slices.Delete(Helper.Channels,idx,idx+1)
			break
		}
	}
	Helper.Mu.Unlock()	

	GetTopics(client)

}