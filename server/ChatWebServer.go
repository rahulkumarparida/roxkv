package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ollama/ollama/api"

	master "github.com/rahulkumarparida/roxkv/agents/Master"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)


var StoreInstance *store.MemoryAlloc
var NameSpaceInsatnce *store.NameSpace
var AgentInstance *api.Client

func WebChatServer(stre *store.MemoryAlloc, namespace *store.NameSpace, agent *api.Client){
	StoreInstance = stre
	NameSpaceInsatnce = namespace
	AgentInstance = agent


	router := mux.NewRouter()
	router.HandleFunc("/healthz", HealthHandler).Methods("GET")
	router.HandleFunc("/api/chat/",callAI).Methods("POST", "OPTIONS")
	fmt.Println("Listening WebChat API at localhost:6972")

	
	err := http.ListenAndServe(":6972",router)	
	if err != nil {
		fmt.Printf("Error in starting the server %v", err)
		return 
	}

}


var payload struct {
		Query string `json:"query"`
}


func callAI(w http.ResponseWriter, r *http.Request){
	
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Error decoding JSON payload:", err)
		return
	}
	defer r.Body.Close()

	hijacker, ok := w.(http.Hijacker)

	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}



	tcpConn , rw , err := hijacker.Hijack()
	
	if err != nil || tcpConn == nil {
			fmt.Println("Not able to hijack the connection:", err)
			fmt.Println("Connections is nil :", tcpConn == nil)
			return
	}
	defer tcpConn.Close()


	client := utils.CreateClient(tcpConn,"client")	
	utils.TotalConnecntions = append(utils.TotalConnecntions, client)
	
	
	



	query := payload.Query

	data := master.MasterAgent(query,StoreInstance,client,NameSpaceInsatnce, AgentInstance)	
	

	fmt.Println("Data:",data , client)

	rw.WriteString("HTTP/1.1 200 OK\r\n")
	rw.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	rw.WriteString("Cache-Control: no-cache\r\n")
	rw.WriteString("Access-Control-Allow-Origin: *\r\n")
	rw.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(fmt.Sprintf("%v",data))))
	rw.WriteString("\r\n") // End of headers
	rw.Flush()


	rw.WriteString(fmt.Sprintf("%v",data))
	rw.Flush()
	
}
