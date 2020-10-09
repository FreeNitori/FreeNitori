package communication

import (
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"github.com/dgraph-io/badger/v2"
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
