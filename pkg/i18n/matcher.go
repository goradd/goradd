package i18n

import (
	"context"
	"github.com/goradd/goradd/pkg/goradd"
	"github.com/goradd/goradd/pkg/session"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// ServerLanguageEntry is a description of a supported language from the server's perspective. It includes the
// information the server will need to describe the language to the browser and to do the translation.
type ServerLanguageEntry struct {
	// Tag is the language tag for the language
	Tag language.Tag
	// Dict is the corresponding language dictionary, helping us to describe the language to users
	Dict *display.Dictionary
	// LangString is the string to display in the lang attribute of the html tag. Leave it blank to get the default from the Tag.
	LangString string
}

// languages is the list of languages that the application supports. By default we just support English, but you can
// change it to what you want
var languages = []language.Tag{
	language.AmericanEnglish, // first one is the default
}

// dictionaries are corresponding dictionaries
var dictionaries = []*display.Dictionary{
	display.English, // first one is the default
}

// langAttributes are the corresponding lang attributes to put in the html tag. This should be the canonical value of the language.
var langAttributes = []string{
	"en", // first one is the default
}

var matcher = language.NewMatcher([]language.Tag{language.AmericanEnglish})

// SetSupportedLanguages sets up the languages that the application supports. It expects both a list of language
// tags and a matching list of dictionaries. You should only call this during application startup to inject your
// list of supported languages into the application.
func SetSupportedLanguages(l ...ServerLanguageEntry) {
	if len(l) < 1 {
		panic("you must have at least one language")
	}
	languages = make([]language.Tag, len(l))
	dictionaries = make([]*display.Dictionary, len(l))
	langAttributes = make([]string, len(l))

	for i, e := range l {
		languages[i] = e.Tag
		dictionaries[i] = e.Dict
		if e.LangString == "" {
			langAttributes[i] = e.Tag.String()
		} else {
			langAttributes[i] = e.LangString
		}
	}

	// Setup a new matcher. Go doc says that matcher is optimized for runtime at the expense of init time.
	matcher = language.NewMatcher(languages)
}

// SupportedLanguages is returuned by GetSupported
type SupportedLanguages []struct {
	LocalName  string
	NativeName string
}

// GetSupportedLanguages returns a slice of the supported languages, in both the language indicated and the native
// representation of the name of that language. You could use this to present a menu to the user. The order is the
// same as the ServerLanguages and ServerDictionaries
func GetSupportedLanguages(t language.Tag) SupportedLanguages {
	_, i, _ := matcher.Match(t)

	d := dictionaries[i]
	l := d.Languages()
	s := make(SupportedLanguages, len(languages))

	for i, t := range languages {
		s[i].LocalName = l.Name(t)
		s[i].NativeName = display.Self.Name(t)
	}

	return s
}

// SetDefaultLanguage is called by the framework to set up the session variable with a default language if one has not
// yet been set. The default language is based on the "accept-language" header value and the list of languages that
// the application supports.
func SetDefaultLanguage(ctx context.Context, acceptLanguageValue string) int {
	if !session.Has(ctx, goradd.SessionLanguage) {
		tags, _, err := language.ParseAcceptLanguage(acceptLanguageValue)
		if err != nil {
			_, i, _ := matcher.Match(tags...)
			session.SetInt(ctx, goradd.SessionLanguage, i)
			return i
		}
	}
	return 0
}

// Call SetLanguage to set the user's language to a specific language from the list of supported languages.
func SetLanguage(ctx context.Context, i int) {
	if i >= len(languages) || i < 0 {
		panic("invalid language setting")
	}
	session.SetInt(ctx, goradd.SessionLanguage, i)
}

func CurrentLanguageAttribute(ctx context.Context) string {
	v, _ := session.GetInt(ctx, goradd.SessionLanguage)
	return langAttributes[v]
}

// CurrentLanguage returns the ordinal value of the current language, and the canonical value
// If the language setting is not yet set, it returns the default language
func CurrentLanguage(ctx context.Context) (int, string) {
	v, _ := session.GetInt(ctx, goradd.SessionLanguage)
	return v, langAttributes[v]
}

// CanonicalValue will return the canonical value of the language at the given position
func CanonicalValue(i int) string {
	return langAttributes[i]
}

// Tag returns the language tag corresponding to the given language position
func Tag(i int) language.Tag {
	return languages[i]
}
