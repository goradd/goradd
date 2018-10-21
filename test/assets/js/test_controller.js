/**************************************************************************
 * Goradd Test Controller Object
 *
 ****************************************************************************/

jQuery.widget( "goradd.testController",  {
    options: {
    },
    _window:null,
    _err:"",
    _create: function() {
        this._super();
    },
    logLine: function(line) {
        this.element.text(this.element.text() + line  + "\n");
    },
    loadUrl: function(step, url) {
        var self = this;
        if (!this._window || this._window.closed) {
            this._window = window.open(url, "testWindow");
            if (!this._window) {
                this._err += "Opening a popup window was blocked by the browser.";
                this.fireStepEvent(step);
                return;
            }
            this._window.addEventListener("load", function(event) {
                self._windowLoadEvent(event, step)
            });
            this._window.addEventListener("error", function(event) {
                self._windowErrorEvent(event, step)
            });


            /*
            this._on( this._window, {
                "load": function(event) {
                    this._windowLoadEvent(event, step)
                },
                "error": function(event) {
                    this._windowErrorEvent(event, step)
                }
            });*/
        } else {
            this._window.location.href = url;
        }
    },
    _windowLoadEvent: function(event, step) {
        // if we got a goradd form, get the formstate or the generated error
        $formstate = $(this._window.document).find("form[data-grctl=form] #Goradd__FormState");
        if ($formstate.length > 0) {
            goradd.setControlValue(this.element.attr("id"), "formstate", $formstate.val());
        }
        this._fireStepEvent(step, null);
    },
    _windowErrorEvent: function(event, step) {
        this._fireStepEvent(step, "Browser load error:" + event.error.message);
    },
    fireStepEvent(step) {
        err = this._err;
        this._err = "";
        this.element.trigger("goradd.teststep", {Step: step, Err: err});
    },
    changeVal: function(step, id, val) {
        var $control = this._findControl(id);

        if (!$control.length) {
            this._err += "Could not find control " + id;
            this.fireStepEvent(step);
            return;
        }

        $control.val(val);
        $control.trigger("change");
        this._fireStepEvent(step);
    },
    _findControl: function(id) {
        return $(this._window.document).find("#" + id);
    }

});