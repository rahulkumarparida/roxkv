package worker

import (
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
)


func RemoveExpired(ms *store.MemoryAlloc){
	vals := ms.Data

	if len(vals) == 0 {
		return 
	}

	for _, key := range vals {
		if !key.Ttl.IsZero() && key.Ttl.Before(time.Now()) {
			logger.InfoLog("Removed "+key.Key+" it expired.")
			store.DelKv(ms , key.Key)
		}
	}
	
}	

func ExpiryWorker(ms *store.MemoryAlloc) {
	ticker := time.NewTicker(1*time.Second)
	defer ticker.Stop()


	expired := make(chan bool)


	go func(){
		for{

			select{
				case <-expired :
					logger.InfoLog("Worker did some work")
					return
				case t := <-ticker.C :
					RemoveExpired(ms)
					logger.InfoLog("Ticker worked at : "+ t.Format("2006-01-02 15:04:05")) 
					
			}
		}
	}()
}