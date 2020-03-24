

goradd.TestController = goradd.extendWidget({
    constructor: function (element, options) {
        this._super(element);
        var self = this;
        goradd.log("Creating test controller");
        window.addEventListener("message", function (event) {
            self._receiveWindowMessage(event)
        }, false);
        this._window = null;
        this._step = 1;
    },
    _receiveWindowMessage: function(event) {
        goradd.log("Message received", event.data);
        if (event.data.pagestate) {
            this._formLoadEvent(event.data.pagestate);
        } else if (event.data.ajaxComplete) {
            this._fireStepEvent(event.data.ajaxComplete);
        } else if (event.data.testMarker) {
            this._fireMarker(event.data.testMarker);
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
            var localpath = this._window.location.href.substr(this._window.location.origin.length)
            if (localpath === url) {
                this._window.location.reload();
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
    _fireMarker(marker) {
        this.trigger("testmarker", marker);
    },

    changeVal: function(step, id, val) {
        val = JSON.parse(val);
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
        var self = this;
        goradd.log ("checkControl", step, id, val);
        var g = this._getGoraddObj(id);

        if (!g) {
            goradd.log("checkControl: element not found", id);
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        var val2 = g.element.checked;

        if (val !== val2) {
            g.click();
        }
        this._fireStepEvent(step);
    },
    /**
     * checkGroup simulates checking the given values in a group. For checkboxes, it also unchecks whatever is checked
     * prior to this. This will generate change events on whatever was changed.
     * @param step
     * @param groupName
     * @param values
     */
    checkGroup: function(step, groupName, values) {
        goradd.log ("checkGroup", step, groupName, values);
        var form = g$(this._window.goradd.form());

        var el = form.qs("input[name=" + groupName + "]");
        if (!el) {
            this._fireStepEvent(step,  "Could not find group " + groupName);
            return;
        }

        if (el.type === "radio") {
            // Check one radio button. The currently checked one should automatically uncheck.
            var s = "input[name=" + groupName + "][value='" + values[0] + "']";
            el = form.qs(s);
            if (el) {
                g$(el).click();
            }
            this._fireStepEvent(step);
            return;
        }

        // Deal with a list of checkboxes
        goradd.each(form.qa("input[name=" + groupName +"]"), function() {
            var toCheck = goradd.contains(values, this.value);
            if (this.checked && !toCheck) {
                g$(this).click(); // uncheck
            } else if (!this.checked && toCheck) {
                g$(this).click(); // check
            }
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
        g.click(function() {
            self._fireStepEvent(step);
        });
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
    // getHtmlElementInfo returns a specific value from an html object identified by the given selector.
    // TODO: chain values, and respond to functions and array references. "attributes.getNamedItem("width").value" would return the value of the width attribute
    getHtmlElementInfo: function (step, selector, attr) {
        goradd.log("GetHtml", step, selector, attr);

        var item = this._window.goradd.qs(selector);
        var ret = "";
        if (item) {
            ret = item[attr];
        }

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

goradd.registerWidget("goradd.TestController", goradd.TestController);