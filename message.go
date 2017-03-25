package main

type msgType string

const (
	msgCmdType msgType = "A"
	msgResType msgType = "B"
)

type msg struct {
	kind  msgType
	value string
	bytes []byte
}

func parseMsg(b []byte) msg {
	return msg{
		kind:  msgType(b[0]),
		value: string(b[1:]),
		bytes: b,
	}
}

func newMsgCmdType(val string) msg {
	return msg{
		kind:  msgCmdType,
		value: val,
		bytes: []byte(string(msgCmdType) + val),
	}
}

func newMsgResType(val string) msg {
	return msg{
		kind:  msgResType,
		value: val,
		bytes: []byte(string(msgResType) + val),
	}
}
