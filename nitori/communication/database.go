package communication

import (
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"github.com/dgraph-io/badger/v2"
	"strings"
)

func size() int64 {
	lsm, vlog := SuperVisor.Database.Size()
	return lsm + vlog
}

func gc() error {
	var err error
	for {
		err = SuperVisor.Database.RunValueLogGC(0.5)
		if err != nil {
			break
		}
	}
	return err
}

func set(k, v string) error {
	return SuperVisor.Database.Update(func(txn *badger.Txn) (err error) {
		return txn.Set([]byte(k), []byte(v))
	})
}

func get(k string) (string, error) {
	var data string

	err := SuperVisor.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
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

func del(keys []string) error {
	return SuperVisor.Database.Update(func(txn *badger.Txn) error {
		for _, key := range keys {
			err := txn.Delete([]byte(key))
			if err != nil {
				break
			}
		}

		return err
	})
}

func seek(offset string, includeOffset bool, iterator *badger.Iterator) {
	if offset != "" {
		iterator.Seek([]byte(offset))
		if includeOffset && iterator.Valid() {
			iterator.Next()
		}
	} else {
		iterator.Rewind()
	}
}

func hset(hashmap, key, value string) error {
	err := set(hashmap+"/{HASH}/"+key, value)
	return err
}

func hget(hashmap, key string) (string, error) {
	return get(hashmap + "/{HASH}/" + key)
}

func hdel(hashmap string, keys []string) error {
	if len(keys) > 0 {
		for i, key := range keys {
			keys[i] = hashmap + "/{HASH}/" + key
		}
	} else {
		err := iter(false, true, hashmap+"/{HASH}/", hashmap+"/{HASH}/", func(key, _ string) bool {
			keys = append(keys, key)
			return true
		})
		if err != nil {
			return err
		}
	}
	return del(keys)
}

func hgetall(hashmap string) (map[string]string, error) {
	result := map[string]string{}
	err := iter(true, true, hashmap, hashmap,
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

func hkeys(hashmap string) ([]string, error) {
	var result []string
	err := iter(false, true, hashmap, hashmap,
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

func hlen(hashmap string) (int, error) {
	length := 0
	err := iter(false, true, hashmap, hashmap,
		func(_, _ string) bool {
			length++
			return true
		})
	return length, err
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

func iter(prefetch, includeOffset bool, offset, prefix string, handler func(key, value string) bool) error {
	return SuperVisor.Database.View(func(txn *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		options.PrefetchValues = prefetch
		iterator := txn.NewIterator(options)
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
