package redisparser

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
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
	tokens , till , err := DecodeArrayString(client.Buffer,len(client.Buffer))
	if err != nil{
		if err == ErrIncompleteRESP {
			logger.ErrorLog("RESP: incomplete command received")
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
			logger.ErrorLog("RESP: connection read error for client " + fmt.Sprintf("%v", client.Conn.RemoteAddr()) + ": " + err.Error())
			break
		}

		client.Buffer = append(client.Buffer, buf[:n]...)	
		
		for {
			cmdList, readTill ,err := HandleMultipleCommands(client, client.Buffer,n)

			if err != nil &&  err == ErrIncompleteRESP {
				break
			}
			if cmdList == nil {
				break
			}else {
				ExecuteTokens(client,cmdList,readTill)

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
					logger.InfoLog("RESP: executing command " + string(tokens[0]))
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




	logger.InfoLog("RESP: parseCommand received: " + input.Cmd)
	switch strings.ToLower(input.Cmd){
		case "ping":
			logger.InfoLog("RESP: ping")
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
			logger.InfoLog("RESP: client " + input.Args[0])
			data := ExecuteClientname(input.Args,client)
			client.Conn.Write([]byte(data))
		case "set":
			logger.InfoLog("RESP: set key=" + input.Args[0])
			data :=ExecuteSet(input , client)
			client.Conn.Write([]byte(data.(string)))
		case "get":
			logger.InfoLog("RESP: get key=" + input.Args[0])
			data := ExecuteGet(input.Args)
			client.Conn.Write(data)
		case "del":
			logger.InfoLog("RESP: del")
			data := ExecuteDel(input.Args)
			client.Conn.Write([]byte(data))
		case "keys":
			logger.InfoLog("RESP: keys")
			data := Executekeys(input.Args)
			client.Conn.Write([]byte(data))
		case "exists":
			logger.InfoLog("RESP: exists")
			data := ExecuteExists(input.Args)
			client.Conn.Write([]byte(data))
		case "save":
			logger.InfoLog("RESP: save")
			val := ExecuteSave()
			client.Conn.Write([]byte(val))
		case "expire":
			logger.InfoLog("RESP: expire")
			val := ExecuteExpire(input.RawArgs,client)
			client.Conn.Write([]byte(val))
		case "flushdb":
			logger.InfoLog("RESP: flushdb")
			val := ExecuteFlushDB()
			client.Conn.Write([]byte(val))
		case "echo":
			logger.InfoLog("RESP: echo")
			val := ExecuteEcho(input.Args)
			client.Conn.Write([]byte(val))
		case "ttl":
			logger.InfoLog("RESP: ttl")
			val := ExecuteTTL(input.Args)
			client.Conn.Write([]byte(val))
		case "persist":
			logger.InfoLog("RESP: persist")
			val := ExecutePersistance(input.Args, client)
			client.Conn.Write([]byte(val))
		case "dbsize":
			logger.InfoLog("RESP: dbsize")
			val := ExecuteDbSize()
			client.Conn.Write([]byte(val))
		case "randomkey":
			logger.InfoLog("RESP: randomkey")
			val := ExecuteRandomKeys()
			client.Conn.Write([]byte(val))
		case "rename":
			logger.InfoLog("RESP: rename")
			val := ExecuteRenameKey(input.Args, client)
			client.Conn.Write([]byte(val))
		case "uptime":
			logger.InfoLog("RESP: uptime")
			val := ExecuteServerUpTime(input.Args)
			client.Conn.Write([]byte(val))
		case "mget":
			logger.InfoLog("RESP: mget")
			val := ExecuteMget(input.Args)
			client.Conn.Write([]byte(val))
		case "incr", "decr", "incrby", "decrby":
			logger.InfoLog("RESP: integer op " + input.Cmd)
			val := ExecuteIntOpration(input.Cmd, input.Args, client)
			client.Conn.Write([]byte(val))
		case "strlen":
			logger.InfoLog("RESP: strlen")
			val := ExecuteStrLen(input.Args)
			client.Conn.Write([]byte(val))
		case "command":
			logger.InfoLog("RESP: command")
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
				logger.InfoLog("RESP: publish to topic " + input.Args[0])
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
			logger.InfoLog("RESP: quit")
			ExecuteQuit(client)
		case "reset":
			logger.InfoLog("RESP: reset")
			data := ExecuteReset(client)
			client.Conn.Write([]byte(data))
		case "lastsave":
			logger.InfoLog("RESP: lastsave")
			data := ExecuteLastSave()
			client.Conn.Write([]byte(data))
		case "monitor":
			logger.InfoLog("RESP: monitor/history")
			data := ExecuteHistory(input.Args)		
			for _, line := range data {
				data , _ := EncodeBulkString(line)
				time.Sleep(2*time.Second)
				client.Conn.Write([]byte(data))
			}
		case "select":
			logger.InfoLog("RESP: select database space: " + fmt.Sprintf("%v", input.Args))
			// Return OK to satisfy connection setup workflows
			client.Conn.Write([]byte("+OK\r\n"))

		default:
			logger.ErrorLog("RESP: unknown command " + input.Cmd)
			client.Conn.Write([]byte("-No such command found\r\n"))

	}


}
