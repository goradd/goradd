
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
     *
     * @param id
     * @returns {*}
     */
    el: function(id) {
        return document.getElementById(id);
    },
    qs: function(sel) {
        return document.querySelector(sel);
    },
    qa: function(sel) {
        return document.querySelectorAll(sel);
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
        if (Element.prototype.matches) {
            return el.matches(filter);
        } else {
            var matches = goradd.qa(sel),
                i = matches.length;
            while (--i >= 0 && matches.item(i) !== el) {}
            return i > -1;
        }
    },

    /**
     * on attaches an event handler to the given html object.
     * Filtering and potentially supplying data to the event are also included.
     * If data is a function, the function will be called when the event fires and the
     * result of the function will be provided as data to the event. The "this" parameter
     * will be the element with the given targetId, and the function will be provided the event object.
     *
     * @param target {string|object} Either a string id of an html object, or the object itself
     * @param eventName
     * @param eventHandler
     * @param filter
     * @param data
     */
    on: function(target, eventName, eventHandler, filter, data) {
        if (typeof target != "object") {
            target = goradd.el(target);
        }
        target.addEventListener(eventName, function(event) {
            if (filter && !goradd.matches(event.target, filter)) {
                return
            }
            if (data) {
                if (typeof data === "function") {
                    data = data.call(target, event);
                }
                event.grdata = data;
            }
            eventHandler.call(target, event);
        });
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
     * Records that a control has changed in order to synchronize the control with
     * the go version on the next request. Send the formObjChanged event to the control
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
        goradd.on(form, 'formObjChanged', goradd.formObjChanged); // Allow any control, including hidden inputs, to trigger a change and post of its data.
        goradd.on(form, 'submit', function(event) {
            if (!goradd.el('Goradd__Params').value) { // did postBack initiate the submit?
                // if not, prevent implicit form submission. This can happen in the rare case we have a single field and no submit button.
                event.preventDefault();
            }
        });
        goradd._registerControls();
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

        var $objForm = $(goradd.getForm());
        var formId = $objForm.attr("id");

        params.callType = "Server";

        // Notify custom controls that we are about to post
        $objForm.trigger("posting", "Server");

        // Post custom javascript control values
        if (!$.isEmptyObject(goradd._controlValues)) {
            params.controlValues = goradd._controlValues;
        }
        $('#Goradd__Params').val(JSON.stringify(params));

        // have $ trigger the submit event (so it can catch all submit events)
        $objForm.trigger("submit");
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
        var $form = $('#' + params.formId),
            controls = $form.find('input,select,textarea'),
            postData = {};

        // Notify controls we are about to post.
        $form.trigger("posting", "Ajax");

        // We try to ignore controls that have not changed to reduce the amount of data sent in an ajax post.
        controls.each(function() {
            var id = this.id,
                blnForm = (id && (id.substr(0, 8) === 'Goradd__'));


            if (!goradd._inputSupport || // if not oninput support, then post all the controls, rather than just the modified ones, because we might have missed something
            goradd._ajaxError || // Ajax error would mean that _formObjsModified is invalid. We need to submit everything.
            (id && goradd._formObjsModified[id]) ||
             blnForm) {  // all controls with Goradd__ at the beginning of the id are always posted.
                var $ctrl = $(this),
                    strType = $ctrl.prop("type");

                switch (strType) {
                    case "radio":
                        // Radio buttons listen to their name.
                        var n = $ctrl.attr("name");
                        var $sel = $('input:radio[name=' + n + ']:checked');
                        var val = null;
                        if ($sel.length) {
                            val = $sel.val();
                        }
                        postData[n] = val;
                        break;
                    case "checkbox":
                        postData[id] = $ctrl.prop("checked");
                        break;
                    default:
                        // All goradd controls and subcontrols MUST have an id for this to work.
                        // There is a special case for checkbox groups, but they get handled on the server
                        // side differently between ajax and server posts.
                        // Also, the .val() will gather an array for multi-select lists automatically.
                        postData[id] = $ctrl.val();
                        break;
                }

            }
        });



        // Update most of the Goradd__ parameters explicitly here. Others, like the state and form id will have been handled above.
        params.callType = "Ajax";
        if (!$.isEmptyObject(goradd._controlValues)) {
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
        var $objForm = $(goradd.getForm()),
            formAction = $objForm.attr("action"),
            async = params.hasOwnProperty("async");

        if (goradd._blockEvents) {
            return;
        }

        params.formId = $objForm.attr('id');

        goradd.log("postAjax", params);

        // Use an ajax queue so ajax requests happen synchronously
        goradd.ajaxq.enqueue(function() {
            var data = goradd._getAjaxData(params);

            return {
                url: formAction,
                data: data,
                error: function (result) {
                    goradd._displayAjaxError(result);
                    goradd.testStep();
                    return false;
                },
                success: function (json) {
                    goradd.log("Ajax success ", json);

                    if (json.js) {
                        var deferreds = [];
                        // Load all javascript files before attempting to process the rest of the response, in case some things depend on the injected files
                        $.each(json.js, function (i, v) {
                            deferreds.push(goradd.loadJavaScriptFile(v));
                        });
                        goradd._processImmediateAjaxResponse(json, params); // go ahead and begin processing things that will not depend on the javascript files to allow parallel processing
                        $.when.apply($, deferreds).then(
                            function () {
                                goradd._processDeferredAjaxResponse(json);
                                goradd._blockEvents = false;
                            }, // success
                            function () {
                                goradd.log('Failed to load a file');
                                goradd._blockEvents = false;
                            } // failed to load a file. What to do?
                        );
                    } else {
                        goradd._processImmediateAjaxResponse(json, params);
                        goradd._processDeferredAjaxResponse(json);
                        goradd._blockEvents = false;
                    }
                }
            };
        }, async);
    },
    /**
     * Displays the ajax error in either a popup window, or a new web page.
     * @param resultText
     * @param textStatus
     * @param errorThrown
     * @private
     */
    _displayAjaxError: function(resultText) {
        var objErrorWindow;

        goradd._ajaxError = true;
        goradd._blockEvents = false;

        if (resultText.substr(0, 15) === '<!DOCTYPE html>') {
            window.alert("An error occurred.\r\n\r\nThe error response will appear in a new popup.");
            objErrorWindow = window.open('about:blank', 'qcubed_error', 'menubar=no,toolbar=no,location=no,status=no,scrollbars=yes,resizable=yes,width=1000,height=700,left=50,top=50');
            objErrorWindow.focus();
            objErrorWindow.document.write(resultText);
        } else {
            resultText = $('<div>').html(resultText);
            $('<div id="Goradd_AJAX_Error" />')
                .append(resultText)
                .append('<button onclick="$(this).parent().hide()">OK</button>')
                .appendTo('form');
        }
    },
    /**
     * Start me up.
     */
    initialize: function() {
        ////////////////////////////////
        // Browser-related functionality
        ////////////////////////////////

        goradd.loadJavaScriptFile = function(strScript, objCallback) {
            return $.ajax({
                url: strScript,
                success: objCallback,
                dataType: "script",
                cache: true
            });
        };

        goradd.loadStyleSheetFile = function(strStyleSheetFile, strMediaType) {
            if (strMediaType){
                strMediaType = " media="+strMediaType;
            }
            $('head').append('<link rel="stylesheet"'+strMediaType+' href="' + strStyleSheetFile + '" type="text/css" />');
        };

        /////////////////////////////
        // Form-related functionality
        /////////////////////////////
        /*
        $(window).on ("storage", function (o) {
            if (o.originalEvent.key === "goradd.broadcast") {
                goradd.updateForm();
            }
        });*/

        goradd._inputSupport = 'oninput' in document;
        // IE 9 has a major bug in oninput, but we are requiring IE 10+, so no problem.
        // I think the only major browser that does not support oninput is Opera mobile.

        $( document ).ajaxComplete(function( event, request, settings ) {
            if (!goradd.ajaxq.isRunning()) {
                goradd._processFinalCommands();  // If there was no ajax queue, we would have already processed final commands
            }
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
        if (json.controls) {
            $.each(json.controls, function(id) {
                var strControlId = id,
                    $control = $(goradd.getControl(strControlId)),
                    $wrapper = $(goradd.getWrapper(strControlId));

                if (this.value !== undefined) {
                    $control.val(this.value);
                }

                if (this.attributes !== undefined) {
                    $control.attr (this.attributes);
                }

                if (this.html !== undefined) {
                    if ($wrapper.length) {
                        // Control's wrapper was found, so replace the control and the wrapper
                        $wrapper.before(this.html).remove();
                    }
                    else if ($control.length) {
                        // control was found without a wrapper, replace it in the same position it was in.
                        // remove related controls (error, name ...) for wrapper-less controls
                        var relSelector = "[data-grel='" + strControlId + "']",
                            relItems = $(relSelector),
                            $relParent;

                        if (relItems && relItems.length) {
                            // if the control is wrapped in a related control, we move the control outside the related controls
                            // before deleting the related controls
                            $relParent = $control.parents(relSelector).last();
                            if ($relParent.length) {
                                $control.insertBefore($relParent);
                            }
                            relItems.remove();
                        }

                        $control.before(this.html).remove();
                    }
                    else {
                        // control is being injected at the top level, so put it at the end of the form.
                        var $objForm = $(goradd.getForm());
                        $objForm.append(this.html);
                    }
                }
            });
        }

        goradd._registerControls();

        if (json.watcher && params.controlId) {
            goradd.broadcastChange();
        }
        if (json.ss) {
            $.each(json.ss, function (i,v) {
                goradd.loadStyleSheetFile(v, "all");
            });
        }
        if (json.alert) {
            $.each(json.alert, function (i,v) {
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
        if (json.commands) { // commands
            $.each(json.commands, function (index, command) {
                if (command.final &&
                    goradd.ajaxq.isRunning()) {

                    goradd._enqueueFinalCommand(command);
                } else {
                    goradd._processCommand(command);
                }
            });
        }
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
            var c = $("#dbProfilePane");
            if (c.length == 0) {
                c = $("<div id = 'dbProfilePane'></div>");
                $(goradd.getForm()).parent().append(c);
            }
            c.html(json.profileHtml);
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
            var script   = document.createElement("script");
            script.type  = "text/javascript";
            script.text  = command.script;
            document.body.appendChild(script);
        }
        else if (command.selector) {
            params = goradd._unpackArray(command.params);

            if (typeof command.selector === 'string') {
                objs = $(command.selector);
            } else {
                objs = $(command.selector[0], command.selector[1]);
            }

            // apply the function on each jQuery object found, using the found jQuery object as the context.
            objs.each (function () {
                var $item = $(this);
                if ($item[command.func]) {
                    $item[command.func].apply($(this), params);
                }
            });
        }
        else if (command.func) {
            params = goradd._unpackArray(command.params);

            // Find the function by name. Walk an object list in the process.
            objs = command.func.split(".");
            var obj = window;
            var ctx = null;

            $.each (objs, function (i, v) {
                ctx = obj;
                obj = obj[v];
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

        $.each(params, function (index, item){
            if ($.type(item) === 'object') {
                if (item.goraddObject) {
                    item = goradd._unpackObj(item);  // top level special object
                }
                else {
                    // look for special objects inside top level objects.
                    var newItem = {};
                    $.each (item, function (key, obj) {
                        newItem[key] = goradd._unpackObj(obj);
                    });
                    item = newItem;
                }
            }
            else if ($.type(item) === 'array') {
                item = goradd._unpackArray (item);
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
        if ($.type(obj) === 'object' &&
                obj.goraddObject) {

            switch (obj.goraddObject) {
                case 'closure':
                    if (obj.params) {
                        params = [];
                        $.each (obj.params, function (i, v) {
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
                    $.each (vars, function (i, v) {
                        val = val[v];
                    });
                    return val;

                case 'func':
                    // Returns the result of the given function called immediately
                    // Find the function and context starting at the window context.
                    var target = window;
                    var params;
                    if (obj.context) {
                       var objects = obj.context.split(".");
                        $.each (objects, function (i, v) {
                            target = target[v];
                        });
                    }

                    if (obj.params) {
                        params = [];
                        $.each (obj.params, function (i, v) {
                            params.push(goradd._unpackObj(v)); // recurse
                        });
                    }
                    var func = target[obj.func];

                    return func.apply(target, params);
            }
        }
        else if ($.type(obj) === 'object') {
            var newItem = {};
            $.each (obj, function (key, obj2) {
                newItem[key] = goradd._unpackObj(obj2);
            });
            return newItem;
        }
        else if ($.type(obj) === 'array') {
            return goradd._unpackArray(obj);
        }
        return obj; // no change
    },
    _registerControls: function() {
        $('[data-grctl]').not('[data-grctl="form"]').each(function() {
            goradd.registerControl(this);
        });
    },


    /***********************
     * Javascriptable Actions
     ***********************/
    focus: function(id) {
        goradd.getControl(id).focus();
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
                $('#' + strControlId).trigger('timerexpiredevent')
            }, intDeltaTime);
        } else {
            goradd._objTimers[strTimerId] = setTimeout(function() {
                $('#' + strControlId).trigger('timerexpiredevent')
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

};


//////////////////////////////
// Action queue support
//////////////////////////////
/* Javascript/jquery has a problem when two events happen simultaneously. In particular, a click event might also
result in a change event, and under certain circumstances this could cause the click event to be dropped. In particular,
if the change event moves the focus away from the button, the click event will not record. We therefore delay
the processing of all events to try to queue them up before processing.
Its very strange. Something to debug at a future date.
*/

goradd.actionQueue = [];
goradd.queueAction = function(params) {
    if (params.last) {
        var delay = 0;

        goradd.actionQueue.forEach(function(item) {
            if (item.d > delay) {
                delay = item.d;
            }
        });
        params.d = delay + 1;
    }
    goradd.log("queueAction: " + params.name);
    goradd.actionQueue.push(params);
    if (!goradd.hasTimer("goradd.actions")) {
        goradd.setTimer("goradd.actions", goradd.processActions, 10);
    }
};
goradd.processActions = function() {
    while (goradd.actionQueue.length > 0) {
        var params = goradd.actionQueue.shift();
        goradd.log("processAction: " + params.name + " delay: " + params.d);
        if (params.d > 0) {
            setTimeout(params.f, params.d);
        } else {
            params.f();
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

goradd.updateForm = function() {
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
    }
    else {
        return document.getElementById($(mixControl).attr('id') + "_ctl")
    }
};

goradd.getForm = function() {
    return $('form[data-grctl="form"]')[0]
};

goradd.getPageState = function() {
    return document.getElementById("Goradd__PageState").value;
};


/**
 * Radio buttons are a little tricky to set if they are part of a group
 * @param strControlId
 */
goradd.setRadioInGroup = function(strControlId) {
    var $objControl = $('#' + strControlId);
    if ($objControl) {
        var groupName = $objControl.prop('name');
        if (groupName) {
            var $radios = $objControl.closest('form').find('input[type=radio][name=' + groupName + ']');
            $radios.val([strControlId]);  // jquery does the work here of setting just the one control
            $radios.trigger('formObjChanged'); // send the new values back to the form
        }
    }
};

goradd.finalCommands = [];
goradd.currentStep = 0;
goradd.stepFunction = null;

goradd.registerControl = function(objControl) {
    var objWrapper;

    if (!objControl) {
        return;
    }

    var $control = $(objControl);

    if ($control.data('gr-reg') === 'reg') {
        return // this control is already registered
    }

    // detect changes to objects before any changes trigger other events
    if (objControl.type === 'checkbox' || objControl.type === 'radio') {
        // clicks are equivalent to changes for checkboxes and radio buttons, but some browsers send change way after a click. We need to capture the click first.
        $(objControl).on ('click', goradd.formObjChanged);
    }
    $(objControl).on ('change input', goradd.formObjChanged);
    $(objControl).on ('change input', 'input, select, textarea', goradd.formObjChanged);   // make sure we get to bubbled events before later attached handlers


    // Link the Wrapper and the Control together
    objWrapper = goradd.getWrapper(objControl.id);
    if (objWrapper) {
        objWrapper.control = objControl;
    }
    $control.data('gr-reg', 'reg') // mark the control as registered so we don't attach events twice
};

goradd.redirect = function(newLocation) {
    window.location = newLocation
};

})( jQuery );

////////////////////////////////
// Goradd Shortcuts and Initialize
////////////////////////////////

goradd.initialize();
