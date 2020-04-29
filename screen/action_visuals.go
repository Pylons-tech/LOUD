package screen

import (
	"fmt"
	"strings"

	loud "github.com/Pylons-tech/LOUD/data"
)

func (screen *GameScreen) buyLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		screen.goldIcon(),
		fmt.Sprintf("%v", loudValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellLoudDesc(loudValue interface{}, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.goldIcon(),
		fmt.Sprintf("%v", loudValue),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) buyItemDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatItem(activeItem),
	}, "")
	return desc
}

func (screen *GameScreen) buyItemSpecDesc(itemSpec loud.ItemSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatItemSpec(itemSpec),
	}, "")
	return desc
}

func (screen *GameScreen) buyCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatCharacter(activeCharacter),
	}, "")
	return desc
}

func (screen *GameScreen) buyCharacterSpecDesc(charSpec loud.CharacterSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
		"\n  ↓\n",
		formatCharacterSpec(charSpec),
	}, "")
	return desc
}

func (screen *GameScreen) sellItemDesc(activeItem loud.Item, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatItem(activeItem),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellCharacterDesc(activeCharacter loud.Character, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatCharacter(activeCharacter),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellItemSpecDesc(activeItem loud.ItemSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatItemSpec(activeItem),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) sellCharacterSpecDesc(activeCharacter loud.CharacterSpec, pylonValue interface{}) string {
	var desc = strings.Join([]string{
		"\n",
		formatCharacterSpec(activeCharacter),
		"\n  ↓\n",
		screen.pylonIcon(),
		fmt.Sprintf("%v", pylonValue),
	}, "")
	return desc
}

func (screen *GameScreen) tradeTableColorDesc(width int) []string {
	var infoLines = []string{screen.regularFont()(fillSpace("", width))}

	infoLines = append(infoLines, screen.regularFont()(fillSpace(loud.Localize("white trade line desc"), width)))
	infoLines = append(infoLines, screen.blueBoldFont()(fillSpace(loud.Localize("bluebold trade line desc"), width)))
	infoLines = append(infoLines, screen.brownBoldFont()(fillSpace(loud.Localize("brownbold trade line desc"), width)))
	infoLines = append(infoLines, screen.brownFont()(fillSpace(loud.Localize("brown trade line desc"), width)))
	return infoLines
}
