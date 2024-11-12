package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTimeMarshalJSON(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		tt   Time
		want int64
	}{
		{Time{}, 0},
		{Time{now}, now.Unix()},
		{Time{now.Add(24 * time.Hour)}, now.Add(24 * time.Hour).Unix()},
		{Time{now.Add(24 * 30 * 12 * time.Hour)}, now.Add(24 * 30 * 12 * time.Hour).Unix()},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(t *testing.T) {
			b, err := tc.tt.MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			var n int64
			if err = json.Unmarshal(b, &n); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.want, n; got != want {
				t.Errorf("Time.Marshal mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func TestTimeUnmarshalJSON(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		n     int64
		want  Time
		isNil bool
	}{
		{now.Unix(), Time{now}, false},
		{Epoch.Unix() - 0xDEAD, Time{Epoch}, false},
		{Epoch.Unix(), Time{Epoch}, false},
		{Epoch.Unix() + 0xDEAD, Time{Epoch.Add(0xDEAD * time.Second)}, false},
		{0, Time{}, true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("", func(t *testing.T) {
			var n *int64
			if !tc.isNil {
				n = &tc.n
			}
			b, err := json.Marshal(n)
			if err != nil {
				t.Fatal(err)
			}
			var tt Time
			if err = tt.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.want.Unix(), tt.Unix(); got != want {
				t.Errorf("Time.Unmarshal mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}
