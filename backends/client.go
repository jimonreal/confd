package backends

import (
	"errors"
	"strings"

	"github.com/jimonreal/confd/backends/consul"
	"github.com/jimonreal/confd/backends/couchbase"
	"github.com/jimonreal/confd/backends/env"
	"github.com/jimonreal/confd/backends/etcd"
	"github.com/jimonreal/confd/backends/redis"
	"github.com/jimonreal/confd/backends/zookeeper"
	"github.com/jimonreal/confd/log"
)

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// New is used to create a storage client based on our configuration.
func New(config Config) (StoreClient, error) {
	if config.Backend == "" {
		config.Backend = "etcd"
	}
	backendNodes := config.BackendNodes
	log.Info("Backend nodes set to " + strings.Join(backendNodes, ", "))
	log.Info("Edited..")
	switch config.Backend {
	case "consul":
		return consul.New(config.BackendNodes, config.Scheme,
			config.ClientCert, config.ClientKey,
			config.ClientCaKeys)
	case "etcd":
		// Create the etcd client upfront and use it for the life of the process.
		// The etcdClient is an http.Client and designed to be reused.
		return etcd.NewEtcdClient(backendNodes, config.ClientCert, config.ClientKey, config.ClientCaKeys)
	case "zookeeper":
		return zookeeper.NewZookeeperClient(backendNodes)
	case "redis":
		return redis.NewRedisClient(backendNodes)
	case "couchbase":
		return couchbase.NewCouchbaseClient(backendNodes)
	case "env":
		return env.NewEnvClient()
	}
	return nil, errors.New("Invalid backend")
}
