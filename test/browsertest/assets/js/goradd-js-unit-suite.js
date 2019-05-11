/**
 * testsuite contains the unit tests that test the goradd.js javascript file, in conjunction with the jsunit_form go
 * form. Tests start with "test".
 */
goradd.testsuite = {
    /**
     * Reset the test space to its original state
     * @param t
     */
    pretest: function(t) {
        goradd.el("testspace").innerHTML =  '<p id="testP">I am here</p> \
        <div id="testD" data-animal-type="bird" spellcheck>a div</div>';
    },
    testEl: function(t) {
       var el = goradd.el("testP");
       t.assert(el.innerText === "I am here");
    },
    testQS: function(t) {
        var el = goradd.qs("p[id='testP']");
        t.assert(el.innerText === "I am here");
    },
    testQA: function(t) {
        var el = goradd.qa("p[id='testP']");
        t.assert(el[0].innerText === "I am here");
    },
    testIsEmptyObj: function(t) {
        t.assert(!goradd.isEmptyObj({a:"b"}));
        t.assert(goradd.isEmptyObj({}));
    },
    testForm: function(t) {
        t.isSame("JsUnitTestForm", goradd.form().id);
    },
    testMatches: function(t) {
        t.assert(goradd.matches("testP", "p[id='testP']"));
        t.assert(!goradd.matches("testP", "div"));
    },
    testParents: function(t) {
        var p = goradd.parents("testP");
        t.isSame("JsUnitTestForm", p[1].id);
    },
    testAttr: function(t) {
        var a1 = goradd.attr("testD", "spellcheck");
        t.isSame(true, a1);
        goradd.attr("testD", {"spellcheck":null, "class":"a"});
        t.isSame(null, goradd.attr("testD", "spellcheck"));
        t.isSame("a", goradd.attr("testD", "class"));
        goradd.attr("testD", "class", "b");
        t.isSame("b", goradd.attr("testD", "class"));
    },
    testProp: function(t) {
        t.isSame("testP", goradd.prop("testP", "id"))
    },
    testEvent: function(t) {
        goradd.on("testD", "et", function() {
            this.innerText = "tested";
        });
        goradd.trigger("testD", "et");
        t.isSame("tested", goradd.el("testD").innerText)
    },
    testHtmlInserts: function(t) {
        goradd.htmlAfter("testP", "<p id='after'>Inserted After</p>");
        goradd.htmlBefore("testP", "<p id='before'>Inserted Before</p>");
        t.isSame("Inserted After", goradd.el("after").innerText);
        t.isSame("Inserted Before", goradd.el("before").innerText);
        goradd.el("testP").innerText = "There";
        goradd.insertHtml("testP", "Here");
        goradd.appendHtml("testP", "Everywhere");
        t.isSame("HereThereEverywhere", goradd.el("testP").innerText);
    },
    testRemove: function(t) {
        goradd.remove("testP");
        t.isSame(goradd.el("testP"), null);
    },
    testEach: function(t) {
        var s = "";
        goradd.each(goradd.qa("testspace", "p,div"), function(i, v) {
            s += v.innerText;
        });
        t.isSame("I am herea div", s);

        s = "";
        goradd.each(["a", "b", "c"], function(i,v) {
            s += v;
        });
        t.isSame("abc", s);

        s = [];
        goradd.each({a:"a", b:"b", c:"c"}, function(k,v) {
            s.push(k+v);
        });
        t.assert(s.indexOf("aa") !== -1);
    },
    testToSnake: function(t) {
        t.isSame("this-is-me", goradd._toKebab("thisIsMe"));
        t.isSame("a-b-c", goradd._toKebab("aBC"));
    },
    testData: function(t) {
        t.isSame("bird", goradd.data("testD", "animalType"));
        goradd.data("testD", "animalType", "dog");
        t.isSame("dog", goradd.data("testD", "animalType"));
    }
};
