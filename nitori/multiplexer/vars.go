package multiplexer

var (
	EventHandlers     []interface{}
	Router            = New()
	NotTargeted       []interface{}
	GuildMemberAdd    []interface{}
	GuildMemberRemove []interface{}
	GuildDelete       []interface{}
)
