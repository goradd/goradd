/**
 * goradd.js
 *
 * This is the shim between the goradd go code and the browser. It enables ajax and other kinds of
 * communication between the server and the client.
 *
 * Goals:
 *  - Compatible with all current browsers, IE 10+ and Opera Mobile (ES5).
 *  - Provide utility code to javascript widgets and plugins.
 *  - Do not use jQuery or other frameworks, but be compatible if its used by the developer.
 */

if (!function () {
    "use strict";
    return Function.prototype.bind && XMLHttpRequest && !this;
}()) {
    window.location = "/Unsupported.g";
}

var goradd;
var g$;

(function( ) {
    // Put everything in a function so we can have private functions and variables.
    // By convention, things starting with underscore will be private.
"use strict";

/**
 * Private functions and members used by code below.
 */

var _controlValues = {};
var _formObjsModified = {};
var _ajaxError = false;
var _blockEvents = false;
var _inputSupport = true;
var _finalCommands = [];
var _prevUpdateTime = 0;

function _toKebab(s) {
    return  s.replace(/[A-Z]/g, function(m) {
        return "-" + m.toLowerCase();
    });
}

function _fromKebab(s) {
    return s.replace( /-([a-z])/gi, function ( o, i ) { return i.toUpperCase(); } );
}

/**
 * formObjChanged is an event handler that records that a control has changed in order to synchronize the control with
 * the server on the next request. Send the formObjChanged event to the control
 * that changed, and it will bubble up to the form and be caught here.
 * @param event
 */
function _formObjChanged(event) {
    _formObjsModified[event.target.id] = true;
}

/**
 * Gets the data to be sent to an ajax call as post data. This will be called from the ajax queueing function, and
 * will erase the cache of changed objects to prepare for the next call.
 *
 * @param {object} params An object containing the following:
 *  controlId {string}: The control id to post an action to
 *  eventId {int}: The event id
 *  async: If true, process the event asynchronously without waiting for other events to complete
 *  formId: The id of the form getting posted
 *  values {object} (optional): An optional object, that contains values to send with the event
 *      event: The event's action value, if one is provided. This can be any type, including an object.
 *      action: The action's action value, if one is provided. Any type.
 *      control: The control's action value, if one is provided. Any type.
 * @return {object} Post Data
 * @private
 */
function _getAjaxData(params) {
    var form = goradd.form(),
        controls = g$(form).qa("input,select,textarea"),
        postData = {};

    // Notify controls we are about to post.
    g$(form).trigger("posting", "Ajax");

    goradd.each(controls, function(i,c) {
        var id = c.id;
        var blnForm = (id && (id.substr(0, 8) === "Goradd__"));

        if (!_inputSupport || // if not oninput support, then post all the controls, rather than just the modified ones, because we might have missed something
            _ajaxError || // Ajax error would mean that _formObjsModified is invalid. We need to submit everything.
            (id && _formObjsModified[id]) ||  // We try to ignore controls that have not changed to reduce the amount of data sent in an ajax post.
            blnForm) {  // all controls with Goradd__ at the beginning of the id are always posted.

            switch (c.type) {
                case "radio":
                    // Radio buttons listen to their name.
                    var n = c.name;
                    var radio = form.querySelector('input[name=' + n + ']:checked');
                    var val = null;
                    if (radio) {
                        val = radio.value;
                    }
                    postData[n] = val;
                    break;
                case "checkbox":
                    postData[id] = c.checked;
                    break;
                default:
                    // All goradd controls and subcontrols MUST have an id for this to work.
                    // There is a special case for checkbox groups, but they get handled on the server
                    // side differently between ajax and server posts.
                    postData[id] = g$(c).val();
                    break;
            }
        }
    });

    // Update most of the Goradd__ parameters explicitly here. Others, like the state and form id will have been handled above.
    params.callType = "Ajax";
    if (!goradd.isEmptyObj(_controlValues)) {
        params.controlValues = _controlValues;
    }
    postData.Goradd__Params = JSON.stringify(params);

    _ajaxError = false;
    _formObjsModified = {};
    _controlValues = {};
    return postData;
}

/**
 * Displays the ajax error in either a popup window, or a new web page.
 * @param resultText {string}
 * @param err
 * @private
 */
function _displayAjaxError(resultText, err) {
    var objErrorWindow;

    _ajaxError = true;
    _blockEvents = false;

    if (resultText.substr(0, 15) === "<!DOCTYPE html>") {
        window.alert("An error occurred.\r\n\r\nThe error response will appear in a new popup.");
        objErrorWindow = window.open("about:blank", "qcubed_error", "menubar=no,toolbar=no,location=no,status=no,scrollbars=yes,resizable=yes,width=1000,height=700,left=50,top=50");
        objErrorWindow.focus();
        objErrorWindow.document.write(resultText);
    } else {
        if (err) {
            resultText = err.toString();
            if (err.sourceURL) {
                resultText += " File:" + err.sourceURL
            }
            if (err.line) {
                resultText += " Line:" + err.line;
            }
        }
        var el = goradd.tagBuilder("div").attr("id", "Goradd_AJAX_Error").
        html("<button onclick=\"window.goradd.g('Goradd_AJAX_Error').remove()\">OK</button>").
        appendTo(goradd.form());
        goradd.tagBuilder("div").html(resultText).appendTo(el);
    }
}

/**
 * Responds to the part of an ajax response that must be handled serially before other handlers can fire.
 *
 * @param {Object} json     json generated by goradd application
 * @param {Object} params   option parameters
 * @private
 */
function _processImmediateAjaxResponse(json, params) {
    goradd.each(json.controls, function(id) {
        var el = goradd.el(id),
            $ctrl = g$(el),
            wrapper = goradd.el(id + "_ctl");

        if (this.value !== undefined && $ctrl) {
            $ctrl.val(this.value);
        }

        if (this.attributes !== undefined && $ctrl) {
            $ctrl.prop(this.attributes);
        }

        if (this.html !== undefined) {
            if (wrapper !== null) {
                // Control's wrapper was found, so replace the control and the wrapper
                g$(wrapper).htmlBefore(this.html);
                g$(wrapper).remove(wrapper);
            } else if ($ctrl) {
                // control was found without a wrapper, replace it in the same position it was in.
                // remove related controls (error, name ...) for wrapper-less controls
                var relSelector = "[data-grel='" + id + "']",
                    relatedItems = goradd.qa(relSelector);

                var p = $ctrl.parents();
                var relatedParent = p.filter(function(item) {
                    return g$(item).matches(relSelector);
                }).pop();

                if (relatedParent) {
                    relatedParent.insertAdjacentElement("beforebegin", el);
                }
                if (relatedItems && relatedItems.length > 0) {
                    goradd.each(relatedItems, function() {
                        g$(this).remove();
                    });
                }
                $ctrl.htmlBefore(this.html);
                $ctrl.remove();
            }
            else {
                // control is being injected at the top level, so put it at the end of the form.
                g$(goradd.form()).appendHtml(this.html);
            }
        }
    });

    _registerControls();

    if (json.watcher && params.controlId) {
        goradd.broadcastChange();
    }
    if (json.ss) {
        goradd.each(json.ss, function (i,v) {
            goradd.loadStyleSheetFile(v, "all");
        });
    }
    if (json.alert) {
        goradd.each(json.alert, function (i,v) {
            window.alert(v);
        });
    }
}
/**
 * Process the part of an ajax response that can be deferred and so be executed in parallel with other operations.
 *
 * @param {object} json  Json generated by the goradd application.
 * @private
 */
function _processDeferredAjaxResponse(json) {
    goradd.each(json.commands, function (i,command) {
        if (command.final &&
            goradd.ajaxq.isRunning()) {
            _enqueueFinalCommand(command);
        } else {
            goradd.processCommand(command);
        }
    });
    if (json.winclose) {
        window.close();
    }
    if (json.loc) {
        if (goradd._closeWebSocket) {
            goradd._closeWebSocket(1001);
        }
        if (json.loc === "reload") {
            window.location.reload(true);
        } else {
            document.location = json.loc;
        }
    }
    if (json.profileHtml) {
        var c = goradd.el("dbProfilePane");
        if (!c) {
            g$(goradd.form()).htmlAfter("<div id = 'dbProfilePane'></div>");
            c = goradd.el("dbProfilePane");
        }
        c.innerHTML = json.profileHtml;
    }
    goradd.testStep();
}

/**
 * Places the given command in the queue so that it is executed last.
 * @param command
 * @private
 */
function _enqueueFinalCommand(command) {
    _finalCommands.push(command);
}

/**
 * Execute the final commands.
 * @private
 */
function _processFinalCommands() {
    while(_finalCommands.length) {
        var command = _finalCommands.pop();
        goradd.processCommand(command);
    }
}

/**
 * Convert from JSON return value to an actual jQuery object. Certain structures don't work in JSON, like closures,
 * but can be part of a javascript object. We use special codes to piece together functions, closures, dates, etc.
 * @param params
 * @returns {*}
 * @private
 */
function _unpackArray(params) {
    if (!params) {
        return null;
    }
    var newParams = [];

    goradd.each(params, function (index, item){
        if (Array.isArray(item)) {
            item = _unpackArray (item);
        } else if (typeof item === 'object' && item !== null) {
            if (item.goraddObject) {
                item = _unpackObj(item);  // top level special object
            }
            else {
                // look for special objects inside top level objects.
                var newItem = {};
                goradd.each (item, function (key, obj) {
                    newItem[key] = _unpackObj(obj);
                });
                item = newItem;
            }
        }
        newParams.push(item);
    });
    return newParams;
}

/**
 * Given an object coming from goradd, will attempt to decode the object into a corresponding javascript object.
 * @param obj
 * @returns {*}
 * @private
 */
function _unpackObj(obj) {
    var params;

    if (typeof obj === "object" && obj !== null) {
        if (Array.isArray(obj)) {
            return _unpackArray(obj);
        } else if (obj.goraddObject) {
            switch (obj.goraddObject) {
                case 'closure':
                    if (obj.params) {
                        params = [];
                        goradd.each (obj.params, function (i, v) {
                            params.push(_unpackObj(v)); // recurse
                        });

                        return new Function(params, obj.func);
                    } else {
                        return new Function(obj.func);
                    }

                case 'date':
                    return goradd.unpackJsonDate(obj);

                case 'varName':
                    // Find the variable value starting at the window context.
                    var vars = obj.varName.split(".");
                    var val = window;
                    goradd.each (vars, function (i, v) {
                        val = val[v];
                    });
                    return val;

                case 'func':
                    // Returns the result of the given function called immediately
                    // Find the function and context starting at the window context.
                    var target = window;
                    if (obj.context) {
                        var objects = obj.context.split(".");
                        goradd.each (objects, function (i, v) {
                            target = target[v];
                        });
                    }

                    if (obj.params) {
                        params = [];
                        goradd.each (obj.params, function (i, v) {
                            params.push(_unpackObj(v)); // recurse
                        });
                    }
                    var func = target[obj.func];

                    return func.apply(target, params);
            }
        }
        else {
            var newItem = {};
            goradd.each (obj, function (key, obj2) {
                newItem[key] = _unpackObj(obj2);
            });
            return newItem;
        }
    }

    return obj; // no change
}

function _registerControls() {
    var els = goradd.qa('[data-grctl]');
    goradd.each(els, function() {
        _registerControl(this);
    });
}

function _registerControl(ctrl) {
    if (!ctrl) {
        return;
    }

    // get the widget
    var g = g$(ctrl);

    if (g.data('gr-reg') === 'reg') {
        return // this control is already registered
    }

    if (ctrl.tagName === "FORM") {
        return;
    }

    g.data('gr-reg', 'reg'); // mark the control as registered so we don't attach events twice. Has the side effect
                             // of attaching the widget to the control.

    // detect changes to objects before any changes trigger other events
    if (ctrl.type === 'checkbox' || ctrl.type === 'radio') {
        // clicks are equivalent to changes for checkboxes and radio buttons, but some browsers send change way after a click. We need to capture the click first.
        g.on('click', _formObjChanged);
    }
    g.on('change input', _formObjChanged, {capture: true}); // make sure we get these events before later attached events

    // widget support, using declarative methods
    if (goradd.widget.new) {
        var widget;
        var options = {};
        goradd.each(g.attr(), function(k,v) {
            if (k === "data-gr-widget") {
                widget = v;
            } else if (k.substr(0, 12) === "data-gr-opt-") {
                options[_fromKebab(k.substr(12))] = v;
            }
        });
        if (widget) {
            widget = goradd.widget.new(widget, options, ctrl);
            // Replace the control's widget with the new one. There can be only one goradd widget associated with
            // a particular control. We will need some other mechanism for mixins if needed.
            ctrl.goradd.widget = widget;
        }
    }
}


/**
 * g$ is a shortcut for goradd.g(). It wraps an element with additional functions defined here to more easily manipulate
 * the element. el can be either an actual HTMLElement, or the id of one.
 * One main difference between jQuery's wrapper and this one is that jQuery wraps an array of elements, and we only
 * wrap one element. Also, all functions return and HTMLElement or array of elements, not a wrapped element.
 * @param el {string | HTMLElement}
 * @returns {Element.goradd.widget}
 */
g$ = function(el) {
    return goradd.g(el);
};

// noinspection JSUnusedGlobalSymbols
 /**
 * @namespace goradd
 */
goradd = {
    /**
     * General support library. Here we recreate a few useful functions from jquery.
     */

    /**
     * Extend merges keys and values of objects into the target object.
     * This version of extend is primarily for the purpose of adding
     * capabilities to javascript classes. It does not do deep merging, but it will copy the members of plain objects
     * if it finds a plain object. Other kinds of objects are copied by reference.
     * If only one argument is provided, the target is the goradd object itself.
     *
     * @param target... {object}
     * @returns {*}
     */
    extend: function( target ) {
        var input = Array.prototype.slice.call( arguments, 1 ),
            key,
            value;

        if (arguments.length === 1) {
            input = [target];
            target = goradd;
        }

        var inputIndex = 0,
            inputLength = input.length;

        for ( ; inputIndex < inputLength; inputIndex++ ) { // iterate through all arguments in order
            var obj = input[ inputIndex ];
            for ( key in obj ) { // iterate through the keys in each argument
                value = obj[ key ];
                if ( obj.hasOwnProperty( key ) && value !== undefined ) {

                    // Clone plain objects
                    if ( goradd.isPlainObject( value ) ) {
                        target[ key ] = goradd.isPlainObject( target[ key ] ) ?
                            goradd.extend( {}, target[ key ], value ) :

                            // Don't extend strings, arrays, etc. with objects
                            goradd.extend( {}, value );
                    } else { // Copy everything else by reference
                        target[ key ] = value;
                    }
                }
            }
        }
        return target;
    },

    /**
     * el returns the html element t. t can be an id, or an element, and if an element, it will just return the element
     * back. This is used below so that all the functions can pass either an element, or the id of an element. Returns
     * null if not found.
     * @param t {string|object}
     * @returns {Object}
     */
    el: function(t) {
        if (typeof t === "object") {
            return t;
        }
        return document.getElementById(t);
    },
    qs: function(sel) {
        return document.querySelector(sel);
    },
    /**
     * qa is a querySelectorAll call that returns an actual array, and not a NodeList.
     * Returns empty array if selector has no results.
     * @param sel {string} The css selector to find
     * @returns {HTMLElement[]}
     */
    qa: function(sel) {
        return Array.prototype.slice.call(document.querySelectorAll(sel));
    },
    isEmptyObj: function(o) {
        if (!o) return false;
        for (var name in o ) {
            return false;
        }
        return true;
    },
    form: function() {
        // TODO: cache this, it will not change. No reason to do this over and over.
        return goradd.qs('form[data-grctl="form"]');
    },
    /**
     * matches returns true if the given element matches the css selector.
     * @param el
     * @param sel
     * @returns {boolean}
     */
    matches: function(el, sel) {
        return g$(el).matches(sel);
    },
    /**
     * loadJavaScriptFile will dynamically load a javascript file. It is designed to be called during ajax calls or
     * other times when a dynamically loaded javascript file is required.
     * @param strScript
     * @param attributes
     */
    loadJavaScriptFile: function(strScript, attributes) {
        var script = document.createElement("script");
        script.src = strScript;
        script.type = 'text/javascript';
        goradd.extend(script, attributes);

        var head = document.getElementsByTagName('head')[0];
        head.appendChild(script);
    },
    /**
     * loadStyleSheetFile dynamically loads a style sheet file. It is used by the ajax code.
     * @param strStyleSheetFile
     * @param attributes
     */
    loadStyleSheetFile: function(strStyleSheetFile, attributes) {
        var link = document.createElement("link");
        link.rel = "stylesheet";
        link.href = strStyleSheetFile;
        goradd.extend(link, attributes);
        var head = document.getElementsByTagName('head')[0];
        head.appendChild(link);
    },
    /**
     * each is a recreation of the jQuery each function, but for our targeted browsers only. It iterates the given object,
     * calling the function for each item found. If the object is an array, or something array-like, like a nodelist,
     * it will pass the index and the item to the function. For a regular object, it will pass the key and the item.
     * "this" is set to the item each time as well.
     * @param obj {object}
     * @param f {function}
     */
    each: function(obj, f) {
        if (!obj) return;
        if (typeof(obj) !== "object") return;
        var i;

        // isArrayLike needs to return true for nodelists.
        var isArrayLike = Array.isArray(obj) || ("length" in obj && typeof(obj.length) === "number" && (obj.length === 0 || ((obj.length - 1) in obj)));
        if (isArrayLike) {
            if (obj.length === 0) return;
            for (i = 0; i < obj.length; i++) {
                if (f.call( obj[ i ], i, obj[ i ] ) === false) break;
            }
        } else {
            for (i in obj) {
                if ( f.call( obj[ i ], i, obj[ i ] ) === false ) {
                    break;
                }
            }
        }
    },
    /**
     * contains returns true if needle is in the array a
     * @param a {Array}
     * @param needle {*}
     * @returns {boolean}
     */
    contains: function(a, needle) {
        return (a.indexOf(needle) !== -1);
    },
    isPlainObject: function( obj ) {
        var proto, Ctor;

        // Detect obvious negatives
        // Use toString instead of jQuery.type to catch host objects
        if ( !obj || {}.toString.call( obj ) !== "[object Object]" ) {
            return false;
        }

        proto = Object.getPrototypeOf( obj );

        // Objects with no prototype (e.g., `Object.create( null )`) are plain
        if ( !proto ) {
            return true;
        }

        // Objects with prototype are plain iff they were constructed by a global Object function
        Ctor = {}.hasOwnProperty.call( proto, "constructor" ) && proto.constructor;
        return typeof Ctor === "function" && {}.hasOwnProperty.toString.call( Ctor ) === {}.hasOwnProperty.toString.call(Object);
    },
    /**
     * setRadioInGroup is a specialized function called from goradd go code.
     * It sets the given radio button to being checked in a group. Since the goradd code already knows what we want to
     * check, that button is just set. However, the button that gets unset by the browser needs to communicate to the
     * go code that it is getting unset.
     * @param id
     */
    setRadioInGroup: function(id) {
        var el = goradd.el(id);
        if (el.type !== "radio") {
            return; // not a radio button, or not part of a group
        }
        var prevItem;

        if (el.name) {
            prevItem = goradd.qs("input[type='radio'][name='" + el.name +"']:checked");
        }
        el.checked = true;
        if (prevItem) {
            g$(el).trigger('formObjChanged');
        }
    },
    /**
     * Unpacks a date object that was packed by dateTime.DateTime.MarshalJson. If the date represented a
     * timestamp on the server side, it will be a timestamp here, but the time will be in local time.
     * In other words, if the server timezone and browser timezone are different,
     * then they will show different times, but both will correspond to the same world time.
     * If on the server side the date represented simply a date and time in local time,
     * the date will become the same date and time in local time here. If the server timezone and browser
     * timezone are different, they will both show the same time, meaning they will not be the same world time.
     * If it was a zero date on the server, it becomes a null here.
     *
     * This solves some problems inherent in the traditional JSON date format consisting of an ISO8601 string.
     *
     * @param o
     * @returns {null|Date}
     */
    unpackJsonDate(o) {
        if (o.z) {
            return null;
        } else if (o.t) {
            return new Date(Date.UTC(o.y, o.mo, o.d, o.h, o.m, o.s, o.ms));
        } else {
            return new Date(o.y, o.mo, o.d, o.h, o.m, o.s, o.ms);
        }
    },

    /**
     * Adds a value to the next ajax or server post for the specified control. You can either call this ongoing, or
     * call it in response to the "posting" event. This is the preferred way for custom javascript controls to send data
     * to their goradd counterparts.
     *
     * @param {string} strControlId     The controlId of the property to set
     * @param {string} strProperty      The name of the property
     * @param strNewValue               The new value of the property. Can be any type, including string, number, object or array
     */
    setControlValue: function(strControlId, strProperty, strNewValue) {
        if (!_controlValues[strControlId]) {
            _controlValues[strControlId] = {};
        }
        _controlValues[strControlId][strProperty] = strNewValue;
    },
    /**
     * Initializes form related scripts. This is called by injected code on a goradd form.
     * TODO: Combine with initialize and waiting for dom loaded
     */
    initForm: function () {
        var form =  goradd.form();
        g$(form).on('formObjChanged', _formObjChanged); // Allow any control, including hidden inputs, to trigger a change and post of its data.
        g$(form).on('submit', function(event) {
            if (!goradd.el('Goradd__Params').value) { // did postBack initiate the submit?
                // if not, prevent implicit form submission. This can happen in the rare case we have a single field and no submit button.
                event.preventDefault();
            } else {
                // Check html5 validity in case it is being used.
                if (typeof form.reportValidity !== "function" ||
                    form.hasAttribute("novalidate") ||
                    form.reportValidity()) {

                    form.submit();
                }
            }
        });
        _registerControls();
    },

    /**
     * Post the form. ServerActions call this.
     *
     * @param {Object} params An object containing the following:
     *  controlId {string}: The control id to post an action to
     *  eventId {int}: The event id
     *  async: If true, process the event asynchronously without waiting for other events to complete
     *  values {object} (optional): An optional object, that contains values coming to send with the event
     *      event: The event's action value, if one is provided. This can be any type, including an object.
     *      action: The action's action value, if one is provided. Any type.
     *      control: The control's action value, if one is provided. Any type.
     *
     */
    postBack: function(params) {
        if (_blockEvents) {
            return;  // We are waiting for a response from the server
        }

        var form = goradd.form();
        var gForm = g$(form);

        params.callType = "Server";

        // Notify custom controls that we are about to post

        gForm.trigger("posting", "Server");

        // Post custom javascript control values
        if (goradd.isEmptyObj(_controlValues)) {
            params.controlValues = _controlValues;
        }
        goradd.el('Goradd__Params').value = JSON.stringify(params);

        // trigger our own form submission so we can catch it
        gForm.trigger("submit");
    },

    /**
     * Posts an ajax call to the ajax queue. Ajax actions call this.
     *
     * @param {Object} params An object containing the following:
     *  controlId {string}: The control id to post an action to
     *  eventId {number}: The event id
     *  async {boolean}: If true, process the event asynchronously without waiting for other events to complete
     *  values {Object} (optional): An optional object, that contains values coming to send with the event
     *      event: The event's action value, if one is provided. This can be any type, including an object.
     *      action: The action's action value, if one is provided. Any type.
     *      control: The control's action value, if one is provided. Any type.
     *
     * @return {void}
     */
    postAjax: function(params) {
        var form = goradd.form(),
            formAction = g$(form).attr("action"),
            async = params.hasOwnProperty("async");

        if (_blockEvents) {
            return;
        }

        params.formId = form.id;

        goradd.log("postAjax", params);

        // Use an ajax queue so ajax requests happen synchronously
        goradd.ajaxq.enqueue(function() {
            var data = _getAjaxData(params);
            goradd.log("Gathered ajax data: " + JSON.stringify(data));

            return {
                url: formAction,
                data: data,
                /**
                 * @param result {string}
                 * @param err {object}
                 * @returns {boolean}
                 */
                error: function (result, err) {
                    _displayAjaxError(result, err);
                    goradd.testStep();
                    return false;
                },
                /**
                 * @param json {object}
                 */
                success: function (json) {
                    goradd.log("Ajax success ", json);

                    if (json.js) {
                        for (var k in json.js) {
                            goradd.loadJavaScriptFile(k, json.js[k]);
                        }
                    }
                    _processImmediateAjaxResponse(json, params);
                    // TODO: Wait until javascripts above are loaded before proceeding?
                    _processDeferredAjaxResponse(json);
                    _blockEvents = false;
                }
            };
        }, async);
    },
    /**
     * Start me up.
     */
    initialize: function() {
        /*
        $(window).on ("storage", function (o) {
            if (o.originalEvent.key === "goradd.broadcast") {
                goradd.updateForm();
            }
        });*/

        _inputSupport = 'oninput' in document;
        // IE 9 has a major bug in oninput, but we are requiring IE 10+, so no problem.
        // I think the only major browser that does not support oninput is Opera mobile.

        g$(goradd.form()).on("ajaxQueueComplete", _processFinalCommands);

        // TODO: Add a detector of the back button. This detector should ping the server to make sure the pagestate exists on the server. If not,
        // it should reload the form.
    },
    /**
     * Process a single command. This is called both from ajax and javascript.
     * @param {object} command
     */
    processCommand: function(command) {
        var params,
            objs;

        if (command.script) {
            // TODO: clean this up a bit by using ids for inserted scripts. Might have multiple scripts for the same id though.
            var script   = document.createElement("script");
            script.type  = "text/javascript";
            script.text  = command.script;
            document.body.appendChild(script);
        }
        else if (command.selector) {
            params = _unpackArray(command.params);

            if (typeof command.selector === 'string') {
                // general selector
                objs = goradd.qa(command.selector);
            } else {
                // First item is the id to select on
                objs = g$(command.selector[0]).qa(command.selector[1]);
            }

            goradd.each (objs, function (i,v) {
                var $c = g$(v);
                if (typeof $c[command.func] === "function") {
                    $c[command.func].apply($c, params);
                }
            });
        }
        else if (command.func) {
            params = _unpackArray(command.params);

            // Find the function by name. Walk an object list in the process.
            objs = command.func.split(".");
            var obj = window;
            if (command.id) {
                obj = g$(command.id);
            } else if (command.jqueryId) {
                obj = jQuery('#' + command.jqueryId);
            }
            var ctx = null;

            goradd.each (objs, function (i, v) {
                ctx = obj;
                obj = obj[v];
                if (!obj) {
                    var p = Object.getPrototypeOf(ctx);
                    if (p && p[v]) {
                        obj = p[v];
                    }
                }
            });
            // obj is now a function object, and ctx is the parent of the function object
            obj.apply(ctx, params);
        }
    },
    updateForm: function() {
        // call this whenever you generally just need the form to update without a specific action.
        var newTime = new Date().getTime();

        // the following code prevents too many updates from happening in a short amount of time.
        // the default will update no faster than once per minUpdateInterval.
        if (newTime - _prevUpdateTime > goradd.minUpdateInterval) {
            //refresh immediately
            goradd.log("Immediate update");
            _prevUpdateTime = new Date().getTime();
            goradd.postAjax ({});
            goradd.clearTimer('goradd.update');
        } else if (!_objTimers['goradd.update']) {
            // delay to let multiple fast actions only trigger periodic refreshes
            goradd.log("Delayed update");
            goradd.setTimer ('goradd.update', goradd.updateForm, goradd.minUpdateInterval);
        }
        // else we already have a queued update, so no need to queue another one
    },

    /******************************************
     * Stub functions for debugging and testing
     * See goradd-testing.js
     ******************************************/

    testStep: function(event) {
    },
    log: function() {
    },

    /***********************
     * Utility Functions
     ***********************/

    /**
     * Sets a cookie with the given parameters. Potentially called by the goradd app.
     * @param name
     * @param val
     * @param expires
     * @param path
     * @param dom
     * @param secure
     */
    setCookie: function(name, val, expires, path, dom, secure) {
        var cookie = name + "=" + encodeURIComponent(val) + "; ";

        if (expires) {
            cookie += "expires=" + expires.toUTCString() + "; ";
        }

        if (path) {
            cookie += "path=" + path + "; ";
        }
        if (dom) {
            cookie += "domain=" + dom + "; ";
        }
        if (secure) {
            cookie += "secure;";
        }

        document.cookie = cookie;
    },

    //////////////////////////////
    // Action queue support
    //////////////////////////////
    /* Javascript has a problem when two events happen simultaneously. In particular, a click event might also
    result in a change event, and under certain circumstances this could cause the click event to be dropped. In particular,
    if the change event moves the focus away from the button that was clicked, the click event will not record.
    We therefore delay the processing of all events to try to queue them up before processing.
    Its very strange. Something to debug at a future date.
    */

    _actionQueue: [],
    queueAction: function(params) {
        if (params.last) {
            var delay = 0;

            goradd._actionQueue.forEach(function(item) {
                if (item.d > delay) {
                    delay = item.d;
                }
            });
            params.d = delay + 1;
        }
        goradd.log("queueAction: " + params.name);
        goradd._actionQueue.push(params);
        if (!goradd.hasTimer("goradd.actions")) {
            goradd.setTimer("goradd.actions", goradd.processActions, 10);
        }
    },
    processActions: function() {
        while (goradd._actionQueue.length > 0) {
            var params = goradd._actionQueue.shift();
            goradd.log("processAction: " + params.name + " delay: " + params.d);
            if (params.d > 0) {
                setTimeout(params.f, params.d);
            } else {
                params.f();
            }
        }
    },
    //////////////////////////////
    // Goradd Action Support
    //////////////////////////////
    /**
     * These support the various GoraddFunction actions available in the action package.
     */
    msg: function(m) {
        alert(m);
    },
    redirect: function(newLocation) {
        window.location = newLocation
    },
    proxyVal: function(event) {
        var target = event.target;
        if (!!event.goradd && !!event.goradd.match) {
            target = event.goradd.match;
        }
        return g$(target).data("grAv");
    }

};

/**
 * Named timer support. These allow you to create timers without having to keep a copy of the timer around.
 */

var _objTimers = {};

goradd.extend({
    /**
     * Sets the named timer, given an action and a delay.
     * @param strTimerId
     * @param action
     * @param intDelay
     */
    setTimer: function (strTimerId, action, intDelay) {
        goradd.log("setTimer", strTimerId, intDelay);
        _objTimers[strTimerId] = setTimeout(
            function () {
                goradd.clearTimer(strTimerId);
                action();
            }, intDelay);
    },
    hasTimer: function (strTimerId) {
        return !!_objTimers[strTimerId];
    },
    /**
     * Creates a timer that can perform periodic events, and that fires the timerexpiredevent event when it is done.
     * @param strControlId
     * @param intDeltaTime
     * @param blnPeriodic
     */
    startTimer: function (strControlId, intDeltaTime, blnPeriodic) {
        var strTimerId = strControlId + '_ct';
        goradd.stopTimer(strControlId, blnPeriodic);
        if (blnPeriodic) {
            _objTimers[strTimerId] = setInterval(function () {
                g$(strControlId).trigger('timerexpiredevent');
            }, intDeltaTime);
        } else {
            _objTimers[strTimerId] = setTimeout(function () {
                g$(strControlId).trigger('timerexpiredevent');
            }, intDeltaTime);
        }
    },
    /**
     * Stops the named timer, allowing you to specify whether its a periodic timer or not.
     * @param strControlId
     * @param blnPeriodic
     */
    stopTimer: function (strControlId, blnPeriodic) {
        var strTimerId = strControlId + '_ct';
        if (_objTimers[strTimerId]) {
            if (blnPeriodic) {
                clearInterval(_objTimers[strTimerId]);
            } else {
                clearTimeout(_objTimers[strTimerId]);
            }
            _objTimers[strTimerId] = null;
        }
    },
    /**
     * Clears the named timer.
     * @param {string} strTimerId
     */
    clearTimer: function (strTimerId) {
        if (_objTimers[strTimerId]) {
            goradd.log("clearTimer", strTimerId);
            clearTimeout(_objTimers[strTimerId]);
            _objTimers[strTimerId] = null;
        }
    },
});

///////////////////////////////
// Watcher support
///////////////////////////////
goradd.minUpdateInterval = 500; // milliseconds to limit broadcast updates. Feel free to change this.
goradd.broadcastChange = function () {
    if ('localStorage' in window && window.localStorage !== null) {
        var newTime = new Date().getTime();
        localStorage.setItem("goradd.broadcast", newTime); // must change value to induce storage event in other windows
    }
};



/////////////////////////////////
// Testing support
/////////////////////////////////

goradd.getPageState = function() {
    return document.getElementById("Goradd__PageState").value;
};

goradd.currentStep = 0;

/////////////////////////////////
// Tag Builder
/////////////////////////////////

/**
 * tagBuilder returns a TagBuilder. Use it as follows:
 * tag = goradd.tagBuilder("div").attr("class", "myClass").text("I am text").appendTo("objId");
 * @param tag {string}
 * @returns {goradd.TagBuilder}
 */
goradd.tagBuilder = function(tag) {
    return new goradd.TagBuilder(tag);
};
/**
 * TagBuilder uses a builder pattern to create and place html tags.
 *
 * @param tag {string}
 * @constructor
 */
goradd.TagBuilder = function(tag) {
    this.el = document.createElement(tag);
};
/**
 *
 * @type {{appendTo: (function((Object|string)): *), insertInto: (function((Object|string)): *), replace: (function((Object|string)): *), html: (function(string): goradd.TagBuilder), text: (function(string): goradd.TagBuilder), attr: (function(string, string): goradd.TagBuilder), insertAfter: (function((Object|string)): *), insertBefore: (function((Object|string)): *)}}
 */
goradd.TagBuilder.prototype = {
    /**
     * attr sets an attribute and returns the tag builder
     * @param a {string} The name of the attribute
     * @param v {string} The value to set the attribute to
     * @returns {goradd.TagBuilder}
     */
    attr: function(a, v) {
        this.el.setAttribute(a,v);
        return this;
    },
    /**
     * html sets the innerHTML to the given value.
     * @param h {string}
     * @returns {goradd.TagBuilder}
     */
    html: function(h) {
        this.el.innerHTML = h;
        return this;
    },
    /**
     * text sets the innterText to the given value.
     * @param t {string}
     * @returns {goradd.TagBuilder}
     */
    text: function(t) {
        this.el.innerText = t;
        return this;
    },
    /**
     * appendTo ends the builder by inserting the tag into the dom as the last child element of the given element.
     * @param el {object|string}
     */
    appendTo: function(el) {
        el = goradd.el(el);
        el.appendChild(this.el);
        return this.el;
    },
    /**
     * insertInto ends the builder by inserting the tag into the dom as the first child element of the given element.
     * @param el {object|string}
     */
    insertInto: function(el) {
        el = goradd.el(el);
        el.insertChild(this.el);
        return this.el;
    },
    /**
     * insertBefore ends the builder by inserting the tag into the dom as a sibling of the given item, and just before it.
     * @param el {object|string}
     */
    insertBefore: function(el) {
        el = goradd.el(el);
        el.parentElement.insertBefore(this.el, el);
        return this.el;
    },
    /**
     * insertAfter ends the builder by inserting the tag into the dom as a sibling of the given item, and just after it.
     * @param el {object|string}
     */
    insertAfter: function(el) {
        el = goradd.el(el);
        el.insertAdjacentElement("afterend", this.el);
        return this.el;
    },
    /**
     * replace ends the builder by replacing the given element.
     * @param el {object|string}
     */
    replace: function(el) {
        el = goradd.el(el);
        el.parentElement.replaceChild(this.el, el);
        return this.el;
    },
    /**
     * Wrap ends the builder by moving the given tag inside of the builder's tag, and
     * then replacing the tag
     * @param el
     */
    wrap: function(el) {
        el = goradd.el(el);
        el.parentElement.replaceChild(this.el, el);
        this.el.appendChild(el);
        return this.el;
    }
};

/***
 * The goradd widget wrapper contains a series of operations that can be performed on an html object.
 */

/**
 * g Wraps an html object in a goradd widget and returns the widget, so you can call functions on it.
 * It also attaches itself to the object so it doesn't need to recreate itself each time.
 * @param el
 */
goradd.g = function(el) {
    el = goradd.el(el);
    if (!el) {
        return undefined;
    }
    if (el.goradd && el.goradd.widget) {
        // Element has an attached goradd widget, so use it. It is either this object or an extension of this object.
        return el.goradd.widget;
    }
    if (!el.goradd) {
        el.goradd = {};
    }
    if (!this._g) {
        // first time through, allow it to be called without a new
        return new goradd.g(el);
    }
    // this is the actual constructor
    this.element = el;
    el.goradd.widget = this;
};


goradd.g.prototype = {
    _g: 1, // just a marker to help with the constructor
    id: function() {
        return this.element.id;
    },
    get: function(key) {
        var v;
        return (v = this.data(key)) !== undefined ? v :
            (v = this.option(key)) !== undefined ? v:
            (v = this.prop(key))  !== undefined ? v :
                undefined;
    },
    prop: function(key, v) {
        var self = this;
        if (arguments.length === 1) {
            if (typeof key === "object") {
                // setting group of keys and values
                goradd.each(key, function(k,v) {
                    self.element[k] = v;
                });
                return
            }
            return this.element[key];
        } else if (arguments.length === 2) {
            this.element[key] = v;
        }
    },
    option: function(key) {
        return this._options[key];
    },
    qs: function(sel) {
        return this.element.querySelector(sel);
    },
    /**
     * qa is a querySelectorAll call that returns an actual array of HTML elements, and not a NodeList.
     * By returning an array, you can call ES5 array functions on it, like forEach.
     * Returns empty array if selector has no results.
     * @param sel {string} The css selector to find
     * @returns {HTMLElement[]}
     */
    qa: function(sel) {
        return Array.prototype.slice.call(this.element.querySelectorAll(sel));
    },
    /**
     * matches returns true if the given element matches the css selector.
     * @param sel
     * @returns {boolean}
     */
    matches: function(sel) {
        if (Element.prototype.matches) {
            return this.element.matches(sel);
        } else {
            var matches = document.querySelectorAll(sel),
                i = matches.length;
            while (--i >= 0 && matches.item(i) !== this.element) {}
            return i > -1;
        }
    },
    /**
     * parents returns the parent nodes, not including the window.
     * @returns {Array}
     */
    parents: function() {
        var a = [];
        var el = this.element;
        while (el.parentElement && el.parentElement !== window) {
            a.push(el.parentElement);
            el = el.parentElement;
        }
        return a;
    },
    /**
     * closest returns the first parent node that matches the given selector, or null
     * @param sel
     */
    closest: function(sel) {
        var el = this.element;
        while (el.parentElement && el.parentElement !== window) {
            if (g$(el).matches(sel)) {
                return el;
            }
            el = el.parentElement;
        }
        return null;
    },
    /**
     * attr gets and sets attributes on a dom object. Remember that attributes are not the same as properties, but can be related.
     * To access properties, use prop. These specifically access the attributes defined in html, but not anything set
     * afterwards.
     *
     * Returns undefined if the attribute does not exist. If no arguments are given, returns an object with all the
     * elements attributes.
     *
     * Note that when returning all attributes, attributes set as an empty string will be returned as a "true".
     * Html has no way of differentiating between an attribute that is an empty string, and an attribute that is
     * set with no value at all, which is common for boolean attributes. Since setting an attribute with an empty
     * string is unusual, and setting an attribute with no value to represent true is more common, we will return
     * the boolean value. If you need an empty string, you will need to look for the boolean value and switch it.
     *
     * @param k [string] The attribute name.
     * @param v [string] The attribute value to set. If no value is given, it just returns the value indicated by k.
     *                   If you pass undefined, null, or false, the attribute will be removed. Passing true here will
     *                   set the attribute with a blank value, which in html indicates a value of true.
     * @returns {null|boolean|*}
     */
    attr: function(k,v) {
        var t = this.element;
        if (arguments.length === 0) {
            // Return an object mapping all the attributes of the html object
            if (t.hasAttributes()) {
                var attr = {};
                // Apparently IE has a quirk where it returns all possible attributes, and not just set attributes.
                goradd.each(this.element.attributes, function(v,n) {
                    n = n.nodeName || n.name;
                    if (t.hasAttribute(n)) {
                        var v2 = t.getAttribute(n);
                        if (v2 === "") { // empty string. Is it a true, or really an empty string? The world may never know.
                            v2 = true;
                        }
                        attr[n] = v2;
                    }
                });
                return attr;
            }
            return undefined; // no attributes are set
        } else if (arguments.length === 1) {
            // get value
            if (!t.hasAttribute(k)) {
                return undefined;
            }
            v = t.getAttribute(k);
            if (v === "true" || v === "") {
                return true; // A boolean attribute, it just exists with no value or with "true"
            } else if (v === "false") {
                return false;
            } else {
                return v;
            }
        } else {
            if (v === undefined || v === null || v === false) {
                t.removeAttribute(k);
                return;
            }
            if (v === true) {
                v = ""; // per the standard for boolean attributes
            }
            t.setAttribute(k,v);
        }
    },
    /**
     * class returns the value of the class, or sets the class, and returns the new class.
     * Prefix the class with a "+" to add the class(es). Prefix with "-" to remove the given classes.
     * Separate class names with a space.
     * @param c
     */
    class: function(c) {
        var el = this.element;
        if (arguments.length === 0) {
            return el.className || el.class;
        }
        if (c.substr(0,1) === "+") {
            // Support: Opera Mini does not support multiple classes, so we do them one at a time
            goradd.each(c.substr(1).split(" "), function(i,v) {
                if (v !== "") {
                    el.classList.add(v);
                }
            });
        } else if (c.substr(0,1) === "-") {
            // Support: Opera Mini does not support multiple classes, so we do them one at a time
            goradd.each(c.substr(1).split(" "), function (i, v) {
                if (v !== "") {
                    el.classList.remove(v);
                }
            });
        } else {
            el.className = c;
        }
        return el.className || el.class;
    },
    /**
     * Returns true if the give class is on the element.
     * @param c {string} class to check for
     * @returns {boolean}
     */
    hasClass: function(c) {
        return this.element.classList.contains(c);
    },
    /**
     * Toggles the given classes. Returns the final class list.
     * @param c
     * @returns {string}
     */
    toggleClass: function(c) {
        var el = this.element;
        goradd.each(c.split(" "), function(i,v) {
            if (v !== "") {
                el.classList.toggle(v);
            }
        });
        return el.className || el.class;
    },
    /**
     * css sets or gets the given css property.
     * @param p {string} Property to set or get
     * @param v [string] Optional value. If ommitted, no setting will happen
     * @returns {string} The value of the property.
     */
    css: function(p, v) {
        var el = this.element;
        if (arguments.length >= 2) {
            el.style[p] = v;
            return v;
        }
        var styles = window.getComputedStyle(el); // TODO: since this is live, should we stash this in the object so we don't have the overhead?
        if (styles.hasOwnProperty(p)) {
            return styles[p];
        }
        if (el.style && el.style.hasOwnProperty(p)) {
            return el.style[p];
        }
        return undefined;
    },

    /**
     * on attaches an event handler to the given html object.
     * Filtering and potentially supplying data to the event are also included.
     *
     * If data is a function, the function will be called when the event fires and the
     * result of the function will be provided as data to the event.
     *
     * The "this" parameter provided to the handler will be the wrapper object that you are attaching the handler to ...
     * essentially the same as the "this" for the on function when you call it.
     *
     * In the event returned to the handler, "target" is the element receiving the event, and "currentTarget" is the element
     * listening for the event.
     *
     * If using a selector,
     *
     * @param {string} eventNames  One or more event names separated by spaces
     * @param {string} [selector] An optional css selector to filter bubbled events. This is here because jQuery does it this way too.
     * @param {function} handler
     * @param {object} [options] Optional additional options as follows:
     *      selector: {string} Same as selector above, just specified in options
     *      bubbles: {boolean} When used with a selector, determines if selector filters parent elements (true), or just
     *        the target. If true, and the event passes the filter, the attached goradd.match object will be the element
     *        that is the first matching selector that the event encountered as it bubbled.
     *      capture: {boolean} Whether to fire event during the capture phase. See addEventListener doc for how this works
     *      data: {*} Data to provide into the goradd.data item attached to the event. If this is a function, the function
     *        will be executed when the event fires, and the result provided to the event. The "this" of the function
     *        will be the "this" of the on call, unless of course you bind a different "this".
     */
    on: function(eventNames, selector, handler, options) {
        // TODO: This code breaks the built-in addEventListener ability to prevent multiple adds of the same handler.
        // However, that code only works when the handler is not anonymous.
        //  We could potentially add code here that would prevent this as well if needed.
        // We could put a "handleEvent" function on ourselves, and then make that the handler. We would then need to do
        // our own management of the attached handlers. We could implement a mechanism where the handler provides a
        // unique id, and so we can prevent multiple adds of the same anonymous function too.
        var self = this;
        if (!eventNames) {
            goradd.log("on must specify an event");
            return;
        }
        // Sort out the arguments
        if (typeof selector !== "string") {
            options = handler;
            handler = selector;
            selector = undefined;
        }
        if (typeof handler !== "function") {
            goradd.log("on must have a handler that is a function");
            return;
        }

        var capture = false;
        var target = self;
        if (!!options) {
            if (typeof options !== "object") {
                goradd.log("options must be an object if it is defined");
                return;
            }
            if (options.capture) {
                capture = true;
            }
            if (!!options.selector) {
                selector = options.selector;
            }
            if (!!options.handlerTarget) {
                target = options.handlerTarget;
            }
        }
        var el = this.element;
        var events = eventNames.split(" ");
        goradd.each(events, function(i,eventName) {
            el.addEventListener(eventName, function (event) {
                goradd.log("triggered: " + event.type);
                if (!!selector) {
                    if (!!options && options.bubbles) {
                        var check = event.target;
                        var match;
                        if (g$(check).matches(selector)) {
                            match = check;
                        }
                        while (!match && !!check && check !== event.currentTarget) {
                            check = check.parentElement;
                            if (g$(check).matches(selector)) {
                                match = check;
                            }
                        }
                        if (match) {
                            if (!event.goradd) {
                                event.goradd = {};
                            }
                            event.goradd.match = match;
                            event.goradd.selector = selector;
                        } else {
                            return;
                        }
                    } else {
                        if (!g$(event.target).matches(selector)) {
                            return;
                        }
                        if (!event.goradd) {
                            event.goradd = {};
                        }
                        event.goradd.selector = selector;
                        event.goradd.match = event.target;
                    }
                }
                var data;
                if (options && options.data) {
                    data = options.data;
                    if (typeof options.data === "function") {
                        data = options.data.call(self, event);
                    }
                    if (!event.goradd) {
                        event.goradd = {};
                    }
                    event.goradd.data = data;
                }
                if (event.detail) {
                    data = event.detail;
                }

                if (data) {
                    handler.call(target, event, data); // add extra item to event handler
                } else {
                    handler.call(target, event);
                }
            }, capture);
        });
    },
    click: function(postFunc) {
        var event;
        // Include extra information as part of the click.
        if (typeof window.Event === "object") {
            goradd.log ("init MouseEvent");
            // Event for browsers which don't natively support the Constructor method
            event = document.createEvent('MouseEvent');
            event.initEvent("click", true, true);
            if (postFunc) {
                event.goradd = {postFunc: postFunc};
            }
        } else {
            goradd.log("new MouseEvent");
            event = new MouseEvent("click", {bubbles: true, cancelable: true, composed: true});
            if (postFunc) {
                event.goradd = {postFunc: postFunc};
            }
        }
        this.element.dispatchEvent(event);
    },
    trigger: function(eventName, extra) {
        var el = this.element;
        var event;

        if (eventName === "click") {
            el.click();
        } else if (eventName === "change") {
            if (typeof window.Event === "object") {
                // Event for browsers which don't natively support the Constructor method
                event = document.createEvent('HTMLEvents');
                event.initEvent(eventName, true, true);
            } else {
                event = new Event(eventName, {bubbles: true})
            }
            // Note that extra information is not supported for the change event. If needed, we can add it
            // in a special area on the event, like in grDetail, and then unpack that in the on handler.
        } else {
            // assume custom event
            if (typeof window.CustomEvent === "object") {
                // CustomEvent for browsers which don't natively support the Constructor method
                event = document.createEvent('CustomEvent');
                event.initCustomEvent(eventName, true, true, extra);
            } else {
                event = new CustomEvent(eventName, {bubbles: true, cancelable: true, composed: true, detail: extra})
            }
        }
        el.dispatchEvent(event);
    },
    /**
     * htmlAfter adds the html after the given element.
     * @param html
     */
    htmlAfter: function(html) {
        this.element.insertAdjacentHTML("afterend", html);
    },
    /**
     * htmlBefore inserts the html before the given element.
     * @param html
     */
    htmlBefore: function(html) {
        this.element.insertAdjacentHTML("beforebegin", html);
    },
    /**
     * insertHtml inserts the given html in the inner html of the given element, but before any other html that is
     * already there.
     * @param html
     */
    insertHtml: function(html) {
        this.element.insertAdjacentHTML("afterbegin", html);
    },
    /**
     * appendHtml inserts the given html into the inner html of the given element, but after any other html that is
     * already there.
     * @param html
     */
    appendHtml: function(html) {
        this.element.insertAdjacentHTML("beforeend", html);
    },
    /**
     * Remove removes the given element from the dom. It returns the removed element.
     * @returns {*}
     */
    remove: function() {
        var el = this.element;
        el.parentElement.removeChild(el);
        return el;
    },
    /**
     * Value sets or gets the value of a goradd control. This is primarily used by the ajax processing code, but
     * external tools can use this too. See below for what each kind of control will return. Note that the actual "value"
     * attribute is not always returned.
     * @param v
     * @returns {*}
     */
    val: function(v) {
        var el = this.element;
        var type = g$(el).prop("type");
        if (arguments.length === 1) {
            // Setting the value
            switch (type) {
                case "select-multiple":
                    // Multi-select selections will attempt to set all items in the given array to the value
                    var opts = g$(el).qa('option');
                    goradd.each(opts, function(i, opt) {
                        opt.selected = v.includes(opt.value);
                    });
                    break;
                case "checkbox":
                    if (typeof v === "boolean") {
                        el.checked = v;
                    } else if (typeof v === "number") {
                        el.checked = v !== 0;
                    } else if ("value" in el) {
                        el.checked = el.value === v;
                    } else {
                        el.checked = false;
                    }
                    break;

                case "radio":
                    if (typeof v === "boolean") {
                        el.checked = v;
                    } else {
                        el.checked = el.value == v;
                    }
                    break;
                default:
                    //if ("value" in el) {
                        el.value = v;
                    //}
                    break;
            }
            return el;
        } else {
            switch (type) {
                case "select-multiple":
                    // Multi-select selections will return an array of selected values
                    var sels = g$(el).qa(':checked');
                    return sels.map(function(s){return s.value});
                case "checkbox":
                case "radio":
                    // Checkboxes and radios will return the value, or true, if checked, and null if not checked.
                    if (el.checked) {
                        if (!("value" in el)) { // if the checkbox has no value, just return true;
                            return true;
                        } else {
                            return el.value;
                        }
                    }
                    break;
                default:
                    if ("value" in el) {
                        // This works for textboxes, textarea (possible problem losing newlines though), and single selects.
                        // Custom controls can add a "value" getter or override the val() mehtod.
                        return el.value;
                    }
                    break;
            }

            return null;
        }
    },
    /**
     * data gets or sets custom data that we assign to an element. If getting the data, we will check our private area
     * first for the data, and then check for an attribute if we have not overridden the attribute with private data.
     * Getting data attached as a "data-*" attribute uses the camelCase version of the name.
     * Private data is stored in the "goradd.data" object attached to the element.
     * @param key
     * @param v
     * @returns {*}
     */
    data: function(key, v) {
        var el = this.element;
        if (arguments.length === 1) {
            // Get the data
            if (el.goradd.data && el.goradd.data.hasOwnProperty(key)) {
                return el.goradd.data[key]; // Use our private data area if its there
            }
            // Otherwise try to get the data from the attribute on the element
            if (el.dataset) { // modern browsers
                // use the key as is
                return el.dataset[key];
            } else {
                // IE 10 or opera mini. Gotta get this from the attribute itself.
                key = _toKebab(key);
                return el.getAttribute("data-" + key);
            }
        } else {
            // We are setting data. We do not alter the attribute (use goradd.attr() if you need that). Instead,
            // we put the data in our private area for later collection.
            if (!el.goradd.data) {
                el.goradd.data = {};
            }
            el.goradd.data[key] = v;
        }
    },
    focus: function() {
        this.element.focus();
    },
    blur: function() {
        this.element.blur();
    },
    selectAll: function() {
        this.element.select();
        // Note, setSelectionRange, etc. appears to NOT be supported in opera mini.
    },
    text: function(t) {
        if (arguments.length === 0) {
            return this.element.innerText;
        } else {
            this.element.innerText = t;
        }
    },
    html: function(t) {
        if (arguments.length === 0) {
            return this.element.innerHTML;
        } else {
            this.element.innerHTML = t;
        }
    },

    /**
     * f calls the named function, with the given parameters, on the goradd widget.
     * @param name
     * @param params
     */
    f: function(name, params) {
        var f = this[name];
        if (typeof f === "function") {
            return f.apply(this, params);
        }
    },
    /**
     * clientWidth returns the bounding width of the client area, which includes content, padding, border, but not margin.
     * @returns {number}
     */
    clientWidth: function() {
        return this.element.clientWidth;
    },
    /**
     * clientHeight returns the bounding height of the client area, which includes content, padding, border, but not margin.
     * @returns {number}
     */
    clientHeight: function() {
        return this.element.clientHeight;
    },

};

/**
 * This is a recreation of the jQuery UI widget factory, with fewer features and specifically supporting IE 10+
 * and Opera Mini.
 *
 * It takes the given prototype, makes it an extension of the base object, and then puts it at the given named
 * spot under the window object. The name can be separated with dots to work down the hierarchy. Start the name
 * with "goradd." to add it to the goradd hierarchy.
 *
 * Note that this name means two things. First, that the prototype will be placed at that location off the goradd global
 * hierarchy, and that the actual object created will be placed at the location off of the goradd object attached
 * to the html object.
 *
 * @param name  The namespaced name of the prototype.
 * @param base  The base object. If not included, goradd.Widget will be used as the base object.
 * @param prototype The prototype to use. Functions will become part of the function prototype, and other objects will
 *                  become static global objects. Instance methods should be placed in the "options" object, or
 *                  simply declared and initialized in the "_create" function.
 */
goradd.widget = function(name, base, prototype) {
    // Use goradd.Widget if there is no base
    if ( !prototype ) {
        prototype = base;
        base = goradd.Widget;
    }

    // make sure we put the prototype on the goradd global object, and the instance on the goradd item attached to the html object.
    var names = name.split( "." );
    if (names[0] !== "goradd") {
        names.unshift("goradd");
    }

    if (names.length === 1) {
        goradd.log("You cannot create a widget at 'goradd'");
        return;
    }

    if (names[0] === "goradd" && names[1] === "data") {
        goradd.log("goradd.data is a reserved location");
        return;
    }

    var obj = window;
    var ctx = null;

    for (var i = 0; i < names.length - 1; i++) {
        var v = names[i];
        ctx = obj;
        if (!obj[v]) {
            obj[v] = {};
        }
        obj = obj[v];
    }
    var loc = names[names.length -1];
    if (obj[loc]) {
        goradd.log(name + " is already defined.");
        return;
    }

    var constructor = obj[loc] = function(options, element) {
        if (this._createWidget) {
            this._createWidget(options, element);
        }
    };

    var basePrototype = new base();
    // Copy the options object
    basePrototype.options = goradd.extend( {}, basePrototype.options );

    var proxiedPrototype = {};
    goradd.each( prototype, function( prop, value ) {
        if (typeof value !== "function" ||
            !base.prototype[ prop ]) { // only create override if there is a base function
            proxiedPrototype[ prop ] = value;
            return;
        }
        proxiedPrototype[ prop ] = ( function() {
            function _super() {
                return base.prototype[ prop ].apply( this, arguments );
            }

            function _superApply( args ) {
                return base.prototype[ prop ].apply( this, args );
            }

            return function() {
                var __super = this._super;
                var __superApply = this._superApply;
                var returnValue;

                this._super = _super;
                this._superApply = _superApply;

                returnValue = value.apply( this, arguments );

                this._super = __super;
                this._superApply = __superApply;

                return returnValue;
            };
        } )();
    } );

    var namespace = names.slice(0, names.length - 2).join("."),
        widgetName = names[names.length - 1],
        widgetFullName = names.join(".");

    constructor.prototype = goradd.extend( basePrototype, proxiedPrototype, {
        constructor: constructor,
        namespace: namespace,
        widgetName: widgetName,
        widgetFullName: widgetFullName
    } );

};

/**
 * widget.new creates and initializes a new widget with the given named constructor.
 * @param constructor
 * @param options
 * @param element
 * @returns {*}
 */
goradd.widget.new = function(constructor, options, element) {
    if (typeof constructor === "string") {
        var names = constructor.split( "." );
        var obj = window;
        var ctx = null;
        goradd.each (names, function (i, v) {
            ctx = obj;
            obj = obj[v];
        });
        constructor = obj;
    }
    return new constructor(options, element);
};


/**
 * This is the definition of the Widget class, which serves as the base class for other widgets. It itself is based
 * on the "g" class, which is a jQuery like wrapper. In other words, all the functions on the g class are available
 * to widgets through the "this" variable, and can be overridden. One important function to override might be the
 * "val" function, which provides the value that will be used by ajax calls. If your widget only works through Ajax,
 * then that is sufficient to keep the go side of things updated.
 */
goradd.widget("goradd.Widget", goradd.g, {
    /**
     * _createWidget acts as the constructor of all widgets. It can be overridden by the widget if needed, but
     * you normally do not need to. Instead, implement _create() to make a private constructor for your widget,
     * which this default constructor will call.
     * @param options
     * @param element
     * @private
     */
    _createWidget: function(options, element) {
        var self = this;
        this.element = goradd.el(element);

        // Merge options into default options

        this.options = goradd.extend( {}, this.options); // copy defaults
        goradd.each(options, function(k,v) {
            if (self.options.hasOwnProperty(k) &&
                typeof(self.options[k]) === "string" &&
                v === true) {
                // deal with special situation where we are trying to set an option to a blank string, instead of a boolean
                // html cannot differentiate between the two when the option is coming from an attribute
                self.options[k] = "";
            } else {
                self.options[k] = v;
            }
        });

        if (this.element) { // if no element, this may be created just to get to its prototype
            this._create();
            this.trigger("create");
            this._init();
        }
    },
    /**
     * _create is the constructor of each individual widget. Call this._super() to call the superclass's constructor too.
     * @private
     */
    _create: function() {
    },
    _init: function() {
    },
    option: function(key, value) {
        this._setOption(key,value);
    },
    _setOption: function(key, value) {
        this.options[key] = value;
    }
});

})(  );


/**
 * Ajax Queue
 *
 * This used to be handled with a jquery plugin, but since we are trying to get away from jquery, and working
 * towards an OperaMini compatible version, we are rolling our own.
 */
(function() {
    var _q = [],
        _currentRequests= {},
    _idCounter= 0;

        goradd.ajaxq = {


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
            _q.push(f);
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
        return _currentRequests.length === 0;
    },
    _dequeue: function() {
        var f = _q.shift();
        if (f) {
            this._do1(f);
        }
    },
    _do1(f) {
        var self = this;
        var opts = f();
        _idCounter++;
        var ajaxID = _idCounter;

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

                delete _currentRequests[ajaxID];
                if (_q.length === 0 && !self.isRunning()) {
                    g$(goradd.form()).trigger("ajaxQueueComplete");
                }
                self._dequeue(); // do the next ajax event in the queue
            }
        };
        _currentRequests[ajaxID] = objRequest;
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
})();


////////////////////////////////
// Goradd Shortcuts and Initialize
////////////////////////////////

goradd.initialize();
