/**************************************************************************
 * Goradd Test Controller Object
 *
 ****************************************************************************/

jQuery.widget( "goradd.testController",  {
    options: {
    },
    _window:null,
    _err:"",
    _step:0,
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
            this.fireStepEvent(step);
            return;
        }/*
        $(this._window).on("load", function(event) {
            self._windowLoadEvent(event, step)
        });*/
        this._window.addEventListener("error", function(event) {
            self._windowErrorEvent(event, step)
        });
    },
    _formLoadEvent: function(formstate) {
        goradd.setControlValue(this.element.attr("id"), "formstate", formstate);
        this.fireStepEvent(this._step, null);
    },
    /*
    _windowLoadEvent: function(event, step) {
        this.fireStepEvent(step, null);
    },*/
    _windowErrorEvent: function(event, step) {
        this.fireStepEvent(step, "Browser load error:" + event.error.message);
    },
    fireStepEvent(step) {
        err = this._err;
        this._err = "";
        this.element.trigger("goradd.teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        var control = this._findControl(id);

        if (!control) {
            this._err += "Could not find control " + id;
            return;
        }

        $(control).val(val);

        // Note that jQuery is very quirky about calling events in another window, because it attaches its own events to the current window.
        var event = new Event('change', { 'bubbles': true });
        control.dispatchEvent(event);
    },
    _findControl: function(id) {
        return this._window.document.getElementById(id);
    },
    click: function (step, id) {
        var control = this._findControl(id);
        if (!control) {
            this._err += "Could not find control " + id;
            return;
        }
        var event = new MouseEvent('click', {
            view: window,
            bubbles: true,
            cancelable: true
        });
        control.dispatchEvent(event);
    }

});