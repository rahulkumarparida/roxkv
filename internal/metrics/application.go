package metrics

import (
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/pubsub"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

func GetUptime() time.Duration {
	// time.start at a begning call
	// to the time asked query
	duration := time.Since(utils.ServerStarted)

	return duration
}

func GetConnectedClients() []utils.ClientMetaData{
	// Total clients Connected and Metadata from the server
	var clientsConnected []utils.ClientMetaData
	mu := sync.Mutex{}
	mu.Lock()
	connections := make([]*utils.NewClient, len(utils.TotalConnecntions))
	copy(connections, utils.TotalConnecntions)
	mu.Unlock()
	for _, client := range connections {
		client := utils.ClientMetaData{
			Id:	client.ID.(string),
			LastUsed: client.LastUsed,
			ConnectedAt: client.ConnectedAt,
			Interactions: client.Interactions,
		}
		clientsConnected = append(clientsConnected, client)
	}

	return clientsConnected
}

func GetTotalKeys(store *store.MemoryAlloc) []string{
	// All the keys keyskv
	store.Mu.Lock()
	keys := store.Keys()
	store.Mu.Unlock()	
	
	return keys
}

func GetCommandCount() int{
	// write a counter that is publicaly avaliable strcut with mutexes
	mu := sync.Mutex{}
	mu.Lock()
	input := utils.TotalInputs
	mu.Unlock()
	return input
}

func GetTopicCount(client *utils.NewClient) (int,[]string){
	// len(topics)
	mu := sync.Mutex{}
	mu.Lock()
	topics := pubsub.GetTopics(client)
	mu.Unlock()
	return len(topics),topics
}
