// Badger database backend.
package badger

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"strings"
)

var err error

var Database Badger

type Badger struct {
	DB *badger.DB
}

func (db *Badger) DBType() string {
	return "Badger"
}

func (db *Badger) Open(path string) error {
	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.Logger = log.Logger
	opts.ValueDir = path
	opts.Truncate = true
	opts.SyncWrites = false
	opts.TableLoadingMode = options.MemoryMap
	opts.ValueLogLoadingMode = options.FileIO
	opts.NumMemtables = 2
	opts.MaxTableSize = 10 << 20
	opts.NumLevelZeroTables = 2
	opts.ValueThreshold = 1

	db.DB, err = badger.Open(opts)

	if err != nil {
		return err
	}

	go (func() {
		for db.DB.RunValueLogGC(0.5) == nil {
		}
	})()

	return nil
}

func (db *Badger) Close() error {
	return db.DB.Close()
}

func (db *Badger) Size() int64 {
	lsm, vlog := db.DB.Size()
	return lsm + vlog
}

func (db *Badger) GC() error {
	var err error
	for {
		err = db.DB.RunValueLogGC(0.5)
		if err != nil {
			break
		}
	}
	return err
}

func (db *Badger) Set(key, value string) error {
	return db.DB.Update(func(txn *badger.Txn) (err error) {
		return txn.Set([]byte(key), []byte(value))
	})
}

func (db *Badger) Get(key string) (string, error) {
	var data string

	err := db.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		data = string(val)

		return nil
	})

	return data, err
}

func (db *Badger) Del(keys []string) error {
	return db.DB.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			err := txn.Delete([]byte(key))
			if err != nil {
				break
			}
		}

		return err
	})
}

func (db *Badger) HSet(hashmap, key, value string) error {
	err := db.Set(hashmap+"/{HASH}/"+key, value)
	return err
}

func (db *Badger) HGet(hashmap, key string) (string, error) {
	return db.Get(hashmap + "/{HASH}/" + key)
}

func (db *Badger) HDel(hashmap string, keys []string) error {
	if len(keys) > 0 {
		for i, key := range keys {
			keys[i] = hashmap + "/{HASH}/" + key
		}
	} else {
		err := db.Iter(false, true, hashmap+"/{HASH}/", hashmap+"/{HASH}/", func(key, _ string) bool {
			keys = append(keys, key)
			return true
		})
		if err != nil {
			return err
		}
	}
	return db.Del(keys)
}

func (db *Badger) HGetAll(hashmap string) (map[string]string, error) {
	result := map[string]string{}
	err := db.Iter(true, true, hashmap+"/{HASH}/", hashmap+"/{HASH}/",
		func(key, value string) bool {
			fields := strings.SplitN(key, "/{HASH}/", 2)
			if len(fields) < 2 {
				return true
			}
			result[fields[1]] = value
			return true
		})
	return result, err
}

func (db *Badger) HKeys(hashmap string) ([]string, error) {
	var result []string
	err := db.Iter(false, true, hashmap+"/{HASH}/", hashmap+"/{HASH}/",
		func(key, _ string) bool {
			fields := strings.SplitN(key, "/{HASH}/", 2)
			if len(fields) < 2 {
				return true
			}
			result = append(result, fields[1])
			return true
		})
	return result, err
}

func (db *Badger) HLen(hashmap string) (int, error) {
	length := 0
	err := db.Iter(false, true, hashmap+"/{HASH}/", hashmap+"/{HASH}/",
		func(_, _ string) bool {
			length++
			return true
		})
	return length, err
}

func seek(offset string, includeOffset bool, iterator *badger.Iterator) {
	if offset == "" {
		iterator.Rewind()
	} else {
		iterator.Seek([]byte(offset))
		if !includeOffset && iterator.Valid() {
			iterator.Next()
		}
	}
}

func validate(prefix string, iterator *badger.Iterator) bool {
	if !iterator.Valid() {
		return false
	}
	if prefix != "" && !iterator.ValidForPrefix([]byte(prefix)) {
		return false
	}
	return true
}

func (db *Badger) Iter(prefetch, includeOffset bool, offset, prefix string, handler func(key, value string) bool) error {
	return db.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = prefetch
		iterator := txn.NewIterator(opts)
		defer iterator.Close()
		for seek(offset, includeOffset, iterator); validate(prefix, iterator); iterator.Next() {
			var key, value []byte
			item := iterator.Item()
			key = item.KeyCopy(nil)

			if prefetch {
				value, _ = item.ValueCopy(nil)
			}

			if !handler(string(key), string(value)) {
				break
			}
		}
		return nil
	})
}
