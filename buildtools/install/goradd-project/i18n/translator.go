// This i18n package provides shell routines and interfaces to implement a translation scheme for your application that is not opinionated.
// Translating an application is tricky and laborious enough without forcing you to do your translations a particular way.
// The one "opinionated" thing it does do is use a Context to discover the locale. GoRadd is a web framework, and so the
// assumption is that you will want the locale to be specific to each user, not application wide.
//
// This package is part of the local application, giving you the opportunity to modify your translation process as
// you see fit.
package i18n

import "context"

// A Translater interface represents an object that is translatable using the current context
type Translater interface {
	Translate(ctx context.Context) string
}

func T(ctx context.Context, key string) string {
	// TODO: extract the locale from the context and hand it off to the translator
	return key
}
