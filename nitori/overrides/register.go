package overrides

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

var (
	simpleEntries  []SimpleConfigurationEntry
	complexEntries []ComplexConfigurationEntry
	customEntries  []CustomConfigurationEntry
)

// SimpleConfigurationEntry is a configuration entry with one item.
type SimpleConfigurationEntry struct {
	// Name is the name of the entry.
	Name string
	// FriendlyName is the friendly name of the entry.
	FriendlyName string
	// Description is the description of the entry.
	Description string
	// DatabaseKey is the database key of the entry.
	DatabaseKey string
	// Cleanup is the function called when the entry is reset.
	Cleanup func(context *multiplexer.Context)
	// Validate is the function called before setting the value, returning a validate success and an ok value.
	Validate func(context *multiplexer.Context, input *string) (bool, bool)
	// Format is the function called when generating preview embed, returning a title, description field and an ok value.
	Format func(context *multiplexer.Context, value string) (string, string, bool)
}

// ComplexConfigurationEntry is a configuration entry with multiple items.
type ComplexConfigurationEntry struct {
	// Name is the name of the entry.
	Name string
	// FriendlyName is the friendly name of the entry.
	FriendlyName string
	// Description is the description of the entry.
	Description string
	// Entries is the slice of simple entries of the entry.
	Entries []SimpleConfigurationEntry
	// CustomEntries is the slice of custom entries of the entry.
	CustomEntries []CustomConfigurationEntry
}

// CustomConfigurationEntry is a configuration entry that directly handles the configuration.
type CustomConfigurationEntry struct {
	// Name is the name of the entry.
	Name string
	// Description is the description of the entry.
	Description string
	// Handler is the handler function of the entry.
	Handler func(context *multiplexer.Context)
}

// GetSimpleEntries returns a slice of SimpleConfigurationEntry registered.
func GetSimpleEntries() []SimpleConfigurationEntry {
	return simpleEntries
}

// RegisterSimpleEntry registers a SimpleConfigurationEntry.
func RegisterSimpleEntry(entry SimpleConfigurationEntry) {
	simpleEntries = append(simpleEntries, entry)
}

// GetComplexEntries returns a slice of ComplexConfigurationEntry registered.
func GetComplexEntries() []ComplexConfigurationEntry {
	return complexEntries
}

// RegisterComplexEntry registers a ComplexConfigurationEntry.
func RegisterComplexEntry(entry ComplexConfigurationEntry) {
	complexEntries = append(complexEntries, entry)
}

// GetCustomEntries returns a slice of CustomConfigurationEntry registered.
func GetCustomEntries() []CustomConfigurationEntry {
	return customEntries
}

// RegisterCustomEntry registers a CustomConfigurationEntry.
func RegisterCustomEntry(entry CustomConfigurationEntry) {
	customEntries = append(customEntries, entry)
}
