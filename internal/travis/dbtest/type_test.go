package dbtest

import (
	"context"
	"github.com/goradd/goradd/pkg/datetime"
	"github.com/stretchr/testify/assert"
	"goradd-project/gen/goradd/model"
	"testing"
	"time"
)

func TestDateTimeType(t *testing.T) {
	ctx := context.Background()

	rec := model.LoadTypeTest(ctx, "1")

	assert.True(t, rec.Date().Equal(datetime.DateOnly(2019,datetime.January, 2)))
	assert.True(t, rec.Time().Equal(datetime.Time(6,17,28,0)))

	dt := rec.DateTime()
	assert.False(t, dt.IsTimestamp())

	ts := rec.Ts()
	assert.True(t, ts.IsTimestamp())

	// Make sure a timestamp that goes into the database in a different timezone will not change
	loc,err := time.LoadLocation("America/New_York")
	if err != nil {
		panic("Cannot load location") // might need to turn this off on windows
	}
	ts2 := datetime.Date(2002,datetime.July, 2, 10,4,3,0,loc)
	rec.SetTs(ts2)
	rec.Save(ctx)

	rec2 := model.LoadTypeTest(ctx, "1")
	ts3 := rec2.Ts()
	assert.True(t, ts3.IsTimestamp())
	assert.True(t, ts3.Equal(ts2))

	rec2.SetTs(ts)
	rec2.Save(ctx)
}

