# GoRADD
A rapid Web application development framework for Go.

GoRADD is a monolithic web development framework for rapidly creating a web application from a concept
in your mind and then allowing you to as easily as possible maintain that application through
all the twists and turns of the change process. It is ideal for prototyping, intranet websites,
websites that are very data intensive with many forms to gather information from users,
websites that require the security and speed of a compiled language, websites with thousands of
simultaneous users, and websites being maintained by one or a small group of developers. It is
particularly good for developers new to GO and/or new to web development.

## Installation
See the [Quick Start](doc/quickstart.md) guide to get started.

## Requirements
- A supported database up and running on your local development computer. 
Current supported databases are:
    - Mysql

### For Developing GoRADD itself
- Sass (to build the css files from the scss source)

## Goals
1) 80-20 rule, where out of the box GoRADD will do most of the hard work of building a website, and
will quickly get you a working website, but not necessarily one that you will want to ship.
GoRADD is architected to allow you to make changes and plug in other open-source 
software as you need.
1) Incremental changes. Most development processes go through a lengthy requirement analysis process,
followed by a design process, and a lengthy build process, only to find out that what you built wasn't 
really what was needed. Instead, GoRADD gets you a working website quickly, and then lets you build out
your application incrementally. It
tries to make it easy to restructure your website, even
at the data structure level, and have your changes filter through your application quickly without
requiring a complete rewrite. 
1) Layered development. GoRADD has its code, you have your code, and then there is an in-between
interface that changes over time. GoRADD uses code generation to create this interface, and clearly
delineates the code that you can change to modify the interface, vs. code that it will generate as you
change your data model. The result is a product that is easy to change as your world and
requirements change.
1) Most development happens in GO. What the user does in the browser is mirrored on the
server, which allows you to work in a way that feels like you are building a desktop application. This
makes your developers more productive and it allows you to build your app using common GO tools like
the built-in unit test environment and documentation server. You can still work in javascript if you
want to or need to do custom UI work, but often you don't have to.
1) Stability. We want to build applications that real people use, and that means reliance on tried
and true technologies that work on a broad range of browsers and servers, rather than technologies
that require lots of Polyfills for emerging standards. JQuery is currently required, and partly because we want
to make it easy to create Bootstrap based applications. However, Bootstrap has announced
that they are removing reliance on JQuery, and we will attempt to do so as well.
1) Progressive enhancement. If you use the provided widgets, you can create a website
that works even if the client turns off Javascript. All major browsers are currently supported,
but we hope to support Opera Mini as well.
1) Rich libraries of widgets. GoRADD provides standard widgets corresponding to 
basic html controls, and also provides Bootstrap widgets. If you have a particular
css or javascript widget library you want to support, building the GoRADD
interface is fairly easy to do, and the Bootstrap library included gives you a 
model to follow.
1) Scalability. GoRADD is architected for scalability. All user state information is serializable
to key-value stores. You might need to build the interface to the particular key-value store you
are interested in, but that is not difficult. Some specific issues to consider:
    1. GoRADD requires a MySQL database at this point for your main data store. 
        SQL is great for creating most common data
           structures, is great when you need to change your structure without destroying data, and
           is fast enough for most applications. However, all data access is done through a common API,
           so switching an application that is already written to another SQL database like Postgresql, Oracle, or any
            other database is very straight-forward
           and is just a matter of implementing the database layer. In fact, the database layer is generic
           enough that you could switch to a NoSQL implementation
           as your product matures and you need scalability at speed.
    2. GoRADD maintains the state of each user of the website in something we call the *pagestate*.
       The pagestate is serializable to any key-value store. Currently, only an in-memory store is
       provided, but writing an interface to any common key-value store is easy.
    3. Live updates work through a pub/sub mechanism. Goradd provides a single-server in-memory 
       system out of the box, but its easy to switch to any other pub/sub mechanism, including 
       distributed systems like pubnub, ally, google cloud messaging, etc. There are no payloads
       with the messages and traffic is minimal.

### Future Goals
* WebComponents. WebComponent architecture fits particularly well with goradd's architecture. However,
WebComponents are not fully supported by all major browsers. As WebComponents gain traction, we hope
to use them for future browser widgets. In the mean-time, we support many JQuery based widgets.
* Matching GraphQL or GraphQL like interface. The ORM architecture has many similarities to
GraphQL, and could potentially auto-generate a GraphQL interface to make it easy to integrate
a mobile app interface. 

### Anti-patterns
1) GoRADD's html server is not microservice based. 
While you can create microservices that serve parts of your application, at its
core goradd is a monolithic framework that includes an ORM, an MVC architecture, and a basic control
and form management library. 
2) Object-oriented. Some of goradd uses a code pattern that mirrors traditional object-oriented
inheritance and gets around some of GO's limitations in this area, including implementing 
virtual functions. We have found this particularly useful in the control library.
If you hate inheritance, goradd is not for you. If you don't mind it, but you still
like object composition too, this is your place.
3) Code generation. GoRADD relies heavily on code generation, and in particular uses the
related github.com/goradd/got template engine to generate code.

## Acknowledgements
GoRADD is a port of the PHP framework [QCubed](https://github.com/qcubed/qcubed). QCubed itself was a 
fork of the PHP framework written by Mike Ho called [QCodo](https://github.com/qcodo/qcodo).
Mike is the original mastermind of many of the concepts in GoRADD, like:
- A code-generated ORM
- The use of "nodes" to describe database entities AND the relationships between them.
- Code-generated CRUD forms to get you started.
- Scaffolding that separates code-generated code from developer code so that code-generation
can continue throughout the life of the project.
- A lightweight javascript layer for processing events and actions through ajax.
- The formstate engine to mirror the state of html and javascript widgets on the server-side
so that the server-side engineer has complete control over what is happening in the html without
needing to write javascript.

GoRADD relies on a number of other open-source projects, including:
- The [Shurcool Github Markdown Library](https://github.com/shurcooL/github_flavored_markdown)
- Alex Edward's [SCS Session Manager Library](https://github.com/alexedwards/scs)
- Akeda Bagus' [Inflector Library](https://github.com/gedex/inflector)
- Gorilla's [Websocket Library](https://github.com/gorilla/websocket)
- Kenneth Shaw's [Snaker Library](https://github.com/knq/snaker)
- Stretchr's [Testify Testing Library](https://github.com/stretchr/testify)

GoRADD was created and is maintained by [Shannon Pekary](https://github.com/spekary)

### Thanks To
[JetBrains](https://www.jetbrains.com/go) for use of the GoLand Go Editor

![BrowserStack](https://d3but80xmlhqzj.cloudfront.net/production/images/static/header/bstack-logo.svg) 
[BrowserStack](http://browserstack.com) for automated browser testing tools
