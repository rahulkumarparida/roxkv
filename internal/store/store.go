package store

import (
	"sync"
)

// The structure of the value to be stored
type Item struct {
	Key string
	Val any
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

	ma.Data[kv.Key] = Item{
		kv.Key,
		kv.Val}

}


// This 
func (ma *MemoryAlloc) Get(key string) (Item,bool){
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()

	data , exist := ma.Data[key]

	if !exist {
		return Item{} ,false
	}

	return data , true

}


func (ma *MemoryAlloc) Del(key string) (bool){
	ma.Mu.RLock()
	defer ma.Mu.RUnlock()

	_ , exist  := ma.Data[key]

	if !exist {
		return  true
	}
	// Because  it is a map  it can store multtiple valuse so we dont want to harm other values insteda of the key value
	delete(ma.Data,key)
	
	return  true

}


func (ma *MemoryAlloc) Keys() []string{
	data := ma.Data
	keys := []string{}

	if len(data) == 0 {
		return []string{}
	}

	for k := range data{
		keys = append(keys, k)
	}

	return  keys
}

// This one initializes the making of data returns a memory address ot type MemoryAlloc we can access all the methods belonging to MemoryAlloc
func StoreInMemory() (*MemoryAlloc){

	store := &MemoryAlloc{
		Data: make(map[string]Item),
	}

	return store
}


// Uses the function from MemoryAlloc and kv Item struct to Set a variable
func SetKv(store *MemoryAlloc,kv *Item ){	
	values := Item{kv.Key,kv.Val}

	store.Set(values)

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

