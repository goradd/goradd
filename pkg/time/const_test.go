package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConst_Time(t *testing.T) {
	assert.True(t, Zero.Time().IsZero())
	assert.True(t, time.Since(Now.Time()).Seconds() == 0)
	assert.True(t, time.Since(Current.Time()).Seconds() == 0)
}
