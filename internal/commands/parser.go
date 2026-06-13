package commands

import (
	"fmt"
	"slices"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func ParseCommands(store *store.MemoryAlloc,input []string) any{
	

	if len(input)  == 0{
		fmt.Println("Check Man page")
		
	}
	data := input[1:]

	switch input[0] {
	case "GET","get","Get":
		val := GetCommand(store,data)
		fmt.Println(val)
	case "SET","set","Set":
		 val:= SetCommand(store,data)
		 if val {
			fmt.Println("Added the KeyValue")
		}else{
			fmt.Println("Some Error occured")
		}
	case "DEL","del","Del":
		val := DelCommand(store,data)
		if val {
			fmt.Println("Removed the value")
		}else{
			fmt.Println("No Such value found")
		}
	case "KEYS","keys","Keys":
		dataItems := KeysCommand(store)
		fmt.Println("The List of keys avaliable: ", dataItems)
	case "SAVE","Save","save":
		datamsg := SaveCommand(store)
		fmt.Println(datamsg)
	case "LOAD","Load","load":
		count := LoaderCommand(store)
		fmt.Println("Restored: ", count , " Keys")
	case "HISTORY","History","history":
		history:= HistoryCommand()
		fmt.Println(history)
	default:
		logger.ErrorLog(input[0]+" command not found")
		fmt.Println("Command Not Found,Check the man page")
	}

	return ""
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
			Val: data[1:sliceFrom],
			Expiry: true,
			Ttl: time.Now(),
		}

	}else{
		inpData = data[1:]
		stripTtl = []string{"",""} 
		dataItems = store.Item{
			Key: data[0],
			Val: inpData,
			Expiry: false,
			Ttl: time.Now(),
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
	for _, key := range keys {
		
		rawdata := store.GetKv(stre,key)

		val := persistence.StoreToJson(rawdata)
		fmt.Println("Saving: ", val)
	}
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

func HistoryCommand() string{

	data , err := logger.ReadFromFile()

	if utils.HandleError("Error while reading from file",err) {
		return "\n"
	}

	return  data
}