package gone

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMacAddr(t *testing.T) {
	if _, err := MacAddr(); err != nil {
		t.Error(err)
		return
	}
}

func TestFormatBytesString(t *testing.T) {
	t.Log(KB, MB, GB, TB, PB, EB)
	t.Log(kb, mb, gb, tb, pb, eb)
	var dir string
	dir, _ = os.Getwd()
	dir = filepath.Join(dir, "tmp")
	var tests = []struct {
		Name string
		Size int64
		Res  string
	}{
		{Name: "900b", Size: 900 * Bytes, Res: "900 bytes"},
		{Name: "1020b", Size: 1020 * Bytes, Res: "1.02 KB"},
		{Name: "900kb", Size: 900 * KB, Res: "921.6 KB"},
		{Name: "1000kb", Size: 1000 * KB, Res: "1.02 MB"},
		{Name: "1029kb", Size: 1029 * KB, Res: "1.05 MB"},
		{Name: "900MB", Size: 900 * MB, Res: "943.72 MB"},
		{Name: "1000MB", Size: 1000 * MB, Res: "1.05 GB"},
		{Name: "1020MB", Size: 1020 * MB, Res: "1.07 GB"},
	}
	for i := range tests {
		b := make([]byte, tests[i].Size)
		_, err := rand.Read(b)
		if err != nil {
			t.Error(err)
		}
		s := FormatBytesString(tests[i].Size)
		if s != tests[i].Res {
			t.Errorf("%s format err : get %s ,but want %s", tests[i].Name, s, tests[i].Res)
		}
		//t.Run(tests[i].Name, func(t *testing.T) {
		//	temp, err := os.CreateTemp(dir, r2)
		//	if err != nil {
		//		t.Error(err)
		//	}
		//	t.Log(tests[i].Name, r)
		//	temp.Write(b)
		//	defer temp.Close()
		//})
	}
}
func BenchmarkFormatBytesString(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = Bytes + rng.Float64()*(EB-Bytes)
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		FormatBytesString(int64(in[i]))
	}
}
