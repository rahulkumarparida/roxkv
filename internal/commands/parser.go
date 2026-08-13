package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func CheckInputLength(input []string, paramsrequired int) bool {
	return len(input) == paramsrequired
}

func ParseCommands(store *store.MemoryAlloc, input []string, client *utils.NewClient) any {

	if len(input) == 0 {
		logger.InfoLog("ParseCommands: empty input received")
	}
	data := input[1:]

	switch input[0] {
	case "GET", "get", "Get":
		val := GetCommand(store, data)
		if b, ok := val.([]byte); ok {
			return string(b)
		}
		return val
	case "SET", "set", "Set":
		
		var bytesdata  [][]byte 
		for _,d := range data{
			singleData := []byte(d)
			bytesdata = append(bytesdata, singleData)
		} 

		val := SetCommand(client, store, [][]byte(bytesdata))
		logger.InfoLog("SET command executed: success=" + fmt.Sprintf("%v", val))
		return val
	case "DEL", "del", "Del":
		val := DelCommand(store, data)
		logger.InfoLog("DEL command executed: success=" + fmt.Sprintf("%v", val))
		return val
	case "KEYS", "keys", "Keys":
		dataItems := KeysCommand(store)
		logger.InfoLog("KEYS command executed")
		return dataItems
	case "SAVE", "Save", "save":
		datamsg := SaveCommand(store)
		logger.InfoLog("SAVE command executed: " + fmt.Sprintf("%v", datamsg))
		return datamsg
	case "LOAD", "Load", "load":
		count := LoaderCommand(store)
		logger.InfoLog("LOAD command executed: restored " + fmt.Sprintf("%d", count) + " keys")
		return count
	case "HISTORY", "History", "history":
		history := HistoryCommand(data)
		return history
	case "SUBSCRIBE", "Subscribe", "subscribe":
		valid := CheckInputLength(data, 1)
		if !valid {
			return false
		}
		val := SubscriberCommand(client, data)
		return val
	case "PUBLISH", "Publish", "publish":
		valid := CheckInputLength(data, 2)
		if !valid {
			return false
		}
		PublishCommand(client, data)
	case "UNSUBSCRIBE", "Unsubscribe", "unsubscribe":
		valid := CheckInputLength(data, 1)
		if !valid {
			return false
		}
		UnsubscribeCommand(client, data)
	case "TOPICS", "Topics", "topics":
		TopicsCommand(client)
	case "CLOSECHANNEL", "CloseChannel", "closechannel":
		valid := CheckInputLength(data, 1)
		if !valid {
			return false
		}
		RemoveTopicCommand(client, data)
	default:
		logger.ErrorLog(input[0] + " command not found")
		return "+RUN\r\n"
	}

	return ""

}

// func ParseInput(data []string) []string {
// 	var result [][]byte
// 	var current []rune

// 	var inQuote bool
// 	var quoteChar rune

// 	for _, token := range data {
// 		for _, ch := range token {

// 			switch {
// 			case !inQuote && (ch == '"' || ch == '\'' || ch == '`'):
// 				inQuote = true
// 				quoteChar = ch

// 			case inQuote && ch == quoteChar:
// 				inQuote = false

// 			case !inQuote && ch == ' ':
// 				if len(current) > 0 {
// 					result = append(result, []byte(current))
// 					current = current[:0]
// 				}

// 			default:
// 				current = append(current, ch)
// 			}
// 		}

// 		if inQuote {
// 			current = append(current, ' ')
// 		} else if len(current) > 0 {
// 			result = append(result, []byte(current))
// 			current = current[:0]
// 		}
// 	}

// 	if len(current) > 0 {
// 		result = append(result, []byte(current))
// 	}

// 	return result
// }

func evaluatetime(ttl int, dur string) time.Time {

	switch strings.ToLower(dur) {
	case "second", "sec":
		return time.Now().Add(time.Duration(ttl) * time.Second)
	case "minutes", "min":
		return time.Now().Add(time.Duration(ttl) * time.Minute)
	case "hour", "hh":
		return time.Now().Add(time.Duration(ttl) * time.Hour)
	default:
		return time.Now()
	}
}

func SetCommand(client *utils.NewClient, stre *store.MemoryAlloc, data [][]byte) bool {
	var dataItems store.Item

	if len(data[1:]) <= 0 {
		logger.ErrorLog("Provided empty value in the key value pair")
		return false
	}

	if string(data[0]) == "--ttl" {
		ttlVals := data[1:3]
		dataVals := data[3:]
		
		if len(ttlVals) > 2 || string(ttlVals[0]) == "" || string(ttlVals[1]) == "" {
			logger.ErrorLog("2 Args after the --ttl flag")
			return false
		}
		value := bytes.Join(dataVals[1:],[]byte(" "))
		sizeOfValue := len(value)
		var timetoAdd, err = strconv.Atoi(string(ttlVals[0])) // time like 12 ,13 ,14

		if utils.HandleError("Error while parsing time provided", err) {
			logger.ErrorLog("Parsing failed while parsing the time. Integre required")
			return false
		}
		futuretime := evaluatetime(timetoAdd, string(ttlVals[1]))

		meta := store.Metadata{
			TTL:           futuretime,
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:          int64(sizeOfValue),
		}

		dataItems = store.Item{
			Key:  string(dataVals[0]),
			Val:  value,
			Meta: meta,
		}

		if time.Until(futuretime) < time.Duration(time.Until(time.Now().Add(24*time.Hour))) {
			store.TTLMetricsContainer.ExpiresToday += 1
		}

		store.TTLMetricsContainer.ActiveTTLKeys += 1
		

	} else {
		value := bytes.Join(data[1:],[]byte(" "))
		sizeOfValue := len(value)
		store.TTLMetricsContainer.PermanentKeys += 1
		meta := store.Metadata{
			TTL:           time.Time{},
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:	int64(sizeOfValue),
		}

		dataItems = store.Item{
			Key:  string(data[0]),
			Val:  value,
			Meta: meta,
		}
	}

	store.SetKv(stre, &dataItems)

	return true
}

