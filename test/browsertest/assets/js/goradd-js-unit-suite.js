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

        el = goradd.g("testspace").qs("p[id='testP']");
        t.assert(el.innerText === "I am here");
    },
    testQA: function(t) {
        var el = goradd.qa("p[id='testP']");
        t.assert(el[0].innerText === "I am here");

        el = goradd.g("testspace").qa("p[id='testP']");
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
        t.assert(goradd.g("testP").matches("p[id='testP']"));
        t.assert(!goradd.g("testP").matches("div"));
    },
    testParents: function(t) {
        var p = goradd.g("testP").parents();
        t.isSame("JsUnitTestForm", p[1].id);
    },
    testAttrProp: function(t) {
        var a1 = goradd.g("testD").attr("spellcheck");
        t.isSame(true, a1);
        goradd.g("testD").prop({"spellcheck":false, "class":"a"});
        t.isSame(false, goradd.g("testD").prop("spellcheck"));
        t.isSame("a", goradd.g("testD").prop("class"));
        goradd.g("testD").prop("class", "b");
        t.isSame("b", goradd.g("testD").prop("class"));
    },
    testClass: function(t) {
        goradd.g("testD").class("b c");
        goradd.g("testD").class("-c");
        t.isSame("b", goradd.g("testD").class());

        goradd.g("testD").class("+a");
        goradd.g("testD").class("-b");
        t.isSame("a", goradd.g("testD").class());
    },
    testEvent: function(t) {
        goradd.g("testD").on("et", function() {
            this.innerText = "tested";
        });
        goradd.g("testD").trigger("et");
        t.isSame("tested", goradd.el("testD").innerText)
    },
    testHtmlInserts: function(t) {
        var p = goradd.g("testP");
        p.htmlAfter("<p id='after'>Inserted After</p>");
        p.htmlBefore("<p id='before'>Inserted Before</p>");
        t.isSame("Inserted After", goradd.el("after").innerText);
        t.isSame("Inserted Before", goradd.el("before").innerText);
        goradd.el("testP").innerText = "There";
        p.insertHtml("Here");
        p.appendHtml("Everywhere");
        t.isSame("HereThereEverywhere", goradd.el("testP").innerText);
    },
    testRemove: function(t) {
        goradd.g("testP").remove();
        t.isSame(goradd.el("testP"), null);
    },
    testEach: function(t) {
        var s = "";
        goradd.each(goradd.g("testspace").qa("p,div"), function(i, v) {
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
        var d = goradd.g("testD");
        t.isSame("bird", d.data("animalType"));
        d.data("animalType", "dog");
        t.isSame("dog", d.data("animalType"));
    }
};
