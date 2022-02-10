package log

import (
	"strings"
)

type Level int

const (
	LOG_OFF Level = iota
	LOG_FATAL
	LOG_ERROR
	LOG_WARN
	LOG_INFO
	LOG_DEBUG
)

var levelStrings = [...]string{"[O]", "[F]", "[E]", "[W]", "[I]", "[D]"}
func (lv Level) String() string {
	if lv < 0 || int(lv) > 5 {
		return "[ ]"
	}
	return levelStrings[int(lv)]
}

func ParseLogLevel(val string) Level {
	val = strings.ToUpper(val)
	switch val {
	case "FATAL":
		return LOG_FATAL
	case "ERROR":
		return LOG_ERROR
	case "WARN":
		return LOG_WARN
	case "INFO":
		return LOG_INFO
	case "DEBUG":
		return LOG_DEBUG
	default:
		return LOG_OFF
	}
}