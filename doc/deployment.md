# Deployment

## Build Tags
Use build tags to build various flavors of your application. Goradd makes use of the following tags internally:
* **release** - Sets the `config.Release` variable to `true`. This is used throughout the framework to turn off code 
that should not be included in a release version of a product, or that might even be dangerous to include, like unit
tests, code generators, etc. It also directs the server to look for assets in a specified asset directory, rather
than your development system. Feel free to use it to change directory paths, database credentials, or whatever
might be different on your deployment server vs. your development system.
* **nodebug** - Sets the `config.Debug` variable to false. This is used throughout the framework to turn on debug
specific features, like logging, profiling, etc. This allows you to create a release version that might be used
by your testers and deployed on a mirror image of your deployment server, but still has particular debug features so 
your testers can recreate issues and deliver usable information to your developers.

## Building a Release Version
See the goradd-project/build directory for sample scripts to build the final version of
your app. There are examples there to build a standalone application and a dockerized container of
your app.

