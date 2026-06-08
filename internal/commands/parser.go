package commands

import (
	"fmt"

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
		return false
	case "KEYS","keys","Keys":
		return false
	default:
		return "Command Not Found"
	}

	return "Check Man page"
} 


func SetCommand(stre *store.MemoryAlloc ,data []string) bool{
	if len(data) != 2  {
		fmt.Println("Check Man page")
	}
	dataItems := store.Item{
		Key: data[0],
		Val: data[1]}

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