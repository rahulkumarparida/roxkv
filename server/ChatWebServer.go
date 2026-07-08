package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ollama/ollama/api"
	master "github.com/rahulkumarparida/roxkv/agents/Master"
	"github.com/rahulkumarparida/roxkv/internal/store"
)


var StoreInstance *store.MemoryAlloc
var NameSpaceInsatnce *store.NameSpace
var AgentInstance *api.Client
func WebChatServer(stre *store.MemoryAlloc, namespace *store.NameSpace, agent *api.Client){
	router := mux.NewRouter()
	router.HandleFunc("/api/chat/",callAI).Methods("POST")


}

func callAI(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var query string

	err := json.NewDecoder(r.Body).Decode(&query)

	if err != nil {
		fmt.Println("Query Not found")
		return
	}

	master.MasterAgent(query,StoreInstance, nil,NameSpaceInsatnce, AgentInstance)
	
}