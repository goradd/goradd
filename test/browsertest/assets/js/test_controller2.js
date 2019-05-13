goradd.widget = function() {

};

goradd.testController = function(el) {
    el = goradd.el(el);

    return new goradd.TestController(el);
};

goradd.TestController = function(el) {
//        this._super();
    this.element = el;
    var self = this;
    goradd.log("Creating test controller");
    window.addEventListener("message", function(event) {
        self._receiveWindowMessage(event)
    } , false);
};

goradd.TestController.prototype = {
    _window:null,
    _step:1,
    _receiveWindowMessage: function(event) {
        goradd.log("Message received", event.data);
        if (event.data.pagestate) {
            this._formLoadEvent(event.data.pagestate);
        } else if (event.data.ajaxComplete) {
            this._fireStepEvent(event.data.ajaxComplete);
        }
    },
    logLine: function(line) {
        this.element.text(this.element.text() + line  + "\n");
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
        goradd.trigger(this.element, "teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        goradd.log ("changeVal", step, id, val);
        var control = this._findElement(id);

        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        goradd.value(control, val);

        goradd.trigger(control, "change");
        this._fireStepEvent(step);
    },
    checkControl: function(step, id, val) {
        goradd.log ("checkControl", step, id, val);
        var control = this._findElement(id);

        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        control.checked = val;

        goradd.trigger(control, "change");
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
            goradd.trigger(this, "change");
        });

        // check whatever needs to be checked
        goradd.each(values, function() {
            var val = this;
            var elements = goradd.qa("input[name=" + id + "][value=" + val + "]");
            goradd.each(elements, function() {
                this.checked = true;
                goradd.trigger(this, "change");
            });
        });

        this._fireStepEvent(step);
    },
    _findElement: function(id) {
        return this._window.document.getElementById(id);
    },
    closeWindow: function(step) {
        this._window.close();
        this._fireStepEvent(step);
    },
    click: function (step, id) {
        goradd.log("click", step, id);
        var control = this._findElement(id);
        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        control.click();
        this._fireStepEvent(step);
    },
    callGoraddElementFunction: function (step, id, f, params) {
        goradd.log("GoraddF", step, id, f, params);
        var ret;

        var control = this._findElement(id);
        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        var func = goradd[f];
        if (!func) {
            this._fireStepEvent(step, "Could not find function " + f + " on goradd");
            return;
        }

        params.unshift(control);

        ret = func.apply(goradd, params);
        goradd.setControlValue(this.element.id, "jsvalue", ret);
        this._fireStepEvent(step);
    },
    typeChars: function (step, id, chars) {
        var control = this._findElement(id);
        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        goradd.value(control, chars);
        this._fireStepEvent(step);
    },
    focus: function (step, id) {
        goradd.log("focus", step, id);
        this._findElement(id).focus();
        this._fireStepEvent(step);
    }

};