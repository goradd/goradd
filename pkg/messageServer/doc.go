/*
The messageServer package implements a general purpose messaging platform based on the gorilla
websocket implementation. The platform is similar to other messaging platforms like pubnub,
in that it is channel based.

The current implementation is limited to communication between goradd forms and objects on
goradd forms, and the server. It uses the pagestate as an authentication token to ensure
that traffic is coming from a valid user of the system.
*/

package messageServer
