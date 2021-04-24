// Package badger implements database backend based on badger.
package badger

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/dgraph-io/badger/v3"
	"os"
	"strings"
	"time"
)

var err error
var backupTicker *time.Ticker

// Database is an instance of the database backend.
var Database Badger

// Badger represents a badger database instance.
type Badger struct {
	DB *badger.DB
}

// DBType returns the name of the database as a string.
func (db *Badger) DBType() string {
	return "Badger"
}

// Open opens the database.
func (db *Badger) Open(path string) error {
	opts := badger.DefaultOptions(path)
	opts.Dir = path
	opts.Logger = log.Instance
	opts.ValueDir = path
	opts.SyncWrites = false
	opts.NumMemtables = 2
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

	if config.Config.System.BackupInterval > 0 {
		log.Infof("Periodical database backup interval %v second(s).", config.Config.System.BackupInterval)
		if _, err := os.Stat("backup"); os.IsNotExist(err) {
			err = os.Mkdir("backup", 0700)
			if err != nil {
				return err
			}
		}
		backupTicker = time.NewTicker(time.Duration(config.Config.System.BackupInterval) * time.Second)
		go func() {
			for {
				select {
				case <-backupTicker.C:
					intermediate, err := os.Create("backup/.intermediate")
					if err != nil {
						log.Errorf("Unable to create intermediate backup file, %s", err)
						continue
					}
					ver, err := db.DB.Backup(intermediate, 0)
					if err != nil {
						log.Errorf("Unable to generate intermediate backup, %s", err)
						continue
					}
					err = os.Rename("backup/.intermediate", fmt.Sprintf("backup/%v", ver))
					if err != nil {
						log.Errorf("Unable to rename intermediate backup, %s", err)
						err = os.Remove("backup/.intermediate")
						if err != nil {
							log.Errorf("Unable to remove intermediate file, %s", err)
							log.Warnf("Backup has been disabled, please resolve the issue above and restart FreeNitori.")
							break
						}
						continue
					}
					log.Infof("Successfully backed up database, version %v", ver)
				}
			}
		}()
	} else {
		log.Infof("Periodical database backup is not enabled.")
	}

	return nil
}

// Close closes the database.
func (db *Badger) Close() error {
	if backupTicker != nil {
		backupTicker.Stop()
	}
	return db.DB.Close()
}

// Size returns the size of the database.
func (db *Badger) Size() int64 {
	lsm, vlog := db.DB.Size()
	return lsm + vlog
}

// GC triggers a value log garbage collection.
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

// Set adds a key-value pair to the database.
func (db *Badger) Set(key, value string) error {
	return db.DB.Update(func(txn *badger.Txn) (err error) {
		return txn.Set([]byte(key), []byte(value))
	})
}

// Get gets the value of a key from the database.
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

// Del deletes a key from the database.
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

// HSet adds a key-value pair to a hashmap.
func (db *Badger) HSet(hashmap, key, value string) error {
	err := db.Set(hashmap+"/{HASH}/"+key, value)
	return err
}

// HGet gets the value of a key from a hashmap.
func (db *Badger) HGet(hashmap, key string) (string, error) {
	return db.Get(hashmap + "/{HASH}/" + key)
}

// HDel deletes a key from a hashmap.
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

// HGetAll gets all key-value pairs of a hashmap.
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

// HKeys gets all keys of a hashmap.
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

// HLen gets the length of a hashmap.
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

// Iter iterates through stuff in the database.
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
