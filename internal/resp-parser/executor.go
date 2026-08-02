package redisparser

import (
	// "bytes"
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var mutex = sync.Mutex{}


func RemoveClientsFromTopics(client *utils.NewClient, topic string){
	pubsub.Helper.Mu.Lock()
	defer pubsub.Helper.Mu.Unlock()

	channel, exist := pubsub.FindChannel(&pubsub.Helper, topic)

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

}

func RemoveDeadClientData(client *utils.NewClient){
	
	// remove from TotalConnections
	mutex.Lock()
		utils.TotalConnecntions = slices.DeleteFunc(utils.TotalConnecntions, func(n *utils.NewClient) bool {
			return n.ID == client.ID
		})
	mutex.Unlock()


	client.Mu.Lock()
	topics := client.Mode.Topic
	client.Mu.Unlock()


	for _, topic := range topics{
		RemoveClientsFromTopics(client,topic)		
	}
	
}

// Convert this []string to receiev [][]byte so it should be binary safe trhough out the process without conversion to string add comparison such that runes are compared instead of converting it to string



func ExecuteSet(args *RedisInput, client *utils.NewClient) any{
	var stre = store.StoreHelper		
	if  len(args.Args) < 2 {
		encode , _:=  EncodeSimpleError("ERR wrong number of arguments for 'set' command")
		return  encode
	}
	

	val := commands.SetCommand(client,stre,args.RawArgs)

	if val {
		encode,_ := EncodeSimpleString("OK")	
		return  encode	
	}else{
		encode , _ := EncodeSimpleError(" ERR wrong number of arguments for 'set' command")
		return encode

	}

}

func ExecuteGet(args []string) []byte{
	var stre = store.StoreHelper		

	if  len(args) < 1 {
		 encode ,_ := EncodeSimpleError("ERR wrong number of arguments for 'get' command")
		return  []byte(encode)
	}

	val := commands.GetCommand(stre,args)

	if val == nil {
		encode := EncodeNullValues()
		return []byte(encode)
	}

	// val is now []byte from store
	valBytes, ok := val.([]byte)
	if !ok {
		return []byte(EncodeNullValues())
	}
	return EncodeBulkBytes(valBytes)
}




func Executekeys(args []string) string{
	var stre = store.StoreHelper
	if len(args) > 1 || len(args) < 1{
		encode ,_  := EncodeSimpleError("(ERR wrong number of arguments for 'keys' command")
		
		return encode
	}
	argument := args[0]
	values := store.KeyKv(stre)	

	if len(values) <= 0 {
		return "*0\r\n"
	}

	if len(argument) == 1 {
		encode , _ := EncodeArray(values)
		return encode
	}else{
		req := argument[1:]
		var filteredKeys []string
		for _, v := range values {
			if strings.Contains(v,req) {
				filteredKeys = append(filteredKeys, v)
			}
		}
		encode , _ := EncodeArray(filteredKeys)
		return encode
	}
}

func ExecuteDel(args []string) string{
	stre := store.StoreHelper
	
	var count int64 = 0 
	for _, key := range args {
		val := store.DelKv(stre,key)
		if val {
			count += 1
		}

	}
	encode , _ := EncodeInteger(count)

	return encode
}


func ExecuteSave() string{
	var stre = store.StoreHelper

	commands.SaveCommand(stre)

	encode := "+OK\r\n"
	return  encode
}


func ExecuteExists(args []string) string{
	stre := store.StoreHelper
	
	var count int64 = 0 
	for _, key := range args {
		val := store.GetKv(stre,key)
		if val.Key != "" {
			count += 1
		}
	}
	encode , _ := EncodeInteger(count)
	return encode
}


func ExecuteFlushDB() string{

	var stre = store.StoreHelper
	store.FlushKv(stre)


	return "+OK\r\n"
}

func ExecutePing(args []string) string{
	
	if len(args) >= 2 {
		encode ,_  := EncodeSimpleError("ERR wrong number of arguments for 'ping' command")
		
		return encode
	}

	if len(args) == 1{
		encode , _ := EncodeBulkString(args[0])
		return encode
	}

	encode := "+PONG\r\n"
	return encode
}

func ExecuteExpire(args [][]byte,client *utils.NewClient) string{
	stre := store.StoreHelper
	
	if len(args) < 2 || len(args) > 2 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'expire' command")
		return encode
	}

	key := string(args[0])

	keyData := store.GetKv(stre,key)

	if keyData.Val == nil {
		encode , _ := EncodeInteger(0)
		return encode
	}

	value := keyData.Val
	var byteData [][]byte
	fmtData := []string{"--ttl",string(args[1]),"sec",keyData.Key}
	for _, bd := range fmtData {
		byteData = append(byteData, []byte(bd))
	}
	byteData = append(byteData, value)


	executed := commands.SetCommand(client , stre ,byteData)
	var encode string
	if executed {
		encode , _ = EncodeInteger(1)
	}else{
		encode , _ = EncodeInteger(0)
	}
	return encode
}


