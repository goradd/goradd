/**
 * goradd.js
 *
 * This is the shim between the goradd go code and the browser. It enables ajax and other kinds of
 * communication between the server and the client.
 *
 * Goals:
 *  - Compatible with all current browsers, IE 10+ and Opera Mobile (ES5).
 *  - Small. We want to be compatible in situations with limited bandwidth.
 *    The minimized version should be as small as possible.
 *  - Provide utility code to javascript widgets and plugins.
 *  - Do not use jQuery or other frameworks, but be compatible if its used by the devoloper.
 */

if (!function () {
    "use strict";
    return Function.prototype.bind && XMLHttpRequest && !this;
}()) {
    window.location = "/Unsupported.g";
}


var $j; // Universal shortcut for jQuery so that you can execute jQuery code outside in a no-conflict mode.
var goradd;

(function( $ ) {
"use strict";

$j = $;

/**
 * @namespace goradd
 */
goradd = {
    /**
     * General support library. Here we recreate a few useful functions from jquery.
     */

    /**
     * el returns the html element t. t can be an id, or an element, and if an element, it will just return the element
     * back. This is used below so that all the functions can pass either an element, or the id of an element. Returns
     * null if not found.
     * @param t {string|object}
     * @returns {Element}
     */
    el: function(t) {
        if (typeof t == "object") {
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
        return goradd.qs('form[data-grctl="form"]');
    },
    /**
     * matches returns true if the given element matches the css selector.
     * @param el
     * @param sel
     * @returns {boolean}
     */
    matches: function(el, sel) {
        return goradd.g(el).matches(sel);
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
        goradd.each(attributes, function() {
            script[key] = this[key];
        });

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
        goradd.each(attributes, function() {
            link[key] = this[key];
        });
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

    _toKebab: function(s) {
        var s2 =  s.replace(/[A-Z]/g, function(m, offset, s) {
            var n = "-" + m.toLowerCase();
           return n;
        });
        return s2;
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
            goradd.g(el).trigger('formObjChanged');
        }
    },

    /**
     * Private members
     */
    _controlValues: {},
    _formObjsModified: {},
    _ajaxError: false,
    _blockEvents: false,
    _inputSupport: true,


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
        if (!goradd._controlValues[strControlId]) {
            goradd._controlValues[strControlId] = {};
        }
        goradd._controlValues[strControlId][strProperty] = strNewValue;
    },
    /**
     * formObjChanged is an event handler that records that a control has changed in order to synchronize the control with
     * the server on the next request. Send the formObjChanged event to the control
     * that changed, and it will bubble up to the form and be caught here.
     * @param event
     */
    formObjChanged: function (event) {
        goradd._formObjsModified[event.target.id] = true;
    },
    /**
     * Initializes form related scripts. This is called by injected code on a goradd form.
     */
    initForm: function () {
        var form =  goradd.form();
        goradd.g(form).on('formObjChanged', goradd.formObjChanged); // Allow any control, including hidden inputs, to trigger a change and post of its data.
        goradd.g(form).on('submit', function(event) {
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
        goradd._registerControls();
    },
    _registerControl: function(ctrl) {
        if (!ctrl) {
            return;
        }

        // get the widget
        var g = goradd.g(ctrl);

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
            g.on('click', goradd.formObjChanged);
        }
        g.on('change input', goradd.formObjChanged, null, null, true); // make sure we get these events before later attached events

        // widget support, using declarative methods
        if (goradd.widget.new) {
            var widget;
            var options = {};
            goradd.each(g.attr(), function(k,v) {
                if (k === "data-gr-widget") {
                    widget = v;
                } else if (k.substr(0, 12) === "data-gr-opt-") {
                    options[k.substr(12)] = v;
                }
            });
            if (widget) {
                widget = goradd.widget.new(widget, options, ctrl);
                // Replace the control's widget with the new one. There can be only one goradd widget associated with
                // a particular control. We will need some other mechanism for mixins if needed.
                ctrl.goradd.widget = widget;
            }
        }
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
        if (goradd._blockEvents) {
            return;  // We are waiting for a response from the server
        }

        var form = goradd.form();
        var gForm = goradd.g(form);

        params.callType = "Server";

        // Notify custom controls that we are about to post

        gForm.trigger("posting", "Server");

        // Post custom javascript control values
        if (goradd.isEmptyObj(goradd._controlValues)) {
            params.controlValues = goradd._controlValues;
        }
        goradd.el('Goradd__Params').value = JSON.stringify(params);

        // trigger our own form submission so we can catch it
        gForm.trigger("submit");
    },


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
    _getAjaxData: function(params) {
        var form = goradd.form(),
            controls = goradd.g(form).qa('input,select,textarea'),
            postData = {};

        // Notify controls we are about to post.
        goradd.g(form).trigger("posting", "Ajax");

        goradd.each(controls, function(i,c) {
            var id = c.id;
            var blnForm = (id && (id.substr(0, 8) === 'Goradd__'));

            if (!goradd._inputSupport || // if not oninput support, then post all the controls, rather than just the modified ones, because we might have missed something
                goradd._ajaxError || // Ajax error would mean that _formObjsModified is invalid. We need to submit everything.
                (id && goradd._formObjsModified[id]) ||  // We try to ignore controls that have not changed to reduce the amount of data sent in an ajax post.
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
                        postData[id] = goradd.g(c).value();
                        break;
                }
            }
        });

        // Update most of the Goradd__ parameters explicitly here. Others, like the state and form id will have been handled above.
        params.callType = "Ajax";
        if (!goradd.isEmptyObj(goradd._controlValues)) {
            params.controlValues = goradd._controlValues;
        }
        postData.Goradd__Params = JSON.stringify(params);

        goradd._ajaxError = false;
        goradd._formObjsModified = {};
        goradd._controlValues = {};

        return postData;
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
            formAction = goradd.g(form).attr("action"),
            async = params.hasOwnProperty("async");

        if (goradd._blockEvents) {
            return;
        }

        params.formId = form.id;

        goradd.log("postAjax", params);

        // Use an ajax queue so ajax requests happen synchronously
        goradd.ajaxq.enqueue(function() {
            var data = goradd._getAjaxData(params);

            return {
                url: formAction,
                data: data,
                error: function (result, err) {
                    goradd._displayAjaxError(result, err);
                    goradd.testStep();
                    return false;
                },
                success: function (json) {
                    goradd.log("Ajax success ", json);

                    if (json.js) {
                        for (var k in json.js) {
                            goradd.loadJavaScriptFile(k, json.js[k]);
                        }
                    }
                    goradd._processImmediateAjaxResponse(json, params);
                    // TODO: Wait until javascripts above are loaded before proceeding?
                    goradd._processDeferredAjaxResponse(json);
                    goradd._blockEvents = false;
                }
            };
        }, async);
    },
    /**
     * Displays the ajax error in either a popup window, or a new web page.
     * @param resultText
     * @private
     */
    _displayAjaxError: function(resultText, err) {
        var objErrorWindow;

        goradd._ajaxError = true;
        goradd._blockEvents = false;

        if (resultText.substr(0, 15) === '<!DOCTYPE html>') {
            window.alert("An error occurred.\r\n\r\nThe error response will appear in a new popup.");
            objErrorWindow = window.open('about:blank', 'qcubed_error', 'menubar=no,toolbar=no,location=no,status=no,scrollbars=yes,resizable=yes,width=1000,height=700,left=50,top=50');
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
            var el = goradd.tb("div").attr("id", "Goradd_AJAX_Error").
                html("<button onclick='goradd.remove(\"Goradd_AJAX_Error\")'>OK</button>").
                appendTo(goradd.form());
            goradd.tb("div").html(resultText).appendTo(el);
        }
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

        goradd._inputSupport = 'oninput' in document;
        // IE 9 has a major bug in oninput, but we are requiring IE 10+, so no problem.
        // I think the only major browser that does not support oninput is Opera mobile.

        goradd.g(goradd.form()).on("ajaxQueueComplete", function() {
            goradd._processFinalCommands();
        });

        // TODO: Add a detector of the back button. This detector should ping the server to make sure the pagestate exists on the server. If not,
        // it should reload the form.
    },
    /**
     * Responds to the part of an ajax response that must be handled serially before other handlers can fire.
     *
     * @param {Object} json     json generated by goradd application
     * @param {Object} params   option parameters
     * @private
     */
    _processImmediateAjaxResponse: function(json, params) {
        goradd.each(json.controls, function(id) {
            var el = goradd.el(id),
                $ctrl = goradd.g(el),
                wrapper = goradd.el(id + "_ctl");

            if (this.value !== undefined && $ctrl) {
                $ctrl.value(this.value);
            }

            if (this.attributes !== undefined && $ctrl) {
                $ctrl.prop(this.attributes);
            }

            if (this.html !== undefined) {
                if (wrapper !== null) {
                    // Control's wrapper was found, so replace the control and the wrapper
                    goradd.g(wrapper).htmlBefore(this.html);
                    goradd.g(wrapper).remove(wrapper);
                } else if ($ctrl) {
                    // control was found without a wrapper, replace it in the same position it was in.
                    // remove related controls (error, name ...) for wrapper-less controls
                    var relSelector = "[data-grel='" + id + "']",
                        relatedItems = goradd.qa(relSelector);

                    var p = $ctrl.parents();
                    var relatedParent = p.filter(function(el) {
                        return goradd.g(el).matches(relSelector);
                    }).pop();

                    if (relatedParent) {
                        relatedParent.insertAdjacentElement("beforebegin", el);
                    }
                    if (relatedItems && relatedItems.length > 0) {
                        goradd.each(relatedItems, function(i, el) {
                            goradd.g(el).remove();
                        })
                    }
                    $ctrl.htmlBefore(this.html);
                    $ctrl.remove();
                }
                else {
                    // control is being injected at the top level, so put it at the end of the form.
                    goradd.f(goradd.form()).appendHtml(this.html);
                }
            }
        });

        goradd._registerControls();

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
    },
    /**
     * Process the part of an ajax response that can be deferred and so be executed in parallel with other operations.
     *
     * @param {object} json  Json generated by the goradd application.
     * @private
     */
    _processDeferredAjaxResponse: function(json) {
        goradd.each(json.commands, function (i,command) {
            if (command.final &&
                goradd.ajaxq.isRunning()) {
                goradd._enqueueFinalCommand(command);
            } else {
                goradd._processCommand(command);
            }
        });
        if (json.winclose) {
            window.close();
        }
        if (json.loc) {
            if (goradd._closeWebSocket) {
                goradd._closeWebSocket(1001);
            }
            if (json.loc === 'reload') {
                window.location.reload(true);
            } else {
                document.location = json.loc;
            }
        }
        if (json.profileHtml) {
            var c = goradd.el("dbProfilePane");
            if (!$c) {
                goradd.g(goradd.form()).htmlAfter("<div id = 'dbProfilePane'></div>");
                c = goradd.el("dbProfilePane");
            }
            c.innerHTML = json.profileHtml;
        }
        goradd.testStep();
    },
    /**
     * Process a single command.
     * @param {object} command
     * @private
     */
    _processCommand: function(command) {
        var params,
            objs;

        if (command.script) {
            // TODO: clean this up a bit by using ids for inserted scripts
            var script   = document.createElement("script");
            script.type  = "text/javascript";
            script.text  = command.script;
            document.body.appendChild(script);
        }
        else if (command.selector) {
            params = goradd._unpackArray(command.params);

            if (typeof command.selector === 'string') {
                // general selector
                objs = goradd.qa(command.selector);
            } else {
                objs = goradd.g(command.selector[0]).qa(command.selector[1]);
            }

            goradd.each (objs, function (i,v) {
                var $c = goradd.g(v);
                if (typeof $c[command.func] === "function") {
                    $c[command.func].apply($c, params);
                }
            });
        }
        else if (command.func) {
            params = goradd._unpackArray(command.params);

            // Find the function by name. Walk an object list in the process.
            objs = command.func.split(".");
            var obj = window;
            if (command.id) {
                obj = goradd.g(command.id);
            } else if (command.jqueryId) {
                obj = jQuery(command.jqueryId);
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
    /**
     * Places the given command in the queue so that it is executed last.
     * @param command
     * @private
     */
    _enqueueFinalCommand: function(command) {
        goradd.finalCommands.push(command);
    },
    /**
     * Execute the final commands.
     * @private
     */
    _processFinalCommands: function() {
        while(goradd.finalCommands.length) {
            var command = goradd.finalCommands.pop();
            goradd._processCommand(command);
        }
    },
    /**
     * Convert from JSON return value to an actual jQuery object. Certain structures don't work in JSON, like closures,
     * but can be part of a javascript object. We use special codes to piece together functions, closures, dates, etc.
     * @param params
     * @returns {*}
     * @private
     */
    _unpackArray: function(params) {
        if (!params) {
            return null;
        }
        var newParams = [];

        goradd.each(params, function (index, item){
            if (Array.isArray(item)) {
                item = goradd._unpackArray (item);
            } else if (typeof item === 'object' && item !== null) {
                if (item.goraddObject) {
                    item = goradd._unpackObj(item);  // top level special object
                }
                else {
                    // look for special objects inside top level objects.
                    var newItem = {};
                    goradd.each (item, function (key, obj) {
                        newItem[key] = goradd._unpackObj(obj);
                    });
                    item = newItem;
                }
            }
            newParams.push(item);
        });
        return newParams;
    },

    /**
     * Given an object coming from goradd, will attempt to decode the object into a corresponding javascript object.
     * @param obj
     * @returns {*}
     * @private
     */
    _unpackObj: function (obj) {
        var params;

        if (typeof obj === "object" && obj !== null) {
            if (Array.isArray(obj)) {
                return goradd._unpackArray(obj);
            } else if (obj.goraddObject) {
                switch (obj.goraddObject) {
                    case 'closure':
                        if (obj.params) {
                            params = [];
                            goradd.each (obj.params, function (i, v) {
                                params.push(goradd._unpackObj(v)); // recurse
                            });

                            return new Function(params, obj.func);
                        } else {
                            return new Function(obj.func);
                        }

                    case 'dt':
                        if (obj.z) {
                            return null;
                        } else if (obj.t) {
                            return new Date(Date.UTC(obj.y, obj.mo, obj.d, obj.h, obj.m, obj.s, obj.ms));
                        } else {
                            return new Date(obj.y, obj.mo, obj.d, obj.h, obj.m, obj.s, obj.ms);
                        }

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
                                params.push(goradd._unpackObj(v)); // recurse
                            });
                        }
                        var func = target[obj.func];

                        return func.apply(target, params);
                }
            }
            else {
                var newItem = {};
                goradd.each (obj, function (key, obj2) {
                    newItem[key] = goradd._unpackObj(obj2);
                });
                return newItem;
            }
        }

        return obj; // no change
    },
    _registerControls: function() {
        var els = goradd.qa('[data-grctl]');
        goradd.each(els, function(el) {
            goradd._registerControl(this);
        });
    },
    updateForm: function() {
        // call this whenever you generally just need the form to update without a specific action.
        var newTime = new Date().getTime();

        // the following code prevents too many updates from happening in a short amount of time.
        // the default will update no faster than once per minUpdateInterval.
        if (newTime - goradd._prevUpdateTime > goradd.minUpdateInterval) {
            //refresh immediately
            goradd.log("Immediate update");
            goradd._prevUpdateTime = new Date().getTime();
            goradd.postAjax ({});
            goradd.clearTimer('goradd.update');
        } else if (!goradd._objTimers['goradd.update']) {
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

    /**
     * Named timer support. These allow you to create timers without having to keep a copy of the timer around.
     */

    /**
     * Current timer store, by id
     */
    _objTimers: {},
    /**
     * Clears the named timer.
     * @param {string} strTimerId
     */
    clearTimer: function(strTimerId) {
        if (goradd._objTimers[strTimerId]) {
            goradd.log("clearTimer", strTimerId);
            clearTimeout(goradd._objTimers[strTimerId]);
            goradd._objTimers[strTimerId] = null;
        }
    },
    /**
     * Sets the named timer, given an action and a delay.
     * @param strTimerId
     * @param action
     * @param intDelay
     */
    setTimer: function(strTimerId, action, intDelay) {
        goradd.log("setTimer", strTimerId, intDelay);
        goradd._objTimers[strTimerId] = setTimeout(
            function() {
                goradd.clearTimer(strTimerId);
                action();
            }, intDelay);
    },
    hasTimer: function(strTimerId) {
        return !!goradd._objTimers[strTimerId];
    },
    /**
     * Creates a timer that performs can perform periodic events, and that fires the timerexperedevent event when it is done.
     * @param strControlId
     * @param intDeltaTime
     * @param blnPeriodic
     */
    startTimer: function(strControlId, intDeltaTime, blnPeriodic) {
        var strTimerId = strControlId + '_ct';
        goradd.stopTimer(strControlId, blnPeriodic);
        if (blnPeriodic) {
            goradd._objTimers[strTimerId] = setInterval(function() {
                goradd.g(strControlId).trigger('timerexpiredevent');
            }, intDeltaTime);
        } else {
            goradd._objTimers[strTimerId] = setTimeout(function() {
                goradd.g(strControlId).trigger('timerexpiredevent');
            }, intDeltaTime);
        }
    },
    /**
     * Stops the named timer, allowing you to specify whether its a periodic timer or not.
     * @param strControlId
     * @param blnPeriodic
     */
    stopTimer: function(strControlId, blnPeriodic) {
        var strTimerId = strControlId + '_ct';
        if (goradd._objTimers[strTimerId]) {
            if (blnPeriodic) {
                clearInterval(goradd._objTimers[strTimerId]);
            } else {
                clearTimeout(goradd._objTimers[strTimerId]);
            }
            goradd._objTimers[strTimerId] = null;
        }
    },
    //////////////////////////////
    // Action queue support
    //////////////////////////////
    /* Javascript has a problem when two events happen simultaneously. In particular, a click event might also
    result in a change event, and under certain circumstances this could cause the click event to be dropped. In particular,
    if the change event moves the focus away from the button, the click event will not record. We therefore delay
    the processing of all events to try to queue them up before processing.
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
    }

};


///////////////////////////////
// Watcher support
///////////////////////////////
goradd._prevUpdateTime = 0;
goradd.minUpdateInterval = 500; // milliseconds to limit broadcast updates. Feel free to change this.
goradd.broadcastChange = function () {
    if ('localStorage' in window && window.localStorage !== null) {
        var newTime = new Date().getTime();
        localStorage.setItem("goradd.broadcast", newTime); // must change value to induce storage event in other windows
    }
};



/////////////////////////////////
// Controls-related functionality
/////////////////////////////////

goradd.getControl = function(controlId) {
    return document.getElementById(controlId);
};

goradd.getWrapper = function(mixControl) {
    if (typeof mixControl === 'string') {
        return document.getElementById(mixControl + "_ctl")
    } else {
        return document.getElementById(mixControl.id + "_ctl")
    }
};

goradd.getPageState = function() {
    return document.getElementById("Goradd__PageState").value;
};

goradd.finalCommands = [];
goradd.currentStep = 0;
goradd.stepFunction = null;


goradd.redirect = function(newLocation) {
    window.location = newLocation
};

/**
 * tb returns a TagBuilder. Use it as follows:
 * tag = goradd.tb("div").attr("class", "myClass").text("I am text").appendTo("objId");
 * @param tag {string}
 * @returns {goradd.TagBuilder}
 */
goradd.tb = function(tag) {
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
            if (this.matches(sel)) {
                return el;
            }
            el = el.parentElement;
        }
        return null;
    },
    /**
     * attr gets attributes on a dom object. Remember that attributes are not the same as properties.
     * To access properties, use prop. These specifically access the attributes defined in html, but not anything set
     * afterwards.
     * Returns undefined if the attribute does not exist.
     * @param a (optional) {string} The attribute name to return. Otherwise returns an object that is a map of all defined attributes.
     * @returns {null|boolean|*}
     */
    attr: function() {
        var t = this.element;
        var self = this;
        if (arguments.length === 0) {
            // Return an object mapping all the attributes of the html object
            if (t.hasAttributes()) {
                var attr = {};
                // Apparently IE has a quirk where it returns all possible attributes, and not just set attributes.
                goradd.each(this.element.attributes, function(v,n) {
                    n = n.nodeName || n.name;
                    if (t.hasAttribute(n)) {
                        attr[n] = t.getAttribute(n);
                    }
                });
                return attr;
            }
            return undefined; // no attributes are set
        }
        if (arguments.length === 1) {
            var a = arguments[0];
            // get value
            if (!t.hasAttribute(a)) {
                return undefined;
            }
            var v = t.getAttribute(a);
            if (v === null || v === "true" || v === "") {
                return true; // A boolean attribute, it just exists with no value or with "true"
            } else if (v === "false") {
                return false;
            } else {
                return v;
            }
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
            goradd.each(c.substr(1).split(" "), function(i,v) {
                if (v !== "") {
                    el.classList.add(v);
                }
            });
        } else if (c.substr(0,1) === "-") {
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
     * on attaches an event handler to the given html object.
     * Filtering and potentially supplying data to the event are also included.
     * If data is a function, the function will be called when the event fires and the
     * result of the function will be provided as data to the event. The "this" parameter
     * will be the element with the given targetId, and the function will be provided the event object.
     *
     * @param eventNames {string} One or more event names separated by spaces
     * @param eventHandler
     * @param filter
     * @param data
     * @param capture True to fire this event during initial capture phase. False to wait until it bubbles.
     */
    on: function(eventNames, eventHandler, filter, data, capture) {
        if (!capture) {
            capture = false;
        }
        var el = this.element;
        var events = eventNames.split(" ");
        goradd.each(events, function(i,eventName) {
            el.addEventListener(eventName, function (event) {
                if (filter && !goradd.g(event.target).matches(filter)) {
                    return
                }
                if (data) {
                    if (typeof data === "function") {
                        data = data.call(el, event);
                    }
                    event.grdata = data;
                }
                if (event.detail) {
                    eventHandler.call(el, event, event.detail); // simulate adding extra items to event handler
                } else if (data) {
                    eventHandler.call(el, event, data); // simulate adding extra items to event handler
                } else {
                    eventHandler.call(el, event);
                }
            }, capture);
        });
    },
    click: function() {
        // use the built-in click to simulate a click on an item.
        this.element.click();
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
                event.initCustomEvent(eventName, true, true, extra);
            } else {
                event = new Event(eventName, {bubbles: true, detail: extra})
            }
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
     * @param el {object|string}
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
    value: function(v) {
        var el = this.element;
        var type = goradd.g(el).prop("type");
        if (arguments.length === 1) {
            // Setting the value
            switch (type) {
                case "select-multiple":
                    // Multi-select selections will attempt to set all items in the given array to the value
                    var opts = goradd.qa(el,'option');
                    goradd.each(opts, function(i, opt) {
                        opt.checked = (opt.value in v);
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
                    if ("value" in el) {
                        el.value = v;
                    }
                    break;
            }
            return el;
        } else {
            switch (type) {
                case "select-multiple":
                    // Multi-select selections will return an array of selected values
                    var sels = goradd.qa(el,':checked');
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
                        // Custom controls can add a "value" getter as well and this will pick that up too.
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
                key = goradd._toKebab(key);
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
    text: function(t) {
        if (arguments.length === 0) {
            return this.element.innerText;
        } else {
            this.element.innerText = t;
        }
    },
    html: function(t) {
        if (arguments.length === 0) {
            return this.element.innerHtml;
        } else {
            this.element.innerHtml = t;
        }
    },

    /**
     * f calls the named function, with the named parameters, on the goradd widget first, and if not found, will attempt
     * to call this on the element.
     * @param name
     * @param params
     */
    /*
    f: function(name, params) {
        var f = this[name];
        if (typeof f === "function") {
            return f.apply(params);
        } else {
            f = this.element[name];
            if (typeof f === "function") {
                return f.apply(params);
            }
        }
    }*/
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
    basePrototype.options = goradd.widget.extend( {}, basePrototype.options );

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

    constructor.prototype = goradd.widget.extend( basePrototype, proxiedPrototype, {
        constructor: constructor,
        namespace: namespace,
        widgetName: widgetName,
        widgetFullName: widgetFullName
    } );

};

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

goradd.widget.extend = function( target ) {
    var input = Array.prototype.slice.call( arguments, 1 );
    var inputIndex = 0;
    var inputLength = input.length;
    var key;
    var value;

    for ( ; inputIndex < inputLength; inputIndex++ ) {
        for ( key in input[ inputIndex ] ) {
            value = input[ inputIndex ][ key ];
            if ( input[ inputIndex ].hasOwnProperty( key ) && value !== undefined ) {

                // Clone objects
                if ( goradd.isPlainObject( value ) ) {
                    target[ key ] = goradd.isPlainObject( target[ key ] ) ?
                        goradd.widget.extend( {}, target[ key ], value ) :

                        // Don't extend strings, arrays, etc. with objects
                        goradd.widget.extend( {}, value );

                    // Copy everything else by reference
                } else {
                    target[ key ] = value;
                }
            }
        }
    }
    return target;
};

/**
 * This is the definition of the Widget class, which serves as the base class for other widgets. It itself is based
 * on the "g" class, which is a jQuery like wrapper. In other words, all the functions on the g class are available
 * to widgets throught the "this" variable, and can be overridden. One important function to override might be the
 * "value" function, which provides the value that will be used by ajax calls. If your widget only works through Ajax,
 * then that is sufficient to keep the go side of things updated.
 */
goradd.widget("goradd.Widget", goradd.g, {
    /**
     * _createWidget acts as the constructor of all widgets. It can be overridden by the widget if needed, but
     * you normally do not need to. Implement _create() to make a private constructor.
     * @param options
     * @param element
     * @private
     */
    _createWidget: function(options, element) {
        this.element = goradd.el(element);

        this.options = goradd.widget.extend( {},
            this.options,
            options );

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
});

})( jQuery );

////////////////////////////////
// Goradd Shortcuts and Initialize
////////////////////////////////

goradd.initialize();
