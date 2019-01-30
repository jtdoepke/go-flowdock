package flowdock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	flowdockTime := Time{}
	json := []byte("1385546251160")
	err := flowdockTime.UnmarshalJSON(json)
	assert.NoError(t, err, "Time.UnmarshalJSON returned error: %v", err)

	want := time.Date(2013, time.November, 27, 9, 57, 31, 0, time.UTC)
	assert.Equal(t,
		want.Local(), flowdockTime.Local(),
		"Time.UnmarshalJSON set time to %v, wanted %v", flowdockTime.Local(), want.Local(),
	)
}