func ExecuteHistory(flag []string) []string{

	data := commands.HistoryCommand(flag)

	arrOfData := strings.Split(data,"\n")
	return arrOfData

}


func ExecuteEcho(args []string) string{
	data := strings.Join(args," ")
	encode , _ := EncodeBulkString(data)
	return encode
} 

func ExecuteTTL(args []string) string{
	stre := store.StoreHelper
	

	if len(args) > 1 || len(args) == 0 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for the 'TTL' command")
		return encode
	}

	keyData := store.GetKv(stre, args[0])

	if keyData.Val == nil {
		encode , _ := EncodeSimpleError("ERR key does not exists")
		return encode
	}

	if keyData.Meta.TTL.IsZero() {
		encode , _ := EncodeInteger(-2)
		return encode
	}

	timLeft := time.Until(keyData.Meta.TTL)
	encode , _ := EncodeInteger(int64(timLeft.Seconds()))
	return encode
}

// removes the expiry from the key
func ExecutePersistance(args []string, client *utils.NewClient) string{
stre := store.StoreHelper
	

	if len(args) > 1 || len(args) == 0 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for the 'TTL' command")
		return encode
	}

	keyData := store.GetKv(stre, args[0])

	if keyData.Val == nil {
		encode , _ := EncodeInteger(0)
		return encode
	}

	if keyData.Meta.TTL.IsZero() {
		encode , _ := EncodeInteger(-1)
		return encode
	}

	var databyte = [][]byte{[]byte(keyData.Key),keyData.Val}

	executed := commands.SetCommand(client,stre,databyte)
	var encode string
	if executed {
		encode , _ = EncodeInteger(1)
		
	}else{
		encode , _ = EncodeSimpleError("ERR something occured while setting up the key")
	}
	return encode

}


func ExecuteDbSize() string{
	var stre = store.StoreHelper

	data := store.KeyKv(stre)
	encode , _ := EncodeInteger(int64(len(data)))
	return encode
}


func ExecuteRandomKeys() string{
	
	var stre = store.StoreHelper
	data := store.KeyKv(stre)
	randomIndex := rand.IntN(len(data))
	randomwKey := data[randomIndex]
	encode , _ := EncodeBulkString(randomwKey)

	return encode
}


func ExecuteRenameKey(args []string, client *utils.NewClient) string{
	if len(args) > 2 || len(args) <= 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'rename' command")
		return encode
	}
	
	var stre = store.StoreHelper

	data := store.GetKv(stre,args[0])

	if data.Val == nil {
		encode , _ := EncodeSimpleError("ERR no such key found")
		return encode	
	}


	executed := commands.SetCommand(client,stre,[][]byte{[]byte(args[1]),data.Val})

	store.DelKv(stre,args[0])

	var encode string
	if executed {
		encode , _ = EncodeSimpleString("OK")
	}else{
		encode,_ = EncodeSimpleError("ERR parsing the data")
	}
		
	return encode
}



func ExecuteServerUpTime(args []string)string{
	if len(args) > 1 || len(args) <= 0 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'rename' command")
		return encode
	}

	starttime := utils.ServerStarted

	uptimeseconds := strconv.FormatInt(int64(time.Since(starttime).Seconds()),10)
	uptimemilliseconds := strconv.FormatInt(int64(time.Since(starttime).Milliseconds()),10)
	


	encode, _ := EncodeArray([]any{uptimeseconds,uptimemilliseconds})
	return encode
}



func ExecuteMget(args []string) string{
	if  len(args) <= 0 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'rename' command")
		return encode
	}
	var stre = store.StoreHelper	

	var valueArray []any

	for _, key := range args {
			
		data := store.GetKv(stre,key)

		if data.Val == nil {
			valueArray = append(valueArray, nil)
			continue
		}
		valueArray = append(valueArray, data.Val)	
	}

	encode , _ := EncodeArray(valueArray)

	return encode
}


