package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddClass(t *testing.T) {
	js := AddClass("a", "b").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').class("+b");`, js)
}

func TestBlur(t *testing.T) {
	js := Blur("a").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').blur();`, js)
}

func TestConfirm(t *testing.T) {
	js := Confirm("a").RenderScript(RenderParams{})
	assert.Equal(t, `if (!window.confirm("a")) return false;`, js)
}

func TestCss(t *testing.T) {
	js := Css("a", "b", "c").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').css("b","c");`, js)

	js = Css("a", "b", 2).RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').css("b",2);`, js)
}

func TestFocus(t *testing.T) {
	js := Focus("a").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').focus();`, js)
}

func TestGoraddFunction(t *testing.T) {
	js := GoraddFunction("a").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.a();`, js)
	js = GoraddFunction("a", "b", 3).RenderScript(RenderParams{})
	assert.Equal(t, `goradd.a("b",3);`, js)
}

func TestHide(t *testing.T) {
	js := Hide("a").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').hide();`, js)
}

func TestJavascript(t *testing.T) {
	js := Javascript("a+b=c").RenderScript(RenderParams{})
	assert.Equal(t, `a+b=c;`, js)
}

func TestMessage(t *testing.T) {
	js := Message("a").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.msg("a");`, js)
}

func TestRedirect(t *testing.T) {
	js := Redirect("http://a/b/c").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.redirect("http://a/b/c");`, js)
}

func TestRefresh(t *testing.T) {
	js := Refresh("a").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.refresh("a");`, js)
}

func TestRemoveClass(t *testing.T) {
	js := RemoveClass("a", "b").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').class("-b");`, js)
}

func TestSelect(t *testing.T) {
	js := Select("a").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').selectAll();`, js)
}

func TestSetControlValue(t *testing.T) {
	js := SetControlValue("a", "b", "c").RenderScript(RenderParams{})
	assert.Equal(t, `goradd.setControlValue("a","b","c");`, js)
}

func TestShow(t *testing.T) {
	js := Show("a").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').show();`, js)
}

func TestToggleClass(t *testing.T) {
	js := ToggleClass("a", "b").RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').toggleClass("b");`, js)
}

func TestTrigger(t *testing.T) {
	js := Trigger("a", "b", 2).RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').trigger("b",2);`, js)
}

func TestWidgetFunction(t *testing.T) {
	js := WidgetFunction("a", "b", 2).RenderScript(RenderParams{})
	assert.Equal(t, `g$('a').b(2);`, js)
}
