package utils

import (
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type NewClient struct {
	ID          any
	Conn        net.Conn
	LastUsed    time.Time
	ConnectedAt time.Time
	Mu          sync.RWMutex
}


func CreateClient(conn net.Conn) *NewClient {

	return &NewClient{
		ID:          conn.RemoteAddr().String(),
		Conn:        conn,
		LastUsed:    time.Now(),
		ConnectedAt: time.Now(),
		Mu:          sync.RWMutex{},
	}
}



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