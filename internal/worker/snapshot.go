package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/persistence"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

type Snapstruct struct{
	Data []store.Item
	Client []ClientMetaData
	CreatedAt time.Time
}

type ClientMetaData struct{
	Id string
	LastUsed time.Time
	ConnectedAt time.Time
	Interactions int
}

func TakeSnapShot(ms *store.MemoryAlloc, snaptime time.Time,clientsconnected []*utils.NewClient){

	filename := time.Now().Format("2006-01-02 15:04:05")+".json"
	dirPath := utils.SnapshotFolder()	

		
		var gatherData Snapstruct
		var kvdata []store.Item
		var clientData []ClientMetaData
		keys := store.KeyKv(ms)
		
		
		for _, key := range keys {
			data := store.GetKv(ms,key)
			fmt.Println("Keys:",data)
			kvdata = append(kvdata, data)
		}

		for _, client := range clientsconnected {
			client := ClientMetaData{client.ID.(string),client.LastUsed,client.ConnectedAt,client.Interactions}
			clientData = append(clientData, client)


		}

		

		gatherData.Data = kvdata
		gatherData.Client = clientData
		gatherData.CreatedAt = snaptime

		persistence.StoreToJson(dirPath,filename,gatherData)
		
}

func SnapshotWorker(ctx context.Context,ms *store.MemoryAlloc,mu *sync.RWMutex,clients []*utils.NewClient){
	ticker := time.NewTicker(5*time.Minute)
	defer ticker.Stop()

	
	for {
		select{
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			mu.Lock()
				clientCopy :=make([]*utils.NewClient, len(clients))
				copy(clientCopy,clients)
			mu.Unlock()
			

			TakeSnapShot(ms,t,clientCopy)
			fmt.Println("Taken Snapss----------------->")			
		}


	}


}