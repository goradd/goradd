/**
 * goradd.js
 *
 * This is the shim between the goradd GO code and the browser. It enables ajax and other kinds of
 * communication between the server and the client.
 *
 * Goals:
 *  - Compatible with all current browsers, IE 10+ Opera Mobile (ES5).
 *  - Provide utility code to javascript widgets and plugins.
 *  - Do not use jQuery or other frameworks, but be compatible if external frameworks are used by the developer.
 */

if (!function () {
    "use strict";
    return Function.prototype.bind && XMLHttpRequest && !this;
}()) {
    var body = document.getElementsByTagName("body")[0];
    body.innerText = "This browser is not supported. Try a different browser.";
}

/**
 * goradd holds the global namespace for all the goradd support javascript.
 */
var goradd;
/**
 * g$ is a shortcut for goradd.g(). It wraps an element with additional functions defined here to more easily manipulate
 * the element. el can be either an actual HTMLElement, or the id of one.
 * One main difference between jQuery's wrapper and this one is that jQuery wraps an array of elements, and we only
 * wrap one element. Also, functions that return an array of elements will return HTMLElements, not an array of wrapped elements.
 * @param {string | Element | goradd.Control} el
 * @returns {goradd.Control}
 */
var g$ = function(el) {
    return goradd.g(el);
};


