package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConst_Time(t *testing.T) {
	assert.True(t, Zero.Time().IsZero())
	assert.Greater(t, 1.0, time.Since(Now.Time()).Seconds())
	assert.Greater(t, 1.0, time.Since(Current.Time()).Seconds())
}
