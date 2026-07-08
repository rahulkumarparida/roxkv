package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/metrics"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

var StoreHelper *store.MemoryAlloc

// Accepting the HTTP request as HTTP
func WebServer(store *store.MemoryAlloc) {

	StoreHelper = store

	fmt.Println("Listening Webserver at localhost:6971")

	router := mux.NewRouter()
	router.HandleFunc("/api/events/{name}", SseHandler).Methods("GET")
	err := http.ListenAndServe(":6971", router)

	if err != nil {
		log.Fatal(err)
	}

}

func SseHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	toolrequiredName := params["name"]

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	tcpConn, rw, herr := hijacker.Hijack()
	if herr != nil {
		http.Error(w, herr.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteString("HTTP/1.1 200 OK\r\n")
	rw.WriteString("Content-Type: text/event-stream\r\n")
	rw.WriteString("Cache-Control: no-cache\r\n")
	rw.WriteString("Connection: keep-alive\r\n")
	rw.WriteString("Access-Control-Allow-Origin: *\r\n")
	rw.WriteString("\r\n")
	rw.Flush()

	client := utils.CreateClient(tcpConn, "user")

	clientGone := r.Context().Done()

	// rc := http.NewResponseController(w)
	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-clientGone:
			tcpConn.Close()
			fmt.Println("Client Dsiconnected: ", client)
			return
		case <-t.C:

			data := HandleResponseData(toolrequiredName, client, StoreHelper)

			stringifiedData, jerr := json.Marshal(data)

			if jerr != nil {
				tcpConn.Close()
				return
			}

			rw.WriteString("data: ")
			rw.Write(stringifiedData)
			rw.WriteString("\n\n")

			if err := rw.Flush(); err != nil {
				tcpConn.Close()
				return
			}

		}
	}

}

func HandleResponseData(name string, client *utils.NewClient, store *store.MemoryAlloc) any {

	switch name {
	case "cpu":
		data := GetCPUData()
		return data
	case "ram":
		data := GetRAMData()
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
	case "pubsubtopics":
		data := GetTopicList(client)
		return data
	case "dbdata":
		data := DataBaseData()
		return data
	case "livehistory":
		data := LiveHistoryData()
		return data
	default:
		return "No such endpoint found"

	}

}

// Machine Health
func GetServreUptimeData() time.Duration {
	data := metrics.GetServerUptime()
	return data
}

func GetCPUData() string {
	data := metrics.GetCPUUsage()
	return data
}

func GetRAMData() metrics.RAM {
	data := metrics.GetRAMUsage()
	return data
}

func GetComputerInformation() utils.MonitorComputeStat {
	data := metrics.GetComputerUsage()
	return data
}

func GetDiskData() metrics.DISK {
	data := metrics.GetDiskUsage("/") // Figure the path situation
	return data
}

func GetRuntimeData() metrics.RuntimeStats {
	data := metrics.GetRuntimeStats()
	return data
}

// Snapshot Info
func SnapShotData() metrics.PersistenceHealth {
	data := metrics.GetPersistenceHealth()
	return data
}

type DataBaseInfo struct {
	SavedDataSize     int64            `json:"savedDataSize"`
	SavedSnapShotSize int64            `json:"savedSnapShotSize"`
	TotalKeysSize     int64            `json:"totalKeysSize"`
	LastSnapShotTime  time.Time        `json:"lastSnapShotTime"`
	TtlMetrics        *store.TTLMetrics `json:"ttlMetrics"`
}

func DataBaseData() DataBaseInfo {
	dbFolder := utils.DbFolder()
	snapshotFolder := utils.SnapshotFolder()

	dbinfo, err := os.Stat(dbFolder)
	snapinfo, err := os.Stat(snapshotFolder)

	if err != nil {
		fmt.Println("Directories not found")
		return DataBaseInfo{}
	}

	dbSize := dbinfo.Size()
	snapshotSize := snapinfo.Size()
	latestSnapShotTime := metrics.GetLatestSnapshot().ModifiedAt

	ttlMetrics := metrics.GetTTLMetrics()

	allKeys := metrics.LiveItems(StoreHelper)

	var TotalKeysSize int64
	for _, key := range allKeys {
		TotalKeysSize += key.Meta.Size
	}

	return DataBaseInfo{
		SavedDataSize:     dbSize,
		SavedSnapShotSize: snapshotSize,
		TotalKeysSize:     TotalKeysSize,
		LastSnapShotTime:  latestSnapShotTime,
		TtlMetrics:        &ttlMetrics,
	}

}

// Pubsub
type TopicList struct {
	Topics      []string `json:"topics"`
	TotalTopics int      `json:"total"`
}

func GetTopicList(client *utils.NewClient) TopicList {
	cnt, list := metrics.GetTopicCount(client)

	data := TopicList{Topics: list, TotalTopics: cnt}

	return data

}



// History Info
func LiveHistoryData() []string{
	data , err := logger.ReadFromFile("--all")
	if err != nil{
		return  []string{"Not Data Found","Chcek after a few seconds"}
	}

	dataArr := strings.Split(data,"\n")
	slices.Reverse(dataArr)
	return dataArr[1:11]
}
