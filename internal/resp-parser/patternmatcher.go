package redisparser

import (
	"fmt"
	"strings"

	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)


type PatternRegistry struct{
	Pattern string
	Clients []*utils.NewClient

}
var PatternRegister map[string]PatternRegistry

//  Takes the Published channelname and returns all the pattern matched to the channel
func PtopicChecker(topic *pubsub.SubrChannel) ([]string,bool){
	
	var data []string		
		

	for _, item := range PatternRegister{
		if count , _ :=PatternChecker(item.Pattern,[]string{topic.Topic}); count >= 1{
			fmt.Println("Matched: ", item.Pattern)
			data = append(data, item.Pattern)
		}
	}

	
	return data, true
}


func PatternChecker(input string, topics []string) (int,[]string){
	


	if strings.Contains(input,"*"){
		return AsteriskPattern(input,topics)
	}


	return  0,nil
}


func AsteriskPattern(input string , topics []string) (int,[]string){
	// all the patterns of channel that have similar close to chars input
	var similarChannels []string
	var inputTopics []string
	var separator string
	if strings.Contains(input,"."){
		separator = "."
		inputTopics = strings.Split(input,".")	
	}else if strings.Contains(input,":"){
		separator = ":"
		inputTopics = strings.Split(input,":")
	}else {
		// IF only * then star* , start*start , match*ng
		idx := strings.Index(input,"*")
		length := len(input)
		if idx+1 == length {
			inputTopics = append(inputTopics, input[:idx])
			inputTopics = append(inputTopics , input[idx+1:])
		}
	}
	var check bool = true
	for _, topic := range topics {
		topicArr := strings.Split(topic,separator)
		if inputTopics[0] == "*" || inputTopics[0] ==  topicArr[0]{

			for idx,subtopic := range topicArr{

					if  inputTopics[idx] == "*"{
						continue
					}else if inputTopics[idx] != subtopic{
						check = false
						break
					}
			}
			
			if check {
				similarChannels = append(similarChannels, topic)
			}else {
				check = true
			}

		} 
	}


	return len(similarChannels),similarChannels

}

