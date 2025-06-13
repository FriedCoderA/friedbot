package kinds

const (
	MessageTypePrivate = "private"
	MessageTypeGroup   = "group"

	MessageSubTypeOther     = "other"
	MessageSubTypeFriend    = "friend"
	MessageSubTypeGroup     = "group"
	MessageSubTypeNormal    = "normal"
	MessageSubTypeAnonymous = "anonymous"
	MessageSubTypeNotice    = "notice"
)

const (
	SessionStateNormal int8 = iota
	SessionStatePause
)
