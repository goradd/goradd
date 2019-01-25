package datetime

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFromSqlDateTime(t *testing.T) {
	d,err := FromSqlDateTime("2008-11-11 13:23:44")
	require.NoError(t, err)
	d2 := Date(2008, November, 11, 13, 23, 44, 0, time.UTC)
	assert.True(t, d.Equal(d2))
	d3 := Date(2008, November, 11, 13, 23, 44, 1, time.UTC)
	assert.False(t, d.Equal(d3))
	_,err = FromSqlDateTime("2008-11-11 13:23:44adf")
	assert.Error(t, err)
}

