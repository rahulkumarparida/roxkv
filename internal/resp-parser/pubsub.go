package redisparser

import (
	"fmt"
	"slices"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

// 127.0.0.1:6379(subscribed mode)> publish game name
// (error) ERR Can't execute 'publish': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context

func SubscribeHandler(args []string,client *utils.NewClient){
	if len(args) < 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'subscribe' command")
		client.Conn.Write([]byte(encode))
		return
	}

	client.Mode = utils.ModeSubsriber	
	var encode string
	for _, arg := range args {
		val := pubsub.HandleSubscribers(client,arg)
		if val {
			client.Mode.Topic = append(client.Mode.Topic, arg)			
			data := []any{"subscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(data)
			fmt.Println("Encoded Subscribe:",encode)
			client.Conn.Write([]byte(encode))
		}
	}


	fmt.Println("Client Mode:",client.Mode)


	
}


func PublisheHandler(args []string,client *utils.NewClient) string{

	if len(args) != 2  {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'publish' command")
		return encode
	}
		

	// if client.Mode.Name == utils.ModeSubsriber.Name{
	// 	encode , _ := EncodeSimpleError(" ERR Can't execute 'publish': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context")
	// 	return encode
	// }

	topicname := args[0]
	message := args[1]

	pubsub.Helper.Mu.Lock()
	topic ,exist := pubsub.FindChannel(&pubsub.Helper,topicname)
	pubsub.Helper.Mu.Unlock()

	if !exist {
		SubscribeHandler([]string{topicname},client)
		topic , _ = pubsub.FindChannel(&pubsub.Helper,topicname)
	}

	pubsub.Helper.Mu.Lock()
	defer pubsub.Helper.Mu.Unlock()

	historymsg := pubsub.TopicHistory{
			Message:   message,
			Size:      int64(len(message)),
			CreatedAt: time.Now(),
	}


	topic.UpdatedAt = time.Now()
	topic.LastPublisher = client
	topic.PublishCount += 1
	topic.TotalSize += int64(len(message))
	topic.History = append(topic.History, historymsg)

	subs := make([]*utils.NewClient, len(topic.Subscribers))
	copy(subs, topic.Subscribers)
	data := []string{"message",topicname,message}
			
	encode , _ := EncodeArray(data)
	for _, sub := range subs {
		if sub == client {
			continue
		}else{
			topic.Wg.Add(1)
			go pubsub.DeliverMessage(sub,encode,&topic.Wg)
		}

	}
	topic.Wg.Wait()

	encode , _ = EncodeInteger(int64(len(subs)))
	
	return encode
}


func UnsubscribeHandler(args []string,client *utils.NewClient){
	
	if len(args) < 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'unsubscribe' command")
		client.Conn.Write([]byte(encode))
		return
	}
	
	
	topics := pubsub.GetTopics(client)
	
	var encode string
	for idx, arg := range args {
		
		if slices.Contains(topics,arg) && slices.Contains(client.Mode.Topic,arg) {
			//remove from the client 
			client.Mu.Lock()
			client.Mode.Topic = slices.Delete(client.Mode.Topic,idx,idx+1)
			client.LastUsed = time.Now()
			client.Mu.Unlock()
			// remove from the Helper Channel if no subscribers left
			

			pubsub.Helper.Mu.Lock()
			pubsub.Helper.Channels = slices.Delete(pubsub.Helper.Channels,idx,idx+1)
			pubsub.Helper.ChannelNames =slices.Delete(pubsub.Helper.ChannelNames,idx,idx+1)
			pubsub.Helper.Mu.Unlock()
			
			data := []any{"unsubscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(data) 
			fmt.Println("Encoded Unsubscribe:",encode)
			client.Conn.Write([]byte(encode))

		}else{

			data := []any{"unsubscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(data) 
			client.Conn.Write([]byte(encode))
			continue
		}



	}

	fmt.Println("Client Mode:",client.Mode)

	if len(client.Mode.Topic) == 0{
		client.Mode = utils.ModeDefault
		return
	}
			



}
// Implement  the rest of the pubsub and test 