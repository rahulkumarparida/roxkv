package redisparser

import (
	"fmt"
	"slices"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
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
	
	var encode string
	for _, arg := range args {
		val := pubsub.HandleSubscribers(client,arg)
		if val != nil {
			client.Mode.Topic = append(client.Mode.Topic, arg)			
			data := []any{"subscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(data)
			logger.InfoLog("SubscribeHandler: client " + fmt.Sprintf("%v", client.ID) + " subscribed to " + arg)
			client.Conn.Write([]byte(encode))
		}
	}


	logger.InfoLog("SubscribeHandler: client " + fmt.Sprintf("%v", client.ID) + " mode=" + client.Mode.Name)


	
}


func PublishHandler(args []string,client *utils.NewClient) string{

	if len(args) != 2  {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'publish' command")
		return encode
	}
		

	topicname := args[0]
	message := args[1]
	
	pubsub.Helper.Mu.Lock()
	topic ,exist := pubsub.FindChannel(&pubsub.Helper,topicname)
	pubsub.Helper.Mu.Unlock()	

	if !exist {
		val := pubsub.HandleSubscribers(client,topicname)
		if val != nil {
			client.Mode.Topic = append(client.Mode.Topic, topicname)	
		}
		topic = val	

		



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
	


	
	for _, sub := range subs {
		if sub != client {
			data := []string{"message",topicname,message}
			encode , _ := EncodeArray(data)
			topic.Wg.Add(1)
			go pubsub.DeliverMessage(sub,encode,&topic.Wg)
		}
	topic.Wg.Wait()
	}

	// for psubscribes
	var psubs = 0
	ptopics,exists:=PtopicChecker(topic)
	if exists{
		for _, ptopic := range ptopics {
			pChannel := PatternRegister[ptopic]
			for _, sub := range pChannel.Clients{
				psubs += 1
				data := []string{"pmessage",ptopic,topicname,message}
				encode , _ := EncodeArray(data)
				topic.Wg.Add(1)
				go pubsub.DeliverMessage(sub,encode,&topic.Wg)
			}
			
		}
	}



	encode , _ := EncodeInteger(int64(len(subs)+psubs))
	
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
	for _, arg := range args {
		
		if slices.Contains(topics,arg) && slices.Contains(client.Mode.Topic,arg) {
			//remove from the client 
			for idx, topic := range client.Mode.Topic {
				
				if topic == arg {
					client.Mu.Lock()
			        client.Mode.Topic = slices.Delete(client.Mode.Topic,idx,idx+1)
			        client.LastUsed = time.Now()
			        client.Mu.Unlock()
				}
			}
            pubsub.Helper.Mu.Lock()
			data := pubsub.Helper.Channels[arg] 
            pubsub.Helper.Mu.Unlock()

			for idx, sub := range data.Subscribers {
				if sub == client {

				pubsub.Helper.Mu.Lock()
                   pubsub.Helper.Channels[arg].Subscribers = slices.Delete(pubsub.Helper.Channels[arg].Subscribers,idx,idx+1)    
				pubsub.Helper.Mu.Unlock()
				}
			}
			
			// remove from the Helper Channel if no subscribers left delete it

		
			
			endata := []any{"unsubscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(endata) 
			logger.InfoLog("UnsubscribeHandler: client " + fmt.Sprintf("%v", client.ID) + " unsubscribed from " + arg)
			client.Conn.Write([]byte(encode))

		}else{

			data := []any{"unsubscribe",arg,len(client.Mode.Topic)}
			encode , _ = EncodeArray(data) 
			client.Conn.Write([]byte(encode))
			continue
		}



	}

	logger.InfoLog("UnsubscribeHandler: client " + fmt.Sprintf("%v", client.ID) + " mode=" + client.Mode.Name)

	if len(client.Mode.Topic) == 0{
		client.Mode = utils.ModeDefault
		return
	}
			

}



// Implement  the  pubsub and test only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE are allowed in this context
// Psubcribe subscribe to all the channels matching as many chars after * but should include the rest of the string and ? a single charachter is the ? is here
func PSubscribeHandler(args []string,client *utils.NewClient){
	if len(args) < 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'psubscribe' command")
		client.Conn.Write([]byte(encode))
		return
	}


	for _, ptopic := range args {

	
		PatternRegister[ptopic] = PatternRegistry{
			Pattern: ptopic,
			Clients: append(PatternRegister[ptopic].Clients, client),
		}


		client.Mode.PTopic = append(client.Mode.PTopic, ptopic)
		data := []any{"psubscribe",ptopic,len(client.Mode.PTopic)}
		encode , _ := EncodeArray(data)
		logger.InfoLog("PSubscribeHandler: client " + fmt.Sprintf("%v", client.ID) + " psubscribed to " + ptopic)
		client.Conn.Write([]byte(encode))
	}



}


func PUnsubscribeHandler(args []string, client *utils.NewClient) {

	//  registry remove and check if the topic has no more client then delete the registry

	if len(args) < 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'punsubscribe' command")
		client.Conn.Write([]byte(encode))
		return
	}


	for _, arg := range args {

		data , exist := PatternRegister[arg]

		if !exist {
			endata := []any{"unsubscribe",arg,0}
			encode , _ := EncodeArray(endata)
			client.Conn.Write([]byte(encode))
			return
		}
		

		client.Mu.Lock()
		for idx, ptopic := range client.Mode.PTopic {
			if arg == ptopic {
				ptopics := slices.Delete(client.Mode.PTopic,idx,idx+1)
				client.Mode.PTopic = ptopics
				break
			}
		}
		client.LastUsed = time.Now()
		client.Interactions += 1
		if len(client.Mode.PTopic) == 0 && len(client.Mode.Topic) == 0{
			client.Mode = utils.ModeDefault
		}
		client.Mu.Unlock()
        

		if len(data.Clients) <= 1 && data.Clients[0] == client{
			
			delete(PatternRegister,arg)
		
		}else{
			
			for idx, sub := range data.Clients {

				if sub == client {
					data.Clients = slices.Delete(data.Clients,idx,idx+1)
					break
				}
				
			}
			PatternRegister[arg] = data
		
		}
         
		endata := []any{"punsubscribe",arg,len(client.Mode.PTopic)}
		encode , _ := EncodeArray(endata)
		client.Conn.Write([]byte(encode))
		
	}



}