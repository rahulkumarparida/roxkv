package store

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)




// TTl metrics
type TTLInfo struct{
	Key string
	ExpiresIn time.Duration
}

type TTLMetrics struct {
    ActiveTTLKeys      int
    PermanentKeys      int
    TotalExpiredKeys   uint64
    ExpiresToday       uint64
    NextExpiringKeys   []TTLInfo
	Mu sync.RWMutex
}

var TTLMetricsContainer *TTLMetrics = &TTLMetrics{
    NextExpiringKeys: make([]TTLInfo, 0), // Initialize slices too!
}
// The structure of the value to be stored
type Metadata struct {
    TTL time.Time
    UpdatedAt time.Time	
    CreatedAt time.Time
    LastAcessedBy *utils.NewClient
    KeyAccessCount int64
    Size int64
    Namespace string
}

type Item struct {
	Key string
	Val any // should be converted to []byte only
	Meta Metadata
}

// This one struturizes on how the data wil be stored and allocated
type MemoryAlloc struct{
	Mu sync.RWMutex
	Data map[string]Item
}

// This refrencfes to memory alloc any one can use it by refrencing to MemoryAlloc
func (ma *MemoryAlloc) Set(kv Item) {

	ma.Mu.Lock()
	defer ma.Mu.Unlock()
	
	ma.Data[kv.Key] = kv

	logger.SucessLog(kv.Key + " added to the memory")

}


// This 
func (ma *MemoryAlloc) Get(key string) (Item,bool){
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()

	data , exist := ma.Data[key]

	if !exist {
		fmt.Println("Key dows not exsist")
		logger.InfoLog("Key does not exists")
		return Item{} ,false
	}

	return data , true

}


func (ma *MemoryAlloc) Del(key string) (bool){
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()

	_ , exist  := ma.Data[key]

	if !exist {
		logger.InfoLog("Key does not exists")
		return  true
	}
	// Because ma.Data is a array of map  it can store multtiple valuse so we dont want to harm other values insteda of the key value
	delete(ma.Data,key)
	logger.InfoLog("Deleted Key "+key)
	
	return  true

}


func (ma *MemoryAlloc) Keys() []string{
	data := ma.Data
	keys := []string{}

	if len(data) == 0 {
		return []string{}
	}

	for k , d := range data{
		if !d.Meta.TTL.IsZero() && d.Meta.TTL.Before(time.Now()) {
			DelKv(ma, d.Key)
			continue
		}
		keys = append(keys, k)
	}
	logger.InfoLog("All keys summoned")
	return  keys
}

// This one initializes the making of data returns a memory address ot type MemoryAlloc we can access all the methods belonging to MemoryAlloc
func StoreInMemory() (*MemoryAlloc,*NameSpace){

	store := &MemoryAlloc{
		Data: make(map[string]Item),
	}
	ns := &NameSpace{
		Data: make(map[string]int),
	}

	return store,ns
}

// Name space 
type NameSpace struct{
	Mu sync.RWMutex
	Data map[string]int
}

func (ns *NameSpace) SetNameSpace(name string){
	ns.Mu.Lock()
	defer ns.Mu.Unlock()

	if ns.Data == nil {
		ns.Data = make(map[string]int)
	}

	ns.Data[name] += 1


}

func (ns *NameSpace) GetNamespace(name string)(string,int){
	if ns.Data == nil {
		ns.Data = make(map[string]int)
	}

	return name , ns.Data[name]
}





// Uses the function from MemoryAlloc and kv Item struct to Set a variable
func SetKv(store *MemoryAlloc,namespace *NameSpace,kv *Item) bool{	

	
	data ,exist := store.Get(kv.Key)

	var values Item
	var ns string
	
	if exist {

		if strings.Contains(data.Key,":") {
			part := strings.Split(data.Key,":")
			ns = part[0]
		}else{
			ns = data.Key
		}
		

		meta := Metadata{
			kv.Meta.TTL,
			time.Now(),
			data.Meta.CreatedAt,
			kv.Meta.LastAcessedBy,
			data.Meta.KeyAccessCount+1,
			kv.Meta.Size,
			ns,
		}

		values = Item{
			data.Key,
			kv.Val,
			meta,
		}

	}else{
		if strings.Contains(kv.Key,":") {
			part := strings.Split(kv.Key,":")
			ns = part[0]
		}else{
			ns = kv.Key
		}
		
		// meta := Metadata{
		// 	kv.Meta.TTL,
		// 	time.Now(),
		// 	time.Now(),
		// 	kv.Meta.LastAcessedBy,
		// 	1,
		// 	kv.Meta.Size,
		// 	ns,
		// }


		values = *kv

		
	}
	
	store.Set(values)
	namespace.SetNameSpace(kv.Key)
	return  true



}

// Retieves the similar value from the memory and sends it back
func GetKv(store *MemoryAlloc, key string) Item{

	data , exists := store.Get(key)

	if !exists{
		return data
	}


	return data
}



// This deletes the Value  from the memory (deletes it down )
func DelKv(store *MemoryAlloc,key string) bool{

	deleted := store.Del(key)

	return  deleted
}


// This will 
func KeyKv(store *MemoryAlloc) []string{
	
	data := store.Keys()
	

	return data
}
