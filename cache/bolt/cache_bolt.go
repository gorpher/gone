package bolt

//
//import (
//	"encoding/json"
//	bolt "go.etcd.io/bbolt"
//	"path/filepath"
//	"time"
//)
//
//type BoltCache struct {
//	fs string
//}
//
//// NewBoltCache bbolt 需要考虑使用go1.22 以上版本
//func NewBoltCache(absPath string) *BoltCache {
//	return &BoltCache{filepath.Clean(absPath)}
//}
//
//func (b *BoltCache) Get(k string) ([]byte, error) {
//	db, err := bolt.Open(b.fs, 0666, nil)
//	if err != nil {
//		return nil, err
//	}
//	defer db.Close() //nolint
//	tx, err := db.Begin(true)
//	if err != nil {
//		return nil, err
//	}
//	defer tx.Commit() // nolint
//	var bucket *bolt.Bucket
//	bucket, err = tx.CreateBucketIfNotExists([]byte("authed"))
//	if err != nil {
//		return nil, err
//	}
//	return bucket.Get([]byte(k)), nil
//
//}
//
//func (b *BoltCache) Set(key string, value interface{}) error {
//	var v []byte
//	switch va := value.(type) {
//	case string:
//		v = []byte(va)
//	case []byte:
//		v = va
//	default:
//		data, err := json.Marshal(va)
//		if err != nil {
//			return err
//		}
//		v = data
//	}
//	db, err := bolt.Open(b.fs, 0666, nil)
//	if err != nil {
//		return err
//	}
//	defer db.Close() //nolint
//	err = db.Update(func(tx *bolt.Tx) error {
//		var bucket *bolt.Bucket
//		bucket, err = tx.CreateBucketIfNotExists([]byte("authed"))
//		if err != nil {
//			return err
//		}
//		return bucket.Put([]byte(key), v)
//	})
//	return err
//}
//
//func (b *BoltCache) Del(key string) error {
//	db, err := bolt.Open(b.fs, 0666, nil)
//	if err != nil {
//		return err
//	}
//	defer db.Close() //nolint
//	tx, err := db.Begin(true)
//	if err != nil {
//		return err
//	}
//	defer tx.Commit() // nolint
//	var bucket *bolt.Bucket
//	bucket, err = tx.CreateBucketIfNotExists([]byte("authed"))
//	if err != nil {
//		return err
//	}
//	return bucket.Delete([]byte(key))
//}
//
//func (b *BoltCache) SetWithTTL(key string, value string, duration time.Duration) error {
//	_ = time.AfterFunc(duration, func() {
//		b.Del(key) //nolint
//	})
//	return b.Set(key, value)
//}
