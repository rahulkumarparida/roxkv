package master

// import (
// 	"context"
// 	"log"

// 	"github.com/ollama/ollama/api"
// 	"github.com/rahulkumarparida/roxkv/internal/store"
// 	"github.com/rahulkumarparida/roxkv/internal/utils"
// )


// var properties = api.NewToolPropertiesMap()	

// func PropertyMapper(name string,params any ,description string) *api.ToolPropertiesMap{

// properties.Set(name,api.ToolProperty{Items:params,Description: description})

// return properties

// }




// // []api.ToolProperty{
// // 	{
// // 		Type: *store.MemoryAlloc,
// // 		Description: "This is the handler of the key value memory, contains all the operations such as get,set,delete,keys,load,save",
// // 	},
// // 	{
// // 		Type: *utils.NewClient,
// // 		Description: "This is the client instance containing the data and connection of the client over the server",
// // 	},

// // }





// func AgentManager(stre *store.MemoryAlloc,user *utils.NewClient){

// 	client , err := api.ClientFromEnvironment()
// 	if err != nil {
// 		log.Fatalf("Failed to create client: %v", err)
// 	}
// 	ctx  := context.Background()

// 	PropertyMapper("store", stre, "This is the handler of the key value memory, contains all the operations such as get,set,delete,keys,load,save")
// 	PropertyMapper("client",user,"This is the client instance containing the data and connection of the client over the server")
	
// 	kvtoolparams := api.ToolFunctionParameters{
// 		Type: "object",
// 		Properties: properties,
// 		Required: []string{"store","client"},
// 	}


// 	var	KeyValueTool = api.Tool{
// 		Type: "function",
// 		Function: api.ToolFunction{
// 			Name: "KeyValueAgent",
// 			Description: "Manipulations of keys , values in the database or retrieving data and summarizations for the key value pairs",
// 			Parameters:kvtoolparams,
// 		},

// 	}

	

// req := &api.ChatRequest{
// 	Model:  "llama3.1",
// 	Stream: false,
// 	Messages: []api.Message{
// 		{Role: "user", Content: ""},
// 	},
// 	Tools: []api.Tool{KeyValueTool}, 
// }



// }