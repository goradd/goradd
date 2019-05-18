"use strict";

goradd.ajaxq = {
    /**
     * Ajax Queue
     *
     * This used to be handled with a jquery plugin, but since we are trying to get away from jquery, and working
     * towards an OperaMini compatible version, we are rolling our own.
     */
    _q: [],
    _currentRequests: {},
    _idCounter: 0,
    /**
     * Queues an ajax request.
     * A new Ajax request won't be started until the previous queued
     * request has finished.
     * @param {function} f function that returns ajax options.
     * @param {boolean} blnAsync true to launch right away.
     */
    enqueue: function(f, blnAsync) {
        if (!blnAsync) {
            var wasRunning = this.isRunning();
            this._q.push(f);
            if (!wasRunning) {
                this._dequeue();
            }
        } else {
            this._do1(f);
        }
    },
    /**
     * Returns true if there is something in the ajax queue. This would happen if we have just queued an item,
     * or if we are waiting for an item to return a result.
     *
     * @returns {boolean} true if the goradd ajax queue has an item in it.
     */
    isRunning: function() {
        return this._currentRequests.length === 0;
    },
    _dequeue: function() {
        var f = this._q.shift();
        if (f) {
            this._do1(f);
        }
    },
    _do1(f) {
        var self = this;
        var opts = f();
        this._idCounter++;
        var ajaxID = this._idCounter;

        var objRequest = new XMLHttpRequest();

        objRequest.open("POST", opts.url, true);
        objRequest.setRequestHeader("Method", "POST " + opts.url + " HTTP/1.1");
        objRequest.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        objRequest.setRequestHeader("X-Requested-With", "xmlhttprequest");

        objRequest.onreadystatechange = function() {
            if (objRequest.readyState === 4) {
                if (objRequest.status === 200) {
                    try {
                        opts.success(JSON.parse(objRequest.response));
                    } catch(err) {
                        // Goradd returns ajax errors as text
                        opts.error(objRequest.response, err);
                    }
                } else {
                    // This would be a problem with the server or client
                    opts.error("An ajax error occurred: " + objRequest.statusText);
                }

                delete self._currentRequests[ajaxID];
                if (self._q.length === 0 && !self.isRunning()) {
                    goradd.g(goradd.form()).trigger("ajaxQueueComplete");
                }
                self._dequeue(); // do the next ajax event in the queue
            }
        };
        self._currentRequests[ajaxID] = objRequest;
        var encoded = self._encodeData(opts.data);
        objRequest.send(encoded);
    },
    _encodeData(data) {
        var a = [];
        var key;
        for (key in data) {
            var value = data[key];
            var s = encodeURIComponent(key) + "=" +
            encodeURIComponent( value == null ? "" : value );
            a.push(s);
        }
        return a.join("&");
    }



};