func ExecuteIntOpration(cmd string,args []string, client *utils.NewClient) string{

	
	if len(args) < 1 || len(args) > 3 || (cmd == "incr" && len(args) >=2  )||(cmd == "decr" && len(args) >=2  ){
		encode ,_ := EncodeSimpleError("ERR wrong  number of arguments for '"+cmd+"' command")
		return encode
	}
	
	var stre = store.StoreHelper
	key := args[0]


	getData := store.GetKv(stre,key)
	var encode string
	if getData.Val == nil {
		encode ,_ := EncodeSimpleError("ERR no such key found")
		return encode
	}
	value := string(getData.Val)



	

	i64 , err := strconv.ParseInt(value,10,64)

	if err != nil {
		encode ,_ := EncodeSimpleError("ERR value is not an integer or out of range")
		return encode		
	}
	

	switch strings.ToLower(cmd) {
	case "incr":

		newval := i64+1
		key := getData.Key
		data := [][]byte{[]byte(key),[]byte(strconv.FormatInt(newval,10))}
		commands.SetCommand(client,stre,data)
		encode , _= EncodeInteger(newval)
	case "decr":
		newval := i64-1
		key := getData.Key
		data := [][]byte{[]byte(key),[]byte(strconv.FormatInt(newval,10))}
		commands.SetCommand(client,stre,data)
		encode , _= EncodeInteger(newval)

	case "decrby":
		argInt , err := strconv.ParseInt(args[1],10,64)
		if err != nil {
			encode ,_ := EncodeSimpleError("ERR value is not an integer or out of range")
			return encode	
		}
		newval := i64-argInt
		key := getData.Key
		data := [][]byte{[]byte(key),[]byte(strconv.FormatInt(newval,10))}
		commands.SetCommand(client,stre,data)
		encode , _= EncodeInteger(newval)

	case "incrby":
		argInt , err := strconv.ParseInt(args[1],10,64)
		if err != nil {
			encode ,_ := EncodeSimpleError("ERR value is not an integer or out of range")
			return encode	
		}
		newval := i64+argInt
		key := getData.Key
		data := [][]byte{[]byte(key),[]byte(strconv.FormatInt(newval,10))}
		commands.SetCommand(client,stre,data)
		encode , _= EncodeInteger(newval)
		

	default:
		encode,_ = EncodeSimpleError("ERR no such commands found")
	}


	return encode
}


func ExecuteStrLen(args []string)string{
	if len(args) > 1 || len(args) <= 0 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'strlen' command")
		return encode
	}
	
	var stre = store.StoreHelper

	data := store.GetKv(stre,args[0])

	if data.Val == nil {
		encode , _ := EncodeSimpleError("ERR no such key found")
		return encode	
	}
	
	length := data.Meta.Size

	encode , _ := EncodeInteger(length)
	return encode
}




func ExecuteQuit(client *utils.NewClient){
	RemoveDeadClientData(client)

	client.Conn.Write([]byte("+OK\r\n"))
	client.Conn.Close()
}

func ExecuteReset(client *utils.NewClient) string{
	client.Mu.Lock()
	topics := client.Mode.Topic
	client.Mu.Unlock()


	for _, topic := range topics{
		RemoveClientsFromTopics(client,topic)
	}

	client.Mode = utils.ModeDefault
	return "+RESET\r\n"

}



func ExecuteLastSave() string{

	timsstamp := persistence.LastSavedFileTimeStamp()

	encode ,_ := EncodeInteger(timsstamp.Unix())

	return encode

}


func ExecuteClientname(args []string,client *utils.NewClient) string{
	if strings.ToLower(args[0]) == "getname"{
		if len(args) < 1 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'client|getname' command")
			return encode 	
		}
		fmt.Println("Asked Client:", client.ID, " Response: ", client.Name)
		if client.Name == "" {
			encode := EncodeNullValues()
			return encode
		}
		encode , _ := EncodeBulkString(client.Name)
		return  encode
	}else if strings.ToLower(args[0]) == "setname" {
		if len(args) < 2 {
		encode , _ := EncodeSimpleError("ERR wrong number of argguments for 'client|setname' command")
		return encode 	
		}
		client.Mu.Lock()
		client.Name = args[1]
		client.Mu.Unlock()
		encode , _ := EncodeSimpleString("OK")
		return encode	
	}else if len(args) == 3 && strings.ToLower(args[0]) == "setinfo"{
		if strings.ToLower(args[1]) == "lib-name" {
			client.Mu.Lock()
			client.Library_Name = args[2]
			client.Mu.Unlock()
			
			return "+OK\r\n"
			
		}else if strings.ToLower(args[1]) == "lib-ver"{
			client.Mu.Lock()
			client.Library_Ver = args[2]
			client.Mu.Unlock()
			return "+OK\r\n"
		}
	}
	encode := EncodeNullValues()
	return encode
}

func ExecuteCommand() string {
	// Properly formatted RESP Array of 3 sub-arrays
	return "*3\r\n*6\r\n$4\r\nping\r\n:2\r\n*1\r\n$7\r\nstaleok\r\n:0\r\n:0\r\n:0\r\n*6\r\n$6\r\nclient\r\n:-2\r\n*1\r\n$5\r\nadmin\r\n:0\r\n:0\r\n:0\r\n*6\r\n$7\r\ncommand\r\n:-1\r\n*2\r\n$7\r\nloading\r\n$7\r\nstaleok\r\n:0\r\n:0\r\n:0\r\n"
}

func ExecuteCommandDocs() string {
	// A map/array payload containing 3 key-value command definitions
	return "*6\r\n$4\r\nping\r\n*2\r\n$7\r\nsummary\r\n$21\r\nReturns PONG if alive\r\n$6\r\nclient\r\n*2\r\n$7\r\nsummary\r\n$21\r\nConnection management\r\n$7\r\ncommand\r\n*2\r\n$7\r\nsummary\r\n$22\r\nReturns command matrix\r\n"
}
