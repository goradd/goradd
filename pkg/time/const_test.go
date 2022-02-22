package time

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConst_Time(t *testing.T) {
	assert.True(t, Zero.Time().IsZero())
	assert.Equal(t, Now.Time().Second(), time.Now().Second())
	assert.Equal(t, Current.Time().Second(), time.Now().Second())
}
