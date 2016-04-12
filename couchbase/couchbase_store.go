package couchbase

import (
	"github.com/bushwood/caddyshack"
	"github.com/bushwood/caddyshack/model"
	"github.com/couchbase/gocb"
)

type CouchbaseStore struct {
	name           string
	model          *model.Definition
	host           string
	bucketName     string
	bucketPassword string
	bucket         *gocb.Bucket
}

func NewCouchbaseStore(host, bucketName string) *CouchbaseStore {
	return &CouchbaseStore{
		name:       "couchbase-store",
		host:       host,
		bucketName: bucketName,
	}
}

func (c *CouchbaseStore) ConnectBucket() error {
	cluster, err := gocb.Connect(c.host)
	if err != nil {
		return err
	}

	b, err := cluster.OpenBucket(c.bucketName, c.bucketPassword)
	if err != nil {
		return err
	}

	c.bucket = b
	return nil
}

func (c *CouchbaseStore) ShutdownBucket() {
	c.bucket.Close()
}

func (c *CouchbaseStore) GetName() string {
	return c.name
}

func (c *CouchbaseStore) SetName(name string) error {
	c.name = name
	return nil
}

// Model will be defined in the future, so perform error checking then.
func (c *CouchbaseStore) Init(model *model.Definition) (error, *CouchbaseStore) {
	c.model = model
	return nil, c
}

func (c *CouchbaseStore) Create(obj *CouchbaseObject) error {
	id, err := c.bucket.Insert(obj.key, obj.data, obj.expiry)
	if err != nil {
		return err
	}

	// This is the unique ID returned by Couchbase.
	obj.objectId = id
	return nil
}

func (c *CouchbaseStore) ReadOne(key string) (error, *CouchbaseObject) {
	var data interface{}
	id, err := c.bucket.Get(key, &data)
	if err != nil {
		return err, nil
	}

	obj := &CouchbaseObject{
		key:      key,
		data:     data,
		objectId: id,
	}

	return nil, obj
}

func (c *CouchbaseStore) UpdateOne(obj *CouchbaseObject) error {
	id, err := c.bucket.Replace(obj.key, obj.data, obj.objectId, obj.expiry)
	if err != nil {
		return err
	}

	// Set the new ID for the document.
	obj.objectId = id
	return nil
}

func (c *CouchbaseStore) DestroyOne(key string) error {
	// We do not need to keep the ID that this returns.
	_, err := c.bucket.Remove(key, 0)
	if err != nil {
		return err
	}

	return nil
}

func (c *CouchbaseStore) Read(q caddyshack.Query) (error, []caddyshack.StoreObject) {
	return q.Execute()
}
