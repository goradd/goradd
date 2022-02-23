package dbtest

import (
	time2 "github.com/goradd/goradd/pkg/time"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goraddUnit/model"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	ctx := getContext()
	rec := model.LoadTypeTest(ctx, "1")
	assert.True(t, rec.Date().Equal(time2.NewDate(2019, 1, 2)))
}

func TestTime(t *testing.T) {
	ctx := getContext()
	rec := model.LoadTypeTest(ctx, "1")
	assert.True(t, rec.Time().Equal(time2.NewTime(6, 17, 28, 0)))
}

func TestTimestamp(t *testing.T) {
	ctx := getContext()
	rec := model.LoadTypeTest(ctx, "1")

	ts := rec.Ts()

	// Make sure a timestamp that goes into the database in a different timezone will be saved in UTC
	loc,err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Cannot load location") // might need to turn this off on windows
	}
	ts2 := time.Date(2002, time.July, 2, 10,4,3,0,loc)
	rec.SetTs(ts2)
	zone,_ := rec.Ts().Zone()
	assert.Equal(t, "EDT", zone)
	rec.Save(ctx)

	rec2 := model.LoadTypeTest(ctx, "1")
	ts3 := rec2.Ts()
	zone,_ = rec2.Ts().Zone()
	assert.Equal(t, "UTC", zone)

	assert.True(t, ts3.Equal(ts2))

	// restore database
	rec2.SetTs(ts)
	rec2.Save(ctx)

	// Make sure setting null value sets to default of current time
	rec.SetTs(nil)
	assert.True(t, time.Since(rec.Ts()).Seconds() == 0)
}

func TestDateTime(t *testing.T) {
	ctx := getContext()
	rec := model.LoadTypeTest(ctx, "1")
	assert.True(t, rec.Time().Equal(time2.NewTime(6, 17, 28, 0)))

	ts := rec.DateTime()

	// Make sure a timestamp that goes into the database in a different timezone will be saved in UTC
	loc,err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Cannot load location") // might need to turn this off on windows
	}
	ts2 := time.Date(2002, time.July, 2, 10,4,3,0,loc)
	rec.SetDateTime(ts2)
	zone,_ := rec.DateTime().Zone()
	assert.Equal(t, "EDT", zone)
	rec.Save(ctx)

	rec2 := model.LoadTypeTest(ctx, "1")
	ts3 := rec2.DateTime()
	zone,_ = rec2.DateTime().Zone()
	assert.Equal(t, "UTC", zone)

	assert.True(t, ts3.Equal(ts2))

	// restore database
	rec2.SetDateTime(ts)
	rec2.Save(ctx)

	// Make sure setting null value sets to default
	rec.SetDateTime(nil)
	assert.True(t, rec.DateTime().IsZero())

}




