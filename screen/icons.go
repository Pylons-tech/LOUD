package screen

import (
	"strings"
	"unicode/utf8"
)


func (screen *GameScreen) pylonIcon() string {
	// return "🔶"
	return "🔷"
}

func (screen *GameScreen) goldIcon() string {
	return "💰"
}

func NumberOfSpaces(message string) int {
	customUnicodes := map[string]string{
		"💰":  "xx",
		"🔶":  "xx",
		"🔷":  "xx",
		"🥺":  "xx",
		"🗡️":  "x",
		"🦘":  "xx",
		"⟳":  "x",
		"📋":  "xx",
		"🥇":  "xx",
		"❦":   "x",
		"↓":   "x",
	}
	for k, v := range customUnicodes {
		message = strings.ReplaceAll(message, k, v)
	}
	return utf8.RuneCountInString(message)
}