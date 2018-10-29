/*

Goradd Testing Additions

This file attaches some code used by the test harness to drive browser-based tests. It is only loaded in debug mode.

*/

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
            }, last: true});
    } else {
        goradd._postTestStep(event);
    }
};


goradd._postTestStep = function(event) {
    if (event) {
        if (!goradd.ajaxQueueIsRunning()) {
            window.opener.postMessage({ajaxComplete: event.detail}, "*");
        } else {
            goradd._testStepPending = true;
            goradd.currentStep = event.detail;
        }
    } else {
        // We are being notified that an ajax action has completed
        if (goradd._testStepPending) {
            goradd._testStepPending = false;
            window.opener.postMessage({ajaxComplete: goradd.currentStep}, "*");
        }
    }
};
