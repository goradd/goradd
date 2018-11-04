/*

Goradd Testing Additions

This file attaches some code used by the test harness to drive browser-based tests. It is only loaded in debug mode.

*/
(function() {

    goradd.initFormTest = function() {
        if (window.opener) { //
            // This next line is a potential security hole, so its important that this code NOT be loaded by the release version.
            window.opener.postMessage({formstate: $('#Goradd__FormState').val()}, "*");
            goradd.getForm().addEventListener ('teststep', goradd.testStep);
        }
    };

    goradd._testStepPending= false;

    goradd.testStep = function(event) {
        if (goradd.actionQueue.length > 0) {
            goradd.queueAction({f: function() {
                    goradd._postTestStep(event);
                }, last: true, name: "testStep"});
        } else {
            goradd._postTestStep(event);
        }
    };


    goradd._postTestStep = function(event) {
        if (event) {
            if (!goradd.ajaxQueueIsRunning()) {
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


}) ();