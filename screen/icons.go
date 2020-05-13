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

// icons list

func NumberOfSpaces(message string) int {
	customUnicodes := map[string]string{
		"💰":  "xx", // Gold
		"🔶":  "xx", //
		"🔷":  "xx", // pylon
		"👀":  "xx", // leave emoji 👀
		"🗡️": "x", // sword
		"🐧":  "xx", // character emoji 🐧
		"⟳":  "x", // refresh emoji
		"📋":  "xx", // copy emoji
		"🥇":  "xx", // metal emoji
		"❦":  "x", // 
		"↓":  "x", // arrow emoji
		"🐉":  "xx", // Undead dragon
		"🦈":  "xx", // Ice dragon 🦈 
		"🦐":  "xx", // Fire dragon
		"🐊":  "xx", // Acid dragon 🐊 
		"🗿":  "xx", // Giant
		"👺":  "xx", // Goblin
		"🐺":  "xx", // Wolf
		"👻":  "xx", // Troll
		"🌊":  "xx", // Ice special
		"🔥":  "xx", // Fire special
		"🥗":  "xx", // Acid special
	}
	for k, v := range customUnicodes {
		message = strings.ReplaceAll(message, k, v)
	}
	return utf8.RuneCountInString(message)
}