// SetCommandBytes stores a raw []byte value for the given key.
// Used by the RESP server path where values are already binary.
func SetCommandBytes(client *utils.NewClient, stre *store.MemoryAlloc, key string, val []byte, ttlSeconds int, ttlUnit string) bool {
	var dataItems store.Item
	sizeOfValue := len(val)

	if ttlSeconds > 0 && ttlUnit != "" {
		futuretime := evaluatetime(ttlSeconds, ttlUnit)
		meta := store.Metadata{
			TTL:           futuretime,
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:          int64(sizeOfValue),
		}
		dataItems = store.Item{
			Key:  key,
			Val:  val,
			Meta: meta,
		}
		store.TTLMetricsContainer.ActiveTTLKeys += 1
	} else {
		store.TTLMetricsContainer.PermanentKeys += 1
		meta := store.Metadata{
			TTL:           time.Time{},
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:          int64(sizeOfValue),
		}
		dataItems = store.Item{
			Key:  key,
			Val:  val,
			Meta: meta,
		}
	}

	store.SetKv(stre, &dataItems)
	return true
}

func GetCommand(stre *store.MemoryAlloc, data []string) any {
	if len(data) != 1 {
		logger.InfoLog("GetCommand: expected 1 arg, got " + fmt.Sprintf("%d", len(data)))
	}
	// fmt.Println("Store:", stre)

	dataVal := store.GetKv(stre, data[0])
	if dataVal.Val == nil {
		return  nil
	}

	return dataVal.Val
}

func DelCommand(stre *store.MemoryAlloc, data []string) bool {
	if len(data) != 1 {
		logger.InfoLog("DelCommand: expected 1 arg, got " + fmt.Sprintf("%d", len(data)))
	}

	dataVal := store.DelKv(stre, data[0])

	return dataVal
}

func KeysCommand(stre *store.MemoryAlloc) []string {
	data := store.KeyKv(stre)

	return data

}

func SaveCommand(stre *store.MemoryAlloc) string {

	keys := store.KeyKv(stre)
	logger.InfoLog("SaveCommand: saving " + fmt.Sprintf("%d", len(keys)) + " keys")
	var allData []store.Item
	for _, key := range keys {

		rawdata := store.GetKv(stre, key)
		if rawdata.Meta.TTL.IsZero() {
			allData = append(allData, rawdata)
		}

	}
	dbpath := utils.DbFolder()
	filename := time.Now().Format("2006-01-02") + ".json"
	val := persistence.StoreToJson(dbpath, filename, allData)
	logger.SucessLog("All keys available in RAM till now are saved to DB")
	if !val {
		logger.ErrorLog("SaveCommand: persistence failed")
	}
	return "Saved"
}

func LoaderCommand(stre *store.MemoryAlloc) int {
	dbpath := utils.DbFolder()
	logmsg := "All keys avaliable in DB are loaded to RAM"
	logger.SucessLog(logmsg)
	count := persistence.LoadJsons(dbpath, stre)

	store.TTLMetricsContainer.PermanentKeys += count
	return count
}

func HistoryCommand(args []string) string {

	if len(args) == 0 || len(args) > 1 {

		data, err := logger.ReadFromFile("--all")

		if utils.HandleError("Error while reading from file", err) {
			return "\n"
		}

		return data

	}

	data, err := logger.ReadFromFile(args[0])

	if utils.HandleError("Error while reading from file", err) {
		return "\n"
	}

	return data

}

func SubscriberCommand(client *utils.NewClient, input []string) bool {

	if input[0] == "" {
		return false
	}

	val := pubsub.HandleSubscribers(client, input[0])
	return val != nil
}

func PublishCommand(client *utils.NewClient, input []string) any {
	if input[0] == "" || input[1] == "" {
		return "Channel name and then data is required, plase check thee help page or man for information."
	}
	pubsub.Broker(client, input[0], input[1])
	return ""
}

func UnsubscribeCommand(client *utils.NewClient, input []string) string {

	if input[0] == "" {
		return "Requires the topic name"
	}

	pubsub.HandleUnsubscribes(client, input[0])

	return ""
}

func TopicsCommand(client *utils.NewClient) {

	data := pubsub.GetTopics(client)

	for idx, topic := range data {
		client.Conn.Write([]byte(strconv.Itoa(idx) + ". " + topic + "\n"))
	}

}

func RemoveTopicCommand(client *utils.NewClient, input []string) {
	if input[0] == "" {
		client.Conn.Write([]byte("Topic name required\n"))
		return
	}
	pubsub.CloseChannel(client, input[0])
}
