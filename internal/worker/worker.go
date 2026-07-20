package worker

import (
	"context"
	"slices"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/store"
)


func RemoveExpired(t time.Time,ms *store.MemoryAlloc){
	ms.Mu.Lock()
	vals := ms.Data
	ms.Mu.Unlock()	


	if len(vals) == 0 {
		return 
	}

	for _, key := range vals {
		if !key.Meta.TTL.IsZero() && key.Meta.TTL.Before(time.Now()) {
			logger.InfoLog("Removed "+key.Key+" it expired.")
			store.DelKv(ms , key.Key)
			
			store.TTLMetricsContainer.TotalExpiredKeys +=1
			store.TTLMetricsContainer.ExpiresToday += 1
			var index int
			for idx , ele := range store.TTLMetricsContainer.NextExpiringKeys{
				if ele.Key == key.Key {
					index = idx	
					break
				}
			}

			store.TTLMetricsContainer.NextExpiringKeys = slices.Delete(store.TTLMetricsContainer.NextExpiringKeys,index,index)

			


			logger.InfoLog("Ticker worked at : "+ t.Format("2006-01-02 15:04:05")) 
		}
	}
	
}	

func ExpiryWorker(ctx context.Context,ms *store.MemoryAlloc) {
	ticker := time.NewTicker(1*time.Second)
	defer ticker.Stop()		
		for{

			select{
				case <-ctx.Done():
					return
				case t := <-ticker.C :
					RemoveExpired(t,ms)
					
			}
		}
	
}

