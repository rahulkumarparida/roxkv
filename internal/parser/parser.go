package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func Command(input string) bool {

	if strings.ToLower(input) == "set" {
		return true
	}
	if strings.ToLower(input) == "get" {
		return true
	}
	if strings.ToLower(input) == "del" {
		return true
	}
	if strings.ToLower(input) == "keys" {
		return true
	}

	return false
}

// Divide all the things from the string like cmd key valuesssss --ttl 23 sec
func Chunking(input string, ch string) (string, bool) {

	if ch == " " {
		return input, true
	}

	return input, false

}
func GetCommad(input string) (string, string) {
	var currVal string
	var cmdFound bool = false
	var cmd string
	var lastIdx int

	for idx, ch := range input {
		currVal += string(ch)

		if !cmdFound {
			Val, found := Chunking(currVal, string(ch))
			if found || len(input)-1 == idx {
				cmdFound = true
				cmd = Val
				currVal = ""
				lastIdx = idx
				break
			}
		}
	}

	fmt.Println("cmd:", cmd, "Len:", len(cmd), "\nIdx:", lastIdx)

	input = input[lastIdx+1:]
	fmt.Println("New:", input)
	return strings.ToLower(strings.TrimSpace(cmd)), strings.TrimSpace(input)
}

func flagExtractor(data string){
	
}


var singleArgCommandRegistry = map[string]func(store *store.MemoryAlloc,client *utils.NewClient) any{

		"keys":func(store *store.MemoryAlloc, client *utils.NewClient) any{
			keys := commands.KeysCommand(store)
			return keys
		},
		"save":func(store *store.MemoryAlloc, client *utils.NewClient) any{
			keys := commands.SaveCommand(store)
			return keys
		},
		"load":func(store *store.MemoryAlloc, client *utils.NewClient) any{
			keys := commands.LoaderCommand(store)
			return keys
		},
		"topics":func(store *store.MemoryAlloc, client *utils.NewClient) any{
			commands.TopicsCommand(client)
			return nil
		},
	
} 

var doubleArgCommandRegistry = map[string]func(store *store.MemoryAlloc,client *utils.NewClient,data string)any{
	"get":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {
		if len(data) == 0{
			return ""
		}
		val := store.GetKv(stre,data)
		
		return val
	},
	"del":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {
		if len(data) == 0{
			return ""
		}
		val := store.DelKv(stre,data)
		
		return val
	},
	"history":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {	
		flag := "--all"
		if len(data) != 0 {
			flag = data
		}

		val , err:= logger.ReadFromFile(flag)
		if utils.HandleError("Error while reading from file",err) {
			return "\n"
		}
		return val
		
	},
	"subscribe":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {
		if len(data) == 0{
			return ""
		}
		val := pubsub.HandleSubscribers(client,data)

		return val
	},
	"unsubscribe":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {
		if len(data) == 0{
			return ""
		}
		pubsub.HandleUnsubscribes(client,data)

		return ""
	
	},
	"closechannel":func(stre *store.MemoryAlloc, client *utils.NewClient, data string) any {
		if len(data) == 0{
			client.Conn.Write([]byte("Topic name required\n"))
			return ""
		}
		pubsub.CloseChannel(client,data)
		return ""
	},

}

var tripleArgCommandRegistry = map[string]func (store *store.MemoryAlloc,client *utils.NewClient,data string) any{
	"publish":func(store *store.MemoryAlloc, client *utils.NewClient, data string) any{
		arg1 , msg  := GetCommad(data)
		if len(arg1) != 0 || len(msg) != 0 {
				return "specify clearly"
		}

		pubsub.Broker(client,arg1,msg)

		return ""
	},
}

var setSlice []string = []string{"set"}
var multiValueArgCommandRegistry = map[string]func(store *store.MemoryAlloc,client *utils.NewClient,data string) any{

	"set":func(store *store.MemoryAlloc, client *utils.NewClient, data string) any {
		return nil
	},

}

func Parser(store *store.MemoryAlloc,input []string , client *utils.NewClient) any{

	

	rl, err := readline.New("roxkv> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	line, err := rl.Readline()
	if err != nil { // io.EOF
		log.Fatal(err)
	}

	cmd,data := GetCommad(line)

	
		// Single args command executions
		// keys
		// save
		// load
		// topics
		fmt.Println("cmd:", cmd)
		if fn,exists := singleArgCommandRegistry[cmd]; exists {
				data := fn(store,client)
				return data
		}
	
		// Two args command executions
		// get key
		// del key
		// unsubcribe topicname
		// subscribe topicname
		// closechannel topicname
		// history --all||--info||--error||--success
		fmt.Println("cmd:", cmd)	
		if fn,exists := doubleArgCommandRegistry[cmd]; exists {
			data := fn(store,client,data)
			return data
		}


	
		// Three args command executions
		// publish topicname values
		fmt.Println("three:", cmd)
		if fn,exists := tripleArgCommandRegistry[cmd]; exists {
			data := fn(store,client,data)
			return data
		}
	

	if slices.Contains(setSlice, cmd) {
		// set args command executions
		// set key value
		// set key value --ttl 10 min
		fmt.Println("cmd:", cmd)
	}


	return false

}


