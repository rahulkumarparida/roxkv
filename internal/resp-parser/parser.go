package redisparser

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/utils"
)


type RedisInput struct{
	Cmd     string
	Args    []string   // string args for command names, keys, flags
	RawArgs [][]byte   // raw binary args from RESP (values stay as bytes)
}

var Mu = sync.Mutex{}
var tmpCommandData [][][]byte

func HandleMultipleCommands(client *utils.NewClient,buf []byte, n int) ([][][]byte,int,error){
	fmt.Println("Iter: ", len(tmpCommandData))
	tokens , till , err := DecodeArrayString(client.Buffer,len(client.Buffer))
	if err != nil{
		if err == ErrIncompleteRESP {
			fmt.Println("Erro:",err)
			return nil, 0 , err
		}
		return nil, till, errors.New("Incorrect RESP command") 
	}
	if tokens != nil{
	tmpCommandData = append(tmpCommandData, tokens)			
	}	
	if tokens == nil || till >= len(buf) {
		Mu.Lock()
		command := tmpCommandData
		tmpCommandData = [][][]byte{}
		Mu.Unlock()
		return  command, till , err
	}
	
	
	
	return 	HandleMultipleCommands(client,buf[till:],n)
}


func ReadAndHandleConnection(client *utils.NewClient) {

	var buf []byte =make([]byte, 4096)
	for {
		n, err := client.Conn.Read(buf)
		if err != nil {
			fmt.Println("Err:", err.Error(), " Client: ", client.Conn)
			break
		}

		client.Buffer = append(client.Buffer, buf[:n]...)	
		
		for {
			datalist, readTill ,err := HandleMultipleCommands(client, client.Buffer,n)

			if err != nil &&  err == ErrIncompleteRESP {
				break
			}
			if datalist == nil {
				break
			}else {
				ExecuteTokens(client,datalist,readTill)

				client.Mu.Lock()
				client.Buffer = client.Buffer[readTill:]
				client.Mu.Unlock()
					
			}
			
			if len(client.Buffer) == 0 {
				break
			}
		
		}
					
	}

	
}

// Should be only to execute the commands not decode it should only receieve decoded [][]string // buffer
func ExecuteTokens(client *utils.NewClient,allTokenList [][][]byte , n int) {


		
			for _, tokens := range allTokenList {
				// fmt.Println("Token:", tokens)/
				
				if  len(tokens)== 0{
					data := EncodeNullValues()
					client.Conn.Write([]byte(data))
					break
				}

				var input *RedisInput

				if len(tokens) >= 2 {
					fmt.Println("Token 0:", string(tokens[0]), tokens[0])
					input = &RedisInput{
						Cmd:     string(tokens[0]),
						Args:    bytesSliceToStrings(tokens[1:]),
						RawArgs: tokens[1:],
					}
				}else {
					input = &RedisInput{
						Cmd:     string(tokens[0]),
						Args:    nil,
						RawArgs: nil,
					}
				}

				
				parseCommand(input, client)		
				allTokenList= allTokenList[1:]
			}
		
			allTokenList = [][][]byte{}	

}



// bytesSliceToStrings converts [][]byte to []string for command/key/flag use.
// Only use for known textual data, not for values.
func bytesSliceToStrings(bs [][]byte) []string {
	ss := make([]string, len(bs))
	for i, b := range bs {
		ss[i] = string(b)
	}
	return ss
}

func parseCommand(input *RedisInput, client *utils.NewClient) {
	if PatternRegister == nil{
		PatternRegister = make(map[string]PatternRegistry)
	}

	subscriberModeCommands :=[]string{"subscribe","ping","unsubscribe","psubscribe","punsubscribe","quit","reset"}

	if client.Mode.Name == utils.ModeSubsriber.Name{
		if !slices.Contains(subscriberModeCommands,strings.ToLower(input.Cmd)) {

			encode , _ := EncodeSimpleError(" ERR Can't execute 'publish': only (P|S)SUBSCRIBE / (P|S)UNSUBSCRIBE / PING / QUIT / RESET are allowed in this context")
			client.Conn.Write([]byte(encode))
			return
		}


	}




	fmt.Println("Input:", input)	
	switch strings.ToLower(input.Cmd){
		case "ping":
			fmt.Println("Executing ping", input.Args)
			if client.Mode.Name == utils.ModeSubsriber.Name{
					if len(input.Args) != 1 {
						encode , _:= EncodeSimpleError("ERR wrong number of arguments for 'ping' command")
						client.Conn.Write([]byte(encode))
						return
					}
				encode , _:= EncodeStringArray([]string{"pong",input.Args[0]})
				client.Conn.Write([]byte(encode))
				return
			}
			val := ExecutePing(input.Args)
			client.Conn.Write([]byte(val))
		case "client":
			fmt.Println("Executing Client ",input.Args[0])
			data := ExecuteClientname(input.Args,client)
			client.Conn.Write([]byte(data))
		case "set":
			fmt.Println("Executing set , args: ", input.RawArgs)
			data :=ExecuteSet(input , client)
			client.Conn.Write([]byte(data.(string)))
		case "get":
			fmt.Println("Executing get")
			data := ExecuteGet(input.Args)
			client.Conn.Write(data)
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
			val := ExecuteExpire(input.RawArgs,client)
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
		case "command":
			fmt.Println("Executed command")
			if len(input.Args) <= 1 {
				if len(input.Args) == 1 && strings.ToLower(input.Args[0]) == "docs"{
					data := ExecuteCommandDocs()
					client.Conn.Write([]byte(data))
					
				}else if len(input.Args) == 0{
					data := ExecuteCommand()
					client.Conn.Write([]byte(data))

				}else{
					encode , _:= EncodeSimpleError("ERR wrong number of arguments for 'command' command")
					client.Conn.Write([]byte(encode))
				}
			}	
			return
		case "subscribe","psubscribe","publish","unsubscribe","punsubscribe":
			
			if input.Cmd =="subscribe" {
				if client.Mode.Name != utils.ModeSubsriber.Name {
					client.Mode = utils.ModeSubsriber
				}
				SubscribeHandler(input.Args,client)
			}
			if input.Cmd =="unsubscribe" {
				UnsubscribeHandler(input.Args,client)
			}
			if input.Cmd =="publish"{
				fmt.Println("Pubslihing shit:", input.Args)
				data :=PublishHandler(input.Args,client)
				client.Conn.Write([]byte(data))
			}
			if input.Cmd == "psubscribe"{
				if client.Mode.Name != utils.ModeSubsriber.Name {
					client.Mode = utils.ModeSubsriber
				}
				PSubscribeHandler(input.Args,client)
			}

			if input.Cmd == "punsubscribe" {
				PUnsubscribeHandler(input.Args,client)
			}

		case "quit":
			fmt.Println("executing quit")
			ExecuteQuit(client)
		case "reset":
			fmt.Println("executing reset")
			data := ExecuteReset(client)
			client.Conn.Write([]byte(data))
		case "lastsave":
			fmt.Println("Executing LastSave")
			data := ExecuteLastSave()
			client.Conn.Write([]byte(data))
		case "monitor":
			fmt.Println("Executing history")
			data := ExecuteHistory(input.Args)		
			for _, line := range data {
				data , _ := EncodeBulkString(line)
				time.Sleep(2*time.Second)
				client.Conn.Write([]byte(data))
			}
		case "select":
			fmt.Println("Executing select database space:", input.Args)
			// Return OK to satisfy connection setup workflows
			client.Conn.Write([]byte("+OK\r\n"))

		default:
			fmt.Println("None Command found")
			client.Conn.Write([]byte("-No such command found\r\n"))

	}


}
