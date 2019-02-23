/**************************************************************************
 * Goradd Test Controller Object
 *
 ****************************************************************************/

jQuery.widget( "goradd.testController",  {
    options: {
    },
    _window:null,
    _step:1,
    _create: function() {
        goradd.log("Creating test controller");
        var self = this;
        this._super();
        window.addEventListener("message", function(event) {
            self._receiveWindowMessage(event)
        } , false);
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
        this.element.text(this.element.text() + line  + "\n");
    },
    loadUrl: function(step, url) {
        goradd.log("loadUrl", step, url);
        var self = this;

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
        goradd.setControlValue(this.element.attr("id"), "pagestate", pagestate);
        this._fireStepEvent(this._step);
    },
    _windowErrorEvent: function(event, step) {
        this._fireStepEvent(step,  "Browser load error:" + event.error.message);
    },
    _fireStepEvent(step, err) {
        this.element.trigger("goradd.teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        goradd.log ("changeVal", step, id, val);
        var control = this._findElement(id);

        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }

        $(control).val(val);

        // Note that jQuery is very quirky about calling events in another window, because it attaches its own events to the current window.
        // So, we instead use native javascript to fire off these events.
        var event = new Event('change', { 'bubbles': true });
        control.dispatchEvent(event);
        event = new CustomEvent('teststep', { bubbles: true, detail: step });
        control.dispatchEvent(event);
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
        var event = new MouseEvent('click', {
            view: window,
            bubbles: true,
            cancelable: true
        });
        control.dispatchEvent(event);
        event = new CustomEvent('teststep', { bubbles: true, detail: step });
        control.dispatchEvent(event);

    },
    callJqueryFunction: function (step, id, f, params) {
        goradd.log("jqValue", step, id, f, params);
        var ret;

        var control = this._findElement(id);
        if (!control) {
            this._fireStepEvent(step,  "Could not find element " + id);
            return;
        }
        var $control = $(control);
        var func = $control[f];
        if (!func) {
            this._fireStepEvent(step, "Could not find function " + f + " on jQuery element " + id);
            return;
        }

        ret = func.apply($control, params);
        goradd.setControlValue(this.element.attr("id"), "jsvalue", ret);
        this._fireStepEvent(step);
    },
    typeChars: function (step, id, chars) {
        //KeyEvent.simulate(chars, [], this._findElement(id));
        $(this._findElement(id)).val(chars);
    },
    focus: function (step, id) {
        goradd.log("focus", step, id);
        this._findElement(id).focus();
    }



});

