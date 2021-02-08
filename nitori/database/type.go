package database

// Backend represents a database backend.
type Backend interface {
	// DBType returns the name of the database as a string.
	DBType() string
	// Open opens the database.
	Open(path string) error
	// Close closes the database.
	Close() error
	// Size returns the size of the database.
	Size() int64
	// GC triggers a value log garbage collection.
	GC() error
	// Set adds a key-value pair to the database.
	Set(key, value string) error
	// Get gets the value of a key from the database.
	Get(key string) (string, error)
	// Del deletes a key from the database.
	Del(keys []string) error
	// HSet adds a key-value pair to a hashmap.
	HSet(hashmap, key, value string) error
	// HGet gets the value of a key from a hashmap.
	HGet(hashmap, key string) (string, error)
	// HDel deletes a key from a hashmap.
	HDel(hashmap string, keys []string) error
	// HGetAll gets all key-value pairs of a hashmap.
	HGetAll(hashmap string) (map[string]string, error)
	// HKeys gets all keys of a hashmap.
	HKeys(hashmap string) ([]string, error)
	// HLen gets the length of a hashmap.
	HLen(hashmap string) (int, error)
	// Iter iterates through stuff in the database.
	Iter(prefetch, includeOffset bool, offset, prefix string, handler func(key, value string) bool) error
}

// Database holds the currently used database backend.
var Database Backend
