/*
The i18n package provides support services for translation of goradd web pages and widgets.

Overview
Internationalization is a complex beast. Beyond all the previous attempts at doing it, GO is actively developing
the golang.org/x/text library to provide additional translation and internationalization support.

This new GO library is unfortunately a bit opinionated. It provides an extractor that searches for strings in your application
and dumps them to a json file ready for translation. However, it will only look for strings that are being processed
by the message.Sprintf and similar functions. It does some code analysis to extract the strings in ways that provide
additional information to translators and the reassembly process, which is great, but a bit limited.

This 18n package works with the page.control.T() and similar functions to try to provide some GO like functionality in
a less opinionated way.

Translators
At the top level are the translators, which are contained in a global map of TranslatorI interfaces that implement the Translate function.
These translators are keyed by a domain name, not an http domain name, but a domain borrowed from the historical gettext
process. It simply lets you specify different translators for different strings, and it is primarily useful by allowing
libraries and the goradd core to provide their own translations and translators of their own string. Also, you can replace those and
provide your own translations of these same strings if you want. Library makers, and the goradd framework,
simply need to register their translators during their own package initialization, and you can replace those
during the local application package initialization process.

As far as implementing translation for your own application, you can either do it the GO way using the message.Printer
functionality provided by golang.org/x/text/message (see https://godoc.org/golang.org/x/text for details and a reference
to a very helpful youtube video on the subject), or you can use goradd's translator service.

Goradd Translation
To send a string to your application translator, simply call T("message") on any control or form, and it will return the translation.
You can add some annotations by adding an i18n.ID() call or i18n.Comment() call to the call, like so:

  newMessage := ctrl.T("my message", i18n.ID("Use this ID for additional context"), i18n.Comment("This can become an extracted comment))

The code generated forms and controls automatically call this function to translate strings.

The extractor for this is not yet built, but it should not be too hard, as a great example is already in the go/text library.

Since translation is provided by an interface, you can handle translation however you want by simply creating an object
that implements the TranslateI interface, and then passing it to RegisterTranslator with the GoraddProject domain. There are
a huge variety of libraries available for managing translations with .po files, with online utilities like Google's own
Translation Toolkit, with databases, static linking, etc. Its up to you.

You can see examples of how the framework itself does it in the source.
 */
package i18n
