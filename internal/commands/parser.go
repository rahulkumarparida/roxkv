package commands

import (
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
		fmt.Println("Check Man page")

	}
	data := input[1:]

	switch input[0] {
	case "GET", "get", "Get":
		val := GetCommand(store, data)
		// fmt.Println(val)
		return val
	case "SET", "set", "Set":
		val := SetCommand(client, store, data)
		if val {
			fmt.Println("Added the KeyValue")
			return val
		} else {
			fmt.Println("Some Error occured")
			return val
		}
	case "DEL", "del", "Del":
		val := DelCommand(store, data)
		if val {
			fmt.Println("Removed the value")
			return val
		} else {
			fmt.Println("No Such value found")
			return val
		}
	case "KEYS", "keys", "Keys":
		dataItems := KeysCommand(store)
		fmt.Println("The List of keys avaliable: ", dataItems)
		return dataItems
	case "SAVE", "Save", "save":
		datamsg := SaveCommand(store)
		fmt.Println(datamsg)
		return datamsg
	case "LOAD", "Load", "load":
		count := LoaderCommand(store)
		fmt.Println("Restored: ", count, " Keys")
		return count
	case "HISTORY", "History", "history":
		history := HistoryCommand(data)
		fmt.Println(history)
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
		// fmt.Println("Command Not Found,Check the man page")
		fmt.Println("Recieved: ", input)
		return "+RUN\r\n"
	}

	return ""

}

func ParseInput(data []string) []string {
	var result []string
	var current []rune

	var inQuote bool
	var quoteChar rune

	for _, token := range data {
		for _, ch := range token {

			switch {
			case !inQuote && (ch == '"' || ch == '\'' || ch == '`'):
				inQuote = true
				quoteChar = ch

			case inQuote && ch == quoteChar:
				inQuote = false

			case !inQuote && ch == ' ':
				if len(current) > 0 {
					result = append(result, string(current))
					current = current[:0]
				}

			default:
				current = append(current, ch)
			}
		}

		if inQuote {
			current = append(current, ' ')
		} else if len(current) > 0 {
			result = append(result, string(current))
			current = current[:0]
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}

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

func SetCommand(client *utils.NewClient, stre *store.MemoryAlloc, data []string) bool {
	var dataItems store.Item

	if len(data[1:]) <= 0 {
		logger.ErrorLog("Provided empty value in the key value pair")
		fmt.Println("Empty value provided")
		return false

	}

	if data[0] == "--ttl" {
		ttlVals := data[1:3]
		dataVals := data[3:]
		
		if len(ttlVals) > 2 || ttlVals[0] == "" || ttlVals[1] == "" {
			logger.ErrorLog("2 Args after the --ttl flag")
			fmt.Println("2 args after --ttl")
			return false
		}
		value := dataVals[1:]
		sizeOfValue := len(value)
		var timetoAdd, err = strconv.Atoi(ttlVals[0]) // time like 12 ,13 ,14

		if utils.HandleError("Error while parsing time provided", err) {
			logger.ErrorLog("Parsing failed while parsing the time. Integre required")
			return false
		}
		futuretime := evaluatetime(timetoAdd, ttlVals[1])

		meta := store.Metadata{
			TTL:           futuretime,
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:          int64(sizeOfValue),
		}

		dataItems = store.Item{
			Key:  dataVals[0],
			Val:  value,
			Meta: meta,
		}

		if time.Until(futuretime) < time.Duration(time.Until(time.Now().Add(24*time.Hour))) {
			store.TTLMetricsContainer.ExpiresToday += 1
		}

		store.TTLMetricsContainer.ActiveTTLKeys += 1
		

	} else {
		value := ParseInput(data[1:])
		sizeOfValue := len(value)
		store.TTLMetricsContainer.PermanentKeys += 1
		meta := store.Metadata{
			TTL:           time.Time{},
			UpdatedAt:     time.Now(),
			LastAcessedBy: client,
			Size:	int64(sizeOfValue),
		}

		dataItems = store.Item{
			Key:  data[0],
			Val:  value,
			Meta: meta,
		}
	}

	store.SetKv(stre, &dataItems)

	return true
}

func GetCommand(stre *store.MemoryAlloc, data []string) any {
	if len(data) != 1 {
		fmt.Println("Check Man page")
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
		fmt.Println("Check Man page")
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
	fmt.Println("Keys:", keys)
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
	fmt.Println("Saving: ", val)
	logmsg := "All keys avaliable in RAM till now are saved to DB"
	logger.SucessLog(logmsg)
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
	return val
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
