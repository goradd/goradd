/*

Goradd Websocket Client

This file attaches a websocket client to the current goradd form. The client replies to one message,
which will cause the form to do an ajax update.

The point of all this is to allow the server to simulate a push update of the form. Essentially,
the server is just telling the client that its time to pull updates. This simplifies the architecture,
and makes it easier to implement this using a messaging server instead as the application expands.

To implement this differently, you would need to substitute this js with a system that subscribes
to a channel using the formstate as the channel id. On the server end, you would send to the channel
the update message. Just replace the initMessageClient function with your own.

You can piggyback on this and add your own websocket messages by simply adding event listeners
to the goradd._ws object.

 */

goradd.initMessagingClient = function(wsPort, wssPort) {
    if (window.WebSocket) {
        var d = new Date();
        var con;
        var port;

        if (location.protocol === 'https:') {
            con = "wss://";
            port = wssPort;
        } else {
            con = "ws://";
            port = wsPort;
        }
        con += window.location.hostname + ":" + port + "/ws?id=" + goradd.getFormState() + "&ch=form-" + goradd.getFormState();

        goradd._ws = new WebSocket(con);
        goradd._ws.addEventListener("message", goradd._handleMessage);
        // we purposefully do not use goradd._ws.onmessage = ... so that we can add multiple event listeners.
    }
};

/*
The default message handler. It treats the message data as a JSON object, looks for a 'grup' item
there, and then updates the form if found. You can therefore send any other messages you want
to your own handlers. If you message is a JSON object and has a grup item, your message will also
update the form. Otherwise, it will be ignored here.
 */
goradd._handleMessage = function(e) {
    var message = JSON.parse(e.data);
    console.log("message");
    if (message.grup) {
        console.log("update " + goradd.getFormState());
        goradd.updateForm();
    }
};
