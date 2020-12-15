package multiplexer

var (
	EventHandlers     []interface{}
	Router            = New()
	Commands          []*Route
	NotTargeted       []interface{}
	GuildMemberAdd    []interface{}
	GuildMemberRemove []interface{}
	GuildDelete       []interface{}
)
