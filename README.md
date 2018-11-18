# goradd
A rapid Web application development framework for Go, based on PHP QCubed.

Goradd is a monolithic web development framework for rapidly creating a web application from a concept
in your mind, and then allowing you to as easily as possible maintain that application through
all the twists and turns of the change process. It is ideal for prototyping, intranet websites,
websites that are very data intensive with many forms to gather information from users,
websites that require the security and speed of a compiled language, websites with thousands of
simultaneous users, and websites being maintain by one or a small group of developers. It is
particularly good for developers new to GO and/or new to web development.

## Goals
Goradd trys to achieve the following goals:
1) 80-20 rule, where out of the box it will do most of the hard work of building a website, and
will quickly get you a working website, but not necessarily one that you will want to ship. Speed is
not always the goal in goradd, but rather the ability for you to incorporate changes that fix bottlenecks
you find along the way. 
1) Fail fast. Most development processes go through a lengthy requirement analysis process,
follow by a design process, and a lengthy build process, only to find out that what you built wasn't 
really what was needed. Instead, Goradd gets you a working website quickly, and then lets you build out
your application incrementally. It
tries to make it easy to restructure your website, even
at the data structure level, and have your changes filter through your application quickly without
requiring a complete rewrite. It achieves this by:
1) Layered development. Goradd has its code, you have your code, and then there is an in-between
interface that changes over time. Goradd uses code generation to create this interface, and clearly
delineates the code that you can change to modify the interface, vs. code that it will generate as you
change your data model. The result is a product that is easy to change as your world and
requirements change.
1) Most development happens in GO. What is happening at the user in the client is mirrored on the
server, which allows you to work in a way that feels like you are building desktop application. This
makes your developers more productive and it allows you to build your app using common go tools like
the built-in unit test environment and documentation server. You can still work in javascript if you
want to or need to do custom UI work, but often you don't have to.
1) Stability. We want to build applications that real people use, and that means reliance on tried
and true technologies that work on a broad range of browsers and servers, rather than technologies
that require lots of Polyfills for emerging standards. JQuery is currently heavily relied on, and partly because we want
to make it easy to create Bootstrap based applications (Bootstrap is not required though)


### Future Goals
* Scalability. Goradd is architected for scalability, but the specific features required for
scalability have not been built and are not scheduled for version 1. Some of those things are
    1. Goradd requires a SQL database at this point. SQL is great for creating most common data
    structures, is great when you need to change your structure without destroying data, and
    is fast enough for most applications. Goradd is architected so that you could switch to a NoSQL implementation
    as your product matures and you need scalability at speed, but the NoSQL drivers are not currently built. 
    They would not be difficult to do, but they will need to rely on a separate data structure definition. The
    plan is that a SQL database would be able to generate a schema that would be used by the NoSQL drivers to
    continue to maintain the data.
    2. Goradd maintains the state of each user of the website in internal memory we call the *pagestate*.
    Since its in memory, each user is currently bound to one server. Go is incredibly fast, so one server should be
    able to manage thousands of users with a reasonable amount of RAM, but to grow beyond this, some work would
    need to be done on serializing the pagestate into an off-site database. This effort is in process.
    3. Live-updates. Live updates in a multi-user environment can be particularly difficult at the data model
    level. However, browser technologies also make them difficult at the client too. The browser world is in rapid flux
    around this topic, with different browsers supporting a variety of technologies 
    (WebSockets, Server-side events, and now fetch streaming), with Microsoft browsers generally 
    lagging the pack. Hopefully things will settle down and we can implement this in a sane way.
* WebComponents. WebComponent architecture fits particularly well with goradd's architecture. However,
WebComponents are not fully supported by all the browsers. As WebComponents gain traction, we hope
to use them for future browser widgets. In the mean-time, we support many JQuery based widgets.

### Anti-patterns
Go purists will be particularly unhappy with the following:
1) No microservices. While you can create microservices that serve parts of your application, at its
core goradd is a monolithic framework that includes an ORM, an MVC architecture, and a basic control
and form management library. If you are trying to build an application for millions of users, goradd is not
for you at this time. However, it is architected to eventually allow parts of it to be handled off-line by
other servers, so it is (or will be soon) scalable.
2) Object-oriented. Some of goradd uses a code pattern that mirrors traditional object-oriented
inheritance and gets around some of GO's limitations in this area, including implementing 
virtual functions. If you hate inheritance, goradd is not for you. If you don't mind it, but you still
like object composition too, this is your place.
3) Code generation. Goradd relies heavily on code generation, and in particular uses the
related github.com/spekary/got template engine to generate code.

## Installation
### For Go 1.10 and below:
1. Create a new directory and set your GOPATH environment variable to it, if needed.
1. Make sure the GOPATH/bin directory is in your execution path, or execute commands from there.
1. Execute ```go get github.com/spekary/goradd```
1. Execute ```goradd install```

### For Go 1.11 and above using modules:
1. If you just installed go, make sure your GOPATH/bin directory is in your execution path.
1. Create a new directory *outside* of your GOPATH and cd to that new directory.
1. Execute ```go get github.com/spekary/goradd```
1. Execute ```goradd install```

##Requirements
### For Development
- Sass (to build the css files from the scss source)