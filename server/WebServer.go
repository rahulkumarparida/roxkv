package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)





func WebServer(){

	router := mux.NewRouter()
	router.HandleFunc("/api/events/{name}",SseHandler).Methods("GET")
	
	err := http.ListenAndServe(":6971",router)
	fmt.Println("Listening at localhost:6971")

	if err != nil {
		log.Fatal(err)
	}


}


func SseHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	toolrequiredName := params["name"]

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)
	t := time.NewTicker(5*time.Second)
	defer t.Stop()

	for{
		select{
		case <-clientGone:
			fmt.Println("Client Dsiconnected: ", r.RemoteAddr)
			return
		case <-t.C:
			 
			data := HandleResponseData(toolrequiredName)

			// stringifiedData , jerr:= json.Marshal(data)

			// if jerr != nil{
			// 	return
			// }

			json.NewEncoder(w).Encode(data)

           
            err := rc.Flush()
            if err != nil {
                return
            }

		}
	}


}



func HandleResponseData(name string) any{

	switch name {
	case "cpu":
		data :=GetCPUData()
		return data
	case "ram":
		data :=GetRAMData()
		return data
	case "uptime":
		data := GetServreUptimeData()
		return data
	case "machineinfo":
		data := GetComputerInformation()
		return data
	case "disk":
		data := GetDiskData()
		return data
	case "runtime":
		data := GetRuntimeData()
		return data
	default:
		return "No such endpoint found"
		
	}

}

func GetServreUptimeData() time.Duration{
	data := metrics.GetServerUptime()
	return  data
}

func GetCPUData() string{
	data := metrics.GetCPUUsage()
	return data
}

func GetRAMData() metrics.RAM{
	data := metrics.GetRAMUsage()
	return data
}

func GetComputerInformation() utils.MonitorComputeStat{
	data := metrics.GetComputerUsage()
	return data
}

func GetDiskData() metrics.DISK{
	data := metrics.GetDiskUsage("/") // Figure the path situation
	return  data
}

func GetRuntimeData() metrics.RuntimeStats{
	data := metrics.GetRuntimeStats()
	return data
}


