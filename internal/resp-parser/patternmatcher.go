package redisparser

import (
	"strings"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)


type PatternRegistry struct {
	Pattern string // hello.* now.? g[^e]me
	ListOfAvaliableChannels []string
	clients []utils.NewClient
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