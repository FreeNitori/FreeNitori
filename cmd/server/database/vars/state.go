package vars

// Interface of a database backend
type Backend interface {
	DBType() string
	Open(path string) error
	Close() error
	Size() int64
	GC() error
	Set(key, value string) error
	Get(key string) (string, error)
	Del(keys []string) error
	HSet(hashmap, key, value string) error
	HGet(hashmap, key string) (string, error)
	HDel(hashmap string, keys []string) error
	HGetAll(hashmap string) (map[string]string, error)
	HKeys(hashmap string) ([]string, error)
	HLen(hashmap string) (int, error)
	Iter(prefetch, includeOffset bool, offset, prefix string, handler func(key, value string) bool) error
}

var Database Backend
