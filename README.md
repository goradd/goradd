# goradd
A rapid Web application and REST service development framework for Go.

Goradd is a monolithic web development framework for rapidly creating a web application from a concept
in your mind and then allowing you to as easily as possible maintain that application through
all the twists and turns of the change process. It is ideal for prototyping, intranet websites,
websites that are very data intensive with many forms to gather information from users,
websites that require the security and speed of a compiled language, websites with thousands of
simultaneous users, and websites being maintain by one or a small group of developers. It is
particularly good for developers new to GO and/or new to web development.

## Goals
1) 80-20 rule, where out of the box Goradd will do most of the hard work of building a website, and
will quickly get you a working website, but not necessarily one that you will want to ship.
Goradd is architected to allow you to make changes and plug in other open-source 
software as you need.
1) Fail fast. Most development processes go through a lengthy requirement analysis process,
followed by a design process, and a lengthy build process, only to find out that what you built wasn't 
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
1) Most development happens in GO. What the user does in the browser is mirrored on the
server, which allows you to work in a way that feels like you are building a desktop application. This
makes your developers more productive and it allows you to build your app using common GO tools like
the built-in unit test environment and documentation server. You can still work in javascript if you
want to or need to do custom UI work, but often you don't have to.
1) Stability. We want to build applications that real people use, and that means reliance on tried
and true technologies that work on a broad range of browsers and servers, rather than technologies
that require lots of Polyfills for emerging standards. JQuery is currently heavily relied on, and partly because we want
to make it easy to create Bootstrap based applications (Bootstrap is not required though)
1) Progressivive enhancement. Using the supported widgets, you can create a website
that works even if the client turns off Javascript.
1) Rich libraries of widgets. Goradd provides standard widgets corresponding to 
basic html controls, and also provides Bootstrap widgets. If you have a particular
css or javascript widget library you want to support, building the Goradd
interface is fairly easy to do, and the Bootstrap library included gives you a 
model to follow.

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
* Matching REST service. Goradd will auto-generate a flexible REST service
to use with mobile app development. 

### Anti-patterns
1) Goradd's html server is not microservice based. 
While you can create microservices that serve parts of your application, at its
core goradd is a monolithic framework that includes an ORM, an MVC architecture, and a basic control
and form management library. If you are trying to build an application for millions of users, goradd is not
for you at this time. However, it is architected to eventually allow parts of it to be handled off-line by
other servers, so it is (or will be soon) scalable. Also, it will feature the
ability to quickly build an http REST service to integrate with mobile apps.
2) Object-oriented. Some of goradd uses a code pattern that mirrors traditional object-oriented
inheritance and gets around some of GO's limitations in this area, including implementing 
virtual functions. If you hate inheritance, goradd is not for you. If you don't mind it, but you still
like object composition too, this is your place.
3) Code generation. Goradd relies heavily on code generation, and in particular uses the
related github.com/goradd/got template engine to generate code.

## Installation
See the [Quick Start](doc/quickstart.md) guide to get started.

## Requirements
- A supported database up and running on your local development computer. 
Current supported databases are:
    - Mysql

### For Developing Goradd itself
- Sass (to build the css files from the scss source)