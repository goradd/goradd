

goradd.widget("goradd.testController", {
    _create: function () {
        var self = this;
        goradd.log("Creating test controller");
        window.addEventListener("message", function (event) {
            self._receiveWindowMessage(event)
        }, false);
        this._window = null;
        this._step = 1;
        this._super();
    },
    _receiveWindowMessage: function(event) {
        goradd.log("Message received", event.data);
        if (event.data.pagestate) {
            this._formLoadEvent(event.data.pagestate);
        } else if (event.data.ajaxComplete) {
            this._fireStepEvent(event.data.ajaxComplete);
        }
    },
    logLine: function(line) {
        this.text(this.text() + line  + "\n");
    },
    loadUrl: function(step, url) {
        goradd.log("loadUrl", step, url);
        var self = this;

        if (this._window && this._window.goradd && this._window.goradd._closeWebSocket) {
            this._window.goradd._closeWebSocket(1001);
        }

        this._step = step;

        if (this._window && !this._window.closed) {
            if (this._window.location.pathname == url) {
                this._window.location.reload(true);
            } else {
                this._window.location.assign(url);
            }
        } else {
            this._window = window.open(url);
        }

        if (!this._window) {
            this._fireStepEvent(step, "Opening a popup window was blocked by the browser.");
            return;
        }
        this._window.addEventListener("error", function(event) {
            self._windowErrorEvent(event, step)
        });
    },
    _formLoadEvent: function(pagestate) {
        goradd.setControlValue(this.element.id, "pagestate", pagestate);
        //this._fireStepEvent(this._step);
    },
    _windowErrorEvent: function(event, step) {
        this._fireStepEvent(step,  "Browser load error:" + event.error.message);
    },
    _fireStepEvent(step, err) {
        this.trigger("teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        goradd.log ("changeVal", step, id, val);
        var g = this._getGoraddObj(id);

        if (!g) {
            goradd.log("changeVal: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        g.val(val);
        g.trigger("change");
        this._fireStepEvent(step);
    },
    checkControl: function(step, id, val) {
        goradd.log ("checkControl", step, id, val);
        var g = this._getGoraddObj(id);

        if (!g) {
            goradd.log("checkControl: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        g.element.checked = val;
        g.trigger("change");
        this._fireStepEvent(step);
    },
    checkGroup: function(step, groupName, values) {
        // checks a group of checkbox or radio controls.
        goradd.log ("checkGroup", step, id, values);

        var changeEvent = new Event('change', { 'bubbles': true });

        // uncheck whatever is checked
        var elements = goradd.qa("input[name=" + id +"]:checked");
        goradd.each(elements, function() {
            this.checked = false;
            goradd.g(this).trigger("change");
        });

        // check whatever needs to be checked
        goradd.each(values, function() {
            var val = this;
            var elements = goradd.qa("input[name=" + id + "][value=" + val + "]");
            goradd.each(elements, function() {
                this.checked = true;
                goradd.g(this).trigger("change");
            });
        });

        this._fireStepEvent(step);
    },
    _getGoraddObj: function(id) {
        return this._window.goradd.g(id);
    },
    closeWindow: function(step) {
        this._window.close();
        this._fireStepEvent(step);
    },
    click: function (step, id) {
        var self = this;
        goradd.log("click", step, id);
        var g = this._getGoraddObj(id);
        if (!g) {
            goradd.log("click: element not found", id);
            self._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        g.click({postFunc: function() {
            self._fireStepEvent(step);
        }});
    },
    callWidgetFunction: function (step, id, f, params) {
        goradd.log("WidgetF", step, id, f, params);

        var g = this._getGoraddObj(id);
        if (!g) {
            goradd.log("callWidgetFunction: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        var ret = g.f(f, params);

        goradd.setControlValue(this.element.id, "jsvalue", ret);
        this._fireStepEvent(step);
    },
    typeChars: function (step, id, chars) {
        var g = this._getGoraddObj(id);
        if (!g) {
            goradd.log("typeChars: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        g.val(chars);
        this._fireStepEvent(step);
    },
    focus: function (step, id) {
        var g = this._getGoraddObj(id);
        if (!g) {
            goradd.log("focus: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        g.focus();
        this._fireStepEvent(step);
    }

});