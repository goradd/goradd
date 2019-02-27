# Why Goradd

Goradd is a framework for developing websites. There are a large variety of
frameworks available for free on the Internet, written in a variety of different
languages. There are also a number of advocates of not using a framework but
just coding everything from scratch. Whatever approach you take will require
an investment from you to learn that framework before you can be productive.


So should you use Goradd? Here are the main reasons why you should:
* Goradd will save you lots of time getting to a finished product vs. other 
approaches and other frameworks. Go is easy to learn if you don't know it.
* Go is a compiled language, and Goradd is incredibly fast.
* By using a framework, you are using the work of many other people who
have thought through the critical issues of building a website.

Here is what Goradd has built-in:

* Code generation to automatically create an ORM to easily manipulate the database,
and default html forms, panels and controls to edit that data.
* Goradd uses a layered approach so that when you change your database 
structure, you do not have to rewrite user interface code.
* Profiling
* Dates and times: display, localization and timezone management between database, server and client
* Correct html generation
* Correct ARIA tags for accessibility
* Bootstrap support and hooks to support any other css/js framework
* Progressive enhancement. Goradd forms will work even without javascript enabled,
but will also take advantage of javascript to improve the client experience when it is enabled.
* Security, including Clickjacking, and XSS prevention.

Here is what is on our TODO list. Goradd is architected to accomplish these,
but they are not yet complete. Please join the team and help!
* Goradd will work on proxied browsers like Opera Mini, so is suitable for
websites targeted to Africa and Asia.
* Goradd will be scalable, and in particular will work on Hiroku and Google 
App Engine with automatic scalability to mulitple machines enabled.
* Goradd will work on NoSQL databases using the same schema definition as
SQL databases, but will create relationships in NoSQL specific ways.

