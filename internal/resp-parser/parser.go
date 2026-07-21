package redisparser

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)


type RedisInput struct{
	Cmd string
	Args []string
}



func ReadAndHandleConnection(client *utils.NewClient) {

	var buf []byte =make([]byte, 512)

	n, err := client.Conn.Read(buf[:])

	if err != nil{
		return 
	}

	tokens , err := DecodeArrayString(buf[:n])

	var input *RedisInput

	if len(tokens) >= 2 {

		input =  &RedisInput{
			Cmd: tokens[0],
			Args: tokens[1:],
		}	
	}else {
		input = &RedisInput{
			Cmd: tokens[0],
			Args: nil,
		}
	}

	
	parseCommand(input, client)

}


func parseCommand(input *RedisInput, client *utils.NewClient) {


	subscriberModeCommands :=[]string{"subscribe","ping","unsubscribe","psubscribe","punsubscribe","quit","reset"}

	if client.Mode.Name == utils.ModeSubsriber.Name{
		if !slices.Contains(subscriberModeCommands,strings.ToLower(input.Cmd)) {

			encode , _ := EncodeSimpleError(" ERR Can't execute 'publish': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context")
			client.Conn.Write([]byte(encode))
			return
		}


	}




	switch strings.ToLower(input.Cmd){
		case "ping":
			fmt.Println("Executing ping", input.Args)
			val := ExecutePing(input.Args)
			client.Conn.Write([]byte(val))
		case "set":
			fmt.Println("Executing set , args: ", input.Args)
			data :=ExecuteSet(input.Args , client)
			client.Conn.Write([]byte(data.(string)))
		case "get":
			fmt.Println("Executing get")
			data := ExecuteGet(input.Args)
			client.Conn.Write([]byte(data))
		case "del":
			fmt.Println("Executing delete")
			data := ExecuteDel(input.Args)
			client.Conn.Write([]byte(data))
		case "keys":
			fmt.Println("Executing Keys")
			data := Executekeys(input.Args)
			client.Conn.Write([]byte(data))
		case "exists":
			fmt.Println("Executing exists")
			data := ExecuteExists(input.Args)
			client.Conn.Write([]byte(data))
		case "save":
			fmt.Println("Executing save")
			val := ExecuteSave()
			client.Conn.Write([]byte(val))
		case "expire":
			fmt.Println("executing Expire")
			val := ExecuteExpire(input.Args,client)
			client.Conn.Write([]byte(val))
		case "flushdb":
			fmt.Println("Executing Flush")
			val := ExecuteFlushDB()
			client.Conn.Write([]byte(val))
		case "echo":
			fmt.Println("Executing echo")
			val := ExecuteEcho(input.Args)
			client.Conn.Write([]byte(val))
		case "ttl":
			fmt.Println("Executing ttl")
			val := ExecuteTTL(input.Args)
			client.Conn.Write([]byte(val))
		case "persist":
			fmt.Println("Executing persist")
			val := ExecutePersistance(input.Args, client)
			client.Conn.Write([]byte(val))
		case "dbsize":
			fmt.Println("Executing dbsize")
			val := ExecuteDbSize()
			client.Conn.Write([]byte(val))
		case "randomkey":
			fmt.Println("Executing randomkey")
			val := ExecuteRandomKeys()
			client.Conn.Write([]byte(val))
		case "rename":
			fmt.Println("Executing rename")
			val := ExecuteRenameKey(input.Args, client)
			client.Conn.Write([]byte(val))
		case "uptime":
			fmt.Println("Executing uptime")
			val := ExecuteServerUpTime(input.Args)
			client.Conn.Write([]byte(val))
		case "mget":
			fmt.Println("Executing mget")
			val := ExecuteMget(input.Args)
			client.Conn.Write([]byte(val))
		case "incr", "decr", "incrby", "decrby":
			fmt.Println("Executing integer operation")
			val := ExecuteIntOpration(input.Cmd, input.Args, client)
			client.Conn.Write([]byte(val))
		case "strlen":
			fmt.Println("Executing strlen")
			val := ExecuteStrLen(input.Args)
			client.Conn.Write([]byte(val))
		case "command","docs":
			fmt.Println("command")
			client.Conn.Write([]byte("*0\r\n"))
		case "subscribe","publish","unsubscribe":
			
			if input.Cmd =="subscribe" {
				SubscribeHandler(input.Args,client)
			}
			if input.Cmd =="unsubscribe" {
				UnsubscribeHandler(input.Args,client)
			}
			if input.Cmd =="publish"{
					
				data :=PublisheHandler(input.Args,client)
				client.Conn.Write([]byte(data))
			}

		case "monitor":
			fmt.Println("Executing history")
			data := ExecuteHistory(input.Args)		
			for _, line := range data {
				data , _ := EncodeBulkString(line)
				time.Sleep(2*time.Second)
				client.Conn.Write([]byte(data))

			}
		default:
			fmt.Println("None Command found")
			client.Conn.Write([]byte("-No such command found\r\n"))

	}


}
