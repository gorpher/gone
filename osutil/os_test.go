package osutil

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
	//t.Log(KiB, MiB, GiB, TiB, PiB, EiB)
	//t.Log(KB, MB, GB, TB, PB, EB)
	var dir string
	dir, _ = os.Getwd()
	dir = filepath.Join(dir, "tmp")
	tests := []struct {
		Name string
		Size int64
		Res  string
	}{
		{Name: "900b", Size: 900, Res: "900 B"},
		{Name: "1020b", Size: 1020, Res: "1.02 KB"},
		{Name: "900kb", Size: 900 * KiB, Res: "921.6 KB"},
		{Name: "1000kb", Size: 1000 * KiB, Res: "1.02 MB"},
		{Name: "1029kb", Size: 1029 * KiB, Res: "1.05 MB"},
		{Name: "900MB", Size: 900 * MiB, Res: "943.72 MB"},
		{Name: "1000MB", Size: 1000 * MiB, Res: "1.05 GB"},
		{Name: "1020MB", Size: 1020 * MiB, Res: "1.07 GB"},
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
		//s2 := FormatBinaryDecimal(tests[i].Size)
		//if s2 != tests[i].Res {
		//	t.Errorf("%s format err : get %s ,but want %s", tests[i].Name, s2, tests[i].Res)
		//}
	}
}

func BenchmarkFormatBytesString(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = 1 + rng.Float64()*(EiB-1)
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		FormatBytesString(int64(in[i]))
	}
}
