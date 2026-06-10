package commands

import (
	"fmt"

	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/store"
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
	default:
		return "Command Not Found"
	}

	return ""
} 


func SetCommand(stre *store.MemoryAlloc ,data []string) bool{
	
	dataItems := store.Item{
		Key: data[0],
		Val: data[1:]}

	store.SetKv(stre , &dataItems)
	
	return true
}

func GetCommand(stre *store.MemoryAlloc ,data []string) any{
	if len(data) != 1  {
		fmt.Println("Check Man page")
	}

	dataVal := store.GetKv(stre , data[0])

	return  dataVal
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


	return "Saved"
}