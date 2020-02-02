package loud

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var gameLanguage string = "en"

func localize(key string) string {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("locale/en.json")
	bundle.MustLoadMessageFile("locale/es.json")

	loc := i18n.NewLocalizer(bundle, gameLanguage)

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
