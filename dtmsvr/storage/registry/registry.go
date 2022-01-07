package registry

import (
	"sync"
	"time"

	"github.com/dtm-labs/dtm/dtmsvr/config"
	"github.com/dtm-labs/dtm/dtmsvr/storage"
	"github.com/dtm-labs/dtm/dtmsvr/storage/boltdb"
	"github.com/dtm-labs/dtm/dtmsvr/storage/redis"
	"github.com/dtm-labs/dtm/dtmsvr/storage/sql"
)

var conf = &config.Config

type storeFactory struct {
	once    sync.Once
	store   storage.Store
	factory func() storage.Store
}

var factories map[string]*storeFactory = map[string]*storeFactory{
	"redis": {
		factory: func() storage.Store {
			return &redis.Store{}
		},
	},
	"mysql": {
		factory: func() storage.Store {
			return &sql.Store{}
		},
	},
	"postgres": {
		factory: func() storage.Store {
			return &sql.Store{}
		},
	},
	"boltdb": {
		factory: func() storage.Store {
			return &boltdb.Store{}
		},
	},
}

// GetStore returns storage.Store
func GetStore() storage.Store {
	f := factories[conf.Store.Driver]
	f.once.Do(func() {
		f.store = f.factory()
	})
	return f.store
}

// WaitStoreUp wait for db to go up
func WaitStoreUp() {
	for err := GetStore().Ping(); err != nil; err = GetStore().Ping() {
		time.Sleep(3 * time.Second)
	}
}
