package cache

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

type BadgerCache struct {
	db *badger.DB
}

func NewBadgerCacheDB(db *badger.DB) *BadgerCache {
	return &BadgerCache{
		db: db,
	}
}

func NewBadgerCache(dir string, inMemory bool) (*BadgerCache, error) {
	options := badger.Options{
		Dir:      dir,
		ValueDir: dir,

		InMemory: inMemory,

		ValueLogFileSize:   102400000,
		ValueLogMaxEntries: 100000,
		VLogPercentile:     0.1,

		MemTableSize:                  64 << 20,
		BaseTableSize:                 2 << 20,
		BaseLevelSize:                 10 << 20,
		TableSizeMultiplier:           2,
		LevelSizeMultiplier:           10,
		MaxLevels:                     7,
		NumGoroutines:                 8,
		MetricsEnabled:                true,
		NumCompactors:                 4,
		NumLevelZeroTables:            5,
		NumLevelZeroTablesStall:       15,
		NumMemtables:                  5,
		BloomFalsePositive:            0.01,
		BlockSize:                     4 * 1024,
		SyncWrites:                    false,
		NumVersionsToKeep:             1,
		CompactL0OnClose:              false,
		VerifyValueChecksum:           false,
		BlockCacheSize:                256 << 20,
		IndexCacheSize:                0,
		ZSTDCompressionLevel:          1,
		EncryptionKey:                 []byte{},
		EncryptionKeyRotationDuration: 10 * 24 * time.Hour, // Default 10 days.
		DetectConflicts:               true,
		NamespaceOffset:               -1,
	}

	cache, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return NewBadgerCacheDB(cache), nil
}

func (c *BadgerCache) SetNX(key string, value interface{}) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if errors.Is(err, badger.ErrKeyNotFound) {
			v := []byte(fmt.Sprintf("%v", value))
			return txn.Set([]byte(key), v)
		}
		return err
	})
	return err
}

func (c *BadgerCache) Set(key string, value interface{}) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		v := []byte(fmt.Sprintf("%v", value))
		return txn.Set([]byte(key), v)
	})
	return err
}

func (c *BadgerCache) Del(key string) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}

func (c *BadgerCache) Clean() error {
	return c.db.DropAll()
}

func (c *BadgerCache) Get(key string) ([]byte, error) {
	var result []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			result = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return result, err
}

func (c *BadgerCache) SetWithTTL(key string, value string, duration time.Duration) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		v := []byte(fmt.Sprintf("%v", value))
		e := badger.NewEntry([]byte(key), v).WithTTL(duration)
		return txn.SetEntry(e)
	})
	return err
}

func (c *BadgerCache) PrefixScanKey(prefixStr string) ([]string, error) {
	var res []string
	err := c.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefixStr)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			res = append(res, string(k))
			return nil
		}
		return nil
	})
	return res, err
}
