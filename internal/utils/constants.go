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
type Mode struct {
	Name string
	Topic []string
}

const (
	RoleSystem Role = "system"
	RoleClient Role = "client"
	RoleAdmin  Role = "admin"
)

var ModeDefault = Mode{
	Name: "default",
	Topic: nil,
}

var ModeSubsriber = Mode{
	Name: "subscriber",
	Topic: []string{},
}
type NewClient struct {
	ID           any          `json:"id"`
	Role         string       `json:"role"`
	Conn         net.Conn     `json:"-"`
	LastUsed     time.Time    `json:"lastUsed"`
	ConnectedAt  time.Time    `json:"connectedAt"`
	Mu           sync.RWMutex `json:"-"`
	Interactions int          `json:"interactions"`
	Mode 		 Mode		  `json:"mode"`
}

type ClientMetaData struct {
	Id           string    `json:"id"`
	LastUsed     time.Time `json:"lastUsed"`
	ConnectedAt  time.Time `json:"connectedAt"`
	Interactions int       `json:"interactions"`
}

func CreateClient(conn net.Conn, role string) *NewClient {

	if role != string(RoleAdmin) && role != string(RoleClient) && role != string(RoleSystem) {
		return nil
	}

	return &NewClient{
		ID:           conn.RemoteAddr().String(),
		Role:         role,
		Conn:         conn,
		LastUsed:     time.Now(),
		ConnectedAt:  time.Now(),
		Mu:           sync.RWMutex{},
		Interactions: 1,
		Mode: ModeDefault,
	}
}

// Metrics
type MonitorComputeStat struct {
	Username        string  `json:"username"`
	Os              string  `json:"os"`
	Architecture    string  `json:"architecture"`
	Cpus            int     `json:"cpus"`
	TotalRam        uint64  `json:"totalRam"`
	FreeRam         uint64  `json:"freeRam"`
	UsedRamPercent  float64 `json:"usedRamPercent"`
	TotalGoRoutines int     `json:"totalGoRoutines"`
	ServerUptime	time.Duration `json:"serverUptime"`
}

// Constant Folders section
func DbFolder() string {
	HomePath, err := os.UserHomeDir()
	dbFolder := filepath.Join(HomePath, ".roxkv", "roxdb")
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return dbFolder

}

func LogFolder() string {
	HomePath, err := os.UserHomeDir()
	logFolder := filepath.Join(HomePath, ".roxkv", "roxlogs")
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return logFolder
}

func SnapshotFolder() string {
	HomePath, err := os.UserHomeDir()
	snapFolder := filepath.Join(HomePath, ".roxkv", "roxsnaps")
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return snapFolder
}

func MetricFolder() string {
	HomePath, err := os.UserHomeDir()
	metricFolder := filepath.Join(HomePath, ".roxkv", "roxmetrics")
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}
	return metricFolder
}
