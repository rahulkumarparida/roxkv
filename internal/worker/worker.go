package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
)


func RemoveExpired(t time.Time,ms *store.MemoryAlloc){
	vals := ms.Data

	if len(vals) == 0 {
		return 
	}

	for _, key := range vals {
		if !key.Ttl.IsZero() && key.Ttl.Before(time.Now()) {
			logger.InfoLog("Removed "+key.Key+" it expired.")
			store.DelKv(ms , key.Key)
			logger.InfoLog("Ticker worked at : "+ t.Format("2006-01-02 15:04:05")) 
		}
	}
	
}	

func ExpiryWorker(ctx context.Context,ms *store.MemoryAlloc) {
	ticker := time.NewTicker(1*time.Second)
	defer ticker.Stop()		
	
	
	
		fmt.Println("expirying")
		for{

			select{
				case <-ctx.Done():
					return
				case t := <-ticker.C :
					RemoveExpired(t,ms)
					
			}
		}
	
}

