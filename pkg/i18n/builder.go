package i18n


type translationBuilder struct {
	domain string
	language string
	id string	// same as msgctxt in .PO files. Disambiguates same text. Usually blank.
	message string
	arguments []interface{}
}

func Build() *translationBuilder {
	return new(translationBuilder)
}

func (b *translationBuilder) Domain(domain string) *translationBuilder {
	b.domain = domain
	return b
}

// Lang adds the canonical value to the builder
func (b *translationBuilder) Lang(lang string) *translationBuilder {
	b.language = lang
	return b
}

// ID adds a context to disambiguate strings with the same message id but different meanings
func (b *translationBuilder) ID(id string) *translationBuilder {
	b.id = id
	return b
}

// Comment will add a comment to the extracted translation file, but will otherwise not change the builder
// Use this to add comments directed to the person doing the translation.
func (b *translationBuilder) Comment(comment string) *translationBuilder {
	return b
}

// T ends the builder and performs the translation
func (b *translationBuilder) T(s string) string {
	if b.message == "" {
		return ""
	}
	if b.domain == "" {
		b.domain = ProjectDomain
	}
	if b.language == "" {
		b.language = langAttributes[0]
	}
	b.message = s

	return translators[b.domain].Translate(b)
}

// Sprintf ends the builder and performs the translation
func (b *translationBuilder) Sprintf(s string, params... interface{}) string {
	if b.message == "" {
		return ""
	}
	if b.domain == "" {
		b.domain = ProjectDomain
	}
	if b.language == "" {
		b.language = langAttributes[0]
	}
	b.message = s
	b.arguments = params

	return translators[b.domain].Translate(b)
}



// The following are modifiers to the T() function in page.Control

type id struct {
	id string
}

// ID is a parameter you can add to the page.control.T() function to specify a message id. Usually the message id is the
// same as the string being translated, but when multiple strings are translated that are the same but have different meaning,
// this will be required. This is used as the msgctxt value in PO files, and is combined with the message to make a composite id
// in golang translation files. Adding a comment is helpful in these situations.
func ID(i string) interface{} {
	return id{i}
}

type comment struct {
	comment string
}

// Comment adds a comment to the translation. It is used in extracted files, but does not impact the translator.
func Comment(c string) interface{} {
	return comment{c}
}

func ExtractBuilderFromArguments(args... interface{}) (b *translationBuilder, args2 []interface{}) {
	b = Build()
	for _,a := range args {
		if i,ok := a.(id); ok {
			b.ID(i.id)
		} else if _,ok := a.(comment); ok {
			// do nothing
		} else {
			args2 = append(args2, a)
		}
	}
	return
}
