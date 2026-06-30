package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rahulkumarparida/roxkv/internal/logger"
	"github.com/rahulkumarparida/roxkv/internal/utils"
)

// The structure of the value to be stored
type Item struct {
	Key string
	Val any // should be converted to []byte only
	Ttl time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
	LastAcessedBy *utils.NewClient
	KeyAccessCount int
	Size  int
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
		if !d.Ttl.IsZero() && d.Ttl.Before(time.Now()) {
			DelKv(ma, d.Key)
			continue
		}
		keys = append(keys, k)
	}
	logger.InfoLog("All keys summoned")
	return  keys
}

// This one initializes the making of data returns a memory address ot type MemoryAlloc we can access all the methods belonging to MemoryAlloc
func StoreInMemory() (*MemoryAlloc){

	store := &MemoryAlloc{
		Data: make(map[string]Item),
	}

	return store
}

func evaluatetime(ttl int,dur string) time.Time{
	
	
	switch strings.ToLower(dur) {
		case "second","sec":
			return time.Now().Add(time.Duration(ttl) * time.Second)
		case "minutes","min":
			return time.Now().Add(time.Duration(ttl) * time.Minute)
		case "hour","hh":
			return time.Now().Add(time.Duration(ttl) * time.Hour)
		default:
			return time.Now()
	}
}

// Uses the function from MemoryAlloc and kv Item struct to Set a variable
func SetKv(store *MemoryAlloc,kv *Item, stripTtl []string) bool{	
	var futureTime time.Time 
	if !kv.Ttl.IsZero() {
		if stripTtl[0] == "" || stripTtl[1] == "" {
			logger.ErrorLog("--ttl arguments not provided please check the man page for arguments")
			return  false
		}
			var timetoAdd , err = strconv.Atoi(stripTtl[0]) // time like 12 ,13 ,14

		if utils.HandleError("Error while parsing time provided", err){
			logger.ErrorLog("Parsing failed while parsing the time. Integre required")
			return false
		}

		var duration = stripTtl[1] // duration like sec , min, hr

		 futureTime = evaluatetime(timetoAdd,duration)
	}
	data ,exist := store.Get(kv.Key)

	var values Item
	if exist {
		// store.Mu.Lock()
		values = Item{
			data.Key,
			kv.Val,
			futureTime,
			time.Now(),
			data.CreatedAt,
			kv.LastAcessedBy,
			data.KeyAccessCount+1,
			kv.Size,
		}
		// store.Mu.Unlock()
		return true
	}else{
		values = Item{
			kv.Key,
			kv.Val,
			futureTime,
			time.Now(),
			time.Now(),
			kv.LastAcessedBy,
			1,
			kv.Size,
		}

		
	}

	store.Set(values)
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

