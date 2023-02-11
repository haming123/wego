package gows

const (
	Frame_Null     = -1
	Frame_Continue = 0
	Frame_Text     = 1
	Frame_Binary   = 2
	Frame_Close    = 8
	Frame_Ping     = 9
	Frame_Pong     = 10
)

func isControlFrame(frameType int) bool {
	if frameType == Frame_Close {
		return true
	}
	if frameType == Frame_Ping {
		return true
	}
	if frameType == Frame_Pong {
		return true
	}
	return false
}

func isMessageFrame(frameType int) bool {
	if frameType == Frame_Binary {
		return true
	}
	if frameType == Frame_Text {
		return true
	}
	if frameType == Frame_Continue {
		return true
	}
	return false
}
