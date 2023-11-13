/*

Goradd Debugging Additions

This file attaches some code to display ajax and javascript errors and progress.
It is only loaded in debug mode.

*/
(function() {

    if (window.console) {
        // This lets us easily turn off logging in production without losing any console.log capabilities
        goradd.log = console.log;
    }

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

        if (resultText.substring(0, 15) === "<!DOCTYPE html>") {
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