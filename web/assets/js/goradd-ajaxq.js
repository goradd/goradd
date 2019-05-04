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
        $.ajax()
    },
    /**
     * Returns true if there is something in the ajax queue. This would happen if we have just queued an item,
     * or if we are waiting for an item to return a result.
     *
     * @returns {boolean} true if the goradd ajax queue has an item in it.
     */
    isRunning: function() {
        return goradd._currentRequests.length == 0;
    },
    _dequeue: function() {
        var f = this._q.shift();
        if (f) {
            this._do1(f);
        }
    },
    _do1(f) {
        var opts = f();
        this._idCounter++;
        var ajaxID = this._idCounter;

        var objRequest;
        if (window.XMLHttpRequest) {
            objRequest = new XMLHttpRequest();
        } else if (typeof ActiveXObject != "undefined") {
            objRequest = new ActiveXObject("Microsoft.XMLHTTP");
        }

        if (objRequest) {
            this._currentRequests[ajaxID] = objRequest;
            objRequest.open("POST", opts.url, true);
            objRequest.setRequestHeader("Method", "POST " + opts.url + " HTTP/1.1");
            objRequest.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

            objRequest.onreadystatechange = function() {
                if (objRequest.readystate === 4) {
                    if (objRequest.status === 200) {
                        // success

                    } else {

                    }

                    delete goradd._currentRequests[ajaxID];
                }
            };
            qcodo.ajaxRequest = objRequest;
            objRequest.send(opts.data);
        }
    }



}