/**************************************************************************
 * Goradd Test Controller Object
 *
 ****************************************************************************/

jQuery.widget( "goradd.testController",  {
    options: {
    },
    _window:null,
    _err:"",
    _step:1,
    _create: function() {
        var self = this;
        this._super();
        window.addEventListener("message", function(event) {
            self._receiveWindowMessage(event)
        } , false);
    },
    _receiveWindowMessage: function(event) {
        if (event.data.formstate) {
            this._formLoadEvent(event.data.formstate);
        } else if (event.data.ajaxComplete) {
            this._fireStepEvent(event.data.ajaxComplete);
        }
    },
    logLine: function(line) {
        this.element.text(this.element.text() + line  + "\n");
    },
    loadUrl: function(step, url) {
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
            this._err += "Opening a popup window was blocked by the browser.";
            this._fireStepEvent(step);
            return;
        }
        this._window.addEventListener("error", function(event) {
            self._windowErrorEvent(event, step)
        });

    },
    _formLoadEvent: function(formstate) {
        goradd.setControlValue(this.element.attr("id"), "formstate", formstate);
        this._fireStepEvent(this._step, null);
    },
    _windowErrorEvent: function(event, step) {
        this._fireStepEvent(step, "Browser load error:" + event.error.message);
    },
    _fireStepEvent(step) {
        err = this._err;
        this._err = "";
        this.element.trigger("goradd.teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        var control = this._findElement(id);

        if (!control) {
            this._err += "Could not find element " + id;
            return;
        }

        $(control).val(val);

        // Note that jQuery is very quirky about calling events in another window, because it attaches its own events to the current window.
        // So, we instead using native javascript to fire off these events.
        var event = new Event('change', { 'bubbles': true });
        control.dispatchEvent(event);
        event = new CustomEvent('teststep', { bubbles: true, detail: step });
        control.dispatchEvent(event);
    },
    _findElement: function(id) {
        return this._window.document.getElementById(id);
    },
    click: function (step, id) {
        var control = this._findElement(id);
        if (!control) {
            this._err += "Could not find element " + id;
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
    jqValue: function (step, id, f, params) {
        var ret;

        var control = this._findElement(id);
        if (!control) {
            this._err += "Could not find element " + id;
            this._fireStepEvent(step);
            return;
        }
        var $control = $(control);
        var func = $control[f];
        if (!func) {
            this._err += "Could not find function " + f + " on jQuery element " + id;
            this._fireStepEvent(step);
            return;
        }

        ret = func.apply($control, params);
        goradd.setControlValue(this.element.attr("id"), "jsvalue", ret);
        this._fireStepEvent(step);
    }


});

