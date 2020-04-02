package loud

import (
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var GameLanguage string = "en"

func Localize(key string) string {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("locale/en.json")
	bundle.MustLoadMessageFile("locale/es.json")

	loc := i18n.NewLocalizer(bundle, GameLanguage)

	translate, err := loc.Localize(
		&i18n.LocalizeConfig{
			MessageID:   key,
			PluralCount: 1,
		})
	if err != nil {
		return key
	}
	return translate
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(Localize(format), a...)
}
