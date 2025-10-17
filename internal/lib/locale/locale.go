package locale

import (
	"github.com/BurntSushi/toml"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type M map[string]interface{}

func SetupLocale() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile("resources/en/message.en.toml")
	if err != nil {
		panic(err.Error())
	}
	_, err = bundle.LoadMessageFile("resources/ru/message.ru.toml")
	if err != nil {
		panic(err.Error())
	}
	return bundle
}

func L(loc *i18n.Localizer, id string, data M) string {
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return id
	}
	return msg
}

func LCtx(c *app.RequestContext, id string, data M) string {
	loc := GetLocalizer(c)
	if loc == nil {
		return id
	}

	return L(loc, id, data)
}

func GetLocalizer(c *app.RequestContext) *i18n.Localizer {
	if v, ok := c.Get("localizer"); ok {
		if loc, ok := v.(*i18n.Localizer); ok {
			return loc
		}
	}
	return nil
}
