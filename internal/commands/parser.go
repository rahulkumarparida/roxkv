package commands

import (
	"fmt"
	"regexp"
	"slices"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func CheckInputLength(input []string, paramsrequired int)bool{
	return len(input) == paramsrequired
}

func ParseCommands(store *store.MemoryAlloc,input []string , client *utils.NewClient) any{
	

	if len(input)  == 0{
		fmt.Println("Check Man page")
		
	}
	data := input[1:]

	switch input[0] {
	case "GET","get","Get":
		val := GetCommand(store,data)
		fmt.Println(val)
		return  val
	case "SET","set","Set":
		 val:= SetCommand(store,data)
		 if val {
			fmt.Println("Added the KeyValue")
			return  val
		}else{
			fmt.Println("Some Error occured")
			return  val
		}
	case "DEL","del","Del":
		val := DelCommand(store,data)
		if val {
			fmt.Println("Removed the value")
			return  val
		}else{
			fmt.Println("No Such value found")
			return  val
		}
	case "KEYS","keys","Keys":
		dataItems := KeysCommand(store)
		fmt.Println("The List of keys avaliable: ", dataItems)
		return  dataItems
	case "SAVE","Save","save":
		datamsg := SaveCommand(store)
		fmt.Println(datamsg)
		return  datamsg
	case "LOAD","Load","load":
		count := LoaderCommand(store)
		fmt.Println("Restored: ", count , " Keys")
		return  count
	case "HISTORY","History","history":
		history:= HistoryCommand(data)
		fmt.Println(history)
		return history
	case "SUBSCRIBE","Subscribe","subscribe":
		valid := CheckInputLength(data,1)
		if !valid{
			return false
		}		
		val := SubscriberCommand(client,data)
		return val
	case "PUBLISH","Publish","publish":
		valid := CheckInputLength(data,2)
		if !valid{
			return false
		}
		PublishCommand(client,data)
	case "UNSUBSCRIBE","Unsubscribe","unsubscribe":
		valid := CheckInputLength(data,1)
		if !valid{
			return false
		}
		UnsubscribeCommand(client,data)
	case "TOPICS","Topics","topics":
		TopicsCommand(client)	
	case "CLOSECHANNEL","CloseChannel","closechannel":
		valid := CheckInputLength(data,1)
		if !valid{
			return false
		}
		RemoveTopicCommand(client,data)		
	default:
		logger.ErrorLog(input[0]+" command not found")
		fmt.Println("Command Not Found,Check the man page")
		
		return "Error"
	}

	return ""

} 

func ParseInput(data []string) []string{
	var dataAppended []string
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	var startWord bool = false
	var tempWord string
	for _, v := range data {
		
		if len(v) == 0 {
			continue
		}

		if !startWord && (v[0] == '\'' || v[0] == '"' || v[0] == '`' ) {
				startWord = true
				tempWord += v
				continue
		}
		if startWord || v[len(v)-1] == '\'' || v[len(v)-1] ==  '"' || v[len(v)-1] == '`'  {
				tempWord +=" " + v
				if v[len(v)-1] == '\'' || v[len(v)-1] ==  '"' || v[len(v)-1] == '`'  {
					dataAppended = append(dataAppended, tempWord)
					startWord = false	
					tempWord = ""
				}
				
				continue
		}

		if !startWord && re.MatchString(v) {
				dataAppended = append(dataAppended, v)
				
		}
	}

	return dataAppended

}


func SetCommand(stre *store.MemoryAlloc ,data []string) bool{
	var dataItems store.Item
	var inpData []string
	var stripTtl []string

	if len(data[1:]) <= 0 {
		logger.ErrorLog("Provided empty value in the key value pair")
			fmt.Println("Empty value provided")
			return false
		
	}

	if slices.Contains(data,"--ttl") {
		sliceFrom := slices.Index(data,"--ttl")
		stripTtl = data[sliceFrom+1:]
		if len(stripTtl) > 2 {
			logger.ErrorLog("2 Args after the --ttl flag")
			fmt.Println("Only 2 args after --ttl")
			return false
		}
		dataItems = store.Item{
			Key: data[0],
			Val: ParseInput(data[1:sliceFrom]),
			Ttl: time.Now(),
		}

	}else{
		inpData = data[1:]
		stripTtl = []string{"",""} 
		dataItems = store.Item{
			Key: data[0],
			Val: ParseInput(inpData),
			Ttl: time.Time{},
		}
	}
	
	

	store.SetKv(stre , &dataItems , stripTtl)
	
	return true
}

func GetCommand(stre *store.MemoryAlloc ,data []string) any{
	if len(data) != 1  {
		fmt.Println("Check Man page")
	}

	dataVal := store.GetKv(stre , data[0])

	return  dataVal.Val
}


func DelCommand(stre *store.MemoryAlloc ,data []string) bool{
	if len(data) != 1  {
		fmt.Println("Check Man page")
	}

	dataVal := store.DelKv(stre , data[0])

	return  dataVal
}


func KeysCommand(stre *store.MemoryAlloc) []string{
	data := store.KeyKv(stre)

	return data
	
}

func SaveCommand(stre *store.MemoryAlloc) string{
	
	keys := store.KeyKv(stre)
	fmt.Println("Keys:", keys)
	var allData []store.Item
	for _, key := range keys {
		
		rawdata := store.GetKv(stre,key)
		if rawdata.Ttl.IsZero() {
			allData = append(allData, rawdata)	
		}

		
	}
	val := persistence.StoreToJson(allData)
	fmt.Println("Saving: ", val)
	logmsg:= "All keys avaliable in RAM till are saved to DB"
	logger.SucessLog(logmsg)
	return "Saved"
}

func LoaderCommand(stre *store.MemoryAlloc) int{
	logmsg:= "All keys avaliable in DB are loaded to RAM"
	logger.SucessLog(logmsg)
	count := persistence.LoadJsons(stre)
	return count
}

func HistoryCommand(args []string) string{

	if len(args) == 0 || len(args) > 1 {
		

		data , err := logger.ReadFromFile("--all")

		if utils.HandleError("Error while reading from file",err) {
			return "\n"
		}

		return  data

	}

	data , err := logger.ReadFromFile(args[0])

	if utils.HandleError("Error while reading from file",err) {
			return "\n"
	}

	return  data


}


func SubscriberCommand(client *utils.NewClient, input []string) bool{

	if input[0] == "" {
		return false
	}

	val := pubsub.HandleSubscribers(client,input[0])
	return val
}

func PublishCommand(client *utils.NewClient,input []string) any {
	if input[0] == "" || input[1] == "" {
		return "Channel name and then data is required, plase check thee help page or man for information."
	}
	pubsub.Broker(client,input[0],input[1])
	return ""
}

func UnsubscribeCommand(client *utils.NewClient, input []string) string{

	if input[0] == "" {
		return "Requires the topic name"
	}

	pubsub.HandleUnsubscribes(client,input[0])

	return ""
}


func TopicsCommand(client *utils.NewClient){
	
	pubsub.GetTopics(client)

}


func RemoveTopicCommand(client *utils.NewClient, input []string){
	if input[0] == "" {
		client.Conn.Write([]byte("Topic name required\n"))
		return 
	}
	pubsub.CloseChannel(client,input[0])
}