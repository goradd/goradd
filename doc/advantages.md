# Why Goradd

You have an idea for a web or mobile app. You want to see if its a good idea. You want to
quickly get something working, try it out, get it in front of others,
and grow your idea incrementally. 
 
Goradd is a framework for developing websites and mobile app backends. 
There are a large variety of
frameworks available for free on the Internet, written in a variety of different
languages. There are also a number of advocates of not using a framework but
just coding everything from scratch. Whatever approach you take will require
an investment in time from you to learn that approach before you can be productive.

So should you use Goradd? Here are the main reasons why you should:
* Goradd will save you lots of time getting to a working product vs. other 
approaches and other frameworks.
* The GO language is easy to learn and is a compiled language, so Goradd is incredibly fast.
* By using a framework, you are using the work of many other people who
have thought through the critical issues of building a website.

Here is what Goradd has built-in:

* Code generation to automatically create an object relational model (ORM), 
which lets you access your database through GO code, rather than using a database language
like SQL. It also creates default html forms, panels and controls to edit that data.
* Goradd uses a layered approach so that when you change your database 
structure, you do not have to rewrite user interface code.
* Profiling to locate performance problems when they crop up
* Dates and times: display, localization and timezone management between database, server and client
* Correct html generation
* Correct ARIA tags for accessibility
* Bootstrap support and hooks to support other css/js frameworks
* Progressive enhancement. Goradd forms will work even without javascript enabled,
but will also take advantage of javascript to improve the client experience when it is enabled.
* Security, including CSRF, Clickjacking, and XSS prevention.

Here is what is on our TODO list. Goradd is architected to accomplish these,
but they are not yet complete. Please join the team and help!
* Goradd will work on proxied browsers like Opera Mini, so is suitable for
websites targeted to Africa and Asia.
* Goradd will be scalable, and in particular will work on Hiroku and Google 
App Engine with automatic scalability to mulitple machines enabled.
* Goradd will work on NoSQL databases using the same schema definition as
SQL databases, but will create relationships in NoSQL specific ways.
* Goradd will create a default GraphQL endpoint for your mobile apps to use, and then
allow you to grow that interface as you mobile app needs change.

