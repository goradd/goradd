/**
 * testsuite contains the unit tests that test the goradd.js javascript file, in conjunction with the jsunit_form go
 * form. Tests start with "test".
 */
goradd.testsuite = {
    testEl: function(t) {
       var el = goradd.el("testP");
       t.assert(el.innerText === "I am here");
    },
    testQs: function(t) {
        var el = goradd.qs("p[id='testP']");
        t.assert(el.innerText === "I am here");
    }
};
