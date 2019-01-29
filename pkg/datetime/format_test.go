package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromSqlDateTime2(t *testing.T) {
	d, err := FromSqlDateTime("2008-11-11 13:23:44")
	require.NoError(t, err)
	d2 := Date(2008, November, 11, 13, 23, 44, 0, time.UTC)
	assert.True(t, d.Equal(d2))
	d3 := Date(2008, November, 11, 13, 23, 44, 1, time.UTC)
	assert.False(t, d.Equal(d3))
	_, err = FromSqlDateTime("2008-11-11 13:23:44adf")
	assert.Error(t, err)
}

func TestFromSqlDateTime(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		wantT   DateTime
		wantErr bool
	}{
		{"Date only", "2008-11-11", NewDateTime("2008-11-11T0:00:00Z"), false},
		{"Time only", "15:02:03", NewDateTime("0000-01-01T15:02:03Z"), false},
		{"Date and Time", "2008-11-11 13:23:44", NewDateTime("2008-11-11T13:23:44Z"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotT, err := FromSqlDateTime(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromSqlDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !gotT.Equal(tt.wantT) {
				t.Errorf("FromSqlDateTime() = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}
