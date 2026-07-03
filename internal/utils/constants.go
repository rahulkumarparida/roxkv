package utils

import (
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DB use
var TotalInputs int = 0

// Server Client
var ServerStarted time.Time
var TotalConnecntions []*NewClient

type Role string

const (
	RoleSystem Role = "system"
	RoleClient Role = "client"
	RoleAdmin Role = "admin"	
)

type NewClient struct {
	ID          any
	Role		string
	Conn        net.Conn
	LastUsed    time.Time
	ConnectedAt time.Time
	Mu          sync.RWMutex
	Interactions int
}

type ClientMetaData struct{
	Id string
	LastUsed time.Time
	ConnectedAt time.Time
	Interactions int
}




func CreateClient(conn net.Conn,role string) *NewClient {

	if role != string(RoleAdmin) && role != string(RoleClient) && role != string(RoleSystem) {
		return nil
	}


	return &NewClient{
		ID:          conn.RemoteAddr().String(),
		Role: 		 role,
		Conn:        conn,
		LastUsed:    time.Now(),
		ConnectedAt: time.Now(),
		Mu:          sync.RWMutex{},
		Interactions: 1,
	}
}


// Metrics
type MonitorComputeStat struct{
	Username string
	Os string
	Architecture string
	Cpus int
	TotalRam uint64
	FreeRam uint64
	UsedRamPercent float64
	TotalGoRoutines int
}



//  Constant Folders section
func DbFolder()  string{
	HomePath , err := os.UserHomeDir()
	dbFolder := filepath.Join(HomePath , ".roxkv" , "roxdb") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return dbFolder

}


func LogFolder() string{
	HomePath , err := os.UserHomeDir()
	logFolder := filepath.Join(HomePath , ".roxkv" , "roxlogs") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return logFolder
}

func SnapshotFolder() string{
	HomePath , err := os.UserHomeDir()
	snapFolder := filepath.Join(HomePath , ".roxkv" , "roxsnaps") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return snapFolder	
}

func MetricFolder() string{
	HomePath , err := os.UserHomeDir()
	metricFolder := filepath.Join(HomePath , ".roxkv" , "roxmetrics") 
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return metricFolder
}