(function () {
    // Put everything in a function so we can have private functions and variables.
    // By convention, things starting with underscore will be private.
    "use strict";

    /**
     * Private functions and members used by code below.
     */

    var _controlValues = {};
    var _formObjsModified = {};
    var _ajaxError = false;
    var _inputSupport = true;
    var _finalCommands = [];
    var _watchers = {};
    var _refresh = [];
    var _registeredWidgets = {};

    /*
    function _toKebab(s) {
        return s.replace(/[A-Z]/g, function (m) {
            return "-" + m.toLowerCase();
        });
    }
    */

    function _fromKebab(s) {
        return s.replace(/-([a-z])/gi, function (o, i) {
            return i.toUpperCase();
        });
    }

    /**
     * formObjChanged is an event handler that records that a control has changed in order to synchronize the control with
     * the server on the next request. Send the formObjChanged event to the control
     * that changed, and it will bubble up to the form and be caught here.
     * @param {object} event
     */
    function _formObjChanged(event) {
        _formObjsModified[event.target.id] = true;
    }

    /**
     * Gets the data to be sent to an ajax call as post data. This will be called from the ajax queueing function, and
     * will erase the cache of changed objects to prepare for the next call.
     *
     * @param {object} params       An object containing the following:
     * @param  {string} params.controlId   The control id to post an action to
     * @param  {number} params.eventId     The event id
     * @param  {boolean} params.async      If true, process the event asynchronously without waiting for other events to complete
     * @param  {string} params.formId      The id of the form getting posted
     * @param  {object} [params.values]    An optional object, that contains values to send with the event
     * @param  {*}   [params.values.event]   The event's action value, if one is provided. This can be any type, including an object.
     * @param  {*}   [params.values.action]  The action's action value, if one is provided. Any type.
     * @param  {*}   [params.values.control] The control's action value, if one is provided. Any type.
     * @return {FormData} Post Data
     * @private
     */
    function _getAjaxData(params) {
        var form = goradd.form(),
            controls = g$(form).qa("input,select,textarea"),
            postData = {};

        // Notify controls we are about to post.
        g$(form).trigger("posting", "Ajax");

        goradd.each(controls, function (i, c) {
            var id = c.id;
            var blnForm = (id && (id.substring(0, 8) === "Goradd__"));
            var blnPost = (g$(c).data("grPost") === "");

            if (!_inputSupport || // if not oninput support, then post all the controls, rather than just the modified ones, because we might have missed something
                _ajaxError || // Ajax error would mean that _formObjsModified is invalid. We need to submit everything.
                (id && _formObjsModified[id]) ||  // We try to ignore controls that have not changed to reduce the amount of data sent in an ajax post.
                blnPost || // was the control marked to always post,
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
        //params.callType = "Ajax";
        postData.Goradd__Params = _getParamsValue(params);

        _ajaxError = false;
        _formObjsModified = {};
        _controlValues = {};
        var fd = new FormData();
        for (var k in postData) {
            var obj = postData[k];
            if (obj instanceof FileList) {
                for (var i = 0; i < obj.length; i++) {
                    goradd.log("Ajax: added file to " + k);
                    fd.append(k, obj[i]);
                }
            } else {
                goradd.log("Ajax: added " + postData[k] + " to " + k);
                fd.append(k, postData[k]);
            }
        }
        return fd;
    }

    /**
     * Returns the Goradd__Params value in an encoded form that will be decoded by the server.
     * Encoding is necessary to accommodate Web site monitors that may interpret the clear text as attempts at
     * subterfuge.
     *
     * @param params
     * @returns {string}
     * @private
     */
    function _getParamsValue(params) {
        params.tz = {"o" : -(new Date()).getTimezoneOffset(), "z" : Intl ? Intl.DateTimeFormat().resolvedOptions().timeZone : null};

        if (!goradd.isEmptyObject(_controlValues)) {
            params.controlValues = _controlValues;
        }
        var s = JSON.stringify(params);
        var b = new TextEncoder().encode(s);
        var bs = String.fromCodePoint(...b);
        return btoa(bs);
    }

    /**
     * Responds to the part of an ajax response that must be handled serially before other handlers can fire.
     *
     * @param {object} json     json generated by goradd application
     * @private
     */
    function _processAjax(json) {
        if (json.alert) {
            goradd.each(json.alert, function (i, v) {
                window.alert(v);
            });
        }

        // Delay until all files are loaded
        var loadCount = 1;
        var fWait = function() {
            loadCount--;
            if (loadCount === 0) {
                _processAjax2(json);
            }
        }

        if (json.ss) {
            goradd.each(json.ss, function (i, v) {
                loadCount ++;
                goradd.loadStyleSheetFile(i, v, fWait);
            });
        }

        if (json.js) {
            goradd.each(json.js, function (i, v) {
                loadCount++;
                goradd.loadJavaScriptFile(i, v, fWait);
            });
        }

        fWait(); // if there are no files being loaded, just go straight to the next step
    }

    /**
     * Process the part of an ajax response that may depend on earlier processing.
     *
     * @param {object} json  Json generated by the goradd application.
     * @private
     */
    function _processAjax2(json) {
        goradd.each(json.controls, function (id) {
            var el = goradd.el(id),
                $ctrl = g$(id),
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
                    var relatedParent = p.filter(function (item) {
                        return g$(item).matches(relSelector);
                    }).pop();

                    if (relatedParent) {
                        relatedParent.insertAdjacentElement("beforebegin", el);
                    }
                    if (relatedItems && relatedItems.length > 0) {
                        goradd.each(relatedItems, function () {
                            g$(this).remove();
                        });
                    }
                    $ctrl.htmlBefore(this.html);
                    $ctrl.remove();
                } else {
                    // control is being injected at the top level, so put it at the end of the form.
                    g$(goradd.form()).appendChildHtml(this.html);
                }
            }
        });

        _registerControls();

        goradd.each(json.commands, function (i, command) {
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
     * @param {object} command
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
        while (_finalCommands.length) {
            var command = _finalCommands.pop();
            goradd.processCommand(command);
        }
    }

    /**
     * Convert from JSON return value to an actual jQuery object. Certain structures don't work in JSON, like closures,
     * but can be part of a javascript object. We use special codes to piece together functions, closures, dates, etc.
     * @param {*[]} params
     * @returns {*}
     * @private
     */
    function _unpackArray(params) {
        if (!params) {
            return null;
        }
        var newParams = [];

        goradd.each(params, function (index, item) {
            if (Array.isArray(item)) {
                item = _unpackArray(item);
            } else if (typeof item === 'object' && item !== null) {
                if (item.goraddObject) {
                    item = _unpackObj(item);  // top level special object
                } else {
                    // look for special objects inside top level objects.
                    var newItem = {};
                    goradd.each(item, function (key, obj) {
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
     * @param {object} obj
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
                            goradd.each(obj.params, function (i, v) {
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
                        goradd.each(vars, function (i, v) {
                            val = val[v];
                        });
                        return val;

                    case 'func':
                        // Returns the result of the given function called immediately
                        // Find the function and context starting at the window context.
                        var target = window;
                        if (obj.context) {
                            var objects = obj.context.split(".");
                            goradd.each(objects, function (i, v) {
                                target = target[v];
                            });
                        }

                        if (obj.params) {
                            params = [];
                            goradd.each(obj.params, function (i, v) {
                                params.push(_unpackObj(v)); // recurse
                            });
                        }
                        var func = target[obj.func];

                        return func.apply(target, params);
                }
            } else {
                var newItem = {};
                goradd.each(obj, function (key, obj2) {
                    newItem[key] = _unpackObj(obj2);
                });
                return newItem;
            }
        }

        return obj; // no change
    }

    function _registerControls() {
        var els = goradd.qa('[data-grctlhigh]'); // High priority controls that other controls are dependent on
        goradd.each(els, function () {
            _registerControl(this);
        });
        els = goradd.qa('[data-grctl]');
        goradd.each(els, function () {
            _registerControl(this);
        });
    }

    function _registerControl(ctrl) {
        if (!ctrl) {
            return;
        }

        // get the widget
        var g = g$(ctrl);

        var watches = g.attr("data-gr-watch");
        if (!!watches) {
            goradd.each(watches.split(";"), function () {
                var s = this.split("=");

                _addWatcher(g.id, s[0])
            });
        }

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
        if (window.jQuery) {
            // Some jQuery libraries intercept events, and then fire change using jQuery, which is not interceptable using built in capture methods
            jQuery(ctrl).on("change", _formObjChanged)
        }

        // widget support, using declarative methods
        var widget;
        var options = {};
        goradd.each(g.attr(), function (k, v) {
            if (k === "data-gr-widget") {
                widget = v;
            } else if (k.substring(0, 12) === "data-gr-opt-") {
                options[_fromKebab(k.substring(12))] = v;
            }
        });
        if (widget) {
            g.attachWidget(widget, options);
        }
    }

    function _addWatcher(id, channel) {
        // val is ignored for now. This would be for field watching.
        if (!_watchers[channel]) {
            _watchers[channel] = [id];
        } else if (!goradd.contains(_watchers[channel], id)) {
            _watchers[channel].push(id);
        }
    }

    function _processWatcherMessage(msg) {
        var channel = msg.channel;
        var message = msg.message;

        if (message === "U" && channel === "redraw") {
            // Update the form, no redrawing. Used primarily by the testing framework at this point.
            goradd.updateForm();
            return;
        }
        var watchers = _watchers[channel];
        if (!!watchers) {
            goradd.each(watchers, function () {
                var g = g$(this);
                if (!!g) { // make sure control was not removed from the form
                    _refresh.push(this); // force a refresh of this control
                }
            });

            if (_refresh.length > 0) {
                goradd.updateForm();
            }
        }
    }

    /**
     * The internal goradd event handler object.
     * @param {[string]} events
     * @param {function} handler
     * @param {boolean} capture
     * @private
     */
    function _EventHandler(events, handler, capture) {
        this.events = events;
        this.handler = handler;
        this.capture = capture;
    }

    _EventHandler.prototype = {
        handleEvent: function(evt) {
            this.handler(evt);
        },
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
         * The ajaxTimeout setting controls the number of milliseconds that ajax operations will wait before they
         * timeout. The default here is to have no timeout for debugging purposes. The release version of the app will set this
         * to the config.AjaxTimeout value in GO.
         */
        ajaxTimeout: 0,
        /**
         * blockEvents is set by event handlers to cause that event to block all other incoming events,
         * until the event is processed.
         */
        blockEvents: false,
		/**
         * toArray converts array-like objects to arrays
         * @param a
         * @return {Array}
         */
        toArray: function(a) {
            return Array.prototype.slice.call(a);
        },

        /**
         * Extend merges keys and values of objects into the target object.
         * This version of extend is primarily for the purpose of adding
         * capabilities to javascript classes. It does not do deep merging, but it will copy the members of plain objects
         * if it finds a plain object. Other kinds of objects are copied by reference.
         * If only one argument is provided, the target is the goradd object itself.
         * Matching keys will overwrite the target key.
         *
         * @param {...object} target
         * @returns {*}
         */
        extend: function (target) {
            var input = goradd.toArray(arguments, 1),
                key,
                value;

            if (arguments.length === 1) {
                input = [target];
                target = goradd;
            }

            var inputIndex = 0,
                inputLength = input.length;

            for (; inputIndex < inputLength; inputIndex++) { // iterate through all arguments in order
                var obj = input[inputIndex];
                for (key in obj) { // iterate through the keys in each argument
                    value = obj[key];
                    if (obj.hasOwnProperty(key) && value !== undefined) {

                        // Clone plain objects
                        if (goradd.isPlainObject(value)) {
                            target[key] = goradd.isPlainObject(target[key]) ?
                                goradd.extend({}, target[key], value) :

                                // Don't extend strings, arrays, etc. with objects
                                goradd.extend({}, value);
                        } else { // Copy everything else by reference
                            target[key] = value;
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
         * @param {string|HTMLElement|goradd.Control} t
         * @returns {HTMLElement}
         */
        el: function (t) {
            if (!t) {
                return null;
            }
            if (typeof t === "object") {
                if (!!t.element && !!t.element.goradd) {
                    return t.element; // Its a goradd Control, so get the element from the widget
                }
                return t;
            }
                return document.getElementById(t);
        },
        /**
         * qs is a shortcut for document.querySelector
         * @param {string} sel
         * @returns {*}
         */
        qs: function (sel) {
            return document.querySelector(sel);
        },
        /**
         * qa is a querySelectorAll call that returns an actual array, and not a NodeList.
         * Returns empty array if selector has no results.
         * @param {string} sel The css selector to find
         * @returns {HTMLElement[]}
         */
        qa: function (sel) {
            return goradd.toArray(document.querySelectorAll(sel));
        },
        /**
         * isEmptyObject will test whether the object is empty
         * @param {object} o
         * @returns {boolean}
         */
        isEmptyObject: function (o) {
            if (!o) return false;
            for (var name in o) { // this is a js trick
                return false;
            }
            return true;
        },
        /**
         * form returns the current form on the page
         * @returns {HTMLElement}
         */
        form: function () {
            // TODO: cache this, it will not change. No reason to do this over and over.
            return goradd.qs('form[data-grctl="form"]');
        },
        /**
         * matches returns true if the given element matches the css selector.
         * @param {string|object|HTMLElement}el
         * @param {string} sel
         * @returns {boolean}
         */
        matches: function (el, sel) {
            return g$(el).matches(sel);
        },
        /**
         * loadJavaScriptFile will dynamically load a javascript file. It is designed to be called during ajax calls or
         * other times when a dynamically loaded javascript file is required.
         * @param {string} strScript
         * @param {object} attributes
         * @param {function} fDone The function to run after the file is loaded
         */
        loadJavaScriptFile: function (strScript, attributes, fDone) {
            var script = document.createElement("script");
            script.src = strScript;
            script.type = 'text/javascript';
            script.onload = fDone;
            script.onerror = fDone;
            goradd.extend(script, attributes);

            var head = document.getElementsByTagName('head')[0];
            head.appendChild(script);
        },
        /**
         * loadStyleSheetFile dynamically loads a style sheet file. It is used by the ajax code.
         * @param {string} strStyleSheetFile
         * @param {object} attributes
         * @param {function} fDone The function to run after the file is loaded
         */
        loadStyleSheetFile: function (strStyleSheetFile, attributes, fDone) {
            var link = document.createElement("link");
            link.rel = "stylesheet";
            link.href = strStyleSheetFile;
            link.onload = fDone;
            link.onerror = fDone;
            goradd.extend(link, attributes);
            var head = document.getElementsByTagName('head')[0];
            head.appendChild(link);
        },
        /**
         * each is a recreation of the jQuery each function, but for our targeted browsers only. It iterates the given object,
         * calling the function for each item found. If the object is an array, or something array-like, like a nodelist,
         * it will pass the index and the item to the function. For a regular object, it will pass the key and the item.
         * "this" is set to the item each time as well.
         * @param {object} obj
         * @param {function} f
         */
        each: function (obj, f) {
            if (!obj) return;
            if (typeof (obj) !== "object") return;
            var i;

            // isArrayLike needs to return true for nodelists.
            var isArrayLike = Array.isArray(obj) || ("length" in obj && typeof (obj.length) === "number" && (obj.length === 0 || ((obj.length - 1) in obj)));
            if (isArrayLike) {
                if (obj.length === 0) return;
                for (i = 0; i < obj.length; i++) {
                    if (f.call(obj[i], i, obj[i]) === false) break;
                }
            } else {
                for (i in obj) {
                    if (f.call(obj[i], i, obj[i]) === false) {
                        break;
                    }
                }
            }
        },
        /**
         * contains returns true if needle is in the array a
         * @param {ArrayLike} a
         * @param needle {*}
         * @returns {boolean}
         */
        contains: function (a, needle) {
            return (a.indexOf(needle) !== -1);
        },
        /**
         * isPlainObject will return true if the given object is a POJO.
         * @param {object} obj
         * @returns {boolean}
         */
        isPlainObject: function (obj) {
            var proto, Ctor;

            // Detect obvious negatives
            // Use toString instead of jQuery.type to catch host objects
            if (!obj || {}.toString.call(obj) !== "[object Object]") {
                return false;
            }

            proto = Object.getPrototypeOf(obj);

            // Objects with no prototype (e.g., `Object.create( null )`) are plain
            if (!proto) {
                return true;
            }

            // Objects with prototype are plain iff they were constructed by a global Object function
            Ctor = {}.hasOwnProperty.call(proto, "constructor") && proto.constructor;
            return typeof Ctor === "function" && {}.hasOwnProperty.toString.call(Ctor) === {}.hasOwnProperty.toString.call(Object);
        },
        /**
         * setRadioInGroup is a specialized function called from goradd go code.
         * It sets the given radio button to being checked in a group. Since the goradd code already knows what we want to
         * check, that button is just set. However, the button that gets unset by the browser needs to communicate to the
         * go code that it is getting unset.
         * @param {string} id
         */
        setRadioInGroup: function (id) {
            var el = goradd.el(id);
            if (el.type !== "radio") {
                return; // not a radio button, or not part of a group
            }
            var prevItem;

            if (el.name) {
                prevItem = goradd.qs("input[type='radio'][name='" + el.name + "']:checked");
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
         * @param {object} o
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
        setControlValue: function (strControlId, strProperty, strNewValue) {
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
            var form = goradd.form();
            /*
    $(window).on ("storage", function (o) {
        if (o.originalEvent.key === "goradd.broadcast") {
            goradd.updateForm();
        }
    });*/

            _inputSupport = 'oninput' in document;
            // IE 9 has a major bug in oninput, but we are requiring IE 10+, so no problem.
            // I think the only major browser that does not support oninput is Opera mobile.

            g$(form).on("ajaxQueueComplete", _processFinalCommands);

            // TODO: Add a detector of the back button. This detector should ping the server to make sure the pagestate exists on the server. If not,
            // it should reload the form.

            g$(form).on('formObjChanged', _formObjChanged); // Allow any control, including hidden inputs, to trigger a change and post of its data.
            g$(form).on('submit', function (event) {
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
         * @param {object} params                   An object containing the following:
         * @param {string} params.controlId         The control id to post an action to
         * @param {number} params.eventId           The event id
         * @param {boolean} [params.async = false]  If true, process the event asynchronously without waiting for other events to complete
         * @param {object} [params.values]          An optional object, that contains values coming to send with the event
         * @param {*}    [params.values.event]      The event's action value, if one is provided. This can be any type, including an object.
         * @param {*}    [params.values.action]     The action's action value, if one is provided. Any type.
         * @param {*}    [params.values.control]    The control's action value, if one is provided. Any type.
         *
         * @return {void}
         */
        postBack: function (params) {
            if (goradd.blockEvents) {
                return;  // We are waiting for a response from the server
            }

            var form = goradd.form();
            var gForm = g$(form);

            //params.callType = "Server";

            // Notify custom controls that we are about to post
            gForm.trigger("posting", "Server");

            goradd.el('Goradd__Params').value = _getParamsValue(params);

            // trigger our own form submission so we can catch it
            gForm.trigger("submit");
        },

        /**
         * Posts an ajax call to the ajax queue. Ajax actions call this.
         *
         * @param {object} params                   An object containing the following:
         * @param {string} params.controlId         The control id to post an action to
         * @param {number} params.eventId           The event id
         * @param {boolean} [params.async = false]  If true, process the event asynchronously without waiting for other events to complete
         * @param {object} [params.actionValues]          An optional object, that contains values coming to send with the event
         * @param {*}    [params.actionValues.event]      The event's action value, if one is provided. This can be any type, including an object.
         * @param {*}    [params.actionValues.action]     The action's action value, if one is provided. Any type.
         * @param {*}    [params.actionValues.control]    The control's action value, if one is provided. Any type.
         *
         * @return {void}
         */
        postAjax: function (params) {
            var form = goradd.form(),
                formAction = g$(form).attr("action"),
                async = params.hasOwnProperty("async");

            if (goradd.blockEvents) {
                return;
            }

            //params.formId = form.id;
            params.refresh = _refresh;

            goradd.log("postAjax", params);

            // Use an ajax queue so ajax requests happen synchronously
            goradd.ajaxq.enqueue(function () {
                var data = _getAjaxData(params);

                return {
                    url: formAction,
                    data: data,
                    /**
                     * @param result {string}
                     * @param err {number}
                     * @returns {boolean}
                     */
                    error: function (result, err) {
                        goradd.displayAjaxError(result, err);
                        goradd.testStep();
                        return false;
                    },
                    /**
                     * @param json {object}
                     */
                    success: function (json) {
                        goradd.log("Ajax response received: ", json);
                        _processAjax(json);
                        goradd.blockEvents = false;
                    }
                };
            }, async);
        },
        /**
         * Process a single command. This is called both from ajax and javascript.
         * @param {object} command
         */
        processCommand: function (command) {
            var params,
                objs;

            if (command.script) {
                // TODO: clean this up a bit by using ids for inserted scripts. Might have multiple scripts for the same id though.
                var script = document.createElement("script");
                script.type = "text/javascript";
                script.text = command.script;
                document.body.appendChild(script);
            } else if (command.selector) {
                params = _unpackArray(command.params);

                if (typeof command.selector === 'string') {
                    // general selector
                    objs = goradd.qa(command.selector);
                } else {
                    // First item is the id to select on
                    objs = g$(command.selector[0]).qa(command.selector[1]);
                }

                goradd.each(objs, function (i, v) {
                    var $c = g$(v);
                    if (typeof $c[command.func] === "function") {
                        $c[command.func].apply($c, params);
                    }
                });
            } else if (command.func) {
                params = _unpackArray(command.params);

                // Find the function by name. Walk an object list in the process.
                objs = command.func.split(".");
                var obj = window;
                if (command.id) {
                    obj = g$(command.id);
                    if (!obj) {
                        goradd.log("object for command not found - " + command.id);
                        return;
                    }
                } else if (command.jqueryId) {
                    obj = jQuery('#' + command.jqueryId);
                    if (!obj || obj.length === 0) {
                        goradd.log("object for command not found - " + command.jqueryId);
                        return;
                    }
                }
                var ctx = null;

                goradd.each(objs, function (i, v) {
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
        /**
         * updateForm will cause the form to update without a specific action. Useful if you know that javascript control
         * values have changed and you want to post the changes back to the server and get a response.
         */
        updateForm: function () {
            goradd.throttle('goradd.update', 500, function () {
                goradd.postAjax({});
            });
        },
        /**
         * registerWidget is called by widget definitions to enable attaching a widget to an html element by declaring
         * the widget in the data-gr-widget attribute.
         * @param {string} name
         * @param {Object} classObj
         */
        registerWidget: function(name, classObj) {
            _registeredWidgets[name] = classObj;
        },

        /******************************************
         * Testing support
         * See goradd-debug.js
         ******************************************/

        testStep: function (event) {},
        log: function () {},
        getPageState: function () {
            return document.getElementById("Goradd__PageState").value;
        },
        currentStep:0,
        displayAjaxError: function(){},

        /***********************
         * Utility Functions
         ***********************/

        /**
         * Sets a cookie with the given parameters. Potentially called by the goradd app.
         * @param {string} name
         * @param {*} val
         * @param {Date} expires
         * @param {string} path
         * @param {string} dom
         * @param {boolean} secure
         */
        setCookie: function (name, val, expires, path, dom, secure) {
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
        result in a change event, and under certain circumstances this could cause the click event to be dropped.
        If the change event moves the focus away from the button that was clicked, the click event will not record.
        We therefore delay the processing of all events to try to queue them up before processing.
        Its very strange. Something to debug at a future date.
        */

        _actionQueue: [],
        queueAction: function (params) {
            if (params.last) {
                var delay = 0;

                goradd._actionQueue.forEach(function (item) {
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
        processActions: function () {
            while (goradd._actionQueue.length > 0) {
                var params = goradd._actionQueue.shift();
                goradd.log("processAction: " + params.name + " delay: " + params.d);
                if (params.d > 0) {
                    // If the event is repeated during the delay, the action will not be repeated
                    goradd.setTimer(params.k, params.f, params.d);
                } else {
                    params.f();
                }
            }
        },
        //////////////////////////////
        // Goradd Action Support
        //////////////////////////////
        /**
         * These support the various Goradd function actions available in the action package.
         */

        /**
         * msg puts up a javascript alert
         * @param {string} m
         */
        msg: function (m) {
            alert(m);
        },
        /**
         * redirect will load a different page without posting the form
         * @param {string} newLocation
         */
        redirect: function (newLocation) {
            window.location = newLocation
        },
        /**
         * refresh will cause the given goradd control to refresh
         * @param {string} id
         */
        refresh: function (id) {
            _refresh.push(id);
            goradd.updateForm();
        },

        /**
         * proxyVal returns a value for the Proxy control
         * @param {Event} event
         * @returns {*}
         */
        proxyVal: function (event) {
            var target = event.target;
            if (!!event.goradd && !!event.goradd.match) {
                target = event.goradd.match;
            }
            return g$(target).data("grAv");
        },

        // Watcher support
        subscribeWatchers: function () {
            goradd.subscribe(Object.keys(_watchers), _processWatcherMessage)
        },
        /**
         * findNamedObject will search the through the base hierarchy for object named and return it. The hierarchy
         * is delimited by periods (.). If no base is given, it uses the window as its starting point. If the value is
         * not found in the object, undefined is returned.
         * @param {string} name
         * @param {object} [base]
         */
        findNamedObject: function (name, base) {
            if (name === "") {
                return base;
            }
            var names = name.split(".");

            if (!base) {
                base = window;
            }
            var obj = base;
            for (var i = 0; i < names.length; i++) {
                var v = names[i];
                if (obj[v] === undefined) {
                    return undefined;
                }
                obj = obj[v];
            }
            return obj;
        },
        /**
         * textToHtml replaces items that should be escaped for html display with their html counterparts.
         * @param {string} text
         * @returns {string}
         */
        textToHtml: function (text) {
            text = text.replace("&", "&amp;");
            text = text.replace("<", "&lt;");
            text = text.replace(">", "&gt;");
            text = text.replace('"', "&quot;");
            text = text.replace("'", "&apos;");
            return text;
        }
    };

    /**
     * Named timer support. These allow you to create timers without having to keep a copy of the timer around.
     */
    var _timers = {};

    goradd.extend({
        /**
         * Sets the named timer, given an action and a delay.
         * @param {string} id           Any id string you want to give to identify the timer. Note that "goradd.*" strings are reserved.
         * @param {function} action     The function to execute when the timre fires
         * @param {number} intDelay     The delay in millis before firing the timer, or how often the timer will fire if a periodic timer
         * @param {boolean} [blnPeriodic=false] True for a repeating timer.
         */
        setTimer: function (id, action, intDelay, blnPeriodic) {
            blnPeriodic = !!blnPeriodic;
            goradd.clearTimer(id);
            goradd.log("setTimer", id, intDelay, blnPeriodic);
            _timers[id] = {
                i: new Date().getTime(), // record when the timer was inserted
                d: intDelay, // delay, for info if we need to recover this
                p: !!blnPeriodic, // is this periodic?
                s: 0,
                e: 0,
                t: blnPeriodic ?
                    setInterval(function () {
                        action(_timers[id]);
                    }, intDelay) :
                    setTimeout(function () {
                        _timers[id].s = new Date().getTime(); // record when the function started
                        goradd.clearTimer(id);
                        action(_timers[id]);
                        _timers[id].e = new Date().getTime(); // record when the function ended
                    }, intDelay)
            };
        },
        /**
         * hasTimer returns true if there is an active timer at the give id
         * @param {string} id
         * @returns {boolean}
         */
        hasTimer: function (id) {
            return !!_timers[id] &&
                _timers[id].t != null; // could have been a previous timer that is now expired.
        },
        /**
         * Clears the named timer.
         * @param {string} id
         */
        clearTimer: function (id) {
            goradd.log("clearTimer", id);
            if (goradd.hasTimer(id)) {
                if (_timers[id].p) {
                    clearInterval(_timers[id].t);
                } else {
                    clearTimeout(_timers[id].t);
                }
                _timers[id].t = null; // keep track of timer history.
                // Allows us to know when the last time a timer was fired.
            }
        },
        /**
         * Creates a timer on a control that fires the timerexpiredevent.
         *
         * @param {string} controlID
         * @param {number} delay
         * @param {boolean} periodic
         */
        startControlTimer: function (controlID, delay, periodic) {
            var strTimerId = controlID + '_ct';
            var action = function () {
                g$(controlID).trigger('timerexpiredevent');
            };
            goradd.setTimer(strTimerId, action, delay, periodic);
        },
        /**
         * Stops the control's timer.
         * @param {string} controlID
         */
        stopControlTimer: function (controlID) {
            var strTimerId = controlID + '_ct';
            goradd.clearTimer(strTimerId);
        },
        /*
         * Use throttle to limit the execution of a function that you might execute many times quickly. If the function is
         * already scheduled to execute, it will do nothing, letting the previous call to throttle finish its execution. If
         * the function has not been scheduled since minInterval, it will execute immediately. Otherwise, it will delay
         * the execution until minInterval has passed since the last time the function executed.
         *
         * @param {string} id  - A timer id you assign
         * @param {number} minInterval - The minimum interval between executions
         * @param {function} f - The function to execute
         * @returns {void}
         */
        throttle: function (id, minInterval, f) {
            if (goradd.hasTimer(id)) {
                return; // the function is already scheduled
            }

            var prevTime = 0;

            if (!!_timers[id]) {
                if (_timers[id].e > 0) {
                    prevTime = _timers[id].s;
                }
            }
            var curTime = new Date().getTime();
            var delay = Math.max(minInterval - (curTime - prevTime), 0);
            goradd.setTimer(id, f, delay);
        }
    });


/////////////////////////////////
// Tag Builder
/////////////////////////////////

    /**
     * tagBuilder returns a TagBuilder. Use it as follows:
     * tag = goradd.tagBuilder("div").attr("class", "myClass").text("I am text").appendTo("objId");
     * @type goradd.TagBuilder
     * @param tag {string}
     * @returns {goradd.TagBuilder}
     */
    goradd.tagBuilder = function (tag) {
        return new goradd.TagBuilder(tag);
    };
    /**
     * TagBuilder uses a builder pattern to create and place html tags.
     *
     * @param {string} tag
     * @constructor
     */
    goradd.TagBuilder = function (tag) {
        this.el = document.createElement(tag);
    };
    /**
     *
     * @type {{appendTo: (function((Object|string)): *), insertInto: (function((Object|string)): *), replace: (function((Object|string)): *), html: (function(string): goradd.TagBuilder), text: (function(string): goradd.TagBuilder), attr: (function(string, string): goradd.TagBuilder), insertAfter: (function((Object|string)): *), insertBefore: (function((Object|string)): *)}}
     */
    goradd.TagBuilder.prototype = {
        /**
         * attr sets an attribute and returns the tag builder
         * @param {string} a  The name of the attribute
         * @param {string} v  The value to set the attribute to
         * @returns {goradd.TagBuilder}
         */
        attr: function (a, v) {
            this.el.setAttribute(a, v);
            return this;
        },
        /**
         * css sets the given css inline style property to the given value
         * @param {string} p
         * @param {string} v
         * @returns {goradd.TagBuilder}
         */
        css: function (p, v) {
            this.el.style[p] = v;
            return this;
        },
        /**
         * html sets the innerHTML to the given value.
         * @param {string} h
         * @returns {goradd.TagBuilder}
         */
        html: function (h) {
            this.el.innerHTML = h;
            return this;
        },
        /**
         * text sets the innerText to the given value.
         * @param {string} t
         * @returns {goradd.TagBuilder}
         */
        text: function (t) {
            this.el.innerText = t;
            return this;
        },
        /**
         * id sets the id of the element being built
         * @param {string} i
         * @return {goradd.TagBuilder}
         */
        id: function (i) {
            this.attr("id", i);
            return this;
        },
        /**
         * class sets the class of the element being built
         * @param {string} c
         * @return {goradd.TagBuilder}
         */
        class: function(c) {
            this.attr("class", c);
            return this;
        },

        /**
         * appendTo ends the builder by inserting the tag into the dom as the last child element of the given element.
         * @param {object|string} el
         * @return {HTMLElement}
         */
        appendTo: function (el) {
            el = goradd.el(el);
            el.appendChild(this.el);
            return this.el;
        },
        /**
         * insertInto ends the builder by inserting the tag into the dom as the first child element of the given element.
         * @param {object|string} el
         * @return {HTMLElement}
         */
        insertInto: function (el) {
            el = goradd.el(el);
            el.insertAdjacentElement("afterbegin", this.el);
            return this.el;
        },
        /**
         * insertBefore ends the builder by inserting the tag into the dom as a sibling of the given item, and just before it.
         * @param {object|string} el
         * @return {HTMLElement}
         */
        insertBefore: function (el) {
            el = goradd.el(el);
            el.parentElement.insertBefore(this.el, el);
            return this.el;
        },
        /**
         * insertAfter ends the builder by inserting the tag into the dom as a sibling of the given item, and just after it.
         * @param {object|string} el
         * @return {HTMLElement}
         */
        insertAfter: function (el) {
            el = goradd.el(el);
            el.insertAdjacentElement("afterend", this.el);
            return this.el;
        },
        /**
         * replace ends the builder by replacing the given element.
         * @param {object|string} el
         * @return {HTMLElement}
         */
        replace: function (el) {
            el = goradd.el(el);
            el.parentElement.replaceChild(this.el, el);
            return this.el;
        },
        /**
         * Wrap ends the builder by moving the given tag inside of the builder's tag, and
         * then replacing the tag
         * @param {object|string} el
         * @return {HTMLElement}
         */
        wrap: function (el) {
            el = goradd.el(el);
            el.parentElement.replaceChild(this.el, el);
            this.el.appendChild(el);
            return this.el;
        },
        /**
         * element ends the builder by simply returning the element without inserting it into the DOM.
         * @returns {HTMLElement}
         */
        element: function() {
            return this.el;
        },
        /**
         * asHtml returns the created tag as html
         * @return {string}
         */
        asHtml: function() {
            return this.el.outerHTML;
        }
    };

    /**
     * Goradd Control object
     */

    /**
     * g wraps an html object in a goradd Control and returns the Control object, so you can call functions on it.
     * It is sort of like the $() command in jQuery.
     * It also attaches itself to the object so it doesn't need to recreate itself each time.
     * You can use the shortcut g$(id|element).func() to use these functions, passing either the id of an element,
     * or the element itself.
     * @param {string|Element|goradd.Control} el
     * @returns {goradd.Control}
     */
    goradd.g = function (el) {
        if (!el) {
            return undefined;
        }
        if (el.element && el.element.goradd) {
            return el; // this is already a goradd object
        }
        el = goradd.el(el);
        if (!el) {
            return undefined;
        }
        if (el.goradd && el.goradd.widget) {
            // Element has an attached goradd widget, so use it. It is either this object or an extension of this object.
            return el.goradd.widget;
        }
        return new goradd.Control(el);
    };

    /**
     * goradd.Control is a basic goradd html object that has extended commands.
     * @class
     * @param {HTMLElement} el
     * @constructor
     */
    goradd.Control = function (el) {
        if (!el) {
            // Being called by our widget factory just to make a copy of the prototype
            return
        }
        if (!el.goradd) {
            el.goradd = {};
        }
        this.element = el;
        el.goradd.widget = this;
        if (el.id) {
            this.id = el.id; // hoist id to the goradd object
        }
    };

    goradd.Control.prototype = {
        /**
         * get returns the named parameter. It will look first in the data, then in the options, and then in the element properties.
         * @param {string} key
         * @returns {*}
         */
        get: function (key) {
            var v;
            return (v = this.data(key)) !== undefined ? v :
                (v = this.option(key)) !== undefined ? v :
                    (v = this.prop(key)) !== undefined ? v :
                        undefined;
        },
        /**
         * Prop gets or sets the named property on the element.
         * @param {string} key
         * @param {*} [v]
         * @returns {*}
         */
        prop: function (key, v) {
            var self = this;
            if (arguments.length === 1) {
                if (typeof key === "object") {
                    // setting group of keys and values
                    goradd.each(key, function (k, v) {
                        self.element[k] = v;
                    });
                    return
                }
                return this.element[key];
            } else if (arguments.length === 2) {
                this.element[key] = v;
            }
        },
        /**
         * options returns the named option value
         * @param {string} key
         * @returns {*}
         */
        option: function (key) {
            return this._options[key];
        },
        /**
         * qs returns the results of querySelector on the element.
         * @param {string} sel
         * @returns {HTMLElement}
         */
        qs: function (sel) {
            return this.element.querySelector(sel);
        },
        /**
         * Find performs a querySelector and returns the result as a goradd.Control object.
         * @param {string} sel
         * @returns {goradd.Control}
         */
        find: function (sel) {
            return g$(this.qs(sel))
        },
        /**
         * qa is a querySelectorAll call that returns an actual array of HTML elements, and not a NodeList.
         * By returning an array, you can call ES5 array functions on it, like forEach.
         * Returns empty array if selector has no results.
         * @param {string} sel  The css selector to find
         * @returns {HTMLElement[]}
         */
        qa: function (sel) {
            return goradd.toArray(this.element.querySelectorAll(sel));
        },
        /**
         * findAll is like querySelectorAll, but returns an array of goradd.Control objects
         * @param {string} sel
         * @returns {goradd.Control[]}
         */
        findAll: function (sel) {
            return this.qa(sel).map(function(o) {return g$(o)});
        },
        /**
         * matches returns true if the given element matches the css selector.
         * @param sel
         * @returns {boolean}
         */
        matches: function (sel) {
            if (Element.prototype.matches) {
                return this.element.matches(sel);
            } else {
                var matches = document.querySelectorAll(sel),
                    i = matches.length;
                while (--i >= 0 && matches.item(i) !== this.element) {
                }
                return i > -1;
            }
        },
        /**
         * parents returns the parent nodes, not including the window.
         * @returns {HTMLElement[]}
         */
        parents: function () {
            var a = [];
            var el = this.element;
            while (el.parentElement && el.parentElement !== window) {
                a.push(el.parentElement);
                el = el.parentElement;
            }
            return a;
        },
        /**
         * children returns the child elements of the control as an array of goradd.Controls
         * @return {goradd.Control[]}
         */
        children: function() {
            return goradd.toArray(this.element.children).map(function(o) {g$(o)});
        },
        /**
         * closest returns the first parent node that matches the given selector, or null, as an goradd object
         * @param sel
         * @returns goradd.Control
         */
        closest: function (sel) {
            var el = this.element;
            while (el.parentElement && el.parentElement !== window) {
                if (g$(el).matches(sel)) {
                    return g$(el);
                }
                el = el.parentElement;
            }
            return null;
        },
        /**
         * firstChild returns the first child element of the object, or nil
         * @return {goradd.Control}
         */
        firstChild: function() {
            return g$(this.element.firstElementChild);
        },
        /**
         * prev returns the html item previous sibling to the current item, or null if no item is previous.
         * @return {goradd.Control}
         */
        prev: function() {
            return g$(this.element.previousElementSibling);
        },
        /**
         * next returns the html item that is the next sibling to the current item, or null if no item is next.
         * @return {goradd.Control}
         */
        next: function() {
            return g$(this.element.nextElementSibling);
        },
        /**
         * clone will clone the element and return a new goradd Control for the element.
         * The element will not be attached to the DOM.
         * @return {goradd.Control}
         */
        clone: function() {
            return g$(this.element.cloneNode(true));
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
         * @param {string} k   The attribute name.
         * @param {*} [v] The attribute value to set. If no value is given, it just returns the value indicated by k.
         *                     If you pass undefined, null, or false, the attribute will be removed. Passing true here will
         *                     set the attribute with a blank value, which in html indicates a value of true.
         * @returns {null|boolean|*}
         */
        attr: function (k, v) {
            var t = this.element;
            if (arguments.length === 0) {
                // Return an object mapping all the attributes of the html object
                if (t.hasAttributes()) {
                    var attr = {};
                    // Apparently IE has a quirk where it returns all possible attributes, and not just set attributes.
                    goradd.each(this.element.attributes, function (v, n) {
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
                t.setAttribute(k, v);
            }
        },
        /**
         * class returns the value of the class, or sets the class, and returns the new class.
         * Prefix the class with a "+" to add the class(es). Prefix with "-" to remove the given classes.
         * Separate class names with a space.
         * @param {string} [c]
         */
        class: function (c) {
            var el = this.element;
            if (arguments.length === 0) {
                return el.className || el.class;
            }
            if (c.substring(0, 1) === "+") {
                // Support: Opera Mini does not support multiple classes, so we do them one at a time
                goradd.each(c.substring(1).split(" "), function (i, v) {
                    if (v !== "") {
                        el.classList.add(v);
                    }
                });
            } else if (c.substring(0, 1) === "-") {
                // Support: Opera Mini does not support multiple classes, so we do them one at a time
                goradd.each(c.substring(1).split(" "), function (i, v) {
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
         * @param {string} c  class to check for
         * @returns {boolean}
         */
        hasClass: function (c) {
            return this.element.classList.contains(c);
        },
        /**
         * Toggles the given classes. Returns the final class list.
         * @param {string} c
         * @returns {string}
         */
        toggleClass: function (c) {
            var el = this.element;
            goradd.each(c.split(" "), function (i, v) {
                if (v !== "") {
                    el.classList.toggle(v);
                }
            });
            return el.className || el.class;
        },
        /**
         * css sets or gets the given css property.
         * @param {string} p Property to set or get
         * @param {string} [v] Optional value. If omitted, no setting will happen
         * @returns {string} The previous value of the property.
         */
        css: function (p, v) {
            var el = this.element;
            var oldV = undefined;
            var styles = window.getComputedStyle(el); // TODO: since this is live, should we stash this in the object so we don't have the overhead?
            if (styles.hasOwnProperty(p)) {
                oldV = styles[p];
            } else if (el.style && el.style.hasOwnProperty(p)) {
                oldV = el.style[p];
            }

            if (arguments.length >= 2) {
                if (v === "" || v === null || v === undefined) {
                    el.style.removeProperty(p);
                } else {
                    el.style.setProperty(p, v);
                }
            }

            return oldV;
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
         * @param {string} eventNames  One or more event names separated by spaces
         * @param {string} [selector] An optional css selector to filter bubbled events. This is here because jQuery does it this way too.
         * @param {function|Array|EventListener} handler The function to execute. If handler is an array, the first item
         *        will become the "this" of the function, and 2nd item is the function itself.
         * @param {object} [options] Optional additional options as follows:
         * @param {string} [options.selector]  Same as selector above, just specified in options
         * @param {boolean} [options.targetElement] True to make the "this" be the element rather than the goradd object.
         * @param {boolean} [options.bubbles] When used with a selector, determines if selector filters parent elements (true), or just
         *        the target. If true, and the event passes the filter, the attached goradd.match object will be the element
         *        that is the first matching selector that the event encountered as it bubbled.
         * @param {boolean}  [options.capture] Whether to fire event during the capture phase. See addEventListener doc for how this works
         * @param {*} [options.data]  Data to provide into the goradd.data item attached to the event. If this is a function, the function
         *        will be executed when the event fires, and the result provided to the event. The "this" of the function
         *        will be the "this" of the on call, unless of course you bind a different "this".
         * @return _EventHandler An object you can use when calling off to remove the event.
         */
        on: function (eventNames, selector, handler, options) {
            // TODO: This code breaks the built-in addEventListener ability to prevent multiple adds of the same handler.
            // However, that code only works when the handler is not anonymous.
            //  We could potentially add code here that would prevent this as well if needed.
            // We could put a "handleEvent" function on ourselves, and then make that the handler. We would then need to do
            // our own management of the attached handlers. We could implement a mechanism where the handler provides a
            // unique id, and so we can prevent multiple adds of the same anonymous function too.
            // also need to think about off.
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
            if (!(handler instanceof Object ) && typeof handler !== "function" && !(Array.isArray(handler) && handler.length === 2)) {
                goradd.log("'on' must have a handler that is an EventListener, function or be a 2 item array");
                return;
            }

            var capture = false;
            var target = self;
            var el = this.element;
            if (options) {
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
                if (options.targetElement) {
                    target = el;
                }
            }

            if (Array.isArray(handler)) {
                target = handler[0];
                handler = handler[1];
            } else if (handler.handleEvent) {
                target = handler;
                handler = target.handleEvent;
            }

            var events = eventNames.split(" ");
            var objEvt = new _EventHandler (events, function (event) {
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

                // This data here is getting set up when "on" is first called.
                var data;
                if (options && options.data !== undefined) {
                    data = options.data;

                    // Calls a specified function when the event is fired to get the value of data
                    if (typeof options.data === "function") {
                        data = options.data.call(self, event);
                    }
                    if (!event.goradd) {
                        event.goradd = {};
                    }
                    event.goradd.data = data; // in case it gets overridden below, we can still get to the data through the event
                }

                // This is data sent through the trigger function at trigger time
                if (event.detail !== undefined) {
                    data = event.detail;
                }

                if (data !== undefined) {
                    handler.call(target, event, data); // add extra item to event handler
                } else {
                    handler.call(target, event);
                }
            }, capture);
            goradd.each(events, function (i, eventName) {
                el.addEventListener(eventName, objEvt, capture);
            });
            return objEvt; // pass this to off
        },
        /**
         * off removes an event that was added by on.
         * @param {_EventHandler} eventObj
         */
        off: function(eventObj) {
            var el = this.element;
            goradd.each(eventObj.events, function() {
                el.removeEventListener(this, eventObj, eventObj.capture)
            });
        },
        /**
         * click will send a click event to the object
         * @param {function} postFunc a function to execute after the click event has fired.
         */
        click: function (postFunc) {
            var event;
            // Include extra information as part of the click.
            if (typeof window.Event === "object") {
                goradd.log("init MouseEvent");
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
        /**
         * trigger
         * @param eventName
         * @param extra
         */
        trigger: function (eventName, extra) {
            var el = this.element;
            var event;

            if (eventName === "click") {
                el.click();
            } else if (eventName === "change") {
                if (typeof window.Event !== "function") {
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
                if (typeof window.CustomEvent !== "function") {
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
         * moveAfter moves the goradd object to be after the given element.
         * @param {string|HTMLElement|goradd.Control} el
         */
        moveAfter: function(el) {
            el = goradd.el(el);
            el.insertAdjacentElement("afterend", this.element);
        },
        /**
         * moveBefore moves the element to be before of the given element.
         * @param {string|HTMLElement|goradd.Control} el
         */
        moveBefore: function(el) {
            el = goradd.el(el);
            el.insertAdjacentElement("beforebegin", this.element);
        },
        /**
         * insertInto puts the element inside the given element before all other child elements.
         * @param {string|HTMLElement|goradd.Control} el
         */
        insertInto: function(el) {
            el = goradd.el(el);
            el.insertAdjacentElement("afterbegin", this.element);
        },
        /**
         * appendTo puts the element inside the given element after all other child elements.
         * @param {string|HTMLElement|goradd.Control} el
         */
        appendTo: function(el) {
            el = goradd.el(el);
            el.insertAdjacentElement("beforeend", this.element);
        },
        /**
         * htmlAfter adds the html after the given element.
         * @param {string} html
         */
        htmlAfter: function (html) {
            this.element.insertAdjacentHTML("afterend", html);
        },
        /**
         * htmlBefore inserts the html before the given element.
         * @param {string} html
         */
        htmlBefore: function (html) {
            this.element.insertAdjacentHTML("beforebegin", html);
        },
        /**
         * insertChildHtml inserts the given html in the inner html of the given element, but before any other html that is
         * already there.
         * @param {string} html
         */
        insertChildHtml: function (html) {
            this.element.insertAdjacentHTML("afterbegin", html);
        },
        /**
         * appendChildHtml inserts the given html into the inner html of the given element, but after any other html that is
         * already there.
         * @param {string} html
         */
        appendChildHtml: function (html) {
            this.element.insertAdjacentHTML("beforeend", html);
        },
        /**
         * Remove removes the given element from the dom. It returns the removed element.
         * @returns {*}
         */
        remove: function () {
            var el = this.element;
            el.parentElement.removeChild(el);
            return el;
        },
        /**
         * Value sets or gets the value of a goradd control. This is primarily used by the ajax processing code, but
         * external tools can use this too. See below for what each kind of control will return. Note that the actual "value"
         * attribute is not always returned.
         * @param {string} [v]      The value to set the control to. If not included, it will just return the value.
         * @returns {*}
         */
        val: function (v) {
            var el = this.element;
            var type = g$(el).prop("type");
            if (arguments.length === 1) {
                // Setting the value
                switch (type) {
                    case "select-multiple":
                        // Multi-select selections will attempt to set all items in the given array to the value
                        var opts = g$(el).qa('option');
                        goradd.each(opts, function (i, opt) {
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
                        return sels.map(function (s) {
                            return s.value
                        });
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
                            // Custom controls can add a "value" getter or override the val() method.
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
         * To get data attached as a "data-*" attribute, use the camelCase version of the name for the key.
         * Private data is stored in the "goradd.data" object attached to the element.
         * @param {string} key
         * @param {*} [v]
         * @returns {*}
         */
        data: function (key, v) {
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
        /**
         * Focus puts the focus on the element.
         */
        focus: function () {
            this.element.focus();
        },
        /**
         * Blur will remove the focus if it is on the element.
         */
        blur: function () {
            this.element.blur();
        },
        /**
         * SelectAll will select all of the element if it is a textbox type element.
         */
        selectAll: function () {
            this.element.select();
            // Note, setSelectionRange, etc. appears to NOT be supported in opera mini.
        },
        /**
         * Text will get or set the innerText of the element.
         * @param {string} [t]
         * @returns {string}
         */
        text: function (t) {
            if (arguments.length === 0) {
                return this.element.innerText;
            } else {
                this.element.innerText = t;
            }
        },
        /**
         * html gets or sets the innerHTML of the element
         * @param {string} [t]
         * @returns {string}
         */
        html: function (t) {
            if (arguments.length === 0) {
                return this.element.innerHTML;
            } else {
                this.element.innerHTML = t;
            }
        },
        /**
         * hide will hide the given element using display:none, and also remember the previous
         * display value. Restore the state using show.
         */
        hide: function () {
            if (!this.element.goradd.hidden) { // do not hide twice
                var d = this.css("display");
                if (!!d) {
                    this.element.goradd.hidden = d;
                }
                this.css("display", "none");
            }
        },
        /**
         * show will show the element if it was hidden using hide() above. It uses inline styles, so if
         * you want to start an item out hidden, you must give it an INLINE style of display:none.
         */
        show: function () {
            if (!!this.element.goradd.hidden) {
                this.css("display", this.element.goradd.hidden);
            } else {
                this.css("display", null);
            }
        },
        /**
         * f calls the named function, with the given parameters, on the goradd widget.
         * @param {string} name
         * @param params {Object}
         */
        f: function (name, params) {
            var f = this[name];
            if (typeof f === "function") {
                return f.apply(this, params);
            }
        },
        /**
         * width returns the width of the object. The box model defines different boundaries for an
         * object, and in the boundary parameter, you can specify one of content, padding, border or margin
         * to use those various boundaries to measure the width of the object. Note that this is the real
         * width of the object as drawn.
         * @param {('content'|'padding'|'border'|'margin')} [boundary = 'content']
         * @returns {number}
         */
        width: function (boundary) {
            var w = this.element.clientWidth; // starting value includes padding but nothing else
            var styles;
            switch (boundary) {
                case undefined:
                case null:
                case "":
                case "content":
                    styles = window.getComputedStyle(this.element);
                    w -= parseFloat(styles.getPropertyValue('padding-left')) + parseFloat(styles.getPropertyValue('padding-right'));
                    break;
                case "padding":
                    break; // clientHeight includes the border
                case "border":
                    w = this.element.offsetWidth;
                    break;
                case "margin":
                    styles = window.getComputedStyle(this.element);
                    w = this.element.offsetHeight + parseFloat(styles.getPropertyValue('margin-left')) + parseFloat(styles.getPropertyValue('margin-right'));
                    break;
            }
            return w;
        },
        /**
         * height returns the height of the object. The box model defines different boundaries for an
         * object, and in the boundary parameter, you can specify one of content, padding, border or margin
         * to use those various boundaries to measure the height of the object. Note that this is the real
         * height of the object as drawn.
         * @param {('content'|'padding'|'border'|'margin')} [boundary = 'content']
         * @returns {number}
         */
        height: function (boundary) {
            var h = this.element.clientHeight; // starting value includes padding but nothing else
            var styles;
            switch (boundary) {
                case undefined:
                case null:
                case "":
                case "content":
                    styles = window.getComputedStyle(this.element);
                    h -= parseFloat(styles.getPropertyValue('padding-top')) + parseFloat(styles.getPropertyValue('padding-bottom'));
                    break;
                case "padding":
                    break; // clientHeight includes the border
                case "border":
                    h = this.element.offsetHeight;
                    break;
                case "margin":
                    styles = window.getComputedStyle(this.element);
                    h = this.element.offsetHeight + parseFloat(styles.getPropertyValue('margin-top')) + parseFloat(styles.getPropertyValue('margin-bottom'));
                    break;
            }
            return h;
        },
        /**
         * columnId is a convenience method to find the column id that the current object resides in.
         * This is specifically from within a goradd Table object, since it marks the colgroup with the column ids.
         */
        columnId: function () {
            var idx;
            if (this.element.tagName === "TH" || this.element.tagName === "TD") {
                idx = this.element.cellIndex;
            } else {
                idx = this.closest("th,td").element.cellIndex;
            }
            var table = this.closest('table');
            var colgroup = table.find('colgroup');
            var ths = colgroup.qa('col');
            return ths[idx].id;
        },
        /**
         * attachWidget attaches the named widget to the element. You can do this via the data-gr-widget attribute,
         * or directly through javascript.
         * @param {string} constructorName
         * @param {object} options
         */
        attachWidget: function(constructorName, options) {
            var classPtr = _registeredWidgets[constructorName];
            if (classPtr) {
                new classPtr(this.element, options);
            } else {
                console.log("Could not find registered widget " + constructorName);
            }
        }

    };

    /**
     * This is a kind of polyfill to mimic the ES6 class extension capabilities. It supports IE 10+
     * and Opera Mini, and allows a similar type of definition of classes in browsers that do not support ES6.
     *
     * It takes the given prototype, makes it an extension of the base object, and then returns it. The main difference between
     * this and ES6 is that sub classes have the variable this._super available to represent the superclass, rather than super.
     *
     * @param {string|goradd.Widget} base  The base object. If not included, goradd.Widget will be used as the base object. You
     *                                     can use the name of a registered widget as well.
     * @param {object} [prototype] The prototype to use. Functions will become part of the function prototype, and other objects will
     *                  become static global objects. Instance methods should be placed in the "options" object, or
     *                  simply declared and initialized in the "_create" function.
     */
    goradd.extendWidget = function (base, prototype) {
        // Use goradd.Widget if there is no base.
        if (!prototype) {
            prototype = base;
            base = goradd.Widget;
        } else if (typeof (base) == "string") {
            // search for the base
            base = goradd.findNamedObject(base, goradd.classes);
        }

        var basePrototype = goradd.extend({}, base.prototype);

        var proxiedPrototype = {};
        // Create override functions
        goradd.each(prototype, function (prop, value) {
            if (typeof value !== "function" ||
                !base.prototype[prop]) { // only create override if there is a base function
                proxiedPrototype[prop] = value;
                return;
            }
            proxiedPrototype[prop] = (function () {
                // Setup super value before calling the subclass function

                return function () {
                    var __super = this._super;
                    var ret;

                    this._super = base;
                    ret = value.apply(this, arguments);
                    this._super = __super;
                    return ret;
                };
            })();
        });

        // Call the class's constructor when _super() is called, mimicking how ES6 works with super()
        var obj = function () {
            if (proxiedPrototype["constructor"]) {
                proxiedPrototype["constructor"].apply(this, arguments);
            }
        };
        obj.prototype = goradd.extend(basePrototype, proxiedPrototype);
        return obj;
    };

    /**
     * extendOptions is a utility function that constructors should use to define the default values of
     * the options that can be set from the server.
     * @param {object} defaultOptions  The default values
     * @param {object} userOptions     The user defined values
     * @returns {object} Returns the default values overridden by any defined values
     */
    goradd.extendOptions = function (defaultOptions, userOptions) {
        var newOptions = goradd.extend({}, defaultOptions, userOptions);
        goradd.each(userOptions, function (k, v) {
            if (defaultOptions.hasOwnProperty(k) &&
                typeof (defaultOptions[k]) === "string" &&
                v === true) {
                // deal with special situation where we are trying to set an option to a blank string, instead of a boolean
                // html cannot differentiate between the two when the option is coming from an attribute
                newOptions[k] = "";
            } else {
                newOptions[k] = v;
            }
        });
        return newOptions;
    };

    /**
     * This is the definition of the Widget class, which serves as the base class for other widgets. It itself is based
     * on the goradd.Control class, which is a jQuery like wrapper. In other words, all the functions on the Control class are available
     * to widgets through the "this" variable, and can be overridden.
     *
     * You can use the Widget class as a base class for either ES5 or ES6 style widgets. From ES5, so this:
     *
     * var MyWidget = goradd.extendWidget(goradd.Widget, {
     *      constructor: function(element, options) {
     *          var optionDefaults = {
     *              myOpt: 0;
     *              myString: "";
     *          };
     *
     *          this._super(element, goradd.extendOptions(optionDefaults, options))
     *      },
     *     // rest of widget definition
     * }
     * goradd.registerWidget("MyWidget", MyWidget);
     *
     * From ES6, just extend the goradd.Widget class
     *
     * class MyWidget extends goradd.Widget {
     *     constructor(element, options) {
     *          var optionDefaults = {
     *              myOpt: 0;
     *              myString: "";
     *          };
     *
     *         super(element, goradd.extendOptions(optionDefaults, options));
     *     }
     *     // rest of widget definition
     * }
     * goradd.registerWidget("MyWidget", MyWidget);
     *
     * One important function to override might be the "val" function, which provides the value that will be used
     * by ajax calls. If your widget only works through Ajax, then that is sufficient to keep the GO side of things updated.
     *
     * The widget above can be attached to an html object by setting the data-gr-widget attribute of that object to "MyWidget".
     * @class
     * @extends goradd.Control
     */
    goradd.Widget = goradd.extendWidget(goradd.Control, {
        constructor: function (element, options) {
            this._super(element); // call the goradd.Control constructor

            this.options = goradd.extend({}, options);
            this.trigger("create");
        },
        option: function (key, value) {
            if (arguments.length === 1) {
                return this.options[key];
            }
            this._setOption(key, value);
        },
        /**
         * _setOption is the protected function that sets a particular option in a widget. Other setOption type calls
         * funnel to here, so if you need to detect the setting of a particular option, you can override this function.
         * @param key
         * @param value
         * @protected
         */
        _setOption: function (key, value) {
            this.options[key] = value;
        },
        setOptions: function(values) {
            if (arguments.length === 0) {
                return this.options;
            }
            var self = this;
            goradd.each(values, function(k,v) {
                self._setOption(k,v);
            });
        }
    });
})();


/**
 * Ajax Queue
 *
 */
(function () {
    var _q = [],
        _currentRequests = {},
        _idCounter = 0;

    goradd.ajaxq = {


        /**
         * Queues an ajax request.
         * A new Ajax request won't be started until the previous queued
         * request has finished.
         * @param {function} f function that returns ajax options.
         * @param {boolean} blnAsync true to launch right away.
         */
        enqueue: function (f, blnAsync) {
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
        isRunning: function () {
            return Object.keys(_currentRequests).length !== 0;
        },
        _dequeue: function () {
            var f = _q.shift();
            if (f) {
                this._do1(f);
            }
        },
        _do1: function (f) {
            var self = this;
            var opts = f();
            _idCounter++;
            var ajaxID = _idCounter;

            var objRequest = new XMLHttpRequest();

            objRequest.open("POST", opts.url, true);
            objRequest.setRequestHeader("Method", "POST " + opts.url + " HTTP/1.1");
            objRequest.setRequestHeader("X-Requested-With", "xmlhttprequest");
            objRequest.timeout = goradd.ajaxTimeout;

            objRequest.onreadystatechange = function () {
                if (objRequest.readyState === 4) {
                    if (objRequest.status === 200) {
                        try {
                            if (objRequest.response.charAt(0) !== "{") {
                                opts.error(objRequest.response, "Json response is not an object.");
                            } else {
                                opts.success(JSON.parse(objRequest.response));
                            }
                        } catch (err) {
                            // Goradd returns ajax errors as text
                            opts.error(objRequest.response, err);
                        }
                    } else {
                        // This would be a problem with the server or client
                        opts.error("An ajax error occurred: " + objRequest.statusText, objRequest.status);
                    }

                    delete _currentRequests[ajaxID];

                    if (_q.length === 0 && !self.isRunning()) {
                        g$(goradd.form()).trigger("ajaxQueueComplete");
                    }
                    self._dequeue(); // do the next ajax event in the queue
                }
            };

            objRequest.onloadend = objRequest.onprogress = function (event) {
                if (event.lengthComputable) {
                    // pass this on to the form to be intercepted there
                    var obj = {
                        loaded: event.loaded,
                        total: event.total
                    };
                    g$(goradd.form()).trigger("progress", obj);
                    goradd.log("progress: " + event.loaded + " out of " + event.total);
                }
            };

            objRequest.onerror = function (event) {
                goradd.log("error");
                opts.error("An error occurred");
            };

            objRequest.ontimeout  = function (event) {
                goradd.log("timeout error");
                opts.error("The server did not respond in time");
            };

            _currentRequests[ajaxID] = objRequest;
            objRequest.send(opts.data);
        },
    };

})();
