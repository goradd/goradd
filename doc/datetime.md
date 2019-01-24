# Dates and Times

In a multi-user, multi-language, client-server product, something so small can be so complicated. 

Consider the following scenarios:
1. You ask a user for a date and time for when an event should happen in the future. What does that
actually mean? If the date is after a change in daylight savings time, should the time move an hour
when the event happens, or should we assume the user knows the time is going to change between now and then?
What if the event is going to happen in a different timezone than the current location of the user?
What if the user changes timezones later?
2. You record a timestamp of an event that happened in the past. But what if someone looks at the event
from a different timezone, what should that user see? What if the event is during a timezone change, 
what is the actual timestamp?
3. You have your database record a timestamp. When your server requests the value, what should be
the timezone of the new value, the server's timezone, or the database's timezone? When you display the
timestamp to the client, how do you convert that value to the client's timezone? How do you localize
what is displayed into the location and language of the client?

Goradd tries to resolve some of these issues in logical ways so that you don't have to figure all of that
out on your own. But first some background...

## Types of date-times

There are a few types of date-times that goradd, and you, need to consider:
1. Timestamps. This is a set moment in historical time. If you view this moment from different timezones,
you should see the hours differently, since 8:00 am in California is also 11:00 am in New York. If the
date and time point to a time that locally is in the middle of a change in daylight savings time, you
might see it as an ambiguous time, but internally, there is no ambiguity. Often, the internal representation
of this is a number of seconds since a particular moment in known time, GMT.
2. DateTimes. A date-time has no timezone information. It means an event at a particular date and time
in whatever timezone you are in. If you change timezones, or daylight savings time changes, the 
displayed date and time do not change. This is what might be in a scheduling application, since if
a user creates an event in the future after a change in daylight savings time, you are assuming that the
user knows there will be a change. For example, if the user sets an alarm for 8:00 in the morning, 
after daylight savings time the alarm should still be at 8:00 in the morning.
3. Date only. Since there is no time information, this could be a DateTime with an empty time.
4. Time only. Since there is no date information, this could be a DateTime with a zero date.

Note that sometimes using a Timezone or DateTime might still need timezone information. 
Consider the following:
1. If you schedule a meeting for a particular date and time, and that meeting is after a 
daylight savings time change, but that meeting might be a video-conference meeting that could be 
attended by people in different timezones, you want to know the relative timezone of everyone, 
but you assume that everyone knows that daylight savings time will change. In this case, you need
to also know the timezone location of the event, and when you display it, you will have to be
very careful to know the location of the viewer compared to the location of the event.
2. If you are showing a timestamped event that happened in the past, you might want to show it
in the local time of when the event occurred, no matter where it is viewed from.

## What is a timezone
There really are two representations known as a timezone.
1. A named timezone, like "America/New_York". This tells us the location in the world where a time
will be displayed. However, the same timestamp might be displayed differently in that location
depending on whether daylight savings time IS in effect when viewing that time, and depending on
whether daylight savings time WAS in effect at the time of the event. For example, if EDT was
in effect at the time of the event, and you are viewing the event from EST, you might want to
see the event in EST time or in EDT time. It depends on the application.
2. An offset from UTC. Abbreviations like EDT, EST, PST, PDT, or +0800, etc. tell us the actual
offset from a daylight savings time agnostic UTC time. As in the above example, when you view
this time from a particular location at a particular date and time, what you want to see will
depend on the application.

## Considerations
### The movement of date-times
The above date-time values will move between different locations. Consider that the client, the server,
and the database may all be in different timezones. Also, while each physically might be in a 
particular timezone, each might store the dates and times in a different timezone, or in
UTC. Sometimes you have control over that, and sometimes not.

### Database capabilities
Some databases can store timezone information, and some can't. Sometimes, it doesn't matter.  

MySQL DateTime objects do not save timezone information, but Postgres gives you the option. 

Both Postgres and MySQL Timestamps are stored internally as
a number of seconds from a known UTC time, but are converted to a timezoned format when the value
is transferred to a server. The GO mysql drivers give you the ability to set what timezone that is,
and it can be different than the server's timezone.

### Browser capabilities
Date objects in browsers internally are represented as UTC Timestamps. 

When creating a Date object, the only way to specify a timezone that is not the current
local timezone is using an ISO-8601 string representation of the date and asking the
constructor to parse it. However, there are lots of warnings in the various descriptions
of the javascript Date constructor to NOT do this, since there are inconsistencies in how browsers
interpret these strings. However, I think the only inconsistency is when you DO NOT specify
a timezone, some browsers assume local time, and others assume UTC. So, if you always specify
a timezone, you should be OK. Other than that, the only way to create a consistent 
date-time is to use the Date.UTC() constructor and specify all values in UTC.

Most modern browsers support the ability to get the local timezone name through the
following call:
```javascript
Intl.DateTimeFormat().resolvedOptions().timeZone
```
Internet Explorer 11 is the exception here. 

Also, most modern browsers can convert UTC date-times to a localized display of the date-time
in the current timezone. They can also display dates and times in UTC time, but they generally
*do not* have the capability of converting a date-time into an arbitrary timezone. GO *does*
have this capability.

## Questions
1. When drawing a date to the browser, should we always just allow the date to be converted
to local time by the browser? Issues include:
- To draw a short version of a date, not just in the current timezone but in the current
language requires support for the Intl.DateTimeFormat object, which is supported
in most major browsers, but not some mobile browsers. See https://caniuse.com/#feat=internationalization.
- Its not possible to convert a date to an arbitrary timezone for display. Browsers only support
conversion back an forth between UTC and the current local timezone.
- Go currently has no international support for time.Time types. If rendering a date or time
from the server, its up to you to get the order of the date-month correct in a short version,
to know whether the current locale uses 12 or 24 hour time, and to translate days of the week.
2. How do we differentiate between a date-time that is agnostic of timezone, and one that cares.
3. How do you do a database search when the date-time that the client is requesting as
the bounding parameters for the search are possibly in a different timezone than the database.

## Solutions
1. Like JavaScript, we will attempt to always store Timestamps and DateTimes internally
in UTC time. 
2. When specifying search parameters on a timestamp, you must convert the client specified
date-time to UTC before searching. You will likely need to do that in javascript.
3. If you are providing the client a way to just search by date, then you should probably
store date-times separately as a Date and Time in the database.
4. If you are making an app like a scheduling app, where the timezone of the created event
matters, you should probably store in the database separate date, time and timezone information
so you can recreate the details needed.
5. If you are making an app that you believe will need to be internationalized, always let
javascript do the drawing. Otherwise you can let the server draw, but you will need to find
out the timezone of the client if the app might be used by people in more than one timezone.
6. If your timestamps need to be drawn in the timezone it was created, you should store the
timezone separately with the timestamp. When drawing, let Javascript do the drawing, but
then append the timezone information you saved.