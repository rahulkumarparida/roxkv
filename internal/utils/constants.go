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
	Name   string
	Topic  []string
	PTopic []string
}

type InitSource struct {
	Source   string
	JoinedAt time.Time
}

const (
	RoleSystem Role = "system"
	RoleClient Role = "client"
	RoleAdmin  Role = "admin"
)

var ModeDefault = Mode{
	Name:   "default",
	Topic:  nil,
	PTopic: nil,
}

var ModeSubsriber = Mode{
	Name:   "subscriber",
	Topic:  []string{},
	PTopic: []string{},
}

var DefaultSource = InitSource{
	Source:   "local",
	JoinedAt: time.Now(),
}

var RedisSource = InitSource{
	Source:   "redis",
	JoinedAt: time.Now(),
}

type NewClient struct {
	ID           any          `json:"id"`
	Name         string       `json:"name"`
	Library_Name string       `json:"library_name"` // I a redis libraby connect we need to save the linraryname and version
	Library_Ver  string       `json:"library_ver"`
	Buffer       []byte       `json:"buffer"` // For client's data to be recie3
	Role         string       `json:"role"`
	Conn         net.Conn     `json:"-"`
	LastUsed     time.Time    `json:"lastUsed"`
	ConnectedAt  time.Time    `json:"connectedAt"`
	Mu           sync.RWMutex `json:"-"`
	Interactions int          `json:"interactions"`
	Mode         Mode         `json:"mode"`
	Initiator    InitSource   `json:"initiator"`
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
		Mode:         ModeDefault,
		Initiator:    DefaultSource,
	}
}

// Metrics
type MonitorComputeStat struct {
	Username        string        `json:"username"`
	Os              string        `json:"os"`
	Architecture    string        `json:"architecture"`
	Cpus            int           `json:"cpus"`
	TotalRam        uint64        `json:"totalRam"`
	FreeRam         uint64        `json:"freeRam"`
	UsedRamPercent  float64       `json:"usedRamPercent"`
	TotalGoRoutines int           `json:"totalGoRoutines"`
	ServerUptime    time.Duration `json:"serverUptime"`
}

func roxKVRoot() string {
	if customRoot := os.Getenv("ROXKV_HOME"); customRoot != "" {
		return customRoot
	}

	homePath, err := os.UserHomeDir()
	if HandleError("Error while fetching Home directory ", err) {
		return ""
	}

	return filepath.Join(homePath, ".roxkv")
}

// Constant Folders section
func DbFolder() string {
	return filepath.Join(roxKVRoot(), "roxdb")

}

func LogFolder() string {
	return filepath.Join(roxKVRoot(), "roxlogs")
}

func SnapshotFolder() string {
	return filepath.Join(roxKVRoot(), "roxsnaps")
}

func MetricFolder() string {
	return filepath.Join(roxKVRoot(), "roxmetrics")
}

func ConfigFolder() string {
	return filepath.Join(roxKVRoot(), "roxconfig")
}
