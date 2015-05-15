package couchbase

import (
	"encoding/json"
	"fmt"
	couchbase "github.com/couchbaselabs/gocb"
	"github.com/jimonreal/confd/log"
	"os"
	"reflect"
	"strings"
)

// Client is a wrapper around the couchbase client
type Client struct {
	client *couchbase.Cluster
	bucket *couchbase.Bucket
	vQuery *couchbase.ViewQuery
}

type viewRow struct {
	Id    string
	Key   string
	Value json.RawMessage
}

// NewCouchbaseClient returns an *couchbase.Client with a connection to named machines.
// It returns an error if a connection to the cluster cannot be made.
func NewCouchbaseClient(machines []string) (*Client, error) {
	var err error
	clusterIPs := strings.Join(machines, ",")
	client, err := couchbase.Connect("http://" + clusterIPs)
	fmt.Println("http://" + clusterIPs)
	if err != nil {
		fmt.Println("Error: NewCouchbaseClient: ", err)
	}

	bucket, err := client.OpenBucket(os.Getenv("CB_BUCKET"), os.Getenv("CB_BUCKET_PASSWORD"))
	if err != nil {
		fmt.Println("Error: Unable to open Bucket: ", err)
	}

	vQuery := couchbase.NewViewQuery(os.Getenv("CB_DDOC"), os.Getenv("CB_VIEW_NAME"))
	vQuery = vQuery.Custom("full_set", "true")
	log.Info("Conected...")

	return &Client{client: client, bucket: bucket, vQuery: vQuery}, err
}

// GetValues queries couchbase for keys prefixed by prefix.
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	log.Debug("GetValues")
	rows := c.bucket.ExecuteViewQuery(c.vQuery)
	row := viewRow{}
	vars := make(map[string]string)
	var err error
	fmt.Printf("Keys: %+v\n\n", keys)

	for rows.Next(&row) {
		log.Debug("Inside the query")
		var tmpDoc Document
		err = json.Unmarshal(row.Value, &tmpDoc)
		if err != nil {
			//Fatal
			fmt.Printf("FATAL: Unable to unmarshal document: %+v\n", row.Value)
			return vars, err
		}
		log.Debug("---------")

		for _, key := range keys {
			tmpKey := strings.Replace(key, "/", "", 1)
			log.Debug("Iterating keys")
			value, ok := getField(tmpDoc, tmpKey)
			if !ok {
				//Fatal no existe la llave
				fmt.Printf("ERROR: Key %s doesn't exist\n")
				continue
			}
			newKey := key + "/" + row.Key
			vars[newKey] = value.(string)
			msg := "Key: " + newKey + " - Value: " + vars[newKey]
			log.Debug(msg)
		}
	}

	if err := rows.Close(); err != nil {
		fmt.Printf("View Query Error: %s", err)
		return vars, err
	}
	fmt.Println("Vars: %+v\n\n", vars)

	return vars, nil
}

// WatchPrefix is not yet implemented.
func (c *Client) WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	<-stopChan
	return 0, nil
}

func getField(from Document, name string) (interface{}, bool) {
	fmt.Printf("DEBUG: getField %s\n", name)
	value := reflect.ValueOf(from).FieldByName(name)
	if !value.IsValid() {
		return nil, false
	}
	return value.Interface(), true
}
