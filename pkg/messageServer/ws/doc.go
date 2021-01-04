/*
The ws package implements a message server based on a Websocket implementation local to the running server. It is appropriate
for single-server application implementations and for development.

Without too much difficulty, it could be modified to run as a separate service that you communicate with using a
rest API or something similar, so that it could accommodate an FCGI server implementation or
a microservice implementation.

It starts two listeners on two ports, for encrypted and unencrypted communication. Unencrypted communication is useful
during development of an application, but stick to encrypted only for deployment.

*/

package ws

