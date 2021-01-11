package overrides

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

var (
	simpleEntries  []SimpleConfigurationEntry
	complexEntries []ComplexConfigurationEntry
	customEntries  []CustomConfigurationEntry
)

type SimpleConfigurationEntry struct {
	Name         string
	FriendlyName string
	Description  string
	DatabaseKey  string
	Cleanup      func(context *multiplexer.Context)
	Validate     func(context *multiplexer.Context, input *string) (bool, bool)
	Format       func(context *multiplexer.Context, value string) (string, string, bool)
}

type ComplexConfigurationEntry struct {
	Name         string
	FriendlyName string
	Description  string
	Entries      []SimpleConfigurationEntry
}

type CustomConfigurationEntry struct {
	Name        string
	Description string
	Handler     func(context *multiplexer.Context)
}

func GetSimpleEntries() []SimpleConfigurationEntry {
	return simpleEntries
}

func RegisterSimpleEntry(entry SimpleConfigurationEntry) {
	simpleEntries = append(simpleEntries, entry)
}

func GetComplexEntries() []ComplexConfigurationEntry {
	return complexEntries
}

func RegisterComplexEntry(entry ComplexConfigurationEntry) {
	complexEntries = append(complexEntries, entry)
}

func GetCustomEntries() []CustomConfigurationEntry {
	return customEntries
}

func RegisterCustomEntry(entry CustomConfigurationEntry) {
	customEntries = append(customEntries, entry)
}
