package errors

const (
	ChannelNotClosedCode uint16 = iota
	TcpAcceptCode
	TcpReadCode
	WsReadMessageCode
	WsWriteMessageCode
	WsWritePingMessageCode
	WsInvalidMessageCode
	WsNoPongCode
	HttpNotFoundCode = 404
)
