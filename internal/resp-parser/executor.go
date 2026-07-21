package redisparser

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)




func ExecuteSet(args []string, client *utils.NewClient) any{
	var stre = store.StoreHelper		
	if  len(args) < 2 {
		encode , _:=  EncodeSimpleError("ERR wrong number of arguments for 'set' command")
		return  encode
	}

	val := commands.SetCommand(client,stre,args)

	fmt.Println("Value :", val)
	if val {
		encode,_ := EncodeSimpleString("OK")
		fmt.Println("returning encoded string:", encode)	
		return  encode	
	}else{
		encode , _ := EncodeSimpleError(" ERR wrong number of arguments for 'set' command")
		return encode

	}

}

func ExecuteGet(args []string) string{
	var stre = store.StoreHelper		

	if  len(args) < 1 {
		 encode ,_ := EncodeSimpleError("ERR wrong number of arguments for 'get' command")
		return  encode
	}


	val := commands.GetCommand(stre,args)


	if val == nil {
		encode := EncodeNullValues()
		return encode
	}

	var encode string
	switch v := val.(type) {
		case string:
			encode, _ = EncodeBulkString(v)
		case []string:
			data := strings.Join(v ," ")
			encode , _ = EncodeBulkString(data)
		default:
			encode = ""
	}

	return encode
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

func ExecuteExpire(args []string,client *utils.NewClient) string{
	stre := store.StoreHelper
	
	if len(args) < 2 || len(args) > 2 {
		encode , _ := EncodeSimpleError("ERR wrong number of arguments for 'expire' command")
		return encode
	}

	key := args[0]

	keyData := store.GetKv(stre,key)

	if keyData.Val == nil {
		encode , _ := EncodeInteger(0)
		return encode
	}

	value := strings.Join(keyData.Val.([]string)," ")	
	
	fmtData := []string{"--ttl",args[1],"sec",keyData.Key,value}

	executed := commands.SetCommand(client , stre ,fmtData)
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

	kv := []string{keyData.Key, strings.Join(keyData.Val.([]string), " ")}

	executed := commands.SetCommand(client,stre,kv)
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
	var value string
	switch v := data.Val.(type) {
		case string:
			value = string(data.Val.(string))
		case []string:
			data := strings.Join(v ," ")
			value = data
		default:
			value = data.Val.(string)
	}


	executed := commands.SetCommand(client,stre,[]string{args[1],value})

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
			valueArray = append(valueArray, data.Val)
			continue
		}
		values := fmt.Sprintf("%v", data.Val)
		
		valueArray = append(valueArray, values)	
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
	var value string
	switch v := getData.Val.(type) {
		case string:
			value = string(getData.Val.(string))
		case []string:
			data := strings.Join(v ," ")
			value = data
		default:
			value = getData.Val.(string)
	}



	

	i64 , err := strconv.ParseInt(value,10,64)

	if err != nil {
		encode ,_ := EncodeSimpleError("ERR value is not an integer or out of range")
		return encode		
	}
	

	switch strings.ToLower(cmd) {
	case "incr":

		newval := i64+1
		key := getData.Key
		data := []string{key,strconv.FormatInt(newval,10)}
		commands.SetCommand(client,stre,data)
		encode , _= EncodeInteger(newval)
	case "decr":
		newval := i64-1
		key := getData.Key
		data := []string{key,strconv.FormatInt(newval,10)}
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
		data := []string{key,strconv.FormatInt(newval,10)}
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
		data := []string{key,strconv.FormatInt(newval,10)}
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



func ExecutePubSub(cmd string,args []string, client *utils.NewClient){
	if strings.ToLower(cmd) == "publish" {
		PublisheHandler(args,client)	
	}
	if strings.ToLower(cmd) == "subscribe" {
		SubscribeHandler(args,client)
	}
	if strings.ToLower(cmd) == "unsubscribe" {
		UnsubscribeHandler(args,client)
	}
}