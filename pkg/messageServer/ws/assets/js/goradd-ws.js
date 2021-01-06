/*

Goradd Websocket Client

This file attaches a websocket client to the current goradd form. The client replies to one message,
which will cause the form to do an ajax update.

The point of all this is to allow the server to simulate a push update of the form. Essentially,
the server is just telling the client that its time to pull updates. This simplifies the architecture,
and makes it easier to implement this using a messaging server instead as the application expands.

To implement this differently, you would need to substitute this js with a system that subscribes
to a channel using the pagestate as the channel id. On the server end, you would send to the channel
the update message. Just replace the initMessageClient function with your own.

You can piggyback on this and add your own websocket messages by simply adding event listeners
to the goradd._ws object.

*/

goradd._channels = {};

goradd.initMessagingClient = function(loc) {
    if (window.WebSocket) {
        var d = new Date();
        var con;
        var port;

        if (location.protocol === 'https:') {
            con = "wss://";
        } else {
            con = "ws://";
        }
        port = location.port ? ":" + location.port : "";
        con += window.location.hostname + port + loc + "?id=" + goradd.getPageState();
        goradd._ws = new WebSocket(con);
        goradd._ws.addEventListener("message", goradd._handleWsMessage);
        goradd._ws.addEventListener("close", goradd._handleWsClose);
        // we purposefully do not use goradd._ws.onmessage = ... so that we can add multiple event listeners.
        goradd._ws.onopen = function() {
            goradd.subscribeWatchers();
            g$(goradd.form()).trigger("messengerReady");
        }
    }
};

// channels is an array of strings indicating the channels to subscribe to
// f is the function to call when a message comes through on that channel
goradd.subscribe = function(channels, f) {
    var msg = {};
    msg["subscribe"] = channels;
    goradd._ws.send(JSON.stringify(msg));
    goradd.each(channels, function() {
        goradd._channels[this] = f;
    });
};

/*
The default message handler. Will route the message to the appropriate channel.
 */
goradd._handleWsMessage = function(e) {
    var messages = JSON.parse(e.data);
    console.log("messages: " + e.data);

    goradd.each(messages, function() {
        var msg = this;
        var f = goradd._channels[msg.channel];
        if (!!f) {
            f(msg);
        }
    });
};

/*
Close the websocket connection.
See https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/close for status codes
*/
goradd._closeWebSocket = function(status) {
    if (goradd._ws) {
        goradd._ws.close(); // Not all browsers support a status code, and attempting one breaks javascript.
    }
};

goradd._handleWsClose = function(e) {
    goradd._channels = {};
    goradd._ws = null;
};

