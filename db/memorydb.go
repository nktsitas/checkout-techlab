package db

import "sync"

type memoryDB struct {
	store_map map[string]interface{}
	mu sync.Mutex
}

func InitMemoryDB() *memoryDB {
	return &memoryDB{
		store_map: make(map[string]interface{}),
	}
}

func (mdb *memoryDB) StoreItem(id string, item interface{}) {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	
	mdb.store_map[id] = item
}

func (mdb *memoryDB) FetchItem(id string) interface{} {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	
	item := mdb.store_map[id]
	return item
}

func (mdb *memoryDB) DeleteItem(id string) {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	
	delete(mdb.store_map, id)
}