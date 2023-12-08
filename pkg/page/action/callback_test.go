package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAjax(t *testing.T) {
	js := Do().ControlID("a").ID(2).RenderScript(RenderParams{
		TriggeringControlID: "b",
		ControlActionValue:  1,
		EventID:             2,
		EventActionValue:    "c",
	})
	assert.Equal(t, `goradd.postAjax({"controlID":"b","eventId":2,"actionValues":{"event":"c","control":1}});`, js)

	js = Do().ControlID("a").ID(2).ActionValue(3).Async().RenderScript(RenderParams{})
	assert.Equal(t, `goradd.postAjax({"controlID":"","eventId":0,"async":true,"actionValues":{"event":eventData,"action":3}});`, js)
}

func TestServer(t *testing.T) {
	js := Do().ControlID("a").ID(2).Post().RenderScript(RenderParams{
		TriggeringControlID: "b",
		ControlActionValue:  1,
		EventID:             2,
		EventActionValue:    "c",
	})
	assert.Equal(t, `goradd.postBack({"controlID":"b","eventId":2,"actionValues":{"event":"c","control":1}});`, js)

	js = Do().ControlID("a").ID(2).Post().ActionValue(3).Async().RenderScript(RenderParams{})
	assert.Equal(t, `goradd.postBack({"controlID":"","eventId":0,"async":true,"actionValues":{"event":eventData,"action":3}});`, js)
}
