/*

Goradd Debugging Additions

This file attaches some code used by the test harness to drive browser-based tests, and other
debug-only code. It is only loaded in debug mode.

*/

(function() {

    goradd.initFormTest = function() {
        if (window.opener) { //
            // This next line is a potential security hole, so its important that this code NOT be loaded by the release version.
            window.opener.postMessage({pagestate: g$('Goradd__PageState').val()}, "*");
            goradd.form().addEventListener ('teststep', goradd.testStep);
            var event = new CustomEvent('teststep', { bubbles: true, detail: -1 });
            goradd.form().dispatchEvent(event);
        }
    };

    goradd._testStepPending= false;

    goradd.testStep = function(event) {
        if (goradd._actionQueue.length > 0) {
            goradd.queueAction({f: function() {
                    goradd._postTestStep(event);
                }, last: true, name: "testStep"});
        } else {
            goradd._postTestStep(event);
        }
    };


    goradd._postTestStep = function(event) {
        if (event) {
            if (!goradd.ajaxq.isRunning()) {
                goradd.log("Posting message: Ajax complete", event.detail);
                window.opener.postMessage({ajaxComplete: event.detail}, "*");
            } else {
                goradd.log("Delaying ajax complete message", event.detail);
                goradd._testStepPending = true;
                goradd.currentStep = event.detail;
            }
        } else {
            // We are being notified that an ajax action has completed
            if (goradd._testStepPending) {
                goradd.log("Reposting delayed message: Ajax complete", goradd.currentStep);
                goradd._testStepPending = false;
                window.opener.postMessage({ajaxComplete: goradd.currentStep}, "*");
            }
        }
    };

    if (window.console) {
        // This lets us easily turn off logging in production without losing any console.log capabilities
        goradd.log = console.log;
    }

    goradd.postMarker = function(marker) {
        if (window.opener) {
            window.opener.postMessage({testMarker: marker}, "*");
        }
    };

    /**
     * Displays the ajax error in either a popup window, or a new web page.
     * @param {string} resultText
     * @param {number} err
     * @private
     */
    goradd.displayAjaxError = function(resultText, err) {
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
            var el = goradd.tagBuilder("div").attr("id", "Goradd_AJAX_Error").html("<button onclick=\"window.goradd.g('Goradd_AJAX_Error').remove()\">OK</button>").appendTo(goradd.form());
            goradd.tagBuilder("div").html(resultText).appendTo(el);
        }
    }


}) ();