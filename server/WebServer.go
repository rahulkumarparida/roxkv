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


// Accepting the HTTP request as HTTP
func WebServer() {


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

	// SSE Headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send headers immediately
	flusher.Flush()

	clientGone := r.Context().Done()

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {

		case <-clientGone:
			fmt.Println("Client disconnected")
			return

		case <-t.C:

			// Remove the client argument later after your Pub/Sub refactor.
			data := HandleResponseData(toolrequiredName, nil, store.StoreHelper)

			stringifiedData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("JSON Error:", err)
				return
			}

			_, err = fmt.Fprintf(w, "data: %s\n\n", stringifiedData)
			if err != nil {
				fmt.Println("Write Error:", err)
				return
			}

			flusher.Flush()
		}
	}
}

func HandleResponseData(name string, client *utils.NewClient, store *store.MemoryAlloc) any {

	switch name {
	case "monitor":
		data := MonitorStatsData()
		return data
	case "storage":
		data := SnapshotData()
		return data
	case "database":
		data := DataBaseData(client)
		return data
	case "activity":
		data := LiveHistoryData()
		return data
	default:
		return "No such endpoint found"

	}

}

// Machine Health
func GetComputerInformation() metrics.RuntimeStats {
	data := metrics.GetRuntimeStats()
	return data
}



type DataBaseInfo struct {
	SavedDataSize     int64            `json:"savedDataSize"`
	SavedSnapShotSize int64            `json:"savedSnapShotSize"`
	TotalKeysSize     int64            `json:"totalKeysSize"`
	LastSnapShotTime  time.Time        `json:"lastSnapShotTime"`
	TtlMetrics        *store.TTLMetrics `json:"ttlMetrics"`
	PubSubTopics      []TopicList         `json:"pubsubTopics"`
}

type MonitorStats struct{
	CpuUsage string `json:"cpuusage"`
	RamUsage metrics.RAM `json:"ramusage"`
	DiskUsage metrics.DISK `json:"diskusage"`
	NetworkUsage metrics.NetworkStat `json:"networkstats"`
	MachineInfo metrics.RuntimeStats `json:"machineinfo"`
}

// Monitor 
func MonitorStatsData() MonitorStats{
	netStats := make(chan metrics.NetworkStat, 0)
	machineinfo := GetComputerInformation()
	cpudata := metrics.GetCPUUsage()
	ramdata := metrics.GetRAMUsage()
	diskdata := metrics.GetDiskUsage("/")
	go metrics.NetworkStatistics(netStats)

	networkData := <-netStats

	return MonitorStats{
		CpuUsage: cpudata,
		RamUsage: ramdata,
		DiskUsage: diskdata,
		NetworkUsage: networkData,
		MachineInfo:machineinfo,
	}

}

// database
func DataBaseData(client *utils.NewClient) DataBaseInfo {
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

	allKeys := metrics.LiveItems(store.StoreHelper)
	topics := GetTopicList(client)
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
		PubSubTopics:      []TopicList{topics},
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


type SnapshotInfo struct {
	TotalSnapShotSize int64     `json:"totalSnapShotSize"`
	TotalSnapshots int       `json:"totalSnapshots"`
	LatestSnapshot time.Time `json:"latestSnapshot"`
	LastCreated    time.Time `json:"lastCreated"`
	NextSnap       float64 `json:"nextSnap"`
}

//Snapshot Info
func SnapshotData() SnapshotInfo{

	snapshotFolder := utils.SnapshotFolder()

	snapinfo, err := os.Stat(snapshotFolder)
		if err != nil {
		fmt.Println("Directories not found")
		return SnapshotInfo{}
	}
	snapshotSize := snapinfo.Size()
	totalSnaps := metrics.GetSnapshotCount()
	lastCreated := metrics.GetLatestSnapshot().ModifiedAt
	latestsnap := metrics.GetLatestSnapshot()
	NextSnap := metrics.NextSnapshotTime()

	return SnapshotInfo{
		TotalSnapShotSize: snapshotSize,
		TotalSnapshots: totalSnaps,
		LatestSnapshot: latestsnap.ModifiedAt,
		LastCreated: lastCreated,
		NextSnap: NextSnap,
	}
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
