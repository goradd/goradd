/**
 * jsUnit is a simple javascript test runner. We use it to unit test the goradd.js code. You can use it for your
 * own unit tests too.
 *
 * @type {{run: goradd.jsUnit.run, _makeTestRunner(*=): *}}
 */
goradd.jsUnit = {
    run: function(suite, result) {
        result = goradd.el(result);
        var curTestInfo = this._makeTestRunner(result);
        for (t in suite) {
            if (t.substr(0, 4) === "test") {
                var f = suite[t];
                if (typeof f === "function") {
                    try {
                        if (suite.pretest) {
                            suite.pretest(t);
                        }
                        curTestInfo.testName = t;
                        f(curTestInfo);
                    } catch(e) {
                        var s = t + " threw an error. " + e + ". " + e.sourceURL + " " + e.line + ":" + e.column;
                        curTestInfo.displayError(s);
                    }
                }
            }
        }
        g$(result).appendHtml("Done");
    },
    _makeTestRunner(result) {
        return new goradd.JsTest(result);
    }
};

/**
 * JsTest is the object passed to each test suite that lets it call asserts and display errors.
 *
 * @param resultDiv {object|string}
 * @constructor
 */
goradd.JsTest = function(resultDiv) {
    this.result = goradd.el(resultDiv);
};
/**
 *
 * @type {{appendTo: (function((Object|string)): *), insertInto: (function((Object|string)): *), replace: (function((Object|string)): *), html: (function(string): goradd.TagBuilder), text: (function(string): goradd.TagBuilder), attr: (function(string, string): goradd.TagBuilder), insertAfter: (function((Object|string)): *), insertBefore: (function((Object|string)): *)}}
 */
goradd.JsTest.prototype = {
    testName:"",
    displayError: function(e) {
        g$(this.result).appendHtml(e);
    },
    assert: function(bVal, msg) {
        if (!msg) {
            msg = "";
        }
        if (!bVal) {
            this.displayError(this._makeError("assert", msg));
        }
    },
    isEqual: function(expected, v, msg) {
        if (!msg) {
            msg = "";
        }
        if (expected != v) {
            this.displayError(this._makeError("isEqual", "Expected: " + expected + " got: " + v + ". " + msg));
        }
    },
    isSame: function(expected, v, msg) {
        if (!msg) {
            msg = "";
        }
        if (expected !== v) {
            this.displayError(this._makeError("isSame", "Expected: " + expected + " got: " + v + ". " + msg));
        }
    },
    _makeError(title, msg) {
        if (!msg) {
            msg = "";
        }
        var s = this._getCallStack(2);
        return this.testName + ": " + title + " failed. " + msg + s[0];
    },
    _getCallStack (offset) {
        var e = new Error();
        if (!e.stack)
            try {
                // Some browsers require a thrown error to get a stack trace.
                throw e;
            } catch (e) {
            }
        var stack = e.stack.toString().split(/\r\n|\n/);
        for (var i = 0; i < offset + 1; i++) {
            stack.shift();
        }
        return stack;
    }
};
